package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/lunancy1992/jianghu-server/internal/ai"
	"github.com/lunancy1992/jianghu-server/internal/cache"
	"github.com/lunancy1992/jianghu-server/internal/config"
	"github.com/lunancy1992/jianghu-server/internal/crawler"
	"github.com/lunancy1992/jianghu-server/internal/crawler/adapter"
	cronpkg "github.com/lunancy1992/jianghu-server/internal/cron"
	"github.com/lunancy1992/jianghu-server/internal/handler"
	"github.com/lunancy1992/jianghu-server/internal/middleware"
	"github.com/lunancy1992/jianghu-server/internal/repo"
	"github.com/lunancy1992/jianghu-server/internal/service"
)

func main() {
	// Load config
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Ensure data directory exists
	if err := os.MkdirAll(cfg.Server.DataDir, 0755); err != nil {
		log.Fatalf("Failed to create data dir: %v", err)
	}

	// Init SQLite
	db, err := repo.InitDB(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := repo.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed")

	// Init cache
	appCache, err := cache.New(cfg.Cache.NumCounters, cfg.Cache.MaxCost)
	if err != nil {
		log.Fatalf("Failed to init cache: %v", err)
	}
	defer appCache.Close()

	// Init repos
	userRepo := repo.NewUserRepo(db)
	newsRepo := repo.NewNewsRepo(db)
	commentRepo := repo.NewCommentRepo(db)
	coinRepo := repo.NewCoinRepo(db)
	eventRepo := repo.NewEventRepo(db)

	// Init AI client
	var aiClient ai.AIClient
	if claudeKey := os.Getenv("ANTHROPIC_API_KEY"); claudeKey != "" {
		aiClient = ai.NewClaudeClient(claudeKey, os.Getenv("ANTHROPIC_MODEL"))
		log.Println("Claude AI client initialized")
	} else if cfg.AI.DeepSeekAPIKey != "" {
		aiClient = ai.NewDeepSeekClient(cfg.AI.DeepSeekAPIKey, cfg.AI.DeepSeekBase, cfg.AI.Model)
		log.Println("DeepSeek AI client initialized")
	} else {
		log.Println("Warning: No AI API key set, AI features disabled")
	}

	// Init services
	auditService := service.NewAuditService(aiClient, commentRepo)
	authService := service.NewAuthService(userRepo, coinRepo, cfg)
	newsService := service.NewNewsService(newsRepo, appCache)
	commentService := service.NewCommentService(commentRepo, coinRepo, newsRepo, auditService)
	coinService := service.NewCoinService(coinRepo, commentRepo)
	eventService := service.NewEventService(eventRepo)

	// Init crawler
	newsCrawler := crawler.NewCrawler(newsRepo, aiClient)
	for _, feed := range cfg.Crawl.Feeds {
		newsCrawler.AddAdapter(adapter.NewRSSAdapter(feed.Name, feed.URL, feed.Enabled))
	}

	// Init handlers
	authHandler := handler.NewAuthHandler(authService)
	newsHandler := handler.NewNewsHandler(newsService)
	commentHandler := handler.NewCommentHandler(commentService)
	coinHandler := handler.NewCoinHandler(coinService)
	eventHandler := handler.NewEventHandler(eventService)
	adminHandler := handler.NewAdminHandler(auditService, newsService, eventService)
	evidenceHandler := handler.NewEvidenceHandler(eventService, cfg.Server.DataDir)

	// Start cron jobs
	scheduler := cronpkg.NewScheduler()

	// News crawl every 30 minutes
	_ = scheduler.AddFunc("0 */30 * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		newsCrawler.Run(ctx)
	})

	// Weekly coin distribution: Monday 00:00
	_ = scheduler.AddFunc("0 0 0 * * 1", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		if err := coinService.WeeklyDistribute(ctx); err != nil {
			log.Printf("[Cron] Weekly distribution error: %v", err)
		}
	})

	// Comment rewards hourly
	_ = scheduler.AddFunc("0 0 * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		if err := coinService.ProcessCommentRewards(ctx); err != nil {
			log.Printf("[Cron] Comment rewards error: %v", err)
		}
	})

	// DB backup every 6 hours
	_ = scheduler.AddFunc("0 0 */6 * * *", func() {
		backupPath := filepath.Join(cfg.Server.DataDir, fmt.Sprintf("backup_%s.db", time.Now().Format("20060102_150405")))
		if _, err := db.Exec("VACUUM INTO ?", backupPath); err != nil {
			log.Printf("[Cron] DB backup error: %v", err)
		} else {
			log.Printf("[Cron] DB backed up to %s", backupPath)
		}
	})

	scheduler.Start()

	// Setup Gin router
	router := gin.Default()

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Global rate limit
	router.Use(middleware.RateLimit(rate.Limit(100), 200))

	// Serve uploaded files
	router.Static("/uploads", filepath.Join(cfg.Server.DataDir, "uploads"))

	// API v1
	v1 := router.Group("/api/v1")

	// Auth routes
	auth := v1.Group("/auth")
	{
		auth.POST("/sms/send", middleware.SMSRateLimit(), authHandler.SendSMS)
		auth.POST("/sms/login", authHandler.LoginWithSMS)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// User routes (authenticated)
	user := v1.Group("/user")
	user.Use(middleware.JWTAuth(cfg.Auth.JWTSecret))
	{
		user.GET("/profile", authHandler.GetProfile)
	}

	// Headlines (public)
	v1.GET("/headlines", newsHandler.GetHeadlines)
	v1.GET("/headlines/history", newsHandler.GetHeadlineHistory)

	// News routes (public with optional auth)
	news := v1.Group("/news")
	news.Use(middleware.OptionalAuth(cfg.Auth.JWTSecret))
	{
		news.GET("", newsHandler.ListNews)
		news.GET("/search", newsHandler.SearchNews)
		news.GET("/:id", newsHandler.GetNews)
		news.GET("/:id/comments", commentHandler.ListComments)
	}

	// News actions (authenticated)
	newsAuth := v1.Group("/news")
	newsAuth.Use(middleware.JWTAuth(cfg.Auth.JWTSecret))
	{
		newsAuth.POST("/:id/read", newsHandler.MarkRead)
		newsAuth.POST("/:id/comments", commentHandler.CreateComment)
	}

	// Like routes (authenticated)
	v1.POST("/like", middleware.JWTAuth(cfg.Auth.JWTSecret), commentHandler.Like)
	v1.DELETE("/like", middleware.JWTAuth(cfg.Auth.JWTSecret), commentHandler.Unlike)

	// Coin routes (authenticated)
	coin := v1.Group("/coin")
	coin.Use(middleware.JWTAuth(cfg.Auth.JWTSecret))
	{
		coin.GET("/balance", coinHandler.GetBalance)
		coin.GET("/transactions", coinHandler.GetTransactions)
	}

	// Event routes (public)
	events := v1.Group("/events")
	events.Use(middleware.OptionalAuth(cfg.Auth.JWTSecret))
	{
		events.GET("", eventHandler.ListEvents)
		events.GET("/:id", eventHandler.GetEvent)
	}

	// Evidence routes (authenticated)
	v1.POST("/evidence", middleware.JWTAuth(cfg.Auth.JWTSecret), evidenceHandler.Upload)
	v1.POST("/events/:id/evidence", middleware.JWTAuth(cfg.Auth.JWTSecret), evidenceHandler.AddEvidence)

	// Admin routes
	admin := v1.Group("/admin")
	admin.Use(middleware.JWTAuth(cfg.Auth.JWTSecret), middleware.RequireAdmin())
	{
		admin.GET("/audit/queue", adminHandler.GetAuditQueue)
		admin.POST("/audit/:id/approve", adminHandler.ApproveComment)
		admin.POST("/audit/:id/reject", adminHandler.RejectComment)
		admin.POST("/headlines", adminHandler.SetHeadlines)
		admin.POST("/events", adminHandler.CreateEvent)
		admin.POST("/events/:id/nodes", adminHandler.CreateEventNode)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	scheduler.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

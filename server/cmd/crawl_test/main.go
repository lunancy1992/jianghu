package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/lunancy1992/jianghu-server/internal/ai"
	"github.com/lunancy1992/jianghu-server/internal/config"
	"github.com/lunancy1992/jianghu-server/internal/crawler"
	"github.com/lunancy1992/jianghu-server/internal/crawler/adapter"
	"github.com/lunancy1992/jianghu-server/internal/repo"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := repo.InitDB(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	defer db.Close()

	if err := repo.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	newsRepo := repo.NewNewsRepo(db)

	var aiClient ai.AIClient
	if claudeKey := os.Getenv("ANTHROPIC_API_KEY"); claudeKey != "" {
		aiClient = ai.NewClaudeClient(claudeKey, os.Getenv("ANTHROPIC_MODEL"))
		log.Println("Claude AI client initialized")
	} else if cfg.AI.DeepSeekAPIKey != "" {
		aiClient = ai.NewDeepSeekClient(cfg.AI.DeepSeekAPIKey, cfg.AI.DeepSeekBase, cfg.AI.Model)
		log.Println("DeepSeek AI client initialized")
	} else {
		log.Println("Warning: No AI key set, AI features disabled")
	}

	c := crawler.NewCrawler(newsRepo, aiClient)
	for _, feed := range cfg.Crawl.Feeds {
		log.Printf("Adding feed: %s (enabled=%v)", feed.Name, feed.Enabled)
		c.AddAdapter(adapter.NewRSSAdapter(feed.Name, feed.URL, feed.Enabled))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	log.Println("Starting crawl...")
	c.Run(ctx)
	log.Println("Crawl done!")
}

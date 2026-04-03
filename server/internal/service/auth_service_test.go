package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/lunancy1992/jianghu-server/internal/config"
	"github.com/lunancy1992/jianghu-server/internal/repo"
)

func newServiceTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:?_foreign_keys=ON&_loc=auto")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	if err := repo.RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}
	return db
}

func validFutureTime() time.Time {
	return time.Now().Add(5 * time.Minute)
}

func TestAuthService_LoginAndRefresh(t *testing.T) {
	db := newServiceTestDB(t)
	userRepo := repo.NewUserRepo(db)
	coinRepo := repo.NewCoinRepo(db)
	cfg := &config.Config{Auth: config.AuthConfig{JWTSecret: "test-secret", JWTExpireHours: 24}}
	service := NewAuthService(userRepo, coinRepo, cfg)
	ctx := context.Background()

	if err := userRepo.SaveVerification(ctx, "13800138000", "123456", validFutureTime()); err != nil {
		t.Fatalf("SaveVerification() error: %v", err)
	}

	token, user, err := service.LoginWithSMS(ctx, "13800138000", "123456")
	if err != nil {
		t.Fatalf("LoginWithSMS() error: %v", err)
	}
	if token == "" {
		t.Fatal("LoginWithSMS() returned empty token")
	}
	if user == nil {
		t.Fatal("LoginWithSMS() returned nil user")
	}
	if user.CoinBalance != 10 {
		t.Fatalf("CoinBalance = %d, want 10", user.CoinBalance)
	}
	if user.Phone != "13800138000" {
		t.Fatalf("Phone = %q, want 13800138000", user.Phone)
	}

	newToken, err := service.RefreshToken(ctx, token)
	if err != nil {
		t.Fatalf("RefreshToken() error: %v", err)
	}
	if newToken == "" {
		t.Fatal("RefreshToken() returned empty token")
	}

	profile, err := service.GetProfile(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetProfile() error: %v", err)
	}
	if profile == nil {
		t.Fatal("GetProfile() returned nil")
	}
	if profile.CoinBalance != 10 {
		t.Fatalf("GetProfile().CoinBalance = %d, want 10", profile.CoinBalance)
	}
}

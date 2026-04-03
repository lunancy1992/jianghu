package repo

import (
	"context"
	"testing"
	"time"

	"github.com/lunancy1992/jianghu-server/internal/model"
)

func setupCoinTest(t *testing.T) (*CoinRepo, int64) {
	t.Helper()
	db := newTestDB(t)
	ur := NewUserRepo(db)
	cr := NewCoinRepo(db)
	ctx := context.Background()

	userID := createTestUser(t, ur)
	err := cr.CreateAccount(ctx, userID, 10)
	if err != nil {
		t.Fatalf("CreateAccount() error: %v", err)
	}
	return cr, userID
}

func TestCoinRepo_CreateAccountAndGetAccount(t *testing.T) {
	cr, userID := setupCoinTest(t)

	acc, err := cr.GetAccount(context.Background(), userID)
	if err != nil {
		t.Fatalf("GetAccount() error: %v", err)
	}
	if acc == nil {
		t.Fatal("GetAccount() returned nil")
	}
	if acc.Balance != 10 {
		t.Errorf("Balance = %d, want 10", acc.Balance)
	}
}

func TestCoinRepo_GetAccount_NotFound(t *testing.T) {
	db := newTestDB(t)
	cr := NewCoinRepo(db)

	acc, err := cr.GetAccount(context.Background(), 9999)
	if err != nil {
		t.Fatalf("GetAccount() error: %v", err)
	}
	if acc != nil {
		t.Error("GetAccount() should return nil for non-existent user")
	}
}

func TestCoinRepo_Deduct(t *testing.T) {
	cr, userID := setupCoinTest(t)
	ctx := context.Background()

	tx, err := cr.BeginTx(ctx)
	if err != nil {
		t.Fatalf("BeginTx() error: %v", err)
	}

	err = cr.Deduct(ctx, tx, userID, 3, "comment_cost", "comment:1", "发表评论")
	if err != nil {
		tx.Rollback()
		t.Fatalf("Deduct() error: %v", err)
	}
	tx.Commit()

	acc, _ := cr.GetAccount(ctx, userID)
	if acc.Balance != 7 {
		t.Errorf("Balance = %d, want 7 (10 - 3)", acc.Balance)
	}
}

func TestCoinRepo_Deduct_InsufficientBalance(t *testing.T) {
	cr, userID := setupCoinTest(t)
	ctx := context.Background()

	tx, _ := cr.BeginTx(ctx)
	err := cr.Deduct(ctx, tx, userID, 100, "comment_cost", "", "should fail")
	if err == nil {
		tx.Rollback()
		t.Fatal("Deduct() should return error for insufficient balance")
	}
	tx.Rollback()

	// Balance should be unchanged
	acc, _ := cr.GetAccount(ctx, userID)
	if acc.Balance != 10 {
		t.Errorf("Balance = %d, want 10 (unchanged)", acc.Balance)
	}
}

func TestCoinRepo_Credit(t *testing.T) {
	cr, userID := setupCoinTest(t)
	ctx := context.Background()

	tx, _ := cr.BeginTx(ctx)
	err := cr.Credit(ctx, tx, userID, 5, "weekly_grant", "", "每周奖励")
	if err != nil {
		tx.Rollback()
		t.Fatalf("Credit() error: %v", err)
	}
	tx.Commit()

	acc, _ := cr.GetAccount(ctx, userID)
	if acc.Balance != 15 {
		t.Errorf("Balance = %d, want 15 (10 + 5)", acc.Balance)
	}
}

func TestCoinRepo_ListTransactions(t *testing.T) {
	cr, userID := setupCoinTest(t)
	ctx := context.Background()

	// Make some transactions
	tx1, _ := cr.BeginTx(ctx)
	cr.Credit(ctx, tx1, userID, 5, "weekly_grant", "", "奖励1")
	tx1.Commit()

	tx2, _ := cr.BeginTx(ctx)
	cr.Deduct(ctx, tx2, userID, 1, "comment_cost", "c:1", "评论")
	tx2.Commit()

	list, total, err := cr.ListTransactions(ctx, userID, 1, 10)
	if err != nil {
		t.Fatalf("ListTransactions() error: %v", err)
	}
	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
	if len(list) != 2 {
		t.Errorf("len(list) = %d, want 2", len(list))
	}
}

func TestCoinRepo_GetActiveUsers(t *testing.T) {
	db := newTestDB(t)
	ur := NewUserRepo(db)
	nr := NewNewsRepo(db)
	cr := NewCoinRepo(db)
	cmr := NewCommentRepo(db)
	ctx := context.Background()

	userID := createTestUser(t, ur)
	newsID := createTestNews(t, nr, "活跃用户测试")

	// Create a comment to mark user as active
	cmr.Create(ctx, &model.Comment{
		NewsID: newsID, UserID: userID, Content: "活跃", Stance: "neutral", Status: 1,
	})

	since := time.Now().Add(-1 * time.Hour)
	activeUsers, err := cr.GetActiveUsers(ctx, since)
	if err != nil {
		t.Fatalf("GetActiveUsers() error: %v", err)
	}
	if len(activeUsers) != 1 {
		t.Fatalf("len(activeUsers) = %d, want 1", len(activeUsers))
	}
	if activeUsers[0] != userID {
		t.Errorf("activeUsers[0] = %d, want %d", activeUsers[0], userID)
	}
}

func TestCoinRepo_CreateAccount_Idempotent(t *testing.T) {
	cr, userID := setupCoinTest(t)
	ctx := context.Background()

	// Create again should not error (INSERT OR IGNORE)
	err := cr.CreateAccount(ctx, userID, 100)
	if err != nil {
		t.Fatalf("CreateAccount() duplicate error: %v", err)
	}

	// Balance should still be original 10
	acc, _ := cr.GetAccount(ctx, userID)
	if acc.Balance != 10 {
		t.Errorf("Balance = %d, want 10 (should not change on duplicate insert)", acc.Balance)
	}
}

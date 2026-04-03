package repo

import (
	"context"
	"testing"
	"time"

	"github.com/lunancy1992/jianghu-server/internal/model"
)

func setupCommentTest(t *testing.T) (*CommentRepo, int64, int64) {
	t.Helper()
	db := newTestDB(t)
	ur := NewUserRepo(db)
	nr := NewNewsRepo(db)
	cr := NewCommentRepo(db)

	userID := createTestUser(t, ur)
	newsID := createTestNews(t, nr, "评论测试新闻")

	return cr, userID, newsID
}

func TestCommentRepo_CreateAndFindByID(t *testing.T) {
	cr, userID, newsID := setupCommentTest(t)
	ctx := context.Background()

	comment := &model.Comment{
		NewsID:  newsID,
		UserID:  userID,
		Content: "这是一条测试评论",
		Stance:  "neutral",
		Status:  1,
	}
	id, err := cr.Create(ctx, comment)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	found, err := cr.FindByID(ctx, id)
	if err != nil {
		t.Fatalf("FindByID() error: %v", err)
	}
	if found == nil {
		t.Fatal("FindByID() returned nil")
	}
	if found.Content != "这是一条测试评论" {
		t.Errorf("Content = %q, want %q", found.Content, "这是一条测试评论")
	}
	if found.Status != 1 {
		t.Errorf("Status = %d, want 1", found.Status)
	}
}

func TestCommentRepo_ListByNews(t *testing.T) {
	cr, userID, newsID := setupCommentTest(t)
	ctx := context.Background()

	// Create approved and pending comments
	cr.Create(ctx, &model.Comment{NewsID: newsID, UserID: userID, Content: "approved", Stance: "agree", Status: 1})
	cr.Create(ctx, &model.Comment{NewsID: newsID, UserID: userID, Content: "pending", Stance: "neutral", Status: 0})

	list, total, err := cr.ListByNews(ctx, newsID, 1, 10)
	if err != nil {
		t.Fatalf("ListByNews() error: %v", err)
	}
	// Only approved comments
	if total != 1 {
		t.Errorf("total = %d, want 1 (only approved)", total)
	}
	if len(list) != 1 {
		t.Errorf("len(list) = %d, want 1", len(list))
	}
}

func TestCommentRepo_UpdateStatus(t *testing.T) {
	cr, userID, newsID := setupCommentTest(t)
	ctx := context.Background()

	id, err := cr.Create(ctx, &model.Comment{NewsID: newsID, UserID: userID, Content: "待审", Stance: "neutral", Status: 0})
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	err = cr.UpdateStatus(ctx, id, 1, "审核通过")
	if err != nil {
		t.Fatalf("UpdateStatus() error: %v", err)
	}

	found, _ := cr.FindByID(ctx, id)
	if found.Status != 1 {
		t.Errorf("Status = %d, want 1", found.Status)
	}
	if found.AuditNote != "审核通过" {
		t.Errorf("AuditNote = %q, want %q", found.AuditNote, "审核通过")
	}
}

func TestCommentRepo_ListPending(t *testing.T) {
	cr, userID, newsID := setupCommentTest(t)
	ctx := context.Background()

	cr.Create(ctx, &model.Comment{NewsID: newsID, UserID: userID, Content: "pending1", Stance: "neutral", Status: 0})
	cr.Create(ctx, &model.Comment{NewsID: newsID, UserID: userID, Content: "pending2", Stance: "neutral", Status: 0})
	cr.Create(ctx, &model.Comment{NewsID: newsID, UserID: userID, Content: "approved", Stance: "neutral", Status: 1})

	list, total, err := cr.ListPending(ctx, 1, 10)
	if err != nil {
		t.Fatalf("ListPending() error: %v", err)
	}
	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
	if len(list) != 2 {
		t.Errorf("len(list) = %d, want 2", len(list))
	}
}

func TestCommentRepo_LikeAndUnlike(t *testing.T) {
	cr, userID, newsID := setupCommentTest(t)
	ctx := context.Background()

	commentID, _ := cr.Create(ctx, &model.Comment{NewsID: newsID, UserID: userID, Content: "点赞测试", Stance: "neutral", Status: 1})

	// Like
	err := cr.CreateLike(ctx, userID, commentID)
	if err != nil {
		t.Fatalf("CreateLike() error: %v", err)
	}

	// Duplicate like should not increment again
	err = cr.CreateLike(ctx, userID, commentID)
	if err != nil {
		t.Fatalf("CreateLike() duplicate error: %v", err)
	}

	found, err := cr.FindByID(ctx, commentID)
	if err != nil {
		t.Fatalf("FindByID() after like error: %v", err)
	}
	if found.LikeCount != 1 {
		t.Errorf("LikeCount = %d, want 1", found.LikeCount)
	}

	// Unlike
	err = cr.DeleteLike(ctx, userID, commentID)
	if err != nil {
		t.Fatalf("DeleteLike() error: %v", err)
	}

	found, err = cr.FindByID(ctx, commentID)
	if err != nil {
		t.Fatalf("FindByID() after unlike error: %v", err)
	}
	if found.LikeCount != 0 {
		t.Errorf("LikeCount = %d, want 0 after unlike", found.LikeCount)
	}
}

func TestCommentRepo_FindTopLiked(t *testing.T) {
	cr, userID, newsID := setupCommentTest(t)
	ctx := context.Background()

	// Create a comment with enough likes
	commentID, _ := cr.Create(ctx, &model.Comment{NewsID: newsID, UserID: userID, Content: "热评", Stance: "neutral", Status: 1})
	// Manually set like_count to 15
	cr.IncrLikeCount(ctx, commentID, 15)

	since := time.Now().Add(-1 * time.Hour)
	top, err := cr.FindTopLiked(ctx, since, 10)
	if err != nil {
		t.Fatalf("FindTopLiked() error: %v", err)
	}
	if len(top) != 1 {
		t.Errorf("len(top) = %d, want 1", len(top))
	}
}

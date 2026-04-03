package repo

import (
	"context"
	"testing"
	"time"

	"github.com/lunancy1992/jianghu-server/internal/model"
)

func createTestUser(t *testing.T, r *UserRepo) int64 {
	t.Helper()
	id, err := r.Create(context.Background(), &model.User{
		Phone: "13800138000", Nickname: "test", Role: "user",
	})
	if err != nil {
		t.Fatalf("create test user: %v", err)
	}
	return id
}

func createTestNews(t *testing.T, r *NewsRepo, title string) int64 {
	t.Helper()
	now := time.Now()
	id, err := r.Create(context.Background(), &model.News{
		Title:       title,
		Summary:     "summary of " + title,
		Content:     "content",
		Source:      "test-source",
		Category:    "武林",
		Section:     "江湖快报",
		Keywords:    "[]",
		Sensitivity: 0,
		Fingerprint: "fp-" + title,
		Status:      1,
		PublishedAt: &now,
	})
	if err != nil {
		t.Fatalf("create test news: %v", err)
	}
	return id
}

func TestNewsRepo_CreateAndFindByID(t *testing.T) {
	db := newTestDB(t)
	r := NewNewsRepo(db)

	id := createTestNews(t, r, "武林大会召开")

	n, err := r.FindByID(context.Background(), id)
	if err != nil {
		t.Fatalf("FindByID() error: %v", err)
	}
	if n == nil {
		t.Fatal("FindByID() returned nil")
	}
	if n.Title != "武林大会召开" {
		t.Errorf("Title = %q, want %q", n.Title, "武林大会召开")
	}
	if n.Status != 1 {
		t.Errorf("Status = %d, want 1", n.Status)
	}
}

func TestNewsRepo_FindByID_NotFound(t *testing.T) {
	db := newTestDB(t)
	r := NewNewsRepo(db)

	n, err := r.FindByID(context.Background(), 9999)
	if err != nil {
		t.Fatalf("FindByID() error: %v", err)
	}
	if n != nil {
		t.Error("FindByID() should return nil for non-existent id")
	}
}

func TestNewsRepo_List(t *testing.T) {
	db := newTestDB(t)
	r := NewNewsRepo(db)

	createTestNews(t, r, "新闻一")
	createTestNews(t, r, "新闻二")
	createTestNews(t, r, "新闻三")

	list, total, err := r.List(context.Background(), 1, 10, "", "")
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if total != 3 {
		t.Errorf("total = %d, want 3", total)
	}
	if len(list) != 3 {
		t.Errorf("len(list) = %d, want 3", len(list))
	}
}

func TestNewsRepo_List_Pagination(t *testing.T) {
	db := newTestDB(t)
	r := NewNewsRepo(db)

	for i := 0; i < 5; i++ {
		createTestNews(t, r, "分页新闻")
	}

	list, total, err := r.List(context.Background(), 1, 2, "", "")
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if total != 5 {
		t.Errorf("total = %d, want 5", total)
	}
	if len(list) != 2 {
		t.Errorf("len(list) = %d, want 2", len(list))
	}
}

func TestNewsRepo_List_FilterByCategory(t *testing.T) {
	db := newTestDB(t)
	r := NewNewsRepo(db)
	ctx := context.Background()

	createTestNews(t, r, "武林新闻")

	now := time.Now()
	r.Create(ctx, &model.News{
		Title: "科技新闻", Category: "科技", Section: "其他",
		Keywords: "[]", Status: 1, PublishedAt: &now, Fingerprint: "fp-tech",
	})

	list, total, _ := r.List(ctx, 1, 10, "武林", "")
	if total != 1 {
		t.Errorf("total = %d, want 1", total)
	}
	if len(list) != 1 {
		t.Errorf("len(list) = %d, want 1", len(list))
	}
}

func TestNewsRepo_Search(t *testing.T) {
	db := newTestDB(t)
	r := NewNewsRepo(db)

	createTestNews(t, r, "少林寺武僧表演")
	createTestNews(t, r, "武当派新掌门")
	createTestNews(t, r, "峨眉山风景")

	list, total, err := r.Search(context.Background(), "武", 1, 10)
	if err != nil {
		t.Fatalf("Search() error: %v", err)
	}
	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
	if len(list) != 2 {
		t.Errorf("len(list) = %d, want 2", len(list))
	}
}

func TestNewsRepo_FindByFingerprint(t *testing.T) {
	db := newTestDB(t)
	r := NewNewsRepo(db)

	createTestNews(t, r, "指纹测试")

	n, err := r.FindByFingerprint(context.Background(), "fp-指纹测试")
	if err != nil {
		t.Fatalf("FindByFingerprint() error: %v", err)
	}
	if n == nil {
		t.Fatal("FindByFingerprint() returned nil")
	}
	if n.Title != "指纹测试" {
		t.Errorf("Title = %q, want %q", n.Title, "指纹测试")
	}

	// Not found
	n2, _ := r.FindByFingerprint(context.Background(), "nonexistent")
	if n2 != nil {
		t.Error("FindByFingerprint() should return nil for unknown fingerprint")
	}
}

func TestNewsRepo_IncrReadCount(t *testing.T) {
	db := newTestDB(t)
	r := NewNewsRepo(db)
	ctx := context.Background()

	id := createTestNews(t, r, "阅读测试")

	r.IncrReadCount(ctx, id)
	r.IncrReadCount(ctx, id)

	n, _ := r.FindByID(ctx, id)
	if n.ReadCount != 2 {
		t.Errorf("ReadCount = %d, want 2", n.ReadCount)
	}
}

func TestNewsRepo_MarkRead(t *testing.T) {
	db := newTestDB(t)
	nr := NewNewsRepo(db)
	ur := NewUserRepo(db)
	ctx := context.Background()

	userID := createTestUser(t, ur)
	newsID := createTestNews(t, nr, "阅读标记测试")

	err := nr.MarkRead(ctx, userID, newsID)
	if err != nil {
		t.Fatalf("MarkRead() error: %v", err)
	}

	// Mark again should not error (INSERT OR IGNORE)
	err = nr.MarkRead(ctx, userID, newsID)
	if err != nil {
		t.Fatalf("MarkRead() duplicate error: %v", err)
	}

	news, err := nr.FindByID(ctx, newsID)
	if err != nil {
		t.Fatalf("FindByID() error: %v", err)
	}
	if news.ReadCount != 1 {
		t.Fatalf("ReadCount = %d, want 1 after duplicate mark", news.ReadCount)
	}

	marks, err := nr.GetReadMarks(ctx, userID, []int64{newsID})
	if err != nil {
		t.Fatalf("GetReadMarks() error: %v", err)
	}
	if !marks[newsID] {
		t.Error("GetReadMarks() should return true for marked news")
	}
}

func TestNewsRepo_GetReadMarks_Empty(t *testing.T) {
	db := newTestDB(t)
	r := NewNewsRepo(db)

	marks, err := r.GetReadMarks(context.Background(), 1, []int64{})
	if err != nil {
		t.Fatalf("GetReadMarks() error: %v", err)
	}
	if len(marks) != 0 {
		t.Errorf("GetReadMarks() should return empty map for empty newsIDs")
	}
}

func TestNewsRepo_Headlines(t *testing.T) {
	db := newTestDB(t)
	r := NewNewsRepo(db)
	ctx := context.Background()

	id1 := createTestNews(t, r, "头条一")
	id2 := createTestNews(t, r, "头条二")

	now := time.Now().UTC() // Use UTC to match SQLite's CURRENT_TIMESTAMP
	err := r.SetHeadlines(ctx, []*model.Headline{
		{NewsID: id1, Rank: 1, Title: "头条一标题", ActiveAt: now.Add(-24 * time.Hour), ExpireAt: now.Add(24 * time.Hour)},
		{NewsID: id2, Rank: 2, Title: "头条二标题", ActiveAt: now.Add(-24 * time.Hour), ExpireAt: now.Add(24 * time.Hour)},
	})
	if err != nil {
		t.Fatalf("SetHeadlines() error: %v", err)
	}

	headlines, err := r.GetHeadlines(ctx)
	if err != nil {
		t.Fatalf("GetHeadlines() error: %v", err)
	}
	if len(headlines) != 2 {
		t.Fatalf("len(headlines) = %d, want 2", len(headlines))
	}
	if headlines[0].Rank != 1 {
		t.Errorf("first headline rank = %d, want 1", headlines[0].Rank)
	}
}

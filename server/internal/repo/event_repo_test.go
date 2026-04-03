package repo

import (
	"context"
	"testing"
	"time"

	"github.com/lunancy1992/jianghu-server/internal/model"
)

func TestEventRepo_CreateAndFindByID(t *testing.T) {
	db := newTestDB(t)
	r := NewEventRepo(db)
	ctx := context.Background()

	event := &model.Event{
		Title:       "某某事件真相",
		Description: "一起关于某某的事实核查",
		Category:    "社会",
		Status:      "ongoing",
	}
	id, err := r.Create(ctx, event)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	found, err := r.FindByID(ctx, id)
	if err != nil {
		t.Fatalf("FindByID() error: %v", err)
	}
	if found == nil {
		t.Fatal("FindByID() returned nil")
	}
	if found.Title != "某某事件真相" {
		t.Errorf("Title = %q, want %q", found.Title, "某某事件真相")
	}
	if found.Status != "ongoing" {
		t.Errorf("Status = %q, want %q", found.Status, "ongoing")
	}
}

func TestEventRepo_FindByID_NotFound(t *testing.T) {
	db := newTestDB(t)
	r := NewEventRepo(db)

	e, err := r.FindByID(context.Background(), 9999)
	if err != nil {
		t.Fatalf("FindByID() error: %v", err)
	}
	if e != nil {
		t.Error("FindByID() should return nil for non-existent id")
	}
}

func TestEventRepo_List(t *testing.T) {
	db := newTestDB(t)
	r := NewEventRepo(db)
	ctx := context.Background()

	eventID, _ := r.Create(ctx, &model.Event{Title: "事件一", Description: "事件摘要", Category: "社会", Status: "ongoing"})
	r.Create(ctx, &model.Event{Title: "事件二", Category: "科技", Status: "resolved"})

	now := time.Now()
	_, _ = r.CreateNode(ctx, &model.EventNode{EventID: eventID, Title: "节点标题", Content: "节点内容", NodeTime: now, Source: "来源", Veracity: 2})

	list, total, err := r.List(ctx, 1, 10)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
	if len(list) != 2 {
		t.Errorf("len(list) = %d, want 2", len(list))
	}
	if list[0].ID != eventID {
		t.Fatalf("first event id = %d, want %d", list[0].ID, eventID)
	}
	if list[0].Summary != "事件摘要" {
		t.Errorf("Summary = %q, want %q", list[0].Summary, "事件摘要")
	}
	if list[0].LatestUpdate != "节点标题" {
		t.Errorf("LatestUpdate = %q, want %q", list[0].LatestUpdate, "节点标题")
	}
	if list[0].Veracity != 2 {
		t.Errorf("Veracity = %d, want 2", list[0].Veracity)
	}
	if list[0].NodeCount != 1 {
		t.Errorf("NodeCount = %d, want 1", list[0].NodeCount)
	}
}

func TestEventRepo_Nodes(t *testing.T) {
	db := newTestDB(t)
	r := NewEventRepo(db)
	ctx := context.Background()

	eventID, _ := r.Create(ctx, &model.Event{Title: "时间线测试", Category: "社会", Status: "ongoing"})

	now := time.Now()
	r.CreateNode(ctx, &model.EventNode{
		EventID: eventID, Title: "节点一", Content: "第一个事件", NodeTime: now.Add(-2 * time.Hour), Source: "来源A",
	})
	r.CreateNode(ctx, &model.EventNode{
		EventID: eventID, Title: "节点二", Content: "第二个事件", NodeTime: now.Add(-1 * time.Hour), Source: "来源B",
	})

	nodes, err := r.ListNodes(ctx, eventID)
	if err != nil {
		t.Fatalf("ListNodes() error: %v", err)
	}
	if len(nodes) != 2 {
		t.Errorf("len(nodes) = %d, want 2", len(nodes))
	}
	// Should be ordered by node_time ASC
	if nodes[0].Title != "节点一" {
		t.Errorf("nodes[0].Title = %q, want %q", nodes[0].Title, "节点一")
	}
}

func TestEventRepo_LinkNews(t *testing.T) {
	db := newTestDB(t)
	er := NewEventRepo(db)
	nr := NewNewsRepo(db)
	ctx := context.Background()

	eventID, _ := er.Create(ctx, &model.Event{Title: "关联新闻测试", Category: "社会", Status: "ongoing"})
	newsID := createTestNews(t, nr, "关联的新闻")

	err := er.LinkNews(ctx, eventID, newsID)
	if err != nil {
		t.Fatalf("LinkNews() error: %v", err)
	}

	// Link again should not error (INSERT OR IGNORE)
	err = er.LinkNews(ctx, eventID, newsID)
	if err != nil {
		t.Fatalf("LinkNews() duplicate error: %v", err)
	}

	linked, err := er.GetLinkedNews(ctx, eventID)
	if err != nil {
		t.Fatalf("GetLinkedNews() error: %v", err)
	}
	if len(linked) != 1 {
		t.Errorf("len(linked) = %d, want 1", len(linked))
	}
	if linked[0].Title != "关联的新闻" {
		t.Errorf("linked[0].Title = %q, want %q", linked[0].Title, "关联的新闻")
	}
}

func TestEventRepo_Evidence(t *testing.T) {
	db := newTestDB(t)
	er := NewEventRepo(db)
	ur := NewUserRepo(db)
	ctx := context.Background()

	eventID, _ := er.Create(ctx, &model.Event{Title: "证据测试", Category: "社会", Status: "ongoing"})
	userID := createTestUser(t, ur)

	// Create approved evidence
	id, err := er.CreateEvidence(ctx, &model.Evidence{
		EventID: eventID, UserID: userID, Type: "image",
		URL: "/uploads/test.jpg", Description: "截图", Status: 1,
	})
	if err != nil {
		t.Fatalf("CreateEvidence() error: %v", err)
	}
	if id <= 0 {
		t.Fatalf("CreateEvidence() returned invalid id: %d", id)
	}

	// Create pending evidence (should not appear in list)
	er.CreateEvidence(ctx, &model.Evidence{
		EventID: eventID, UserID: userID, Type: "document",
		URL: "/uploads/doc.pdf", Description: "文档", Status: 0,
	})

	list, err := er.ListEvidences(ctx, eventID)
	if err != nil {
		t.Fatalf("ListEvidences() error: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("len(list) = %d, want 1 (only approved)", len(list))
	}
}

package repo

import (
	"context"
	"testing"
	"time"

	"github.com/lunancy1992/jianghu-server/internal/model"
)

func TestUserRepo_CreateAndFindByID(t *testing.T) {
	db := newTestDB(t)
	r := NewUserRepo(db)
	ctx := context.Background()

	user := &model.User{
		Phone:    "13800138000",
		Nickname: "测试用户",
		Role:     "user",
		Status:   0,
	}

	id, err := r.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}
	if id <= 0 {
		t.Fatalf("Create() returned invalid id: %d", id)
	}

	found, err := r.FindByID(ctx, id)
	if err != nil {
		t.Fatalf("FindByID() error: %v", err)
	}
	if found == nil {
		t.Fatal("FindByID() returned nil")
	}
	if found.Phone != "13800138000" {
		t.Errorf("Phone = %q, want %q", found.Phone, "13800138000")
	}
	if found.Nickname != "测试用户" {
		t.Errorf("Nickname = %q, want %q", found.Nickname, "测试用户")
	}
}

func TestUserRepo_FindByPhone(t *testing.T) {
	db := newTestDB(t)
	r := NewUserRepo(db)
	ctx := context.Background()

	// Not found
	u, err := r.FindByPhone(ctx, "13900139000")
	if err != nil {
		t.Fatalf("FindByPhone() error: %v", err)
	}
	if u != nil {
		t.Fatal("FindByPhone() should return nil for non-existent phone")
	}

	// Create and find
	r.Create(ctx, &model.User{Phone: "13900139000", Nickname: "test", Role: "user"})
	u, err = r.FindByPhone(ctx, "13900139000")
	if err != nil {
		t.Fatalf("FindByPhone() error: %v", err)
	}
	if u == nil {
		t.Fatal("FindByPhone() should find user")
	}
	if u.Phone != "13900139000" {
		t.Errorf("Phone = %q, want %q", u.Phone, "13900139000")
	}
}

func TestUserRepo_FindByID_NotFound(t *testing.T) {
	db := newTestDB(t)
	r := NewUserRepo(db)

	u, err := r.FindByID(context.Background(), 9999)
	if err != nil {
		t.Fatalf("FindByID() error: %v", err)
	}
	if u != nil {
		t.Error("FindByID() should return nil for non-existent id")
	}
}

func TestUserRepo_Update(t *testing.T) {
	db := newTestDB(t)
	r := NewUserRepo(db)
	ctx := context.Background()

	id, _ := r.Create(ctx, &model.User{Phone: "13700137000", Nickname: "old", Role: "user"})
	user, _ := r.FindByID(ctx, id)

	user.Nickname = "new-name"
	user.Role = "admin"
	err := r.Update(ctx, user)
	if err != nil {
		t.Fatalf("Update() error: %v", err)
	}

	updated, _ := r.FindByID(ctx, id)
	if updated.Nickname != "new-name" {
		t.Errorf("Nickname = %q, want %q", updated.Nickname, "new-name")
	}
	if updated.Role != "admin" {
		t.Errorf("Role = %q, want %q", updated.Role, "admin")
	}
}

func TestUserRepo_OAuth(t *testing.T) {
	db := newTestDB(t)
	r := NewUserRepo(db)
	ctx := context.Background()

	userID, _ := r.Create(ctx, &model.User{Phone: "13600136000", Nickname: "oauth-user", Role: "user"})

	oauth := &model.UserOAuth{
		UserID:     userID,
		Provider:   "wechat",
		ProviderID: "wx_123456",
		UnionID:    "union_abc",
	}
	oauthID, err := r.CreateOAuth(ctx, oauth)
	if err != nil {
		t.Fatalf("CreateOAuth() error: %v", err)
	}
	if oauthID <= 0 {
		t.Fatalf("CreateOAuth() returned invalid id: %d", oauthID)
	}

	found, err := r.FindOAuth(ctx, "wechat", "wx_123456")
	if err != nil {
		t.Fatalf("FindOAuth() error: %v", err)
	}
	if found == nil {
		t.Fatal("FindOAuth() returned nil")
	}
	if found.UserID != userID {
		t.Errorf("UserID = %d, want %d", found.UserID, userID)
	}
}

func TestUserRepo_FindOAuth_NotFound(t *testing.T) {
	db := newTestDB(t)
	r := NewUserRepo(db)

	o, err := r.FindOAuth(context.Background(), "apple", "nonexistent")
	if err != nil {
		t.Fatalf("FindOAuth() error: %v", err)
	}
	if o != nil {
		t.Error("FindOAuth() should return nil for non-existent oauth")
	}
}

func TestUserRepo_Verification(t *testing.T) {
	db := newTestDB(t)
	r := NewUserRepo(db)
	ctx := context.Background()
	phone := "13500135000"

	// Save verification code
	expiresAt := time.Now().Add(5 * time.Minute)
	err := r.SaveVerification(ctx, phone, "123456", expiresAt)
	if err != nil {
		t.Fatalf("SaveVerification() error: %v", err)
	}

	// Find verification
	v, err := r.FindVerification(ctx, phone, "123456")
	if err != nil {
		t.Fatalf("FindVerification() error: %v", err)
	}
	if v == nil {
		t.Fatal("FindVerification() returned nil")
	}
	if v.Phone != phone {
		t.Errorf("Phone = %q, want %q", v.Phone, phone)
	}
	if v.Code != "123456" {
		t.Errorf("Code = %q, want %q", v.Code, "123456")
	}

	// Wrong code
	v2, err := r.FindVerification(ctx, phone, "000000")
	if err != nil {
		t.Fatalf("FindVerification() error: %v", err)
	}
	if v2 != nil {
		t.Error("FindVerification() should return nil for wrong code")
	}

	// Mark used
	err = r.MarkVerificationUsed(ctx, v.ID)
	if err != nil {
		t.Fatalf("MarkVerificationUsed() error: %v", err)
	}

	// Should not find used verification
	v3, err := r.FindVerification(ctx, phone, "123456")
	if err != nil {
		t.Fatalf("FindVerification() error: %v", err)
	}
	if v3 != nil {
		t.Error("FindVerification() should return nil for used verification")
	}
}

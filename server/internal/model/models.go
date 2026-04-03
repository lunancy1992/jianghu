package model

import (
	"time"
)

// User represents a registered user.
type User struct {
	ID        int64     `json:"id" db:"id"`
	Phone     string    `json:"phone" db:"phone"`
	Nickname  string    `json:"nickname" db:"nickname"`
	Avatar    string    `json:"avatar" db:"avatar"`
	Role      string    `json:"role" db:"role"` // user, admin
	Status    int       `json:"status" db:"status"` // 0=active, 1=banned
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserVerification holds SMS verification codes.
type UserVerification struct {
	ID        int64     `json:"id" db:"id"`
	Phone     string    `json:"phone" db:"phone"`
	Code      string    `json:"code" db:"code"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	Used      bool      `json:"used" db:"used"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// UserOAuth stores third-party OAuth info.
type UserOAuth struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	Provider   string    `json:"provider" db:"provider"` // wechat, apple
	ProviderID string    `json:"provider_id" db:"provider_id"`
	UnionID    string    `json:"union_id" db:"union_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// News represents a news article.
type News struct {
	ID          int64          `json:"id" db:"id"`
	Title       string         `json:"title" db:"title"`
	Summary     string         `json:"summary" db:"summary"`
	Content     string         `json:"content" db:"content"`
	Source      string         `json:"source" db:"source"`
	SourceURL   string         `json:"source_url" db:"source_url"`
	Author      string         `json:"author" db:"author"`
	Category    string         `json:"category" db:"category"`
	Section     string         `json:"section" db:"section"` // 江湖快报, 武林秘闻, etc.
	Keywords    string         `json:"keywords" db:"keywords"` // JSON array
	CoverImage  string         `json:"cover_image" db:"cover_image"`
	Sensitivity int            `json:"sensitivity" db:"sensitivity"` // 0-10
	Veracity    int            `json:"veracity" db:"veracity"` // 0=待核实, 1=已证实, 2=已辟谣
	CommentText string         `json:"comment_text" db:"comment_text"` // 编辑点评
	Fingerprint string         `json:"fingerprint" db:"fingerprint"`
	Status      int            `json:"status" db:"status"` // 0=pending, 1=published, 2=rejected, 3=archived
	ReadCount   int64          `json:"read_count" db:"read_count"`
	LikeCount   int64          `json:"like_count" db:"like_count"`
	CommentCnt  int64          `json:"comment_count" db:"comment_count"`
	PublishedAt *time.Time     `json:"published_at" db:"published_at"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

// Headline is a pinned/highlighted news item.
type Headline struct {
	ID        int64     `json:"id" db:"id"`
	NewsID    int64     `json:"news_id" db:"news_id"`
	Rank      int       `json:"rank" db:"rank"`
	Title     string    `json:"title" db:"title"` // override title
	ActiveAt  time.Time `json:"active_at" db:"active_at"`
	ExpireAt  time.Time `json:"expire_at" db:"expire_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// HeadlineWithNews is a headline enriched with its news data.
type HeadlineWithNews struct {
	ID          int64      `json:"id"`
	NewsID      int64      `json:"news_id"`
	Rank        int        `json:"rank"`
	Title       string     `json:"title"`
	Source      string     `json:"source"`
	SourceURL   string     `json:"source_url"`
	PublishedAt *time.Time `json:"published_at"`
	Veracity    int        `json:"veracity"`
	CommentText string     `json:"comment_text"`
}

// UserReadMark tracks which news a user has read.
type UserReadMark struct {
	ID     int64     `json:"id" db:"id"`
	UserID int64     `json:"user_id" db:"user_id"`
	NewsID int64     `json:"news_id" db:"news_id"`
	ReadAt time.Time `json:"read_at" db:"read_at"`
}

// Comment is a user comment on a news article.
type Comment struct {
	ID        int64          `json:"id" db:"id"`
	NewsID    int64          `json:"news_id" db:"news_id"`
	UserID    int64          `json:"user_id" db:"user_id"`
	ParentID  int64  `json:"parent_id" db:"parent_id"`
	Content   string `json:"content" db:"content"`
	Stance    string `json:"stance" db:"stance"` // agree, disagree, neutral
	LikeCount int64  `json:"like_count" db:"like_count"`
	Status    int    `json:"status" db:"status"` // 0=pending, 1=approved, 2=rejected
	AuditNote string `json:"audit_note" db:"audit_note"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`

	// Joined fields (not stored in DB)
	Nickname string `json:"nickname,omitempty" db:"-"`
	Avatar   string `json:"avatar,omitempty" db:"-"`
}

// Like tracks likes on comments.
type Like struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	CommentID int64     `json:"comment_id" db:"comment_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Evidence is an uploaded proof file attached to events.
type Evidence struct {
	ID          int64     `json:"id" db:"id"`
	EventID     int64     `json:"event_id" db:"event_id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	Type        string    `json:"type" db:"type"` // image, document, link
	URL         string    `json:"url" db:"url"`
	Description string    `json:"description" db:"description"`
	Status      int       `json:"status" db:"status"` // 0=pending, 1=approved, 2=rejected
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Report is a user report on content.
type Report struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	TargetType string    `json:"target_type" db:"target_type"` // news, comment
	TargetID   int64     `json:"target_id" db:"target_id"`
	Reason     string    `json:"reason" db:"reason"`
	Status     int       `json:"status" db:"status"` // 0=pending, 1=resolved
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// CoinAccount holds the user's coin balance.
type CoinAccount struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Balance   int64     `json:"balance" db:"balance"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CoinTransaction records a single coin credit/debit.
type CoinTransaction struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Amount    int64     `json:"amount" db:"amount"` // positive=credit, negative=debit
	Type      string    `json:"type" db:"type"` // weekly_grant, comment_cost, comment_reward, evidence_reward, admin_adjust
	RefID     string    `json:"ref_id" db:"ref_id"` // reference to related entity
	Note      string    `json:"note" db:"note"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Event is a fact-checking event/topic.
type Event struct {
	ID          int64     `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Category    string    `json:"category" db:"category"`
	Status      string    `json:"status" db:"status"` // ongoing, resolved, debunked
	CoverImage  string    `json:"cover_image" db:"cover_image"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type EventListItem struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Status       string    `json:"status"`
	UpdatedAt    time.Time `json:"updated_at"`
	Summary      string    `json:"summary"`
	LatestUpdate string    `json:"latest_update"`
	Veracity     int       `json:"veracity"`
	NodeCount    int       `json:"node_count"`
}

// EventNode is a timeline node in an event.
type EventNode struct {
	ID        int64     `json:"id" db:"id"`
	EventID   int64     `json:"event_id" db:"event_id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	NodeTime  time.Time `json:"node_time" db:"node_time"`
	Source    string    `json:"source" db:"source"`
	Veracity  int       `json:"veracity" db:"veracity"` // 0=待核实, 1=已证实, 2=已辟谣
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// EventNews links events to news articles.
type EventNews struct {
	ID      int64 `json:"id" db:"id"`
	EventID int64 `json:"event_id" db:"event_id"`
	NewsID  int64 `json:"news_id" db:"news_id"`
}

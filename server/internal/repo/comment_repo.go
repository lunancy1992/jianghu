package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/lunancy1992/jianghu-server/internal/model"
)

// scanParentID is a helper to scan a nullable parent_id column.
func scanParentID(val sql.NullInt64) int64 {
	if val.Valid {
		return val.Int64
	}
	return 0
}

type CommentRepo struct {
	db *sql.DB
}

func NewCommentRepo(db *sql.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) Create(ctx context.Context, c *model.Comment) (int64, error) {
	now := time.Now()
	var parentID interface{}
	if c.ParentID > 0 {
		parentID = c.ParentID
	}
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO comments (news_id, user_id, parent_id, content, stance, status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		c.NewsID, c.UserID, parentID, c.Content, c.Stance, c.Status, now, now,
	)
	if err != nil {
		return 0, err
	}
	// Increment news comment count
	_, _ = r.db.ExecContext(ctx,
		`UPDATE news SET comment_count = comment_count + 1 WHERE id = ?`, c.NewsID)
	return res.LastInsertId()
}

func (r *CommentRepo) ListByNews(ctx context.Context, newsID int64, page, size int) ([]*model.Comment, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM comments WHERE news_id=? AND status=1`, newsID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	rows, err := r.db.QueryContext(ctx,
		`SELECT c.id, c.news_id, c.user_id, c.parent_id, c.content, c.stance,
		 c.like_count, c.status, c.created_at, u.nickname, u.avatar
		 FROM comments c
		 LEFT JOIN users u ON c.user_id = u.id
		 WHERE c.news_id=? AND c.status=1
		 ORDER BY c.like_count DESC, c.created_at DESC
		 LIMIT ? OFFSET ?`, newsID, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*model.Comment
	for rows.Next() {
		c := &model.Comment{}
		var parentID sql.NullInt64
		if err := rows.Scan(
			&c.ID, &c.NewsID, &c.UserID, &parentID, &c.Content, &c.Stance,
			&c.LikeCount, &c.Status, &c.CreatedAt, &c.Nickname, &c.Avatar,
		); err != nil {
			return nil, 0, err
		}
		c.ParentID = scanParentID(parentID)
		list = append(list, c)
	}
	return list, total, rows.Err()
}

func (r *CommentRepo) UpdateStatus(ctx context.Context, id int64, status int, note string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE comments SET status=?, audit_note=?, updated_at=? WHERE id=?`,
		status, note, time.Now(), id,
	)
	return err
}

func (r *CommentRepo) IncrLikeCount(ctx context.Context, commentID int64, delta int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE comments SET like_count = like_count + ? WHERE id = ?`, delta, commentID)
	return err
}

func (r *CommentRepo) FindTopLiked(ctx context.Context, since time.Time, limit int) ([]*model.Comment, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, news_id, user_id, content, like_count, created_at
		 FROM comments
		 WHERE status=1 AND like_count >= 10 AND created_at >= ?
		 ORDER BY like_count DESC
		 LIMIT ?`, since, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.Comment
	for rows.Next() {
		c := &model.Comment{}
		if err := rows.Scan(&c.ID, &c.NewsID, &c.UserID, &c.Content, &c.LikeCount, &c.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (r *CommentRepo) FindByID(ctx context.Context, id int64) (*model.Comment, error) {
	c := &model.Comment{}
	var parentID sql.NullInt64
	var auditNote sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT id, news_id, user_id, parent_id, content, stance, like_count, status, audit_note, created_at, updated_at
		 FROM comments WHERE id=?`, id).Scan(
		&c.ID, &c.NewsID, &c.UserID, &parentID, &c.Content, &c.Stance,
		&c.LikeCount, &c.Status, &auditNote, &c.CreatedAt, &c.UpdatedAt,
	)
	c.ParentID = scanParentID(parentID)
	c.AuditNote = auditNote.String
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CommentRepo) ListPending(ctx context.Context, page, size int) ([]*model.Comment, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM comments WHERE status=0`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	rows, err := r.db.QueryContext(ctx,
		`SELECT c.id, c.news_id, c.user_id, c.content, c.stance, c.status, c.created_at,
		 u.nickname
		 FROM comments c
		 LEFT JOIN users u ON c.user_id = u.id
		 WHERE c.status=0
		 ORDER BY c.created_at ASC
		 LIMIT ? OFFSET ?`, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*model.Comment
	for rows.Next() {
		c := &model.Comment{}
		if err := rows.Scan(
			&c.ID, &c.NewsID, &c.UserID, &c.Content, &c.Stance, &c.Status, &c.CreatedAt, &c.Nickname,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, c)
	}
	return list, total, rows.Err()
}

// --- Likes ---

func (r *CommentRepo) CreateLike(ctx context.Context, userID, commentID int64) error {
	res, err := r.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO likes (user_id, comment_id, created_at) VALUES (?, ?, ?)`,
		userID, commentID, time.Now(),
	)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return nil
	}
	return r.IncrLikeCount(ctx, commentID, 1)
}

func (r *CommentRepo) DeleteLike(ctx context.Context, userID, commentID int64) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM likes WHERE user_id=? AND comment_id=?`, userID, commentID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n > 0 {
		return r.IncrLikeCount(ctx, commentID, -1)
	}
	return nil
}

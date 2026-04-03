package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lunancy1992/jianghu-server/internal/model"
)

type NewsRepo struct {
	db *sql.DB
}

func NewNewsRepo(db *sql.DB) *NewsRepo {
	return &NewsRepo{db: db}
}

func (r *NewsRepo) Create(ctx context.Context, n *model.News) (int64, error) {
	now := time.Now()
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO news (title, summary, content, source, source_url, author, category, section,
		 keywords, cover_image, sensitivity, fingerprint, status, published_at, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		n.Title, n.Summary, n.Content, n.Source, n.SourceURL, n.Author, n.Category, n.Section,
		n.Keywords, n.CoverImage, n.Sensitivity, n.Fingerprint, n.Status, n.PublishedAt, now, now,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *NewsRepo) FindByID(ctx context.Context, id int64) (*model.News, error) {
	n := &model.News{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, summary, content, source, source_url, author, category, section,
		 keywords, cover_image, sensitivity, veracity, comment_text, fingerprint, status,
		 read_count, like_count, comment_count, published_at, created_at, updated_at
		 FROM news WHERE id = ?`, id).Scan(
		&n.ID, &n.Title, &n.Summary, &n.Content, &n.Source, &n.SourceURL, &n.Author, &n.Category,
		&n.Section, &n.Keywords, &n.CoverImage, &n.Sensitivity, &n.Veracity, &n.CommentText,
		&n.Fingerprint, &n.Status, &n.ReadCount, &n.LikeCount, &n.CommentCnt,
		&n.PublishedAt, &n.CreatedAt, &n.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (r *NewsRepo) List(ctx context.Context, page, size int, category, section string) ([]*model.News, int, error) {
	where := "WHERE status = 1"
	args := []interface{}{}
	if category != "" {
		where += " AND category = ?"
		args = append(args, category)
	}
	if section != "" {
		where += " AND section = ?"
		args = append(args, section)
	}

	var total int
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM news "+where, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	query := fmt.Sprintf(
		`SELECT id, title, summary, source, source_url, category, section, cover_image, veracity,
		 comment_text, read_count, like_count, comment_count, published_at, created_at
		 FROM news %s ORDER BY published_at DESC LIMIT ? OFFSET ?`, where)
	args = append(args, size, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*model.News
	for rows.Next() {
		n := &model.News{}
		if err := rows.Scan(
			&n.ID, &n.Title, &n.Summary, &n.Source, &n.SourceURL, &n.Category, &n.Section,
			&n.CoverImage, &n.Veracity, &n.CommentText, &n.ReadCount, &n.LikeCount, &n.CommentCnt,
			&n.PublishedAt, &n.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, n)
	}
	return list, total, rows.Err()
}

func (r *NewsRepo) Search(ctx context.Context, q string, page, size int) ([]*model.News, int, error) {
	like := "%" + q + "%"
	var total int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM news WHERE status=1 AND (title LIKE ? OR summary LIKE ?)`,
		like, like).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, summary, source, source_url, category, section, cover_image, veracity,
		 comment_text, read_count, like_count, comment_count, published_at, created_at
		 FROM news WHERE status=1 AND (title LIKE ? OR summary LIKE ?)
		 ORDER BY published_at DESC LIMIT ? OFFSET ?`,
		like, like, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*model.News
	for rows.Next() {
		n := &model.News{}
		if err := rows.Scan(
			&n.ID, &n.Title, &n.Summary, &n.Source, &n.SourceURL, &n.Category, &n.Section,
			&n.CoverImage, &n.Veracity, &n.CommentText, &n.ReadCount, &n.LikeCount, &n.CommentCnt,
			&n.PublishedAt, &n.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, n)
	}
	return list, total, rows.Err()
}

func (r *NewsRepo) FindByFingerprint(ctx context.Context, fp string) (*model.News, error) {
	n := &model.News{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, fingerprint FROM news WHERE fingerprint = ?`, fp).Scan(
		&n.ID, &n.Title, &n.Fingerprint,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (r *NewsRepo) IncrReadCount(ctx context.Context, newsID int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE news SET read_count = read_count + 1 WHERE id = ?`, newsID)
	return err
}

// --- Headlines ---

func (r *NewsRepo) SetHeadlines(ctx context.Context, headlines []*model.Headline) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, h := range headlines {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO headlines (news_id, rank, title, active_at, expire_at, created_at)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			h.NewsID, h.Rank, h.Title, h.ActiveAt, h.ExpireAt, time.Now().UTC(),
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *NewsRepo) GetHeadlines(ctx context.Context) ([]*model.HeadlineWithNews, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT h.id, h.news_id, h.rank,
		        CASE WHEN h.title != '' THEN h.title ELSE n.title END,
		        n.source, n.source_url, n.published_at, n.veracity, n.comment_text
		 FROM headlines h
		 JOIN news n ON n.id = h.news_id
		 WHERE h.active_at <= CURRENT_TIMESTAMP AND h.expire_at > CURRENT_TIMESTAMP
		 ORDER BY h.rank ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.HeadlineWithNews
	for rows.Next() {
		h := &model.HeadlineWithNews{}
		if err := rows.Scan(&h.ID, &h.NewsID, &h.Rank, &h.Title, &h.Source, &h.SourceURL, &h.PublishedAt, &h.Veracity, &h.CommentText); err != nil {
			return nil, err
		}
		list = append(list, h)
	}
	return list, rows.Err()
}

func (r *NewsRepo) GetHeadlinesByDate(ctx context.Context, date string) ([]*model.HeadlineWithNews, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT h.id, h.news_id, h.rank,
		        CASE WHEN h.title != '' THEN h.title ELSE n.title END,
		        n.source, n.source_url, n.published_at, n.veracity, n.comment_text
		 FROM headlines h
		 JOIN news n ON n.id = h.news_id
		 WHERE DATE(h.active_at) = ?
		 ORDER BY h.rank ASC`, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.HeadlineWithNews
	for rows.Next() {
		h := &model.HeadlineWithNews{}
		if err := rows.Scan(&h.ID, &h.NewsID, &h.Rank, &h.Title, &h.Source, &h.SourceURL, &h.PublishedAt, &h.Veracity, &h.CommentText); err != nil {
			return nil, err
		}
		list = append(list, h)
	}
	return list, rows.Err()
}

// --- Read Marks ---

func (r *NewsRepo) MarkRead(ctx context.Context, userID, newsID int64) error {
	res, err := r.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO user_read_marks (user_id, news_id, read_at) VALUES (?, ?, ?)`,
		userID, newsID, time.Now(),
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
	return r.IncrReadCount(ctx, newsID)
}

func (r *NewsRepo) GetReadMarks(ctx context.Context, userID int64, newsIDs []int64) (map[int64]bool, error) {
	if len(newsIDs) == 0 {
		return map[int64]bool{}, nil
	}
	query := "SELECT news_id FROM user_read_marks WHERE user_id = ? AND news_id IN ("
	args := []interface{}{userID}
	for i, id := range newsIDs {
		if i > 0 {
			query += ","
		}
		query += "?"
		args = append(args, id)
	}
	query += ")"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64]bool)
	for rows.Next() {
		var nid int64
		if err := rows.Scan(&nid); err != nil {
			return nil, err
		}
		result[nid] = true
	}
	return result, rows.Err()
}

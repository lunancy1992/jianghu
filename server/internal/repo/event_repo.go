package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/lunancy1992/jianghu-server/internal/model"
)

type EventRepo struct {
	db *sql.DB
}

func NewEventRepo(db *sql.DB) *EventRepo {
	return &EventRepo{db: db}
}

func (r *EventRepo) Create(ctx context.Context, e *model.Event) (int64, error) {
	now := time.Now()
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO events (title, description, category, status, cover_image, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		e.Title, e.Description, e.Category, e.Status, e.CoverImage, now, now,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepo) Exists(ctx context.Context, id int64) (bool, error) {
	var exists int
	err := r.db.QueryRowContext(ctx, `SELECT 1 FROM events WHERE id=? LIMIT 1`, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *EventRepo) FindByID(ctx context.Context, id int64) (*model.Event, error) {
	e := &model.Event{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, description, category, status, cover_image, created_at, updated_at
		 FROM events WHERE id=?`, id).Scan(
		&e.ID, &e.Title, &e.Description, &e.Category, &e.Status, &e.CoverImage, &e.CreatedAt, &e.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *EventRepo) List(ctx context.Context, page, size int) ([]*model.EventListItem, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM events`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	rows, err := r.db.QueryContext(ctx, `
		SELECT e.id, e.title, e.status, e.updated_at,
		       COALESCE(NULLIF(e.description, ''), '') AS summary,
		       COALESCE(
		         (
		           SELECT COALESCE(NULLIF(en.title, ''), en.content, '')
		           FROM event_nodes en
		           WHERE en.event_id = e.id
		           ORDER BY en.node_time DESC, en.id DESC
		           LIMIT 1
		         ),
		         ''
		       ) AS latest_update,
		       COALESCE(
		         (
		           SELECT en.veracity
		           FROM event_nodes en
		           WHERE en.event_id = e.id
		           ORDER BY en.node_time DESC, en.id DESC
		           LIMIT 1
		         ),
		         0
		       ) AS veracity,
		       COALESCE((SELECT COUNT(*) FROM event_nodes en WHERE en.event_id = e.id), 0) AS node_count
		FROM events e
		ORDER BY e.updated_at DESC
		LIMIT ? OFFSET ?`, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*model.EventListItem
	for rows.Next() {
		e := &model.EventListItem{}
		if err := rows.Scan(
			&e.ID, &e.Title, &e.Status, &e.UpdatedAt, &e.Summary, &e.LatestUpdate, &e.Veracity, &e.NodeCount,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, e)
	}
	return list, total, rows.Err()
}

func (r *EventRepo) CreateNode(ctx context.Context, n *model.EventNode) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO event_nodes (event_id, title, content, node_time, source, veracity, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		n.EventID, n.Title, n.Content, n.NodeTime, n.Source, n.Veracity, time.Now(),
	)
	if err != nil {
		return 0, err
	}
	_, _ = r.db.ExecContext(ctx, `UPDATE events SET updated_at = ? WHERE id = ?`, time.Now(), n.EventID)
	return res.LastInsertId()
}

func (r *EventRepo) ListNodes(ctx context.Context, eventID int64) ([]*model.EventNode, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, event_id, title, content, node_time, source, veracity, created_at
		 FROM event_nodes WHERE event_id=? ORDER BY node_time ASC`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.EventNode
	for rows.Next() {
		n := &model.EventNode{}
		if err := rows.Scan(&n.ID, &n.EventID, &n.Title, &n.Content, &n.NodeTime, &n.Source, &n.Veracity, &n.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, n)
	}
	return list, rows.Err()
}

func (r *EventRepo) LinkNews(ctx context.Context, eventID, newsID int64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO event_news (event_id, news_id) VALUES (?, ?)`, eventID, newsID)
	return err
}

func (r *EventRepo) GetLinkedNews(ctx context.Context, eventID int64) ([]*model.News, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT n.id, n.title, n.summary, n.source, n.published_at
		 FROM news n
		 INNER JOIN event_news en ON n.id = en.news_id
		 WHERE en.event_id = ?
		 ORDER BY n.published_at DESC`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.News
	for rows.Next() {
		n := &model.News{}
		if err := rows.Scan(&n.ID, &n.Title, &n.Summary, &n.Source, &n.PublishedAt); err != nil {
			return nil, err
		}
		list = append(list, n)
	}
	return list, rows.Err()
}

// --- Evidence ---

func (r *EventRepo) CreateEvidence(ctx context.Context, e *model.Evidence) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO evidences (event_id, user_id, type, url, description, status, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		e.EventID, e.UserID, e.Type, e.URL, e.Description, e.Status, time.Now(),
	)
	if err != nil {
		return 0, err
	}
	_, _ = r.db.ExecContext(ctx, `UPDATE events SET updated_at = ? WHERE id = ?`, time.Now(), e.EventID)
	return res.LastInsertId()
}

func (r *EventRepo) ListEvidences(ctx context.Context, eventID int64) ([]*model.Evidence, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, event_id, user_id, type, url, description, status, created_at
		 FROM evidences WHERE event_id=? AND status=1 ORDER BY created_at DESC`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.Evidence
	for rows.Next() {
		e := &model.Evidence{}
		if err := rows.Scan(&e.ID, &e.EventID, &e.UserID, &e.Type, &e.URL, &e.Description, &e.Status, &e.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

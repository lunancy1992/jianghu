package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/lunancy1992/jianghu-server/internal/model"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, phone, nickname, avatar, role, status, created_at, updated_at
		 FROM users WHERE phone = ?`, phone).Scan(
		&u.ID, &u.Phone, &u.Nickname, &u.Avatar, &u.Role, &u.Status, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) FindByID(ctx context.Context, id int64) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, phone, nickname, avatar, role, status, created_at, updated_at
		 FROM users WHERE id = ?`, id).Scan(
		&u.ID, &u.Phone, &u.Nickname, &u.Avatar, &u.Role, &u.Status, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) Create(ctx context.Context, u *model.User) (int64, error) {
	now := time.Now()
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO users (phone, nickname, avatar, role, status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		u.Phone, u.Nickname, u.Avatar, u.Role, u.Status, now, now,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *UserRepo) Update(ctx context.Context, u *model.User) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET nickname=?, avatar=?, role=?, status=?, updated_at=? WHERE id=?`,
		u.Nickname, u.Avatar, u.Role, u.Status, time.Now(), u.ID,
	)
	return err
}

func (r *UserRepo) FindOAuth(ctx context.Context, provider, providerID string) (*model.UserOAuth, error) {
	o := &model.UserOAuth{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, provider, provider_id, union_id, created_at
		 FROM user_oauth WHERE provider=? AND provider_id=?`, provider, providerID).Scan(
		&o.ID, &o.UserID, &o.Provider, &o.ProviderID, &o.UnionID, &o.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (r *UserRepo) CreateOAuth(ctx context.Context, o *model.UserOAuth) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO user_oauth (user_id, provider, provider_id, union_id, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		o.UserID, o.Provider, o.ProviderID, o.UnionID, time.Now(),
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// SaveVerification stores an SMS verification code.
func (r *UserRepo) SaveVerification(ctx context.Context, phone, code string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO user_verifications (phone, code, expires_at) VALUES (?, ?, ?)`,
		phone, code, expiresAt.UTC().Format("2006-01-02 15:04:05"),
	)
	return err
}

// FindVerification retrieves the latest unused verification code for a phone.
func (r *UserRepo) FindVerification(ctx context.Context, phone, code string) (*model.UserVerification, error) {
	v := &model.UserVerification{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, phone, code, expires_at, used, created_at
		 FROM user_verifications
		 WHERE phone=? AND code=? AND used=0 AND expires_at > CURRENT_TIMESTAMP
		 ORDER BY id DESC LIMIT 1`, phone, code).Scan(
		&v.ID, &v.Phone, &v.Code, &v.ExpiresAt, &v.Used, &v.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return v, nil
}

// MarkVerificationUsed marks a verification record as used.
func (r *UserRepo) MarkVerificationUsed(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE user_verifications SET used=1 WHERE id=?`, id,
	)
	return err
}

package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"time2meet/internal/domain/entity"
	"time2meet/internal/domain/repository"
	"time2meet/internal/domain/valueobject"
	"time2meet/internal/infrastructure/persistence/postgres/dto"
	"time2meet/pkg/apperror"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

var _ repository.UserRepository = (*UserRepo)(nil)

func (r *UserRepo) Create(ctx context.Context, u entity.User) (valueobject.UUID, error) {
	q := `
		INSERT INTO users (email, password_hash, full_name, phone, role, is_active)
		VALUES ($1, $2, $3, NULLIF($4, ''), $5, $6)
		RETURNING id
	`
	var id string
	if err := r.db.QueryRowxContext(ctx, q,
		u.Email.String(),
		u.PasswordHash,
		u.FullName,
		u.Phone,
		string(u.Role),
		u.IsActive,
	).Scan(&id); err != nil {
		return valueobject.Nil, apperror.New(apperror.CodeInternal, "create user failed", err)
	}
	uid, err := valueobject.ParseUUID(id)
	if err != nil {
		return valueobject.Nil, apperror.New(apperror.CodeInternal, "invalid uuid returned from db", err)
	}
	return uid, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id valueobject.UUID) (entity.User, error) {
	q := `SELECT id, email, password_hash, full_name, phone, role, is_active, created_at, updated_at FROM users WHERE id = $1`
	var row dto.UserRow
	if err := r.db.GetContext(ctx, &row, q, id.String()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, apperror.New(apperror.CodeNotFound, "user not found", err)
		}
		return entity.User{}, apperror.New(apperror.CodeInternal, "get user failed", err)
	}
	return mapUserRow(row)
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	q := `SELECT id, email, password_hash, full_name, phone, role, is_active, created_at, updated_at FROM users WHERE email = $1`
	var row dto.UserRow
	if err := r.db.GetContext(ctx, &row, q, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, apperror.New(apperror.CodeNotFound, "user not found", err)
		}
		return entity.User{}, apperror.New(apperror.CodeInternal, "get user by email failed", err)
	}
	return mapUserRow(row)
}

func (r *UserRepo) List(ctx context.Context, limit, offset int) ([]entity.User, error) {
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	q := `
		SELECT id, email, password_hash, full_name, phone, role, is_active, created_at, updated_at
		FROM users
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	var rows []dto.UserRow
	if err := r.db.SelectContext(ctx, &rows, q, limit, offset); err != nil {
		return nil, apperror.New(apperror.CodeInternal, "list users failed", err)
	}
	out := make([]entity.User, 0, len(rows))
	for _, row := range rows {
		u, err := mapUserRow(row)
		if err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, nil
}

func (r *UserRepo) Update(ctx context.Context, u entity.User) error {
	q := `
		UPDATE users
		SET email = $1, password_hash = $2, full_name = $3, phone = NULLIF($4, ''), role = $5, is_active = $6
		WHERE id = $7
	`
	res, err := r.db.ExecContext(ctx, q,
		u.Email.String(),
		u.PasswordHash,
		u.FullName,
		u.Phone,
		string(u.Role),
		u.IsActive,
		u.ID.String(),
	)
	if err != nil {
		return apperror.New(apperror.CodeInternal, "update user failed", err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return apperror.New(apperror.CodeNotFound, "user not found", sql.ErrNoRows)
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id valueobject.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id.String())
	if err != nil {
		return apperror.New(apperror.CodeInternal, "delete user failed", err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return apperror.New(apperror.CodeNotFound, "user not found", sql.ErrNoRows)
	}
	return nil
}

func mapUserRow(row dto.UserRow) (entity.User, error) {
	uid, err := valueobject.ParseUUID(row.ID)
	if err != nil {
		return entity.User{}, apperror.New(apperror.CodeInternal, "invalid user id in db", err)
	}
	em, err := valueobject.ParseEmail(row.Email)
	if err != nil {
		return entity.User{}, apperror.New(apperror.CodeInternal, "invalid email in db", err)
	}
	u := entity.User{
		ID:           uid,
		Email:        em,
		PasswordHash: row.PasswordHash,
		FullName:     row.FullName,
		Phone:        "",
		Role:         entity.UserRole(row.Role),
		IsActive:     row.IsActive,
	}
	if row.Phone.Valid {
		u.Phone = row.Phone.String
	}
	if row.CreatedAt.Valid {
		u.CreatedAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		u.UpdatedAt = row.UpdatedAt.Time
	}
	return u, nil
}

type UserProfileRepo struct {
	db *sqlx.DB
}

func NewUserProfileRepo(db *sqlx.DB) *UserProfileRepo {
	return &UserProfileRepo{db: db}
}

var _ repository.UserProfileRepository = (*UserProfileRepo)(nil)

func (r *UserProfileRepo) Upsert(ctx context.Context, p entity.UserProfile) error {
	social, err := json.Marshal(p.SocialLinks)
	if err != nil {
		return apperror.New(apperror.CodeInternal, "marshal social_links failed", err)
	}
	prefs, err := json.Marshal(p.Preferences)
	if err != nil {
		return apperror.New(apperror.CodeInternal, "marshal preferences failed", err)
	}

	q := `
		INSERT INTO user_profiles (user_id, avatar_url, bio, social_links, preferences)
		VALUES ($1, NULLIF($2, ''), NULLIF($3, ''), $4::jsonb, $5::jsonb)
		ON CONFLICT (user_id) DO UPDATE
		SET avatar_url = EXCLUDED.avatar_url,
		    bio = EXCLUDED.bio,
		    social_links = EXCLUDED.social_links,
		    preferences = EXCLUDED.preferences,
		    updated_at = NOW()
	`
	_, err = r.db.ExecContext(ctx, q, p.UserID.String(), p.AvatarURL, p.Bio, string(social), string(prefs))
	if err != nil {
		return apperror.New(apperror.CodeInternal, "upsert user profile failed", err)
	}
	return nil
}

func (r *UserProfileRepo) GetByUserID(ctx context.Context, userID valueobject.UUID) (entity.UserProfile, error) {
	q := `SELECT id, user_id, avatar_url, bio, social_links, preferences, updated_at FROM user_profiles WHERE user_id = $1`
	var row dto.UserProfileRow
	if err := r.db.GetContext(ctx, &row, q, userID.String()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.UserProfile{}, apperror.New(apperror.CodeNotFound, "user profile not found", err)
		}
		return entity.UserProfile{}, apperror.New(apperror.CodeInternal, "get user profile failed", err)
	}
	pid, err := valueobject.ParseUUID(row.ID)
	if err != nil {
		return entity.UserProfile{}, apperror.New(apperror.CodeInternal, "invalid profile id in db", err)
	}
	uid, err := valueobject.ParseUUID(row.UserID)
	if err != nil {
		return entity.UserProfile{}, apperror.New(apperror.CodeInternal, "invalid profile user_id in db", err)
	}
	p := entity.UserProfile{
		ID:          pid,
		UserID:      uid,
		AvatarURL:   "",
		Bio:         "",
		SocialLinks: map[string]any{},
		Preferences: map[string]any{},
	}
	if row.AvatarURL.Valid {
		p.AvatarURL = row.AvatarURL.String
	}
	if row.Bio.Valid {
		p.Bio = row.Bio.String
	}
	if len(row.SocialLinks) > 0 {
		if err := json.Unmarshal(row.SocialLinks, &p.SocialLinks); err != nil {
			return entity.UserProfile{}, fmt.Errorf("unmarshal social_links: %w", err)
		}
	}
	if len(row.Preferences) > 0 {
		if err := json.Unmarshal(row.Preferences, &p.Preferences); err != nil {
			return entity.UserProfile{}, fmt.Errorf("unmarshal preferences: %w", err)
		}
	}
	if row.UpdatedAt.Valid {
		p.UpdatedAt = row.UpdatedAt.Time
	}
	return p, nil
}

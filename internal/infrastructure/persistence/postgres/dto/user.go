package dto

import (
	"database/sql"
	"encoding/json"
)

type UserRow struct {
	ID           string         `db:"id"`
	Email        string         `db:"email"`
	PasswordHash string         `db:"password_hash"`
	FullName     string         `db:"full_name"`
	Phone        sql.NullString `db:"phone"`
	Role         string         `db:"role"`
	IsActive     bool           `db:"is_active"`
	CreatedAt    sql.NullTime   `db:"created_at"`
	UpdatedAt    sql.NullTime   `db:"updated_at"`
}

type UserProfileRow struct {
	ID          string          `db:"id"`
	UserID      string          `db:"user_id"`
	AvatarURL   sql.NullString  `db:"avatar_url"`
	Bio         sql.NullString  `db:"bio"`
	SocialLinks json.RawMessage `db:"social_links"`
	Preferences json.RawMessage `db:"preferences"`
	UpdatedAt   sql.NullTime    `db:"updated_at"`
}



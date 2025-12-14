package dto

import "database/sql"

type EventRow struct {
	ID              string         `db:"id"`
	OrganizerID     string         `db:"organizer_id"`
	Title           string         `db:"title"`
	Description     sql.NullString `db:"description"`
	Status          string         `db:"status"`
	IsPublic        bool           `db:"is_public"`
	MaxParticipants sql.NullInt64  `db:"max_participants"`
	CoverImage      sql.NullString `db:"cover_image"`
	CreatedAt       sql.NullTime   `db:"created_at"`
	UpdatedAt       sql.NullTime   `db:"updated_at"`
}



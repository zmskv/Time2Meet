package dto

import (
	"database/sql"
	"encoding/json"
)

type VenueRow struct {
	ID           string         `db:"id"`
	Name         string         `db:"name"`
	Address      string         `db:"address"`
	City         string         `db:"city"`
	Country      string         `db:"country"`
	Capacity     int            `db:"capacity"`
	ContactPhone sql.NullString `db:"contact_phone"`
	ContactEmail sql.NullString `db:"contact_email"`
	Website      sql.NullString `db:"website"`
	IsActive     bool           `db:"is_active"`
	CreatedAt    sql.NullTime   `db:"created_at"`
	UpdatedAt    sql.NullTime   `db:"updated_at"`
}

type RoomRow struct {
	ID          string          `db:"id"`
	VenueID     string          `db:"venue_id"`
	Name        string          `db:"name"`
	Capacity    int             `db:"capacity"`
	Floor       sql.NullInt64   `db:"floor"`
	Equipment   json.RawMessage `db:"equipment"`
	HourlyRate  string          `db:"hourly_rate"`
	IsAvailable bool            `db:"is_available"`
	CreatedAt   sql.NullTime    `db:"created_at"`
	UpdatedAt   sql.NullTime    `db:"updated_at"`
}



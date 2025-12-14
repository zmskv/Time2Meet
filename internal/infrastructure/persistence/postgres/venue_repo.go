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

type VenueRepo struct{ db *sqlx.DB }

func NewVenueRepo(db *sqlx.DB) *VenueRepo { return &VenueRepo{db: db} }

var _ repository.VenueRepository = (*VenueRepo)(nil)

func (r *VenueRepo) Create(ctx context.Context, v entity.Venue) (valueobject.UUID, error) {
	q := `
		INSERT INTO venues (name, address, city, country, capacity, contact_phone, contact_email, website, is_active)
		VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''), NULLIF($7, ''), NULLIF($8, ''), $9)
		RETURNING id
	`
	var id string
	if err := r.db.QueryRowxContext(ctx, q,
		v.Name, v.Address, v.City, v.Country, v.Capacity, v.ContactPhone, v.ContactEmail, v.Website, v.IsActive,
	).Scan(&id); err != nil {
		return valueobject.Nil, apperror.New(apperror.CodeInternal, "create venue failed", err)
	}
	vid, err := valueobject.ParseUUID(id)
	if err != nil {
		return valueobject.Nil, apperror.New(apperror.CodeInternal, "invalid uuid returned from db", err)
	}
	return vid, nil
}

func (r *VenueRepo) GetByID(ctx context.Context, id valueobject.UUID) (entity.Venue, error) {
	q := `SELECT id, name, address, city, country, capacity, contact_phone, contact_email, website, is_active, created_at, updated_at FROM venues WHERE id = $1`
	var row dto.VenueRow
	if err := r.db.GetContext(ctx, &row, q, id.String()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Venue{}, apperror.New(apperror.CodeNotFound, "venue not found", err)
		}
		return entity.Venue{}, apperror.New(apperror.CodeInternal, "get venue failed", err)
	}
	return mapVenueRow(row), nil
}

func (r *VenueRepo) List(ctx context.Context, limit, offset int) ([]entity.Venue, error) {
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	q := `
		SELECT id, name, address, city, country, capacity, contact_phone, contact_email, website, is_active, created_at, updated_at
		FROM venues
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	var rows []dto.VenueRow
	if err := r.db.SelectContext(ctx, &rows, q, limit, offset); err != nil {
		return nil, apperror.New(apperror.CodeInternal, "list venues failed", err)
	}
	out := make([]entity.Venue, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapVenueRow(row))
	}
	return out, nil
}

func (r *VenueRepo) Update(ctx context.Context, v entity.Venue) error {
	q := `
		UPDATE venues
		SET name=$1, address=$2, city=$3, country=$4, capacity=$5,
		    contact_phone=NULLIF($6,''), contact_email=NULLIF($7,''), website=NULLIF($8,''), is_active=$9
		WHERE id=$10
	`
	res, err := r.db.ExecContext(ctx, q, v.Name, v.Address, v.City, v.Country, v.Capacity, v.ContactPhone, v.ContactEmail, v.Website, v.IsActive, v.ID.String())
	if err != nil {
		return apperror.New(apperror.CodeInternal, "update venue failed", err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return apperror.New(apperror.CodeNotFound, "venue not found", sql.ErrNoRows)
	}
	return nil
}

func (r *VenueRepo) Delete(ctx context.Context, id valueobject.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM venues WHERE id = $1`, id.String())
	if err != nil {
		return apperror.New(apperror.CodeInternal, "delete venue failed", err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return apperror.New(apperror.CodeNotFound, "venue not found", sql.ErrNoRows)
	}
	return nil
}

func mapVenueRow(row dto.VenueRow) entity.Venue {
	vid, _ := valueobject.ParseUUID(row.ID)
	v := entity.Venue{
		ID:       vid,
		Name:     row.Name,
		Address:  row.Address,
		City:     row.City,
		Country:  row.Country,
		Capacity: row.Capacity,
		IsActive: row.IsActive,
	}
	if row.ContactPhone.Valid {
		v.ContactPhone = row.ContactPhone.String
	}
	if row.ContactEmail.Valid {
		v.ContactEmail = row.ContactEmail.String
	}
	if row.Website.Valid {
		v.Website = row.Website.String
	}
	if row.CreatedAt.Valid {
		v.CreatedAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		v.UpdatedAt = row.UpdatedAt.Time
	}
	return v
}

type RoomRepo struct{ db *sqlx.DB }

func NewRoomRepo(db *sqlx.DB) *RoomRepo { return &RoomRepo{db: db} }

var _ repository.RoomRepository = (*RoomRepo)(nil)

func (r *RoomRepo) Create(ctx context.Context, rm entity.Room) (valueobject.UUID, error) {
	eq, err := json.Marshal(rm.Equipment)
	if err != nil {
		return valueobject.Nil, apperror.New(apperror.CodeInternal, "marshal equipment failed", err)
	}
	q := `
		INSERT INTO rooms (venue_id, name, capacity, floor, equipment, hourly_rate, is_available)
		VALUES ($1, $2, $3, $4, $5::jsonb, $6, $7)
		RETURNING id
	`
	var id string
	var floor any
	if rm.Floor != nil {
		floor = *rm.Floor
	}
	if err := r.db.QueryRowxContext(ctx, q, rm.VenueID.String(), rm.Name, rm.Capacity, floor, string(eq), rm.HourlyRate, rm.IsAvailable).Scan(&id); err != nil {
		return valueobject.Nil, apperror.New(apperror.CodeInternal, "create room failed", err)
	}
	rid, err := valueobject.ParseUUID(id)
	if err != nil {
		return valueobject.Nil, apperror.New(apperror.CodeInternal, "invalid uuid returned from db", err)
	}
	return rid, nil
}

func (r *RoomRepo) GetByID(ctx context.Context, id valueobject.UUID) (entity.Room, error) {
	q := `SELECT id, venue_id, name, capacity, floor, equipment, hourly_rate, is_available, created_at, updated_at FROM rooms WHERE id=$1`
	var row dto.RoomRow
	if err := r.db.GetContext(ctx, &row, q, id.String()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Room{}, apperror.New(apperror.CodeNotFound, "room not found", err)
		}
		return entity.Room{}, apperror.New(apperror.CodeInternal, "get room failed", err)
	}
	return mapRoomRow(row)
}

func (r *RoomRepo) ListByVenueID(ctx context.Context, venueID valueobject.UUID) ([]entity.Room, error) {
	q := `
		SELECT id, venue_id, name, capacity, floor, equipment, hourly_rate, is_available, created_at, updated_at
		FROM rooms
		WHERE venue_id = $1
		ORDER BY id DESC
	`
	var rows []dto.RoomRow
	if err := r.db.SelectContext(ctx, &rows, q, venueID.String()); err != nil {
		return nil, apperror.New(apperror.CodeInternal, "list rooms failed", err)
	}
	out := make([]entity.Room, 0, len(rows))
	for _, row := range rows {
		rm, err := mapRoomRow(row)
		if err != nil {
			return nil, err
		}
		out = append(out, rm)
	}
	return out, nil
}

func (r *RoomRepo) Update(ctx context.Context, rm entity.Room) error {
	eq, err := json.Marshal(rm.Equipment)
	if err != nil {
		return apperror.New(apperror.CodeInternal, "marshal equipment failed", err)
	}
	var floor any
	if rm.Floor != nil {
		floor = *rm.Floor
	}
	q := `
		UPDATE rooms
		SET venue_id=$1, name=$2, capacity=$3, floor=$4, equipment=$5::jsonb, hourly_rate=$6, is_available=$7
		WHERE id=$8
	`
	res, err := r.db.ExecContext(ctx, q, rm.VenueID.String(), rm.Name, rm.Capacity, floor, string(eq), rm.HourlyRate, rm.IsAvailable, rm.ID.String())
	if err != nil {
		return apperror.New(apperror.CodeInternal, "update room failed", err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return apperror.New(apperror.CodeNotFound, "room not found", sql.ErrNoRows)
	}
	return nil
}

func (r *RoomRepo) Delete(ctx context.Context, id valueobject.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM rooms WHERE id=$1`, id.String())
	if err != nil {
		return apperror.New(apperror.CodeInternal, "delete room failed", err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return apperror.New(apperror.CodeNotFound, "room not found", sql.ErrNoRows)
	}
	return nil
}

func mapRoomRow(row dto.RoomRow) (entity.Room, error) {
	rid, err := valueobject.ParseUUID(row.ID)
	if err != nil {
		return entity.Room{}, fmt.Errorf("invalid room id in db: %w", err)
	}
	vid, err := valueobject.ParseUUID(row.VenueID)
	if err != nil {
		return entity.Room{}, fmt.Errorf("invalid venue id in db: %w", err)
	}
	rm := entity.Room{
		ID:          rid,
		VenueID:     vid,
		Name:        row.Name,
		Capacity:    row.Capacity,
		Equipment:   map[string]any{},
		HourlyRate:  row.HourlyRate,
		IsAvailable: row.IsAvailable,
	}
	if row.Floor.Valid {
		v := int(row.Floor.Int64)
		rm.Floor = &v
	}
	if len(row.Equipment) > 0 {
		if err := json.Unmarshal(row.Equipment, &rm.Equipment); err != nil {
			return entity.Room{}, fmt.Errorf("unmarshal equipment: %w", err)
		}
	}
	if row.CreatedAt.Valid {
		rm.CreatedAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		rm.UpdatedAt = row.UpdatedAt.Time
	}
	return rm, nil
}

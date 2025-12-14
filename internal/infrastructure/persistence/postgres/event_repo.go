package postgres

import (
	"context"
	"database/sql"
	"errors"

	"time2meet/internal/domain/entity"
	"time2meet/internal/domain/repository"
	"time2meet/internal/domain/valueobject"
	"time2meet/internal/infrastructure/persistence/postgres/dto"
	"time2meet/pkg/apperror"

	"github.com/jmoiron/sqlx"
)

type EventRepo struct {
	db *sqlx.DB
}

func NewEventRepo(db *sqlx.DB) *EventRepo { return &EventRepo{db: db} }

var _ repository.EventRepository = (*EventRepo)(nil)

func (r *EventRepo) Create(ctx context.Context, e entity.Event) (valueobject.UUID, error) {
	q := `
		INSERT INTO events (organizer_id, title, description, status, is_public, max_participants, cover_image)
		VALUES ($1, $2, NULLIF($3, ''), $4, $5, $6, NULLIF($7, ''))
		RETURNING id
	`
	var id string
	var maxp any
	if e.MaxParticipants != nil {
		maxp = *e.MaxParticipants
	}
	if err := r.db.QueryRowxContext(ctx, q,
		e.OrganizerID.String(),
		e.Title,
		e.Description,
		string(e.Status),
		e.IsPublic,
		maxp,
		e.CoverImage,
	).Scan(&id); err != nil {
		return valueobject.Nil, apperror.New(apperror.CodeInternal, "create event failed", err)
	}
	eid, err := valueobject.ParseUUID(id)
	if err != nil {
		return valueobject.Nil, apperror.New(apperror.CodeInternal, "invalid uuid returned from db", err)
	}
	return eid, nil
}

func (r *EventRepo) GetByID(ctx context.Context, id valueobject.UUID) (entity.Event, error) {
	q := `
		SELECT id, organizer_id, title, description, status, is_public, max_participants, cover_image, created_at, updated_at
		FROM events
		WHERE id = $1
	`
	var row dto.EventRow
	if err := r.db.GetContext(ctx, &row, q, id.String()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Event{}, apperror.New(apperror.CodeNotFound, "event not found", err)
		}
		return entity.Event{}, apperror.New(apperror.CodeInternal, "get event failed", err)
	}
	return mapEventRow(row)
}

func (r *EventRepo) List(ctx context.Context, organizerID *valueobject.UUID, status *string, limit, offset int) ([]entity.Event, error) {
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	q := `
		SELECT id, organizer_id, title, description, status, is_public, max_participants, cover_image, created_at, updated_at
		FROM events
		WHERE ($1::BIGINT IS NULL OR organizer_id = $1)
		  AND ($2::TEXT IS NULL OR status = $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	var org any
	if organizerID != nil {
		org = organizerID.String()
	}
	var st any
	if status != nil {
		st = *status
	}

	var rows []dto.EventRow
	if err := r.db.SelectContext(ctx, &rows, q, org, st, limit, offset); err != nil {
		return nil, apperror.New(apperror.CodeInternal, "list events failed", err)
	}
	out := make([]entity.Event, 0, len(rows))
	for _, row := range rows {
		e, err := mapEventRow(row)
		if err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, nil
}

func (r *EventRepo) Update(ctx context.Context, e entity.Event) error {
	q := `
		UPDATE events
		SET title = $1,
		    description = NULLIF($2, ''),
		    status = $3,
		    is_public = $4,
		    max_participants = $5,
		    cover_image = NULLIF($6, '')
		WHERE id = $7
	`
	var maxp any
	if e.MaxParticipants != nil {
		maxp = *e.MaxParticipants
	}
	res, err := r.db.ExecContext(ctx, q,
		e.Title, e.Description, string(e.Status), e.IsPublic, maxp, e.CoverImage, e.ID.String(),
	)
	if err != nil {
		return apperror.New(apperror.CodeInternal, "update event failed", err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return apperror.New(apperror.CodeNotFound, "event not found", sql.ErrNoRows)
	}
	return nil
}

func (r *EventRepo) Delete(ctx context.Context, id valueobject.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM events WHERE id = $1`, id.String())
	if err != nil {
		return apperror.New(apperror.CodeInternal, "delete event failed", err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return apperror.New(apperror.CodeNotFound, "event not found", sql.ErrNoRows)
	}
	return nil
}

func mapEventRow(row dto.EventRow) (entity.Event, error) {
	eid, err := valueobject.ParseUUID(row.ID)
	if err != nil {
		return entity.Event{}, apperror.New(apperror.CodeInternal, "invalid event id in db", err)
	}
	oid, err := valueobject.ParseUUID(row.OrganizerID)
	if err != nil {
		return entity.Event{}, apperror.New(apperror.CodeInternal, "invalid organizer id in db", err)
	}
	st := valueobject.EventStatus(row.Status)
	if err := st.Validate(); err != nil {
		return entity.Event{}, apperror.New(apperror.CodeInternal, "invalid event status in db", err)
	}
	e := entity.Event{
		ID:          eid,
		OrganizerID: oid,
		Title:       row.Title,
		Description: "",
		Status:      st,
		IsPublic:    row.IsPublic,
		CoverImage:  "",
	}
	if row.Description.Valid {
		e.Description = row.Description.String
	}
	if row.CoverImage.Valid {
		e.CoverImage = row.CoverImage.String
	}
	if row.MaxParticipants.Valid {
		v := int(row.MaxParticipants.Int64)
		e.MaxParticipants = &v
	}
	if row.CreatedAt.Valid {
		e.CreatedAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		e.UpdatedAt = row.UpdatedAt.Time
	}
	return e, nil
}

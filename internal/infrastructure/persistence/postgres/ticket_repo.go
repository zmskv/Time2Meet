package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"time2meet/internal/domain/entity"
	"time2meet/internal/domain/repository"
	"time2meet/internal/domain/valueobject"
	"time2meet/internal/infrastructure/persistence/postgres/dto"
	"time2meet/pkg/apperror"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type TicketTypeRepo struct{ db *sqlx.DB }

func NewTicketTypeRepo(db *sqlx.DB) *TicketTypeRepo { return &TicketTypeRepo{db: db} }

var _ repository.TicketTypeRepository = (*TicketTypeRepo)(nil)

func (r *TicketTypeRepo) Create(ctx context.Context, tt entity.TicketType) (valueobject.UUID, error) {
	q := `
		INSERT INTO ticket_types (event_id, name, price, quantity_total, quantity_sold, sale_start, sale_end, description, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NULLIF($8, ''), $9)
		RETURNING id
	`
	var id string
	var saleStart any
	var saleEnd any
	if tt.SaleStart != nil {
		saleStart = *tt.SaleStart
	}
	if tt.SaleEnd != nil {
		saleEnd = *tt.SaleEnd
	}
	if err := r.db.QueryRowxContext(ctx, q,
		tt.EventID.String(),
		tt.Name,
		tt.Price.Amount.StringFixed(2),
		tt.QuantityTotal,
		tt.QuantitySold,
		saleStart,
		saleEnd,
		tt.Description,
		tt.IsActive,
	).Scan(&id); err != nil {
		return valueobject.Nil, apperror.New(apperror.CodeInternal, "create ticket type failed", err)
	}
	out, err := valueobject.ParseUUID(id)
	if err != nil {
		return valueobject.Nil, apperror.New(apperror.CodeInternal, "invalid uuid returned from db", err)
	}
	return out, nil
}

func (r *TicketTypeRepo) GetByID(ctx context.Context, id valueobject.UUID) (entity.TicketType, error) {
	q := `
		SELECT id, event_id, name, price, currency, quantity_total, quantity_sold, sale_start, sale_end, description, is_active, created_at, updated_at
		FROM ticket_types
		WHERE id = $1
	`
	var row dto.TicketTypeRow
	if err := r.db.GetContext(ctx, &row, q, id.String()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.TicketType{}, apperror.New(apperror.CodeNotFound, "ticket type not found", err)
		}
		return entity.TicketType{}, apperror.New(apperror.CodeInternal, "get ticket type failed", err)
	}
	return mapTicketTypeRow(row)
}

func (r *TicketTypeRepo) ListByEventID(ctx context.Context, eventID valueobject.UUID) ([]entity.TicketType, error) {
	q := `
		SELECT id, event_id, name, price, currency, quantity_total, quantity_sold, sale_start, sale_end, description, is_active, created_at, updated_at
		FROM ticket_types
		WHERE event_id = $1
		ORDER BY id ASC
	`
	var rows []dto.TicketTypeRow
	if err := r.db.SelectContext(ctx, &rows, q, eventID.String()); err != nil {
		return nil, apperror.New(apperror.CodeInternal, "list ticket types failed", err)
	}
	out := make([]entity.TicketType, 0, len(rows))
	for _, row := range rows {
		tt, err := mapTicketTypeRow(row)
		if err != nil {
			return nil, err
		}
		out = append(out, tt)
	}
	return out, nil
}

func (r *TicketTypeRepo) Update(ctx context.Context, tt entity.TicketType) error {
	q := `
		UPDATE ticket_types
		SET name=$1, price=$2, quantity_total=$3, quantity_sold=$4,
		    sale_start=$5, sale_end=$6, description=NULLIF($7,''), is_active=$8
		WHERE id=$9
	`
	var saleStart any
	var saleEnd any
	if tt.SaleStart != nil {
		saleStart = *tt.SaleStart
	}
	if tt.SaleEnd != nil {
		saleEnd = *tt.SaleEnd
	}
	res, err := r.db.ExecContext(ctx, q,
		tt.Name,
		tt.Price.Amount.StringFixed(2),
		tt.QuantityTotal,
		tt.QuantitySold,
		saleStart,
		saleEnd,
		tt.Description,
		tt.IsActive,
		tt.ID.String(),
	)
	if err != nil {
		return apperror.New(apperror.CodeInternal, "update ticket type failed", err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return apperror.New(apperror.CodeNotFound, "ticket type not found", sql.ErrNoRows)
	}
	return nil
}

func (r *TicketTypeRepo) Delete(ctx context.Context, id valueobject.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM ticket_types WHERE id=$1`, id.String())
	if err != nil {
		return apperror.New(apperror.CodeInternal, "delete ticket type failed", err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return apperror.New(apperror.CodeNotFound, "ticket type not found", sql.ErrNoRows)
	}
	return nil
}

func mapTicketTypeRow(row dto.TicketTypeRow) (entity.TicketType, error) {
	id, err := valueobject.ParseUUID(row.ID)
	if err != nil {
		return entity.TicketType{}, fmt.Errorf("invalid ticket_type id in db: %w", err)
	}
	eid, err := valueobject.ParseUUID(row.EventID)
	if err != nil {
		return entity.TicketType{}, fmt.Errorf("invalid ticket_type event_id in db: %w", err)
	}
	amt, err := decimal.NewFromString(row.Price)
	if err != nil {
		return entity.TicketType{}, fmt.Errorf("invalid price in db: %w", err)
	}
	price, err := valueobject.NewMoney(amt)
	if err != nil {
		return entity.TicketType{}, fmt.Errorf("invalid money in db: %w", err)
	}
	tt := entity.TicketType{
		ID:            id,
		EventID:       eid,
		Name:          row.Name,
		Price:         price,
		QuantityTotal: row.QuantityTotal,
		QuantitySold:  row.QuantitySold,
		Description:   "",
		IsActive:      row.IsActive,
	}
	if row.Description.Valid {
		tt.Description = row.Description.String
	}
	if row.SaleStart.Valid {
		t := row.SaleStart.Time
		tt.SaleStart = &t
	}
	if row.SaleEnd.Valid {
		t := row.SaleEnd.Time
		tt.SaleEnd = &t
	}
	if row.CreatedAt.Valid {
		tt.CreatedAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		tt.UpdatedAt = row.UpdatedAt.Time
	}
	return tt, nil
}

type TicketRepo struct{ db *sqlx.DB }

func NewTicketRepo(db *sqlx.DB) *TicketRepo { return &TicketRepo{db: db} }

var _ repository.TicketRepository = (*TicketRepo)(nil)

func (r *TicketRepo) Create(ctx context.Context, t entity.Ticket) (valueobject.UUID, error) {
	q := `
		INSERT INTO tickets (ticket_type_id, buyer_id, purchase_date, status, qr_code, amount_paid, used_at)
		VALUES ($1, $2, COALESCE($3, NOW()), $4, $5, $6, $7)
		RETURNING id
	`
	var id string
	var purchase any
	if !t.PurchaseDate.IsZero() {
		purchase = t.PurchaseDate
	}
	var used any
	if t.UsedAt != nil {
		used = *t.UsedAt
	}
	if err := r.db.QueryRowxContext(ctx, q,
		t.TicketTypeID.String(),
		t.BuyerID.String(),
		purchase,
		string(t.Status),
		t.QRCode,
		t.AmountPaid.Amount.StringFixed(2),
		used,
	).Scan(&id); err != nil {
		return valueobject.Nil, apperror.New(apperror.CodeInternal, "create ticket failed", err)
	}
	out, err := valueobject.ParseUUID(id)
	if err != nil {
		return valueobject.Nil, apperror.New(apperror.CodeInternal, "invalid uuid returned from db", err)
	}
	return out, nil
}

func (r *TicketRepo) GetByID(ctx context.Context, id valueobject.UUID) (entity.Ticket, error) {
	q := `
		SELECT id, ticket_type_id, buyer_id, purchase_date, status, qr_code, amount_paid, used_at, created_at, updated_at
		FROM tickets
		WHERE id = $1
	`
	var row dto.TicketRow
	if err := r.db.GetContext(ctx, &row, q, id.String()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Ticket{}, apperror.New(apperror.CodeNotFound, "ticket not found", err)
		}
		return entity.Ticket{}, apperror.New(apperror.CodeInternal, "get ticket failed", err)
	}
	return mapTicketRow(row)
}

func (r *TicketRepo) ListByBuyerID(ctx context.Context, buyerID valueobject.UUID, limit, offset int) ([]entity.Ticket, error) {
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	q := `
		SELECT id, ticket_type_id, buyer_id, purchase_date, status, qr_code, amount_paid, used_at, created_at, updated_at
		FROM tickets
		WHERE buyer_id = $1
		ORDER BY purchase_date DESC
		LIMIT $2 OFFSET $3
	`
	var rows []dto.TicketRow
	if err := r.db.SelectContext(ctx, &rows, q, buyerID.String(), limit, offset); err != nil {
		return nil, apperror.New(apperror.CodeInternal, "list tickets failed", err)
	}
	out := make([]entity.Ticket, 0, len(rows))
	for _, row := range rows {
		t, err := mapTicketRow(row)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

func (r *TicketRepo) UpdateStatus(ctx context.Context, id valueobject.UUID, status string) error {
	st := valueobject.TicketStatus(status)
	if err := st.Validate(); err != nil {
		return apperror.New(apperror.CodeValidation, "invalid ticket status", err)
	}
	q := `UPDATE tickets SET status = $1 WHERE id = $2`
	res, err := r.db.ExecContext(ctx, q, status, id.String())
	if err != nil {
		return apperror.New(apperror.CodeInternal, "update ticket status failed", err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return apperror.New(apperror.CodeNotFound, "ticket not found", sql.ErrNoRows)
	}
	return nil
}

func (r *TicketRepo) Delete(ctx context.Context, id valueobject.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM tickets WHERE id=$1`, id.String())
	if err != nil {
		return apperror.New(apperror.CodeInternal, "delete ticket failed", err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return apperror.New(apperror.CodeNotFound, "ticket not found", sql.ErrNoRows)
	}
	return nil
}

func mapTicketRow(row dto.TicketRow) (entity.Ticket, error) {
	tid, err := valueobject.ParseUUID(row.ID)
	if err != nil {
		return entity.Ticket{}, apperror.New(apperror.CodeInternal, "invalid ticket id in db", err)
	}
	ttid, err := valueobject.ParseUUID(row.TicketTypeID)
	if err != nil {
		return entity.Ticket{}, apperror.New(apperror.CodeInternal, "invalid ticket_type_id in db", err)
	}
	bid, err := valueobject.ParseUUID(row.BuyerID)
	if err != nil {
		return entity.Ticket{}, apperror.New(apperror.CodeInternal, "invalid buyer_id in db", err)
	}
	st := valueobject.TicketStatus(row.Status)
	if err := st.Validate(); err != nil {
		return entity.Ticket{}, apperror.New(apperror.CodeInternal, "invalid ticket status in db", err)
	}
	amt, err := decimal.NewFromString(row.AmountPaid)
	if err != nil {
		return entity.Ticket{}, fmt.Errorf("invalid amount_paid in db: %w", err)
	}
	paid, err := valueobject.NewMoney(amt)
	if err != nil {
		return entity.Ticket{}, fmt.Errorf("invalid money in db: %w", err)
	}
	t := entity.Ticket{
		ID:           tid,
		TicketTypeID: ttid,
		BuyerID:      bid,
		Status:       st,
		QRCode:       row.QRCode,
		AmountPaid:   paid,
	}
	if row.PurchaseDate.Valid {
		t.PurchaseDate = row.PurchaseDate.Time
	}
	if row.UsedAt.Valid {
		tt := row.UsedAt.Time
		t.UsedAt = &tt
	}
	if row.CreatedAt.Valid {
		t.CreatedAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		t.UpdatedAt = row.UpdatedAt.Time
	}
	return t, nil
}

type RegistrationRepo struct{ db *sqlx.DB }

func NewRegistrationRepo(db *sqlx.DB) *RegistrationRepo { return &RegistrationRepo{db: db} }

var _ repository.RegistrationRepository = (*RegistrationRepo)(nil)

func (r *RegistrationRepo) Create(ctx context.Context, reg entity.Registration) (valueobject.UUID, error) {
	q := `
		INSERT INTO registrations (user_id, event_id, status, registered_at, attendance_confirmed, notes)
		VALUES ($1, $2, $3, COALESCE($4, NOW()), $5, NULLIF($6, ''))
		RETURNING id
	`
	var id string
	var regAt any
	if !reg.RegisteredAt.IsZero() {
		regAt = reg.RegisteredAt
	}
	if err := r.db.QueryRowxContext(ctx, q, reg.UserID.String(), reg.EventID.String(), string(reg.Status), regAt, reg.AttendanceConfirmed, reg.Notes).Scan(&id); err != nil {
		return valueobject.Nil, apperror.New(apperror.CodeInternal, "create registration failed", err)
	}
	out, err := valueobject.ParseUUID(id)
	if err != nil {
		return valueobject.Nil, apperror.New(apperror.CodeInternal, "invalid uuid returned from db", err)
	}
	return out, nil
}

func (r *RegistrationRepo) GetByID(ctx context.Context, id valueobject.UUID) (entity.Registration, error) {
	q := `
		SELECT id, user_id, event_id, status, registered_at, attendance_confirmed, notes, created_at, updated_at
		FROM registrations
		WHERE id = $1
	`
	var row dto.RegistrationRow
	if err := r.db.GetContext(ctx, &row, q, id.String()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Registration{}, apperror.New(apperror.CodeNotFound, "registration not found", err)
		}
		return entity.Registration{}, apperror.New(apperror.CodeInternal, "get registration failed", err)
	}
	return mapRegistrationRow(row)
}

func (r *RegistrationRepo) ListByEventID(ctx context.Context, eventID valueobject.UUID, limit, offset int) ([]entity.Registration, error) {
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	q := `
		SELECT id, user_id, event_id, status, registered_at, attendance_confirmed, notes, created_at, updated_at
		FROM registrations
		WHERE event_id = $1
		ORDER BY registered_at DESC
		LIMIT $2 OFFSET $3
	`
	var rows []dto.RegistrationRow
	if err := r.db.SelectContext(ctx, &rows, q, eventID.String(), limit, offset); err != nil {
		return nil, apperror.New(apperror.CodeInternal, "list registrations failed", err)
	}
	out := make([]entity.Registration, 0, len(rows))
	for _, row := range rows {
		reg, err := mapRegistrationRow(row)
		if err != nil {
			return nil, err
		}
		out = append(out, reg)
	}
	return out, nil
}

func (r *RegistrationRepo) Update(ctx context.Context, reg entity.Registration) error {
	q := `
		UPDATE registrations
		SET status=$1, attendance_confirmed=$2, notes=NULLIF($3,'')
		WHERE id=$4
	`
	res, err := r.db.ExecContext(ctx, q, string(reg.Status), reg.AttendanceConfirmed, reg.Notes, reg.ID.String())
	if err != nil {
		return apperror.New(apperror.CodeInternal, "update registration failed", err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return apperror.New(apperror.CodeNotFound, "registration not found", sql.ErrNoRows)
	}
	return nil
}

func (r *RegistrationRepo) Delete(ctx context.Context, id valueobject.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM registrations WHERE id=$1`, id.String())
	if err != nil {
		return apperror.New(apperror.CodeInternal, "delete registration failed", err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return apperror.New(apperror.CodeNotFound, "registration not found", sql.ErrNoRows)
	}
	return nil
}

func mapRegistrationRow(row dto.RegistrationRow) (entity.Registration, error) {
	rid, err := valueobject.ParseUUID(row.ID)
	if err != nil {
		return entity.Registration{}, apperror.New(apperror.CodeInternal, "invalid registration id in db", err)
	}
	uid, err := valueobject.ParseUUID(row.UserID)
	if err != nil {
		return entity.Registration{}, apperror.New(apperror.CodeInternal, "invalid registration user_id in db", err)
	}
	eid, err := valueobject.ParseUUID(row.EventID)
	if err != nil {
		return entity.Registration{}, apperror.New(apperror.CodeInternal, "invalid registration event_id in db", err)
	}
	st := valueobject.RegistrationStatus(row.Status)
	if err := st.Validate(); err != nil {
		return entity.Registration{}, apperror.New(apperror.CodeInternal, "invalid registration status in db", err)
	}
	reg := entity.Registration{
		ID:                  rid,
		UserID:              uid,
		EventID:             eid,
		Status:              st,
		AttendanceConfirmed: row.AttendanceConfirmed,
		Notes:               "",
	}
	if row.RegisteredAt.Valid {
		reg.RegisteredAt = row.RegisteredAt.Time
	}
	if row.Notes.Valid {
		reg.Notes = row.Notes.String
	}
	if row.CreatedAt.Valid {
		reg.CreatedAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		reg.UpdatedAt = row.UpdatedAt.Time
	}
	return reg, nil
}

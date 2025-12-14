package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"time2meet/internal/application/port/tickettx"
	"time2meet/internal/domain/valueobject"
	"time2meet/pkg/apperror"

	"github.com/jmoiron/sqlx"
)

type TicketTxQueries struct{}

func NewTicketTxQueries() *TicketTxQueries { return &TicketTxQueries{} }

var _ tickettx.Queries = (*TicketTxQueries)(nil)

func (q *TicketTxQueries) LockTicketTypeForUpdate(ctx context.Context, tx *sqlx.Tx, ticketTypeID valueobject.UUID) (qtyTotal int, qtySold int, err error) {
	lockQ := `SELECT quantity_total, quantity_sold FROM ticket_types WHERE id = $1 FOR UPDATE`
	if err := tx.QueryRowxContext(ctx, lockQ, ticketTypeID.String()).Scan(&qtyTotal, &qtySold); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, apperror.New(apperror.CodeNotFound, "ticket type not found", err)
		}
		return 0, 0, apperror.New(apperror.CodeInternal, "lock ticket type failed", err)
	}
	return qtyTotal, qtySold, nil
}

func (q *TicketTxQueries) InsertPaidTicket(ctx context.Context, tx *sqlx.Tx, ticketTypeID, buyerID valueobject.UUID, purchaseDate time.Time, qrCode string, amountPaid string) (valueobject.UUID, error) {
	insQ := `
		INSERT INTO tickets (ticket_type_id, buyer_id, purchase_date, status, qr_code, amount_paid)
		VALUES ($1, $2, $3, 'paid', $4, $5)
		RETURNING id
	`
	var id string
	if err := tx.QueryRowxContext(ctx, insQ, ticketTypeID.String(), buyerID.String(), purchaseDate, qrCode, amountPaid).Scan(&id); err != nil {
		return valueobject.Nil, apperror.New(apperror.CodeInternal, "insert ticket failed", err)
	}
	uid, err := valueobject.ParseUUID(id)
	if err != nil {
		return valueobject.Nil, apperror.New(apperror.CodeInternal, "invalid uuid returned from db", err)
	}
	return uid, nil
}

func (q *TicketTxQueries) MarkTicketUsed(ctx context.Context, tx *sqlx.Tx, ticketID valueobject.UUID, usedAt time.Time) (bool, error) {
	upd := `
		UPDATE tickets
		SET status = 'used', used_at = $1
		WHERE id = $2 AND status IN ('paid', 'used')
	`
	res, err := tx.ExecContext(ctx, upd, usedAt, ticketID.String())
	if err != nil {
		return false, apperror.New(apperror.CodeInternal, "validate ticket failed", err)
	}
	aff, _ := res.RowsAffected()
	return aff > 0, nil
}



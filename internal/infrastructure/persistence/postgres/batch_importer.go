package postgres

import (
	"context"

	"time2meet/internal/application/port/batchimport"
	"time2meet/pkg/apperror"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type BatchImporter struct {
	log *zap.Logger
}

func NewBatchImporter(log *zap.Logger) *BatchImporter { return &BatchImporter{log: log} }

var _ batchimport.Importer = (*BatchImporter)(nil)

func (bi *BatchImporter) ImportUsers(ctx context.Context, tx *sqlx.Tx, items []batchimport.ImportUsersItem, continueOnError bool) (batchimport.Result, error) {
	res := batchimport.Result{Total: len(items), Errors: []batchimport.BatchError{}}
	for i, it := range items {
		if _, err := tx.ExecContext(ctx, "SAVEPOINT sp_item"); err != nil {
			return res, apperror.New(apperror.CodeInternal, "savepoint failed", err)
		}
		q := `
			INSERT INTO users (email, password_hash, full_name, phone, role, is_active)
			VALUES ($1, $2, $3, NULLIF($4,''), $5, true)
		`
		_, err := tx.ExecContext(ctx, q, it.Email, it.PasswordHash, it.FullName, it.Phone, it.Role)
		if err != nil {
			bi.log.Warn("batch import users row failed", zap.Int("index", i), zap.Error(err))
			_, _ = tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT sp_item")
			_, _ = tx.ExecContext(ctx, "RELEASE SAVEPOINT sp_item")
			res.Failed++
			res.Errors = append(res.Errors, batchimport.BatchError{Index: i, Error: err.Error()})
			if !continueOnError {
				return res, apperror.New(apperror.CodeConflict, "batch import users failed", err)
			}
			continue
		}
		_, _ = tx.ExecContext(ctx, "RELEASE SAVEPOINT sp_item")
		res.Success++
	}
	return res, nil
}

func (bi *BatchImporter) ImportEvents(ctx context.Context, tx *sqlx.Tx, items []batchimport.ImportEventsItem, continueOnError bool) (batchimport.Result, error) {
	res := batchimport.Result{Total: len(items), Errors: []batchimport.BatchError{}}
	for i, it := range items {
		if _, err := tx.ExecContext(ctx, "SAVEPOINT sp_item"); err != nil {
			return res, apperror.New(apperror.CodeInternal, "savepoint failed", err)
		}
		var maxp any
		if it.MaxParticipants != nil {
			maxp = *it.MaxParticipants
		}
		q := `
			INSERT INTO events (organizer_id, title, description, status, is_public, max_participants, cover_image)
			VALUES ($1, $2, NULLIF($3,''), $4, $5, $6, NULLIF($7,''))
		`
		_, err := tx.ExecContext(ctx, q, it.OrganizerID, it.Title, it.Description, it.Status, it.IsPublic, maxp, it.CoverImage)
		if err != nil {
			bi.log.Warn("batch import events row failed", zap.Int("index", i), zap.Error(err))
			_, _ = tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT sp_item")
			_, _ = tx.ExecContext(ctx, "RELEASE SAVEPOINT sp_item")
			res.Failed++
			res.Errors = append(res.Errors, batchimport.BatchError{Index: i, Error: err.Error()})
			if !continueOnError {
				return res, apperror.New(apperror.CodeConflict, "batch import events failed", err)
			}
			continue
		}
		_, _ = tx.ExecContext(ctx, "RELEASE SAVEPOINT sp_item")
		res.Success++
	}
	return res, nil
}

func (bi *BatchImporter) ImportTickets(ctx context.Context, tx *sqlx.Tx, items []batchimport.ImportTicketsItem, continueOnError bool) (batchimport.Result, error) {
	res := batchimport.Result{Total: len(items), Errors: []batchimport.BatchError{}}
	for i, it := range items {
		if _, err := tx.ExecContext(ctx, "SAVEPOINT sp_item"); err != nil {
			return res, apperror.New(apperror.CodeInternal, "savepoint failed", err)
		}
		var pd any
		if it.PurchaseDate != nil {
			pd = *it.PurchaseDate
		}
		q := `
			INSERT INTO tickets (ticket_type_id, buyer_id, purchase_date, status, qr_code, amount_paid)
			VALUES ($1, $2, COALESCE($3, NOW()), $4, $5, $6)
		`
		_, err := tx.ExecContext(ctx, q, it.TicketTypeID, it.BuyerID, pd, it.Status, it.QRCode, it.AmountPaid)
		if err != nil {
			bi.log.Warn("batch import tickets row failed", zap.Int("index", i), zap.Error(err))
			_, _ = tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT sp_item")
			_, _ = tx.ExecContext(ctx, "RELEASE SAVEPOINT sp_item")
			res.Failed++
			res.Errors = append(res.Errors, batchimport.BatchError{Index: i, Error: err.Error()})
			if !continueOnError {
				return res, apperror.New(apperror.CodeConflict, "batch import tickets failed", err)
			}
			continue
		}
		_, _ = tx.ExecContext(ctx, "RELEASE SAVEPOINT sp_item")
		res.Success++
	}
	return res, nil
}



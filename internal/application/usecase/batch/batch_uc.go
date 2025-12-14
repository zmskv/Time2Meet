package batch

import (
	"context"

	"time2meet/internal/application/port/auditctx"
	"time2meet/internal/application/port/batchimport"
	"time2meet/internal/application/tx"
	"time2meet/internal/domain/valueobject"

	"github.com/jmoiron/sqlx"
)

type UseCase struct {
	tx    tx.Manager
	audit auditctx.Setter
	imp   batchimport.Importer
}

func New(txm tx.Manager, audit auditctx.Setter, imp batchimport.Importer) *UseCase {
	return &UseCase{tx: txm, audit: audit, imp: imp}
}

type ImportUsersInput struct {
	UserID          valueobject.UUID
	IP              string
	ContinueOnError bool
	Items           []batchimport.ImportUsersItem
}

func (uc *UseCase) ImportUsers(ctx context.Context, in ImportUsersInput) (batchimport.Result, error) {
	var res batchimport.Result
	err := uc.tx.WithTx(ctx, func(ctx context.Context, txx *sqlx.Tx) error {
		if err := uc.audit.Set(ctx, txx, in.UserID, in.IP); err != nil {
			return err
		}
		out, err := uc.imp.ImportUsers(ctx, txx, in.Items, in.ContinueOnError)
		res = out
		return err
	})
	return res, err
}

type ImportEventsInput struct {
	UserID          valueobject.UUID
	IP              string
	ContinueOnError bool
	Items           []batchimport.ImportEventsItem
}

func (uc *UseCase) ImportEvents(ctx context.Context, in ImportEventsInput) (batchimport.Result, error) {
	var res batchimport.Result
	err := uc.tx.WithTx(ctx, func(ctx context.Context, txx *sqlx.Tx) error {
		if err := uc.audit.Set(ctx, txx, in.UserID, in.IP); err != nil {
			return err
		}
		out, err := uc.imp.ImportEvents(ctx, txx, in.Items, in.ContinueOnError)
		res = out
		return err
	})
	return res, err
}

type ImportTicketsInput struct {
	UserID          valueobject.UUID
	IP              string
	ContinueOnError bool
	Items           []batchimport.ImportTicketsItem
}

func (uc *UseCase) ImportTickets(ctx context.Context, in ImportTicketsInput) (batchimport.Result, error) {
	var res batchimport.Result
	err := uc.tx.WithTx(ctx, func(ctx context.Context, txx *sqlx.Tx) error {
		if err := uc.audit.Set(ctx, txx, in.UserID, in.IP); err != nil {
			return err
		}
		out, err := uc.imp.ImportTickets(ctx, txx, in.Items, in.ContinueOnError)
		res = out
		return err
	})
	return res, err
}

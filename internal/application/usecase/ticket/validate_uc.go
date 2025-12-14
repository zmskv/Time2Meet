package ticket

import (
	"context"
	"time"

	"time2meet/internal/application/port/auditctx"
	"time2meet/internal/application/port/tickettx"
	"time2meet/internal/application/tx"
	"time2meet/internal/domain/valueobject"
	"time2meet/pkg/apperror"

	"github.com/jmoiron/sqlx"
)

type ValidateUseCase struct {
	tx    tx.Manager
	audit auditctx.Setter
	q     tickettx.Queries
}

func NewValidate(txm tx.Manager, audit auditctx.Setter, q tickettx.Queries) *ValidateUseCase {
	return &ValidateUseCase{tx: txm, audit: audit, q: q}
}

type ValidateInput struct {
	UserID   valueobject.UUID
	IP       string
	TicketID valueobject.UUID
}

func (uc *ValidateUseCase) Validate(ctx context.Context, in ValidateInput) error {
	if in.UserID == valueobject.Nil {
		return apperror.New(apperror.CodeValidation, "user_id is required", nil)
	}
	if in.TicketID == valueobject.Nil {
		return apperror.New(apperror.CodeValidation, "ticket_id is required", nil)
	}

	return uc.tx.WithTx(ctx, func(ctx context.Context, txx *sqlx.Tx) error {
		if err := uc.audit.Set(ctx, txx, in.UserID, in.IP); err != nil {
			return err
		}
		ok, err := uc.q.MarkTicketUsed(ctx, txx, in.TicketID, time.Now().UTC())
		if err != nil {
			return err
		}
		if !ok {
			return apperror.New(apperror.CodeNotFound, "ticket not found or invalid state", nil)
		}
		return nil
	})
}



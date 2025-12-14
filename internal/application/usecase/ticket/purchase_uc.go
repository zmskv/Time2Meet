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
	"github.com/shopspring/decimal"
)

type PurchaseUseCase struct {
	tx    tx.Manager
	audit auditctx.Setter
	q     tickettx.Queries
}

func NewPurchase(txm tx.Manager, audit auditctx.Setter, q tickettx.Queries) *PurchaseUseCase {
	return &PurchaseUseCase{tx: txm, audit: audit, q: q}
}

type PurchaseInput struct {
	UserID       valueobject.UUID
	IP           string
	TicketTypeID valueobject.UUID
	QRCode       string
	AmountPaid   string // decimal string
	Currency     string
}

type PurchaseOutput struct {
	TicketID valueobject.UUID
}

func (uc *PurchaseUseCase) Purchase(ctx context.Context, in PurchaseInput) (PurchaseOutput, error) {
	if in.UserID == valueobject.Nil {
		return PurchaseOutput{}, apperror.New(apperror.CodeValidation, "user_id is required", nil)
	}
	if in.TicketTypeID == valueobject.Nil {
		return PurchaseOutput{}, apperror.New(apperror.CodeValidation, "ticket_type_id is required", nil)
	}
	if in.QRCode == "" {
		return PurchaseOutput{}, apperror.New(apperror.CodeValidation, "qr_code is required", nil)
	}
	amt, err := decimal.NewFromString(in.AmountPaid)
	if err != nil {
		return PurchaseOutput{}, apperror.New(apperror.CodeValidation, "amount_paid must be decimal string", err)
	}
	money, err := valueobject.NewMoney(amt)
	if err != nil {
		return PurchaseOutput{}, apperror.New(apperror.CodeValidation, "invalid money", err)
	}

	var out PurchaseOutput
	err = uc.tx.WithTx(ctx, func(ctx context.Context, txx *sqlx.Tx) error {
		if err := uc.audit.Set(ctx, txx, in.UserID, in.IP); err != nil {
			return err
		}

		qtyTotal, qtySold, err := uc.q.LockTicketTypeForUpdate(ctx, txx, in.TicketTypeID)
		if err != nil {
			return err
		}
		if qtySold >= qtyTotal {
			return apperror.New(apperror.CodeConflict, "sold out", nil)
		}

		ticketID, err := uc.q.InsertPaidTicket(ctx, txx, in.TicketTypeID, in.UserID, time.Now().UTC(), in.QRCode, money.Amount.StringFixed(2))
		if err != nil {
			return err
		}
		out.TicketID = ticketID
		return nil
	})
	if err != nil {
		return PurchaseOutput{}, err
	}
	return out, nil
}



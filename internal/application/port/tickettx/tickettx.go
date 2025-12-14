package tickettx

import (
	"context"
	"time"

	"time2meet/internal/domain/valueobject"

	"github.com/jmoiron/sqlx"
)

type Queries interface {
	LockTicketTypeForUpdate(ctx context.Context, tx *sqlx.Tx, ticketTypeID valueobject.UUID) (qtyTotal int, qtySold int, err error)

	InsertPaidTicket(ctx context.Context, tx *sqlx.Tx, ticketTypeID, buyerID valueobject.UUID, purchaseDate time.Time, qrCode string, amountPaid string) (valueobject.UUID, error)

	MarkTicketUsed(ctx context.Context, tx *sqlx.Tx, ticketID valueobject.UUID, usedAt time.Time) (bool, error)
}

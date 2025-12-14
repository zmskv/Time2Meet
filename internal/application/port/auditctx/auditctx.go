package auditctx

import (
	"context"

	"time2meet/internal/domain/valueobject"

	"github.com/jmoiron/sqlx"
)

type Setter interface {
	Set(ctx context.Context, tx *sqlx.Tx, userID valueobject.UUID, ip string) error
}

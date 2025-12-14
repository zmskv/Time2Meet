package tx

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Manager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context, tx *sqlx.Tx) error) error
}



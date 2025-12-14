package postgres

import (
	"context"

	"time2meet/internal/application/tx"
	"time2meet/pkg/apperror"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type TxManager struct {
	db  *sqlx.DB
	log *zap.Logger
}

func NewTxManager(db *sqlx.DB, log *zap.Logger) *TxManager {
	return &TxManager{db: db, log: log}
}

var _ tx.Manager = (*TxManager)(nil)

func (m *TxManager) WithTx(ctx context.Context, fn func(ctx context.Context, tx *sqlx.Tx) error) error {
	txx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return apperror.New(apperror.CodeUnavailable, "begin tx failed", err)
	}
	defer func() { _ = txx.Rollback() }()

	if err := fn(ctx, txx); err != nil {
		return err
	}
	if err := txx.Commit(); err != nil {
		return apperror.New(apperror.CodeUnavailable, "commit failed", err)
	}
	return nil
}

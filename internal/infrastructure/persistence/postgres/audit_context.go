package postgres

import (
	"context"

	"time2meet/internal/application/port/auditctx"
	"time2meet/internal/domain/valueobject"
	"time2meet/pkg/apperror"

	"github.com/jmoiron/sqlx"
)

type AuditContextSetter struct{}

func NewAuditContextSetter() *AuditContextSetter { return &AuditContextSetter{} }

var _ auditctx.Setter = (*AuditContextSetter)(nil)

func (s *AuditContextSetter) Set(ctx context.Context, tx *sqlx.Tx, userID valueobject.UUID, ip string) error {
	if userID != valueobject.Nil {
		if _, err := tx.ExecContext(ctx, `SELECT set_config('app.user_id', $1, true)`, userID.String()); err != nil {
			return apperror.New(apperror.CodeInternal, "set audit user_id failed", err)
		}
	}
	if ip != "" {
		if _, err := tx.ExecContext(ctx, `SELECT set_config('app.ip', $1, true)`, ip); err != nil {
			return apperror.New(apperror.CodeInternal, "set audit ip failed", err)
		}
	}
	return nil
}



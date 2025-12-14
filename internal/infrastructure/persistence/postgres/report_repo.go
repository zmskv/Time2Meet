package postgres

import (
	"context"
	"time"

	"time2meet/internal/domain/repository"
	"time2meet/internal/domain/valueobject"
	"time2meet/internal/infrastructure/persistence/postgres/dto"
	"time2meet/pkg/apperror"

	"github.com/jmoiron/sqlx"
)

type ReportRepo struct{ db *sqlx.DB }

func NewReportRepo(db *sqlx.DB) *ReportRepo { return &ReportRepo{db: db} }

var _ repository.ReportRepository = (*ReportRepo)(nil)

func (r *ReportRepo) SalesReport(ctx context.Context, start, end time.Time) ([]repository.SalesReportRow, error) {
	q := `SELECT event_id, event_title, tickets_sold, revenue, unique_buyers FROM get_sales_report($1::date, $2::date)`
	var rows []dto.SalesReportRow
	if err := r.db.SelectContext(ctx, &rows, q, start, end); err != nil {
		return nil, apperror.New(apperror.CodeInternal, "sales report failed", err)
	}
	out := make([]repository.SalesReportRow, 0, len(rows))
	for _, row := range rows {
		eid, err := valueobject.ParseUUID(row.EventID)
		if err != nil {
			return nil, apperror.New(apperror.CodeInternal, "invalid event_id in report row", err)
		}
		out = append(out, repository.SalesReportRow{
			EventID:      eid,
			EventTitle:   row.EventTitle,
			TicketsSold:  row.TicketsSold,
			Revenue:      row.Revenue,
			UniqueBuyers: row.UniqueBuyers,
		})
	}
	return out, nil
}

func (r *ReportRepo) AttendanceStats(ctx context.Context, eventID valueobject.UUID) ([]repository.AttendanceRow, error) {
	q := `SELECT ticket_type, sold, used, attendance_rate FROM get_attendance_stats($1)`
	var rows []dto.AttendanceRow
	if err := r.db.SelectContext(ctx, &rows, q, eventID.String()); err != nil {
		return nil, apperror.New(apperror.CodeInternal, "attendance stats failed", err)
	}
	out := make([]repository.AttendanceRow, 0, len(rows))
	for _, row := range rows {
		out = append(out, repository.AttendanceRow{
			TicketType:     row.TicketType,
			Sold:           row.Sold,
			Used:           row.Used,
			AttendanceRate: row.AttendanceRate,
		})
	}
	return out, nil
}

func (r *ReportRepo) PopularEvents(ctx context.Context, limit int, days int) ([]repository.PopularEventRow, error) {
	q := `SELECT event_id, title, registrations, tickets_sold FROM get_popular_events($1, $2)`
	var rows []dto.PopularEventRow
	if err := r.db.SelectContext(ctx, &rows, q, limit, days); err != nil {
		return nil, apperror.New(apperror.CodeInternal, "popular events failed", err)
	}
	out := make([]repository.PopularEventRow, 0, len(rows))
	for _, row := range rows {
		eid, err := valueobject.ParseUUID(row.EventID)
		if err != nil {
			return nil, apperror.New(apperror.CodeInternal, "invalid event_id in popular row", err)
		}
		out = append(out, repository.PopularEventRow{
			EventID:       eid,
			Title:         row.Title,
			Registrations: row.Registrations,
			TicketsSold:   row.TicketsSold,
		})
	}
	return out, nil
}

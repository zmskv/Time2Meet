package report

import (
	"context"
	"time"

	"time2meet/internal/domain/repository"
	"time2meet/internal/domain/valueobject"
)

type UseCase struct {
	reports repository.ReportRepository
}

func New(reports repository.ReportRepository) *UseCase { return &UseCase{reports: reports} }

func (uc *UseCase) Sales(ctx context.Context, start, end time.Time) ([]repository.SalesReportRow, error) {
	return uc.reports.SalesReport(ctx, start, end)
}

func (uc *UseCase) Attendance(ctx context.Context, eventID valueobject.UUID) ([]repository.AttendanceRow, error) {
	return uc.reports.AttendanceStats(ctx, eventID)
}

func (uc *UseCase) Popular(ctx context.Context, limit, days int) ([]repository.PopularEventRow, error) {
	return uc.reports.PopularEvents(ctx, limit, days)
}



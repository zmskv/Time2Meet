package repository

import (
	"context"
	"time"

	"time2meet/internal/domain/valueobject"
)

type SalesReportRow struct {
	EventID      valueobject.UUID
	EventTitle   string
	TicketsSold  int64
	Revenue      string
	UniqueBuyers int64
}

type AttendanceRow struct {
	TicketType     string
	Sold           int64
	Used           int64
	AttendanceRate string
}

type PopularEventRow struct {
	EventID       valueobject.UUID
	Title         string
	Registrations int64
	TicketsSold   int64
}

type ReportRepository interface {
	SalesReport(ctx context.Context, start, end time.Time) ([]SalesReportRow, error)
	AttendanceStats(ctx context.Context, eventID valueobject.UUID) ([]AttendanceRow, error)
	PopularEvents(ctx context.Context, limit int, days int) ([]PopularEventRow, error)
}

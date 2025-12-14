package dto

type SalesReportRow struct {
	EventID      string `db:"event_id"`
	EventTitle   string `db:"event_title"`
	TicketsSold  int64  `db:"tickets_sold"`
	Revenue      string `db:"revenue"`
	UniqueBuyers int64  `db:"unique_buyers"`
}

type AttendanceRow struct {
	TicketType     string `db:"ticket_type"`
	Sold           int64  `db:"sold"`
	Used           int64  `db:"used"`
	AttendanceRate string `db:"attendance_rate"`
}

type PopularEventRow struct {
	EventID       string `db:"event_id"`
	Title         string `db:"title"`
	Registrations int64  `db:"registrations"`
	TicketsSold   int64  `db:"tickets_sold"`
}



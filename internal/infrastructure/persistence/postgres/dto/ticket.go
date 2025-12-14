package dto

import "database/sql"

type TicketTypeRow struct {
	ID            string         `db:"id"`
	EventID       string         `db:"event_id"`
	Name          string         `db:"name"`
	Price         string         `db:"price"`
	Currency      string         `db:"currency"`
	QuantityTotal int            `db:"quantity_total"`
	QuantitySold  int            `db:"quantity_sold"`
	SaleStart     sql.NullTime   `db:"sale_start"`
	SaleEnd       sql.NullTime   `db:"sale_end"`
	Description   sql.NullString `db:"description"`
	IsActive      bool           `db:"is_active"`
	CreatedAt     sql.NullTime   `db:"created_at"`
	UpdatedAt     sql.NullTime   `db:"updated_at"`
}

type TicketRow struct {
	ID           string       `db:"id"`
	TicketTypeID string       `db:"ticket_type_id"`
	BuyerID      string       `db:"buyer_id"`
	PurchaseDate sql.NullTime `db:"purchase_date"`
	Status       string       `db:"status"`
	QRCode       string       `db:"qr_code"`
	AmountPaid   string       `db:"amount_paid"`
	UsedAt       sql.NullTime `db:"used_at"`
	CreatedAt    sql.NullTime `db:"created_at"`
	UpdatedAt    sql.NullTime `db:"updated_at"`
}

type RegistrationRow struct {
	ID                  string         `db:"id"`
	UserID              string         `db:"user_id"`
	EventID             string         `db:"event_id"`
	Status              string         `db:"status"`
	RegisteredAt        sql.NullTime   `db:"registered_at"`
	AttendanceConfirmed bool           `db:"attendance_confirmed"`
	Notes               sql.NullString `db:"notes"`
	CreatedAt           sql.NullTime   `db:"created_at"`
	UpdatedAt           sql.NullTime   `db:"updated_at"`
}



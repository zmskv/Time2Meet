package batchimport

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type BatchError struct {
	Index int    `json:"index"`
	Error string `json:"error"`
}

type Result struct {
	Total   int          `json:"total"`
	Success int          `json:"success"`
	Failed  int          `json:"failed"`
	Errors  []BatchError `json:"errors"`
}

type ImportUsersItem struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	FullName     string `json:"full_name"`
	Phone        string `json:"phone"`
	Role         string `json:"role"`
}

type ImportEventsItem struct {
	OrganizerID     string `json:"organizer_id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Status          string `json:"status"`
	IsPublic        bool   `json:"is_public"`
	MaxParticipants *int   `json:"max_participants"`
	CoverImage      string `json:"cover_image"`
}

type ImportTicketsItem struct {
	TicketTypeID string     `json:"ticket_type_id"`
	BuyerID      string     `json:"buyer_id"`
	PurchaseDate *time.Time `json:"purchase_date"`
	Status       string     `json:"status"`
	QRCode       string     `json:"qr_code"`
	AmountPaid   string     `json:"amount_paid"`
}

type Importer interface {
	ImportUsers(ctx context.Context, tx *sqlx.Tx, items []ImportUsersItem, continueOnError bool) (Result, error)
	ImportEvents(ctx context.Context, tx *sqlx.Tx, items []ImportEventsItem, continueOnError bool) (Result, error)
	ImportTickets(ctx context.Context, tx *sqlx.Tx, items []ImportTicketsItem, continueOnError bool) (Result, error)
}

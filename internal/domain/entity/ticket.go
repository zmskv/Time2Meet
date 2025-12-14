package entity

import (
	"time"

	"time2meet/internal/domain/valueobject"
)

type TicketType struct {
	ID            valueobject.UUID
	EventID        valueobject.UUID
	Name           string
	Price          valueobject.Money
	QuantityTotal  int
	QuantitySold   int
	SaleStart      *time.Time
	SaleEnd        *time.Time
	Description    string
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Ticket struct {
	ID           valueobject.UUID
	TicketTypeID valueobject.UUID
	BuyerID      valueobject.UUID
	PurchaseDate time.Time
	Status       valueobject.TicketStatus
	QRCode       string
	AmountPaid   valueobject.Money
	UsedAt       *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Registration struct {
	ID                  valueobject.UUID
	UserID              valueobject.UUID
	EventID             valueobject.UUID
	Status              valueobject.RegistrationStatus
	RegisteredAt         time.Time
	AttendanceConfirmed bool
	Notes               string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}


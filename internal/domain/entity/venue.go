package entity

import (
	"time"

	"time2meet/internal/domain/valueobject"
)

type Venue struct {
	ID           valueobject.UUID
	Name         string
	Address      string
	City         string
	Country      string
	Capacity     int
	ContactPhone string
	ContactEmail string
	Website      string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Room struct {
	ID          valueobject.UUID
	VenueID     valueobject.UUID
	Name        string
	Capacity    int
	Floor       *int
	Equipment   map[string]any
	HourlyRate  string
	IsAvailable bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}


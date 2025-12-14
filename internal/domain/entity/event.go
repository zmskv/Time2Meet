package entity

import (
	"time"

	"time2meet/internal/domain/valueobject"
)

type Event struct {
	ID              valueobject.UUID
	OrganizerID     valueobject.UUID
	Title           string
	Description     string
	Status          valueobject.EventStatus
	IsPublic        bool
	MaxParticipants *int
	CoverImage      string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Category struct {
	ID          valueobject.UUID
	Name        string
	Description string
	Icon        string
	Slug        string
	ParentID    *valueobject.UUID
	SortOrder   int
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type EventSchedule struct {
	ID        valueobject.UUID
	EventID   valueobject.UUID
	RoomID    valueobject.UUID
	StartTime time.Time
	EndTime   time.Time
	Status    string
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}


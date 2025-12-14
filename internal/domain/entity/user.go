package entity

import (
	"time"

	"time2meet/internal/domain/valueobject"
)

type UserRole string

const (
	UserRoleAdmin     UserRole = "admin"
	UserRoleOrganizer UserRole = "organizer"
	UserRoleAttendee  UserRole = "attendee"
)

type User struct {
	ID           valueobject.UUID
	Email        valueobject.Email
	PasswordHash string
	FullName     string
	Phone        string
	Role         UserRole
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserProfile struct {
	ID          valueobject.UUID
	UserID      valueobject.UUID
	AvatarURL   string
	Bio         string
	SocialLinks map[string]any
	Preferences map[string]any
	UpdatedAt   time.Time
}


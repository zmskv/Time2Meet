package valueobject

import "fmt"

type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusPublished EventStatus = "published"
	EventStatusCancelled EventStatus = "cancelled"
	EventStatusCompleted EventStatus = "completed"
)

func (s EventStatus) Validate() error {
	switch s {
	case EventStatusDraft, EventStatusPublished, EventStatusCancelled, EventStatusCompleted:
		return nil
	default:
		return fmt.Errorf("invalid event status: %q", s)
	}
}

type TicketStatus string

const (
	TicketStatusPaid     TicketStatus = "paid"
	TicketStatusRefunded TicketStatus = "refunded"
	TicketStatusVoid     TicketStatus = "void"
	TicketStatusUsed     TicketStatus = "used"
)

func (s TicketStatus) Validate() error {
	switch s {
	case TicketStatusPaid, TicketStatusRefunded, TicketStatusVoid, TicketStatusUsed:
		return nil
	default:
		return fmt.Errorf("invalid ticket status: %q", s)
	}
}

type RegistrationStatus string

const (
	RegistrationStatusRegistered RegistrationStatus = "registered"
	RegistrationStatusCancelled  RegistrationStatus = "cancelled"
	RegistrationStatusAttended   RegistrationStatus = "attended"
	RegistrationStatusNoShow     RegistrationStatus = "no_show"
)

func (s RegistrationStatus) Validate() error {
	switch s {
	case RegistrationStatusRegistered, RegistrationStatusCancelled, RegistrationStatusAttended, RegistrationStatusNoShow:
		return nil
	default:
		return fmt.Errorf("invalid registration status: %q", s)
	}
}

package repository

import (
	"context"

	"time2meet/internal/domain/entity"
	"time2meet/internal/domain/valueobject"
)

type TicketTypeRepository interface {
	Create(ctx context.Context, tt entity.TicketType) (valueobject.UUID, error)
	GetByID(ctx context.Context, id valueobject.UUID) (entity.TicketType, error)
	ListByEventID(ctx context.Context, eventID valueobject.UUID) ([]entity.TicketType, error)
	Update(ctx context.Context, tt entity.TicketType) error
	Delete(ctx context.Context, id valueobject.UUID) error
}

type TicketRepository interface {
	Create(ctx context.Context, t entity.Ticket) (valueobject.UUID, error)
	GetByID(ctx context.Context, id valueobject.UUID) (entity.Ticket, error)
	ListByBuyerID(ctx context.Context, buyerID valueobject.UUID, limit, offset int) ([]entity.Ticket, error)
	UpdateStatus(ctx context.Context, id valueobject.UUID, status string) error
	Delete(ctx context.Context, id valueobject.UUID) error
}

type RegistrationRepository interface {
	Create(ctx context.Context, r entity.Registration) (valueobject.UUID, error)
	GetByID(ctx context.Context, id valueobject.UUID) (entity.Registration, error)
	ListByEventID(ctx context.Context, eventID valueobject.UUID, limit, offset int) ([]entity.Registration, error)
	Update(ctx context.Context, r entity.Registration) error
	Delete(ctx context.Context, id valueobject.UUID) error
}

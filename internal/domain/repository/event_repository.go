package repository

import (
	"context"

	"time2meet/internal/domain/entity"
	"time2meet/internal/domain/valueobject"
)

type EventRepository interface {
	Create(ctx context.Context, e entity.Event) (valueobject.UUID, error)
	GetByID(ctx context.Context, id valueobject.UUID) (entity.Event, error)
	List(ctx context.Context, organizerID *valueobject.UUID, status *string, limit, offset int) ([]entity.Event, error)
	Update(ctx context.Context, e entity.Event) error
	Delete(ctx context.Context, id valueobject.UUID) error
}

type CategoryRepository interface {
	Create(ctx context.Context, c entity.Category) (valueobject.UUID, error)
	GetByID(ctx context.Context, id valueobject.UUID) (entity.Category, error)
	List(ctx context.Context, limit, offset int) ([]entity.Category, error)
	Update(ctx context.Context, c entity.Category) error
	Delete(ctx context.Context, id valueobject.UUID) error
}

type EventScheduleRepository interface {
	Create(ctx context.Context, s entity.EventSchedule) (valueobject.UUID, error)
	GetByID(ctx context.Context, id valueobject.UUID) (entity.EventSchedule, error)
	ListByEventID(ctx context.Context, eventID valueobject.UUID) ([]entity.EventSchedule, error)
	Update(ctx context.Context, s entity.EventSchedule) error
	Delete(ctx context.Context, id valueobject.UUID) error
}

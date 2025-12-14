package repository

import (
	"context"

	"time2meet/internal/domain/entity"
	"time2meet/internal/domain/valueobject"
)

type VenueRepository interface {
	Create(ctx context.Context, v entity.Venue) (valueobject.UUID, error)
	GetByID(ctx context.Context, id valueobject.UUID) (entity.Venue, error)
	List(ctx context.Context, limit, offset int) ([]entity.Venue, error)
	Update(ctx context.Context, v entity.Venue) error
	Delete(ctx context.Context, id valueobject.UUID) error
}

type RoomRepository interface {
	Create(ctx context.Context, r entity.Room) (valueobject.UUID, error)
	GetByID(ctx context.Context, id valueobject.UUID) (entity.Room, error)
	ListByVenueID(ctx context.Context, venueID valueobject.UUID) ([]entity.Room, error)
	Update(ctx context.Context, r entity.Room) error
	Delete(ctx context.Context, id valueobject.UUID) error
}

package venue

import (
	"context"

	"time2meet/internal/domain/entity"
	"time2meet/internal/domain/repository"
	"time2meet/internal/domain/valueobject"
	"time2meet/pkg/apperror"
)

type UseCase struct {
	venues repository.VenueRepository
	rooms  repository.RoomRepository
}

func New(venues repository.VenueRepository, rooms repository.RoomRepository) *UseCase {
	return &UseCase{venues: venues, rooms: rooms}
}

type CreateVenueInput struct {
	Name         string
	Address      string
	City         string
	Country      string
	Capacity     int
	ContactPhone string
	ContactEmail string
	Website      string
}

func (uc *UseCase) CreateVenue(ctx context.Context, in CreateVenueInput) (valueobject.UUID, error) {
	if in.Name == "" || in.Address == "" || in.City == "" {
		return valueobject.Nil, apperror.New(apperror.CodeValidation, "name/address/city are required", nil)
	}
	if in.Capacity < 0 {
		return valueobject.Nil, apperror.New(apperror.CodeValidation, "capacity must be >= 0", nil)
	}
	v := entity.Venue{
		Name:         in.Name,
		Address:      in.Address,
		City:         in.City,
		Country:      in.Country,
		Capacity:     in.Capacity,
		ContactPhone: in.ContactPhone,
		ContactEmail: in.ContactEmail,
		Website:      in.Website,
		IsActive:     true,
	}
	return uc.venues.Create(ctx, v)
}

func (uc *UseCase) GetVenue(ctx context.Context, id valueobject.UUID) (entity.Venue, error) {
	return uc.venues.GetByID(ctx, id)
}

func (uc *UseCase) ListVenues(ctx context.Context, limit, offset int) ([]entity.Venue, error) {
	return uc.venues.List(ctx, limit, offset)
}

type UpdateVenueInput struct {
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
}

func (uc *UseCase) UpdateVenue(ctx context.Context, in UpdateVenueInput) error {
	if in.ID == valueobject.Nil {
		return apperror.New(apperror.CodeValidation, "id is required", nil)
	}
	if in.Name == "" || in.Address == "" || in.City == "" {
		return apperror.New(apperror.CodeValidation, "name/address/city are required", nil)
	}
	if in.Capacity < 0 {
		return apperror.New(apperror.CodeValidation, "capacity must be >= 0", nil)
	}
	v := entity.Venue{
		ID:           in.ID,
		Name:         in.Name,
		Address:      in.Address,
		City:         in.City,
		Country:      in.Country,
		Capacity:     in.Capacity,
		ContactPhone: in.ContactPhone,
		ContactEmail: in.ContactEmail,
		Website:      in.Website,
		IsActive:     in.IsActive,
	}
	return uc.venues.Update(ctx, v)
}

func (uc *UseCase) DeleteVenue(ctx context.Context, id valueobject.UUID) error {
	return uc.venues.Delete(ctx, id)
}

type CreateRoomInput struct {
	VenueID     valueobject.UUID
	Name        string
	Capacity    int
	Floor       *int
	Equipment   map[string]any
	HourlyRate  string
	IsAvailable bool
}

func (uc *UseCase) CreateRoom(ctx context.Context, in CreateRoomInput) (valueobject.UUID, error) {
	if in.VenueID == valueobject.Nil {
		return valueobject.Nil, apperror.New(apperror.CodeValidation, "venue_id is required", nil)
	}
	if in.Name == "" {
		return valueobject.Nil, apperror.New(apperror.CodeValidation, "name is required", nil)
	}
	if in.Capacity < 0 {
		return valueobject.Nil, apperror.New(apperror.CodeValidation, "capacity must be >= 0", nil)
	}
	rm := entity.Room{
		VenueID:     in.VenueID,
		Name:        in.Name,
		Capacity:    in.Capacity,
		Floor:       in.Floor,
		Equipment:   in.Equipment,
		HourlyRate:  in.HourlyRate,
		IsAvailable: in.IsAvailable,
	}
	return uc.rooms.Create(ctx, rm)
}

func (uc *UseCase) ListRooms(ctx context.Context, venueID valueobject.UUID) ([]entity.Room, error) {
	if venueID == valueobject.Nil {
		return nil, apperror.New(apperror.CodeValidation, "venue_id is required", nil)
	}
	return uc.rooms.ListByVenueID(ctx, venueID)
}



package event

import (
	"context"

	"time2meet/internal/domain/entity"
	"time2meet/internal/domain/repository"
	"time2meet/internal/domain/valueobject"
	"time2meet/pkg/apperror"
)

type UseCase struct {
	events repository.EventRepository
}

func New(events repository.EventRepository) *UseCase { return &UseCase{events: events} }

type CreateEventInput struct {
	OrganizerID     valueobject.UUID
	Title           string
	Description     string
	Status          string
	IsPublic        bool
	MaxParticipants *int
	CoverImage      string
}

func (uc *UseCase) Create(ctx context.Context, in CreateEventInput) (valueobject.UUID, error) {
	if in.OrganizerID == valueobject.Nil {
		return valueobject.Nil, apperror.New(apperror.CodeValidation, "organizer_id is required", nil)
	}
	if in.Title == "" {
		return valueobject.Nil, apperror.New(apperror.CodeValidation, "title is required", nil)
	}
	st := valueobject.EventStatus(in.Status)
	if err := st.Validate(); err != nil {
		return valueobject.Nil, apperror.New(apperror.CodeValidation, "invalid status", err)
	}
	e := entity.Event{
		OrganizerID:     in.OrganizerID,
		Title:           in.Title,
		Description:     in.Description,
		Status:          st,
		IsPublic:        in.IsPublic,
		MaxParticipants: in.MaxParticipants,
		CoverImage:      in.CoverImage,
	}
	return uc.events.Create(ctx, e)
}

func (uc *UseCase) Get(ctx context.Context, id valueobject.UUID) (entity.Event, error) {
	return uc.events.GetByID(ctx, id)
}

func (uc *UseCase) List(ctx context.Context, organizerID *valueobject.UUID, status *string, limit, offset int) ([]entity.Event, error) {
	return uc.events.List(ctx, organizerID, status, limit, offset)
}

type UpdateEventInput struct {
	ID              valueobject.UUID
	Title           string
	Description     string
	Status          string
	IsPublic        bool
	MaxParticipants *int
	CoverImage      string
}

func (uc *UseCase) Update(ctx context.Context, in UpdateEventInput) error {
	if in.ID == valueobject.Nil {
		return apperror.New(apperror.CodeValidation, "id is required", nil)
	}
	if in.Title == "" {
		return apperror.New(apperror.CodeValidation, "title is required", nil)
	}
	st := valueobject.EventStatus(in.Status)
	if err := st.Validate(); err != nil {
		return apperror.New(apperror.CodeValidation, "invalid status", err)
	}
	e := entity.Event{
		ID:              in.ID,
		Title:           in.Title,
		Description:     in.Description,
		Status:          st,
		IsPublic:        in.IsPublic,
		MaxParticipants: in.MaxParticipants,
		CoverImage:      in.CoverImage,
	}
	return uc.events.Update(ctx, e)
}

func (uc *UseCase) Delete(ctx context.Context, id valueobject.UUID) error {
	return uc.events.Delete(ctx, id)
}

func (uc *UseCase) Cancel(ctx context.Context, id valueobject.UUID) error {
	e, err := uc.events.GetByID(ctx, id)
	if err != nil {
		return err
	}
	e.Status = valueobject.EventStatusCancelled
	return uc.events.Update(ctx, e)
}



package ticket

import (
	"context"

	"time2meet/internal/domain/entity"
	"time2meet/internal/domain/repository"
	"time2meet/internal/domain/valueobject"
	"time2meet/pkg/apperror"
)

type TicketUseCase struct {
	tickets repository.TicketRepository
}

func NewTicketUC(tickets repository.TicketRepository) *TicketUseCase {
	return &TicketUseCase{tickets: tickets}
}

func (uc *TicketUseCase) Get(ctx context.Context, id valueobject.UUID) (entity.Ticket, error) {
	return uc.tickets.GetByID(ctx, id)
}

func (uc *TicketUseCase) ListByBuyer(ctx context.Context, buyerID valueobject.UUID, limit, offset int) ([]entity.Ticket, error) {
	if buyerID == valueobject.Nil {
		return nil, apperror.New(apperror.CodeValidation, "buyer_id is required", nil)
	}
	return uc.tickets.ListByBuyerID(ctx, buyerID, limit, offset)
}

func (uc *TicketUseCase) UpdateStatus(ctx context.Context, id valueobject.UUID, status string) error {
	if id == valueobject.Nil {
		return apperror.New(apperror.CodeValidation, "id is required", nil)
	}
	return uc.tickets.UpdateStatus(ctx, id, status)
}

func (uc *TicketUseCase) Delete(ctx context.Context, id valueobject.UUID) error {
	return uc.tickets.Delete(ctx, id)
}



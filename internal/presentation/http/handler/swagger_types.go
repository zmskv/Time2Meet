package handler

import (
	"time2meet/internal/domain/entity"
	"time2meet/internal/domain/repository"
	"time2meet/pkg/apperror"
)

type IDResponse struct {
	ID string `json:"id"`
}

type TicketIDResponse struct {
	TicketID string `json:"ticket_id"`
}

type ErrorResponse struct {
	Code    apperror.Code `json:"code"`
	Message string        `json:"message"`
}

type UserSwagger = entity.User
type EventSwagger = entity.Event
type VenueSwagger = entity.Venue
type RoomSwagger = entity.Room
type TicketSwagger = entity.Ticket

type SalesReportRowSwagger = repository.SalesReportRow
type AttendanceRowSwagger = repository.AttendanceRow
type PopularEventRowSwagger = repository.PopularEventRow

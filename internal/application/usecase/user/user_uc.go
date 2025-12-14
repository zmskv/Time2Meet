package user

import (
	"context"

	"time2meet/internal/domain/entity"
	"time2meet/internal/domain/repository"
	"time2meet/internal/domain/valueobject"
	"time2meet/pkg/apperror"
)

type UseCase struct {
	users    repository.UserRepository
	profiles repository.UserProfileRepository
}

func New(users repository.UserRepository, profiles repository.UserProfileRepository) *UseCase {
	return &UseCase{users: users, profiles: profiles}
}

type CreateUserInput struct {
	Email        string
	PasswordHash string
	FullName     string
	Phone        string
	Role         string
}

func (uc *UseCase) Create(ctx context.Context, in CreateUserInput) (valueobject.UUID, error) {
	email, err := valueobject.ParseEmail(in.Email)
	if err != nil {
		return valueobject.Nil, apperror.New(apperror.CodeValidation, "invalid email", err)
	}
	if in.PasswordHash == "" {
		return valueobject.Nil, apperror.New(apperror.CodeValidation, "password_hash is required", nil)
	}
	if in.FullName == "" {
		return valueobject.Nil, apperror.New(apperror.CodeValidation, "full_name is required", nil)
	}
	role := entity.UserRole(in.Role)
	switch role {
	case entity.UserRoleAdmin, entity.UserRoleOrganizer, entity.UserRoleAttendee:
	default:
		return valueobject.Nil, apperror.New(apperror.CodeValidation, "invalid role", nil)
	}
	u := entity.User{
		Email:        email,
		PasswordHash: in.PasswordHash,
		FullName:     in.FullName,
		Phone:        in.Phone,
		Role:         role,
		IsActive:     true,
	}
	return uc.users.Create(ctx, u)
}

func (uc *UseCase) Get(ctx context.Context, id valueobject.UUID) (entity.User, error) {
	return uc.users.GetByID(ctx, id)
}

func (uc *UseCase) List(ctx context.Context, limit, offset int) ([]entity.User, error) {
	return uc.users.List(ctx, limit, offset)
}

type UpdateUserInput struct {
	ID           valueobject.UUID
	Email        string
	PasswordHash string
	FullName     string
	Phone        string
	Role         string
	IsActive     bool
}

func (uc *UseCase) Update(ctx context.Context, in UpdateUserInput) error {
	email, err := valueobject.ParseEmail(in.Email)
	if err != nil {
		return apperror.New(apperror.CodeValidation, "invalid email", err)
	}
	role := entity.UserRole(in.Role)
	switch role {
	case entity.UserRoleAdmin, entity.UserRoleOrganizer, entity.UserRoleAttendee:
	default:
		return apperror.New(apperror.CodeValidation, "invalid role", nil)
	}
	u := entity.User{
		ID:           in.ID,
		Email:        email,
		PasswordHash: in.PasswordHash,
		FullName:     in.FullName,
		Phone:        in.Phone,
		Role:         role,
		IsActive:     in.IsActive,
	}
	return uc.users.Update(ctx, u)
}

func (uc *UseCase) Delete(ctx context.Context, id valueobject.UUID) error {
	return uc.users.Delete(ctx, id)
}



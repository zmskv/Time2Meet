package repository

import (
	"context"

	"time2meet/internal/domain/entity"
	"time2meet/internal/domain/valueobject"
)

type UserRepository interface {
	Create(ctx context.Context, u entity.User) (valueobject.UUID, error)
	GetByID(ctx context.Context, id valueobject.UUID) (entity.User, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	List(ctx context.Context, limit, offset int) ([]entity.User, error)
	Update(ctx context.Context, u entity.User) error
	Delete(ctx context.Context, id valueobject.UUID) error
}

type UserProfileRepository interface {
	Upsert(ctx context.Context, p entity.UserProfile) error
	GetByUserID(ctx context.Context, userID valueobject.UUID) (entity.UserProfile, error)
}

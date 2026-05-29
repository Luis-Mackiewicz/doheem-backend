package domain

import (
	"context"
	"time"
)

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	AvatarURL    *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateUserParams struct {
	Name         string
	Email        string
	PasswordHash string
	AvatarURL    *string
}

type UpdateUserParams struct {
	Name      *string
	Email     *string
	AvatarURL *string
}

type UserRepository interface {
	GetByID(ctx context.Context, id string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	List(ctx context.Context) ([]User, error)
	Create(ctx context.Context, params CreateUserParams) (User, error)
	Update(ctx context.Context, id string, params UpdateUserParams) (User, error)
	UpdatePassword(ctx context.Context, id string, passwordHash string) error
	Delete(ctx context.Context, id string) error
}

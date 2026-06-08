package user

import (
	"context"
	"errors"
	"time"
)

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	AvatarURL    *string
	Phone        *string
	Document     *string
	BirthDate    *time.Time
	Cep          *string
	IsAdmin      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateUserParams struct {
	Name         string
	Email        string
	PasswordHash string
	AvatarURL    *string
	Phone        *string
	Document     *string
	BirthDate    *time.Time
	Cep          *string
}

type UpdateUserParams struct {
	Name      *string
	Email     *string
	AvatarURL *string
	Phone     *string
	Document  *string
	BirthDate *time.Time
	Cep       *string
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

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrEmailAlreadyExists     = errors.New("email already exists")
	ErrRefreshTokenNotFound   = errors.New("refresh token not found")
	ErrRefreshTokenExpired    = errors.New("refresh token has expired")
	ErrRefreshTokenRevoked    = errors.New("refresh token has been revoked")
)

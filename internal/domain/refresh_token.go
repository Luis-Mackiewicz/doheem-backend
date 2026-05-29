package domain

import (
	"context"
	"time"
)

type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

type CreateRefreshTokenParams struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, params CreateRefreshTokenParams) error
	FindByHash(ctx context.Context, hash string) (RefreshToken, error)
	Revoke(ctx context.Context, hash string) error
	RevokeAllByUser(ctx context.Context, userID string) error
}

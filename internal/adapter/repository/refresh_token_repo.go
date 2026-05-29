package repository

import (
	"context"
	"crypto/sha256"
	"fmt"

	"doheem-backend/internal/adapter/repository/db"
	"doheem-backend/internal/domain"

	"github.com/jackc/pgx/v5"
)

type RefreshTokenRepo struct {
	q *db.Queries
}

func NewRefreshTokenRepo(q *db.Queries) *RefreshTokenRepo {
	return &RefreshTokenRepo{q: q}
}

func (r *RefreshTokenRepo) Create(ctx context.Context, params domain.CreateRefreshTokenParams) error {
	return r.q.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		ID:        uuidFromString(params.ID),
		UserID:    uuidFromString(params.UserID),
		TokenHash: params.TokenHash,
		ExpiresAt: timestamptzFromTime(params.ExpiresAt),
	})
}

func (r *RefreshTokenRepo) FindByHash(ctx context.Context, hash string) (domain.RefreshToken, error) {
	rt, err := r.q.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.RefreshToken{}, domain.ErrRefreshTokenNotFound
		}
		return domain.RefreshToken{}, fmt.Errorf("find refresh token: %w", err)
	}
	return domain.RefreshToken{
		ID:        uuidToString(rt.ID),
		UserID:    uuidToString(rt.UserID),
		TokenHash: rt.TokenHash,
		ExpiresAt: rt.ExpiresAt.Time,
		RevokedAt: timestamptzToTimePtr(rt.RevokedAt),
		CreatedAt: rt.CreatedAt.Time,
	}, nil
}

func (r *RefreshTokenRepo) Revoke(ctx context.Context, hash string) error {
	return r.q.RevokeRefreshToken(ctx, hash)
}

func (r *RefreshTokenRepo) RevokeAllByUser(ctx context.Context, userID string) error {
	return r.q.RevokeAllUserRefreshTokens(ctx, uuidFromString(userID))
}

func hashRefreshToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}

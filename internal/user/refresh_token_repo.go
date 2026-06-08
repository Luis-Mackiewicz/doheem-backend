package user

import (
	"context"
	"crypto/sha256"
	"fmt"

	"doheem-backend/internal/db"

	"github.com/jackc/pgx/v5"
)

type RefreshTokenRepo struct {
	q *db.Queries
}

func NewRefreshTokenRepo(q *db.Queries) *RefreshTokenRepo {
	return &RefreshTokenRepo{q: q}
}

func (r *RefreshTokenRepo) Create(ctx context.Context, params CreateRefreshTokenParams) error {
	return r.q.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		ID:        db.UUIDFromString(params.ID),
		UserID:    db.UUIDFromString(params.UserID),
		TokenHash: params.TokenHash,
		ExpiresAt: db.TimestamptzFromTime(params.ExpiresAt),
	})
}

func (r *RefreshTokenRepo) FindByHash(ctx context.Context, hash string) (RefreshToken, error) {
	rt, err := r.q.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return RefreshToken{}, ErrRefreshTokenNotFound
		}
		return RefreshToken{}, fmt.Errorf("find refresh token: %w", err)
	}
	return RefreshToken{
		ID:        db.UUIDToString(rt.ID),
		UserID:    db.UUIDToString(rt.UserID),
		TokenHash: rt.TokenHash,
		ExpiresAt: rt.ExpiresAt.Time,
		RevokedAt: db.TimestamptzToTimePtr(rt.RevokedAt),
		CreatedAt: rt.CreatedAt.Time,
	}, nil
}

func (r *RefreshTokenRepo) Revoke(ctx context.Context, hash string) error {
	return r.q.RevokeRefreshToken(ctx, hash)
}

func (r *RefreshTokenRepo) RevokeAllByUser(ctx context.Context, userID string) error {
	return r.q.RevokeAllUserRefreshTokens(ctx, db.UUIDFromString(userID))
}

func hashRefreshToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}

package user

import (
	"context"

	"doheem-backend/internal/db"
)

type UserRepo struct {
	q *db.Queries
}

func NewUserRepo(q *db.Queries) *UserRepo {
	return &UserRepo{q: q}
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (User, error) {
	u, err := r.q.GetUserByID(ctx, db.UUIDFromString(id))
	if err != nil {
		return User{}, err
	}
	return toUser(u), nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (User, error) {
	u, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return User{}, err
	}
	return toUser(u), nil
}

func (r *UserRepo) List(ctx context.Context) ([]User, error) {
	users, err := r.q.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	return toUsers(users), nil
}

func (r *UserRepo) Create(ctx context.Context, params CreateUserParams) (User, error) {
	u, err := r.q.CreateUser(ctx, db.CreateUserParams{
		Name:         params.Name,
		Email:        params.Email,
		PasswordHash: params.PasswordHash,
		AvatarUrl:    db.TextFromStringPtr(params.AvatarURL),
	})
	if err != nil {
		return User{}, err
	}
	return toUser(u), nil
}

func (r *UserRepo) Update(ctx context.Context, id string, params UpdateUserParams) (User, error) {
	var name string
	if params.Name != nil {
		name = *params.Name
	}
	var email string
	if params.Email != nil {
		email = *params.Email
	}
	u, err := r.q.UpdateUser(ctx, db.UpdateUserParams{
		ID:        db.UUIDFromString(id),
		Name:      name,
		Email:     email,
		AvatarUrl: db.TextFromStringPtr(params.AvatarURL),
	})
	if err != nil {
		return User{}, err
	}
	return toUser(u), nil
}

func (r *UserRepo) UpdatePassword(ctx context.Context, id string, passwordHash string) error {
	return r.q.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:           db.UUIDFromString(id),
		PasswordHash: passwordHash,
	})
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	return r.q.DeleteUser(ctx, db.UUIDFromString(id))
}

func toUser(u db.User) User {
	return User{
		ID:           db.UUIDToString(u.ID),
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		AvatarURL:    db.TextToStringPtr(u.AvatarUrl),
		CreatedAt:    u.CreatedAt.Time,
		UpdatedAt:    u.UpdatedAt.Time,
	}
}

func toUsers(users []db.User) []User {
	result := make([]User, len(users))
	for i, u := range users {
		result[i] = toUser(u)
	}
	return result
}

func toRefreshToken(rt db.RefreshToken) RefreshToken {
	return RefreshToken{
		ID:        db.UUIDToString(rt.ID),
		UserID:    db.UUIDToString(rt.UserID),
		TokenHash: rt.TokenHash,
		ExpiresAt: rt.ExpiresAt.Time,
		RevokedAt: db.TimestamptzToTimePtr(rt.RevokedAt),
		CreatedAt: rt.CreatedAt.Time,
	}
}

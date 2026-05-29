package repository

import (
	"context"

	"doheem-backend/internal/adapter/repository/db"
	"doheem-backend/internal/domain"
)

type UserRepo struct {
	q *db.Queries
}

func NewUserRepo(q *db.Queries) *UserRepo {
	return &UserRepo{q: q}
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (domain.User, error) {
	u, err := r.q.GetUserByID(ctx, uuidFromString(id))
	if err != nil {
		return domain.User{}, err
	}
	return domainUser(u), nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domainUser(u), nil
}

func (r *UserRepo) List(ctx context.Context) ([]domain.User, error) {
	users, err := r.q.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	return domainUsers(users), nil
}

func (r *UserRepo) Create(ctx context.Context, params domain.CreateUserParams) (domain.User, error) {
	u, err := r.q.CreateUser(ctx, db.CreateUserParams{
		Name:         params.Name,
		Email:        params.Email,
		PasswordHash: params.PasswordHash,
		AvatarUrl:    textFromStringPtr(params.AvatarURL),
	})
	if err != nil {
		return domain.User{}, err
	}
	return domainUser(u), nil
}

func (r *UserRepo) Update(ctx context.Context, id string, params domain.UpdateUserParams) (domain.User, error) {
	var name string
	if params.Name != nil {
		name = *params.Name
	}
	var email string
	if params.Email != nil {
		email = *params.Email
	}
	u, err := r.q.UpdateUser(ctx, db.UpdateUserParams{
		ID:        uuidFromString(id),
		Name:      name,
		Email:     email,
		AvatarUrl: textFromStringPtr(params.AvatarURL),
	})
	if err != nil {
		return domain.User{}, err
	}
	return domainUser(u), nil
}

func (r *UserRepo) UpdatePassword(ctx context.Context, id string, passwordHash string) error {
	return r.q.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:           uuidFromString(id),
		PasswordHash: passwordHash,
	})
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	return r.q.DeleteUser(ctx, uuidFromString(id))
}

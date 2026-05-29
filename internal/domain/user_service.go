package domain

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, params CreateUserParams) (User, error) {
	existing, err := s.repo.GetByEmail(ctx, params.Email)
	if err == nil && existing.ID != "" {
		return User{}, ErrEmailAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(params.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("failed to hash password: %w", err)
	}
	params.PasswordHash = string(hash)

	return s.repo.Create(ctx, params)
}

func (s *UserService) Login(ctx context.Context, email, password string) (User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return User{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return User{}, ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) List(ctx context.Context) ([]User, error) {
	return s.repo.List(ctx)
}

func (s *UserService) Update(ctx context.Context, id string, params UpdateUserParams) (User, error) {
	if params.Email != nil {
		existing, err := s.repo.GetByEmail(ctx, *params.Email)
		if err == nil && existing.ID != "" && existing.ID != id {
			return User{}, ErrEmailAlreadyExists
		}
	}
	return s.repo.Update(ctx, id, params)
}

func (s *UserService) UpdatePassword(ctx context.Context, id, oldPassword, newPassword string) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return s.repo.UpdatePassword(ctx, id, string(hash))
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

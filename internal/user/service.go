package user

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

type UserService struct {
	repo        UserRepository
	refreshRepo RefreshTokenRepository
}

func NewUserService(repo UserRepository, refreshRepo RefreshTokenRepository) *UserService {
	return &UserService{repo: repo, refreshRepo: refreshRepo}
}

func (s *UserService) Register(ctx context.Context, params CreateUserParams) (User, error) {
	if params.Email != "" {
		existing, err := s.repo.GetByEmail(ctx, params.Email)
		if err == nil && existing.ID != "" {
			return User{}, ErrEmailAlreadyExists
		}
	}
	if params.Document != nil && *params.Document != "" {
		existing, err := s.repo.GetByDocument(ctx, *params.Document)
		if err == nil && existing.ID != "" {
			return User{}, ErrDocumentAlreadyExists
		}
	}
	if params.Phone != nil && *params.Phone != "" {
		existing, err := s.repo.GetByPhone(ctx, *params.Phone)
		if err == nil && existing.ID != "" {
			return User{}, ErrPhoneAlreadyExists
		}
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
	if params.AvatarURL != nil {
		maxSize := 5 * 1024 * 1024
		if len(*params.AvatarURL) > maxSize {
			return User{}, fmt.Errorf("arquivo muito grande: máximo %d bytes", maxSize)
		}
	}

	if params.Email != nil {
		existing, err := s.repo.GetByEmail(ctx, *params.Email)
		if err == nil && existing.ID != "" && existing.ID != id {
			return User{}, ErrEmailAlreadyExists
		}
	}
	if params.Document != nil && *params.Document != "" {
		existing, err := s.repo.GetByDocument(ctx, *params.Document)
		if err == nil && existing.ID != "" && existing.ID != id {
			return User{}, ErrDocumentAlreadyExists
		}
	}
	if params.Phone != nil && *params.Phone != "" {
		existing, err := s.repo.GetByPhone(ctx, *params.Phone)
		if err == nil && existing.ID != "" && existing.ID != id {
			return User{}, ErrPhoneAlreadyExists
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

func (s *UserService) StoreRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("falha ao gerar uuid: %w", err)
	}
	return s.refreshRepo.Create(ctx, CreateRefreshTokenParams{
		ID:        id.String(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	})
}

func (s *UserService) RefreshToken(ctx context.Context, refreshTokenStr string) (RefreshToken, error) {
	rt, err := s.refreshRepo.FindByHash(ctx, refreshTokenStr)
	if err != nil {
		return RefreshToken{}, err
	}
	if rt.ExpiresAt.Before(time.Now()) {
		return RefreshToken{}, ErrRefreshTokenExpired
	}
	return rt, nil
}

func (s *UserService) RevokeRefreshToken(ctx context.Context, hash string) error {
	return s.refreshRepo.Revoke(ctx, hash)
}

func (s *UserService) RevokeAllUserRefreshTokens(ctx context.Context, userID string) error {
	return s.refreshRepo.RevokeAllByUser(ctx, userID)
}

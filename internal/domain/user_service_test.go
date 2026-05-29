package domain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(User), args.Error(1)
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(User), args.Error(1)
}

func (m *mockUserRepo) List(ctx context.Context) ([]User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]User), args.Error(1)
}

func (m *mockUserRepo) Create(ctx context.Context, params CreateUserParams) (User, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(User), args.Error(1)
}

func (m *mockUserRepo) Update(ctx context.Context, id string, params UpdateUserParams) (User, error) {
	args := m.Called(ctx, id, params)
	return args.Get(0).(User), args.Error(1)
}

func (m *mockUserRepo) UpdatePassword(ctx context.Context, id string, passwordHash string) error {
	args := m.Called(ctx, id, passwordHash)
	return args.Error(0)
}

func (m *mockUserRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUserService_Register_Success(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetByEmail", ctx, "test@example.com").Return(User{}, nil)
	mockRepo.On("Create", ctx, mock.MatchedBy(func(p CreateUserParams) bool {
		return p.Email == "test@example.com" && p.Name == "Test" && p.PasswordHash != ""
	})).Return(User{ID: "1", Name: "Test", Email: "test@example.com"}, nil)

	user, err := svc.Register(ctx, CreateUserParams{
		Name:         "Test",
		Email:        "test@example.com",
		PasswordHash: "plain-password",
	})

	assert.NoError(t, err)
	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Register_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()

	existing := User{ID: "1", Email: "test@example.com"}
	mockRepo.On("GetByEmail", ctx, "test@example.com").Return(existing, nil)

	user, err := svc.Register(ctx, CreateUserParams{
		Name:         "Test",
		Email:        "test@example.com",
		PasswordHash: "plain-password",
	})

	assert.ErrorIs(t, err, ErrEmailAlreadyExists)
	assert.Empty(t, user.ID)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_Success(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	stored := User{ID: "1", Email: "test@example.com", PasswordHash: string(hash)}
	mockRepo.On("GetByEmail", ctx, "test@example.com").Return(stored, nil)

	user, err := svc.Login(ctx, "test@example.com", "correct-password")

	assert.NoError(t, err)
	assert.Equal(t, "1", user.ID)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_InvalidCredentials(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	stored := User{ID: "1", Email: "test@example.com", PasswordHash: string(hash)}
	mockRepo.On("GetByEmail", ctx, "test@example.com").Return(stored, nil)

	user, err := svc.Login(ctx, "test@example.com", "wrong-password")

	assert.ErrorIs(t, err, ErrInvalidCredentials)
	assert.Empty(t, user.ID)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetByEmail", ctx, "nonexistent@example.com").Return(User{}, assert.AnError)

	user, err := svc.Login(ctx, "nonexistent@example.com", "password")

	assert.ErrorIs(t, err, ErrInvalidCredentials)
	assert.Empty(t, user.ID)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetByID_Success(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetByID", ctx, "1").Return(User{ID: "1", Name: "Test"}, nil)

	user, err := svc.GetByID(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, "1", user.ID)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetByID", ctx, "999").Return(User{}, assert.AnError)

	user, err := svc.GetByID(ctx, "999")

	assert.ErrorIs(t, err, ErrUserNotFound)
	assert.Empty(t, user.ID)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Update_Success(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()

	newEmail := "new@example.com"
	params := UpdateUserParams{Email: &newEmail}
	mockRepo.On("GetByEmail", ctx, newEmail).Return(User{}, nil)
	mockRepo.On("Update", ctx, "1", params).Return(User{ID: "1", Email: newEmail}, nil)

	user, err := svc.Update(ctx, "1", params)

	assert.NoError(t, err)
	assert.Equal(t, newEmail, user.Email)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Update_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()

	newEmail := "other@example.com"
	mockRepo.On("GetByEmail", ctx, newEmail).Return(User{ID: "2", Email: newEmail}, nil)

	user, err := svc.Update(ctx, "1", UpdateUserParams{Email: &newEmail})

	assert.ErrorIs(t, err, ErrEmailAlreadyExists)
	assert.Empty(t, user.ID)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdatePassword_Success(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("old-password"), bcrypt.DefaultCost)
	mockRepo.On("GetByID", ctx, "1").Return(User{ID: "1", PasswordHash: string(hash)}, nil)
	mockRepo.On("UpdatePassword", ctx, "1", mock.MatchedBy(func(h string) bool {
		return bcrypt.CompareHashAndPassword([]byte(h), []byte("new-password")) == nil
	})).Return(nil)

	err := svc.UpdatePassword(ctx, "1", "old-password", "new-password")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdatePassword_InvalidOldPassword(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("actual-password"), bcrypt.DefaultCost)
	mockRepo.On("GetByID", ctx, "1").Return(User{ID: "1", PasswordHash: string(hash)}, nil)

	err := svc.UpdatePassword(ctx, "1", "wrong-password", "new-password")

	assert.ErrorIs(t, err, ErrInvalidCredentials)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Delete(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()

	mockRepo.On("Delete", ctx, "1").Return(nil)

	err := svc.Delete(ctx, "1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

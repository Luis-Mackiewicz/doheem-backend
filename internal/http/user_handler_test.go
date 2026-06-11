package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"doheem-backend/internal/user"

	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (user.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(user.User), args.Error(1)
}
func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (user.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(user.User), args.Error(1)
}
func (m *mockUserRepo) GetByDocument(ctx context.Context, document string) (user.User, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(user.User), args.Error(1)
}
func (m *mockUserRepo) GetByPhone(ctx context.Context, phone string) (user.User, error) {
	args := m.Called(ctx, phone)
	return args.Get(0).(user.User), args.Error(1)
}
func (m *mockUserRepo) List(ctx context.Context) ([]user.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]user.User), args.Error(1)
}
func (m *mockUserRepo) Create(ctx context.Context, params user.CreateUserParams) (user.User, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(user.User), args.Error(1)
}
func (m *mockUserRepo) Update(ctx context.Context, id string, params user.UpdateUserParams) (user.User, error) {
	args := m.Called(ctx, id, params)
	return args.Get(0).(user.User), args.Error(1)
}
func (m *mockUserRepo) UpdatePassword(ctx context.Context, id string, passwordHash string) error {
	args := m.Called(ctx, id, passwordHash)
	return args.Error(0)
}
func (m *mockUserRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockRefreshTokenRepo struct{ mock.Mock }

func (m *mockRefreshTokenRepo) Create(ctx context.Context, params user.CreateRefreshTokenParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *mockRefreshTokenRepo) FindByHash(ctx context.Context, hash string) (user.RefreshToken, error) {
	args := m.Called(ctx, hash)
	return args.Get(0).(user.RefreshToken), args.Error(1)
}

func (m *mockRefreshTokenRepo) Revoke(ctx context.Context, hash string) error {
	args := m.Called(ctx, hash)
	return args.Error(0)
}

func (m *mockRefreshTokenRepo) RevokeAllByUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func newTestJWT(t *testing.T) *JWTService {
	t.Helper()
	return NewJWTService("test-secret", 24*time.Hour, 168*time.Hour)
}

func authCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, UserIDKey, "test-user-id")
}

func readJSON(t *testing.T, body []byte, v any) {
	t.Helper()
	if err := json.Unmarshal(body, v); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
}

func TestRegister_Success(t *testing.T) {
	repo := new(mockUserRepo)
	refreshRepo := new(mockRefreshTokenRepo)
	svc := user.NewUserService(repo, refreshRepo)
	jwt := newTestJWT(t)
	handler := NewUserHandler(svc, jwt)

	repo.On("GetByEmail", mock.Anything, "test@example.com").Return(user.User{}, nil)
	repo.On("GetByDocument", mock.Anything, "52998224725").Return(user.User{}, nil)
	repo.On("GetByPhone", mock.Anything, "11999999999").Return(user.User{}, nil)
	repo.On("Create", mock.Anything, mock.MatchedBy(func(p user.CreateUserParams) bool {
		return p.Email == "test@example.com" && p.Name == "Test" && p.PasswordHash != ""
	})).Return(user.User{ID: "u1", Name: "Test", Email: "test@example.com"}, nil)
	refreshRepo.On("Create", mock.Anything, mock.MatchedBy(func(p user.CreateRefreshTokenParams) bool {
		return p.UserID == "u1" && p.TokenHash != ""
	})).Return(nil)

	body := `{"name":"Test","email":"test@example.com","password":"123456","phone":"(11) 99999-9999","document":"529.982.247-25","birth_date":"1990-05-20","cep":"01001-000"}`
	r := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", resp.StatusCode, w.Body.String())
	}
	var data authResponse
	readJSON(t, w.Body.Bytes(), &data)
	if data.User.ID != "u1" || data.Token == "" {
		t.Fatal("unexpected response")
	}
	repo.AssertExpectations(t)
	refreshRepo.AssertExpectations(t)
}

func TestRegister_ValidationError(t *testing.T) {
	handler := NewUserHandler(user.NewUserService(new(mockUserRepo), new(mockRefreshTokenRepo)), newTestJWT(t))

	body := `{"name":"","email":"bad","password":"12"}`
	r := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
	var errResp map[string]any
	readJSON(t, w.Body.Bytes(), &errResp)
	if _, ok := errResp["fields"]; !ok {
		t.Fatal("expected validation fields in response")
	}
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	repo := new(mockUserRepo)
	svc := user.NewUserService(repo, new(mockRefreshTokenRepo))
	handler := NewUserHandler(svc, newTestJWT(t))

	repo.On("GetByEmail", mock.Anything, "exists@example.com").Return(user.User{ID: "existing"}, nil)
	repo.On("GetByDocument", mock.Anything, "52998224725").Return(user.User{}, nil).Maybe()
	repo.On("GetByPhone", mock.Anything, "11999999999").Return(user.User{}, nil).Maybe()

	body := `{"name":"Test","email":"exists@example.com","password":"123456","phone":"11999999999","document":"529.982.247-25","birth_date":"1990-05-20","cep":"01001-000"}`
	r := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409, got %d", resp.StatusCode)
	}
	repo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	repo := new(mockUserRepo)
	refreshRepo := new(mockRefreshTokenRepo)
	svc := user.NewUserService(repo, refreshRepo)
	handler := NewUserHandler(svc, newTestJWT(t))

	hash, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	repo.On("GetByEmail", mock.Anything, "test@example.com").Return(user.User{ID: "u1", PasswordHash: string(hash)}, nil)
	refreshRepo.On("Create", mock.Anything, mock.MatchedBy(func(p user.CreateRefreshTokenParams) bool {
		return p.UserID == "u1" && p.TokenHash != ""
	})).Return(nil)

	body := `{"email":"test@example.com","password":"correct-password"}`
	r := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, w.Body.String())
	}
	var data authResponse
	readJSON(t, w.Body.Bytes(), &data)
	if data.Token == "" {
		t.Fatal("expected token")
	}
	repo.AssertExpectations(t)
	refreshRepo.AssertExpectations(t)
}

func TestGetProfile_Success(t *testing.T) {
	repo := new(mockUserRepo)
	svc := user.NewUserService(repo, new(mockRefreshTokenRepo))
	handler := NewUserHandler(svc, newTestJWT(t))

	repo.On("GetByID", mock.Anything, "test-user-id").Return(user.User{
		ID: "test-user-id", Name: "Profile", Email: "profile@example.com",
	}, nil)

	r := httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
	r = r.WithContext(authCtx(r.Context()))
	w := httptest.NewRecorder()

	handler.GetProfile(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var data userResponse
	readJSON(t, w.Body.Bytes(), &data)
	if data.ID != "test-user-id" || data.Name != "Profile" {
		t.Fatal("unexpected profile data")
	}
	repo.AssertExpectations(t)
}

func TestGetProfile_NotFound(t *testing.T) {
	repo := new(mockUserRepo)
	svc := user.NewUserService(repo, new(mockRefreshTokenRepo))
	handler := NewUserHandler(svc, newTestJWT(t))

	repo.On("GetByID", mock.Anything, "test-user-id").Return(user.User{}, user.ErrUserNotFound)

	r := httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
	r = r.WithContext(authCtx(r.Context()))
	w := httptest.NewRecorder()

	handler.GetProfile(w, r)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Result().StatusCode)
	}
	repo.AssertExpectations(t)
}

func TestChangePassword_Success(t *testing.T) {
	repo := new(mockUserRepo)
	svc := user.NewUserService(repo, new(mockRefreshTokenRepo))
	handler := NewUserHandler(svc, newTestJWT(t))

	hash, _ := bcrypt.GenerateFromPassword([]byte("old-pass"), bcrypt.DefaultCost)
	repo.On("GetByID", mock.Anything, "test-user-id").Return(user.User{
		ID: "test-user-id", PasswordHash: string(hash),
	}, nil)
	repo.On("UpdatePassword", mock.Anything, "test-user-id", mock.Anything).Return(nil)

	body := `{"old_password":"old-pass","new_password":"new-pass"}`
	r := httptest.NewRequest(http.MethodPut, "/api/users/me/password", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r = r.WithContext(authCtx(r.Context()))
	w := httptest.NewRecorder()

	handler.ChangePassword(w, r)

	if w.Result().StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", w.Result().StatusCode, w.Body.String())
	}
	repo.AssertExpectations(t)
}


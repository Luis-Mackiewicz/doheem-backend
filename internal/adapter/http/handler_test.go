package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"doheem-backend/internal/domain"

	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}
func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(domain.User), args.Error(1)
}
func (m *mockUserRepo) List(ctx context.Context) ([]domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.User), args.Error(1)
}
func (m *mockUserRepo) Create(ctx context.Context, params domain.CreateUserParams) (domain.User, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(domain.User), args.Error(1)
}
func (m *mockUserRepo) Update(ctx context.Context, id string, params domain.UpdateUserParams) (domain.User, error) {
	args := m.Called(ctx, id, params)
	return args.Get(0).(domain.User), args.Error(1)
}
func (m *mockUserRepo) UpdatePassword(ctx context.Context, id string, passwordHash string) error {
	args := m.Called(ctx, id, passwordHash)
	return args.Error(0)
}
func (m *mockUserRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func newTestJWT(t *testing.T) *JWTService {
	t.Helper()
	return NewJWTService("test-secret", 24*time.Hour)
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
	svc := domain.NewUserService(repo)
	jwt := newTestJWT(t)
	handler := NewUserHandler(svc, jwt)

	repo.On("GetByEmail", mock.Anything, "test@example.com").Return(domain.User{}, nil)
	repo.On("Create", mock.Anything, mock.MatchedBy(func(p domain.CreateUserParams) bool {
		return p.Email == "test@example.com" && p.Name == "Test" && p.PasswordHash != ""
	})).Return(domain.User{ID: "u1", Name: "Test", Email: "test@example.com"}, nil)

	body := `{"name":"Test","email":"test@example.com","password":"123456"}`
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
}

func TestRegister_ValidationError(t *testing.T) {
	handler := NewUserHandler(domain.NewUserService(new(mockUserRepo)), newTestJWT(t))

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
	svc := domain.NewUserService(repo)
	handler := NewUserHandler(svc, newTestJWT(t))

	repo.On("GetByEmail", mock.Anything, "exists@example.com").Return(domain.User{ID: "existing"}, nil)

	body := `{"name":"Test","email":"exists@example.com","password":"123456"}`
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
	svc := domain.NewUserService(repo)
	handler := NewUserHandler(svc, newTestJWT(t))

	hash, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	repo.On("GetByEmail", mock.Anything, "test@example.com").Return(domain.User{ID: "u1", PasswordHash: string(hash)}, nil)

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
}

func TestGetProfile_Success(t *testing.T) {
	repo := new(mockUserRepo)
	svc := domain.NewUserService(repo)
	handler := NewUserHandler(svc, newTestJWT(t))

	repo.On("GetByID", mock.Anything, "test-user-id").Return(domain.User{
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
	svc := domain.NewUserService(repo)
	handler := NewUserHandler(svc, newTestJWT(t))

	repo.On("GetByID", mock.Anything, "test-user-id").Return(domain.User{}, domain.ErrUserNotFound)

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
	svc := domain.NewUserService(repo)
	handler := NewUserHandler(svc, newTestJWT(t))

	hash, _ := bcrypt.GenerateFromPassword([]byte("old-pass"), bcrypt.DefaultCost)
	repo.On("GetByID", mock.Anything, "test-user-id").Return(domain.User{
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


package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"doheem-backend/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGroupRepo struct{ mock.Mock }

func (m *mockGroupRepo) GetByID(ctx context.Context, id string) (domain.Group, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Group), args.Error(1)
}
func (m *mockGroupRepo) ListByUserID(ctx context.Context, userID string) ([]domain.Group, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Group), args.Error(1)
}
func (m *mockGroupRepo) Create(ctx context.Context, params domain.CreateGroupParams) (domain.Group, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(domain.Group), args.Error(1)
}
func (m *mockGroupRepo) Update(ctx context.Context, id string, params domain.UpdateGroupParams) (domain.Group, error) {
	args := m.Called(ctx, id, params)
	return args.Get(0).(domain.Group), args.Error(1)
}
func (m *mockGroupRepo) SoftDelete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockGroupRepo) Deactivate(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockGroupRepo) Activate(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockGroupMemberRepo struct{ mock.Mock }

func (m *mockGroupMemberRepo) GetByID(ctx context.Context, id string) (domain.GroupMember, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.GroupMember), args.Error(1)
}
func (m *mockGroupMemberRepo) Get(ctx context.Context, groupID, userID string) (domain.GroupMember, error) {
	args := m.Called(ctx, groupID, userID)
	return args.Get(0).(domain.GroupMember), args.Error(1)
}
func (m *mockGroupMemberRepo) ListByGroup(ctx context.Context, groupID string) ([]domain.GroupMemberWithUser, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).([]domain.GroupMemberWithUser), args.Error(1)
}
func (m *mockGroupMemberRepo) Create(ctx context.Context, groupID, userID, role string) (domain.GroupMember, error) {
	args := m.Called(ctx, groupID, userID, role)
	return args.Get(0).(domain.GroupMember), args.Error(1)
}
func (m *mockGroupMemberRepo) UpdateRole(ctx context.Context, groupID, userID, role string) (domain.GroupMember, error) {
	args := m.Called(ctx, groupID, userID, role)
	return args.Get(0).(domain.GroupMember), args.Error(1)
}
func (m *mockGroupMemberRepo) Remove(ctx context.Context, groupID, userID string) error {
	args := m.Called(ctx, groupID, userID)
	return args.Error(0)
}
func (m *mockGroupMemberRepo) CountActive(ctx context.Context, groupID string) (int64, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).(int64), args.Error(1)
}

func TestGroupCreate_Success(t *testing.T) {
	groupRepo := new(mockGroupRepo)
	memberRepo := new(mockGroupMemberRepo)
	svc := domain.NewGroupService(groupRepo, memberRepo)
	handler := NewGroupHandler(svc)

	groupRepo.On("Create", mock.Anything, mock.Anything).Return(domain.Group{ID: "g1", Name: "My Group", Currency: "BRL"}, nil)
	memberRepo.On("Create", mock.Anything, "g1", "test-user-id", "owner").Return(domain.GroupMember{}, nil)

	body := `{"name":"My Group","currency":"BRL"}`
	r := httptest.NewRequest(http.MethodPost, "/api/groups", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r = r.WithContext(authCtx(r.Context()))
	w := httptest.NewRecorder()

	handler.Create(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", resp.StatusCode, w.Body.String())
	}
	groupRepo.AssertExpectations(t)
	memberRepo.AssertExpectations(t)
}

func TestGroupCreate_ValidationError(t *testing.T) {
	handler := NewGroupHandler(domain.NewGroupService(new(mockGroupRepo), new(mockGroupMemberRepo)))

	body := `{"name":"","currency":"BR"}`
	r := httptest.NewRequest(http.MethodPost, "/api/groups", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r = r.WithContext(authCtx(r.Context()))
	w := httptest.NewRecorder()

	handler.Create(w, r)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Result().StatusCode)
	}
}

func TestGroupGetByID_Success(t *testing.T) {
	groupRepo := new(mockGroupRepo)
	svc := domain.NewGroupService(groupRepo, new(mockGroupMemberRepo))
	handler := NewGroupHandler(svc)

	groupRepo.On("GetByID", mock.Anything, "g1").Return(domain.Group{ID: "g1", Name: "My Group"}, nil)

	r := httptest.NewRequest(http.MethodGet, "/api/groups/g1", nil)
	r = r.WithContext(authCtx(r.Context()))
	r.SetPathValue("id", "g1")
	w := httptest.NewRecorder()

	handler.GetByID(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Result().StatusCode)
	}
	groupRepo.AssertExpectations(t)
}

func TestGroupGetByID_NotFound(t *testing.T) {
	groupRepo := new(mockGroupRepo)
	svc := domain.NewGroupService(groupRepo, new(mockGroupMemberRepo))
	handler := NewGroupHandler(svc)

	groupRepo.On("GetByID", mock.Anything, "999").Return(domain.Group{}, assert.AnError)

	r := httptest.NewRequest(http.MethodGet, "/api/groups/999", nil)
	r = r.WithContext(authCtx(r.Context()))
	r.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	handler.GetByID(w, r)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Result().StatusCode)
	}
	groupRepo.AssertExpectations(t)
}

func TestGroupList_Success(t *testing.T) {
	groupRepo := new(mockGroupRepo)
	svc := domain.NewGroupService(groupRepo, new(mockGroupMemberRepo))
	handler := NewGroupHandler(svc)

	groupRepo.On("ListByUserID", mock.Anything, "test-user-id").Return([]domain.Group{
		{ID: "g1", Name: "Group 1"},
		{ID: "g2", Name: "Group 2"},
	}, nil)

	r := httptest.NewRequest(http.MethodGet, "/api/groups", nil)
	r = r.WithContext(authCtx(r.Context()))
	w := httptest.NewRecorder()

	handler.List(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var groups []domain.Group
	readJSON(t, w.Body.Bytes(), &groups)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	groupRepo.AssertExpectations(t)
}

func TestGroupSoftDelete_Success(t *testing.T) {
	groupRepo := new(mockGroupRepo)
	svc := domain.NewGroupService(groupRepo, new(mockGroupMemberRepo))
	handler := NewGroupHandler(svc)

	groupRepo.On("SoftDelete", mock.Anything, "g1").Return(nil)

	r := httptest.NewRequest(http.MethodDelete, "/api/groups/g1", nil)
	r = r.WithContext(authCtx(r.Context()))
	r.SetPathValue("id", "g1")
	w := httptest.NewRecorder()

	handler.SoftDelete(w, r)

	if w.Result().StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Result().StatusCode)
	}
	groupRepo.AssertExpectations(t)
}

func TestGroupDeactivate_Success(t *testing.T) {
	groupRepo := new(mockGroupRepo)
	svc := domain.NewGroupService(groupRepo, new(mockGroupMemberRepo))
	handler := NewGroupHandler(svc)

	groupRepo.On("Deactivate", mock.Anything, "g1").Return(nil)

	r := httptest.NewRequest(http.MethodPatch, "/api/groups/g1/deactivate", nil)
	r = r.WithContext(authCtx(r.Context()))
	r.SetPathValue("id", "g1")
	w := httptest.NewRecorder()

	handler.Deactivate(w, r)

	if w.Result().StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Result().StatusCode)
	}
	groupRepo.AssertExpectations(t)
}

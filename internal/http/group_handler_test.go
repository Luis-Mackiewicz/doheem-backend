package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"doheem-backend/internal/group"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGroupRepo struct{ mock.Mock }

func (m *mockGroupRepo) GetByID(ctx context.Context, id string) (group.Group, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(group.Group), args.Error(1)
}
func (m *mockGroupRepo) ListByUserID(ctx context.Context, userID string) ([]group.Group, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]group.Group), args.Error(1)
}
func (m *mockGroupRepo) Create(ctx context.Context, params group.CreateGroupParams) (group.Group, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(group.Group), args.Error(1)
}
func (m *mockGroupRepo) Update(ctx context.Context, id string, params group.UpdateGroupParams) (group.Group, error) {
	args := m.Called(ctx, id, params)
	return args.Get(0).(group.Group), args.Error(1)
}
func (m *mockGroupRepo) RegenerateInviteToken(ctx context.Context, id, token string) error {
	args := m.Called(ctx, id, token)
	return args.Error(0)
}
func (m *mockGroupRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockGroupRepo) CountByUserID(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

type mockGroupMemberRepo struct{ mock.Mock }

func (m *mockGroupMemberRepo) GetByID(ctx context.Context, id string) (group.GroupMember, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(group.GroupMember), args.Error(1)
}
func (m *mockGroupMemberRepo) Get(ctx context.Context, groupID, userID string) (group.GroupMember, error) {
	args := m.Called(ctx, groupID, userID)
	return args.Get(0).(group.GroupMember), args.Error(1)
}
func (m *mockGroupMemberRepo) ListByGroup(ctx context.Context, groupID string) ([]group.GroupMemberWithUser, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).([]group.GroupMemberWithUser), args.Error(1)
}
func (m *mockGroupMemberRepo) Create(ctx context.Context, groupID, userID string, isAdmin bool) (group.GroupMember, error) {
	args := m.Called(ctx, groupID, userID, isAdmin)
	return args.Get(0).(group.GroupMember), args.Error(1)
}
func (m *mockGroupMemberRepo) UpdateRole(ctx context.Context, groupID, userID string, isAdmin bool) (group.GroupMember, error) {
	args := m.Called(ctx, groupID, userID, isAdmin)
	return args.Get(0).(group.GroupMember), args.Error(1)
}
func (m *mockGroupMemberRepo) Remove(ctx context.Context, groupID, userID string) error {
	args := m.Called(ctx, groupID, userID)
	return args.Error(0)
}
func (m *mockGroupMemberRepo) Count(ctx context.Context, groupID string) (int64, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockGroupMemberRepo) CountAdmins(ctx context.Context, groupID string) (int64, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).(int64), args.Error(1)
}

func TestGroupCreate_Success(t *testing.T) {
	groupRepo := new(mockGroupRepo)
	memberRepo := new(mockGroupMemberRepo)
	svc := group.NewGroupService(groupRepo, memberRepo, nil)
	handler := NewGroupHandler(svc, nil, nil, nil)

	groupRepo.On("Create", mock.Anything, mock.Anything).Return(group.Group{ID: "g1", Name: "My Group"}, nil)
	memberRepo.On("Create", mock.Anything, "g1", "test-user-id", true).Return(group.GroupMember{}, nil)

	body := `{"name":"My Group"}`
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
	handler := NewGroupHandler(group.NewGroupService(new(mockGroupRepo), new(mockGroupMemberRepo), nil), nil, nil, nil)

	body := `{"name":""}`
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
	svc := group.NewGroupService(groupRepo, new(mockGroupMemberRepo), nil)
	handler := NewGroupHandler(svc, nil, nil, nil)

	groupRepo.On("GetByID", mock.Anything, "g1").Return(group.Group{ID: "g1", Name: "My Group"}, nil)

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
	svc := group.NewGroupService(groupRepo, new(mockGroupMemberRepo), nil)
	handler := NewGroupHandler(svc, nil, nil, nil)

	groupRepo.On("GetByID", mock.Anything, "999").Return(group.Group{}, assert.AnError)

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
	svc := group.NewGroupService(groupRepo, new(mockGroupMemberRepo), nil)
	handler := NewGroupHandler(svc, nil, nil, nil)

	groupRepo.On("ListByUserID", mock.Anything, "test-user-id").Return([]group.Group{
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
	var result paginatedResponse
	readJSON(t, w.Body.Bytes(), &result)
	if result.Total != 2 {
		t.Fatalf("expected total 2, got %d", result.Total)
	}
	groupRepo.AssertExpectations(t)
}

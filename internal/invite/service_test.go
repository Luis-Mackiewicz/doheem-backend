package invite

import (
	"context"
	"testing"
	"time"

	"doheem-backend/internal/group"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockInviteRepo struct {
	mock.Mock
}

func (m *mockInviteRepo) GetByID(ctx context.Context, id string) (Invite, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Invite), args.Error(1)
}

func (m *mockInviteRepo) GetByCode(ctx context.Context, code string) (InviteWithGroup, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(InviteWithGroup), args.Error(1)
}

func (m *mockInviteRepo) ListByGroup(ctx context.Context, groupID string) ([]InviteWithCreator, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).([]InviteWithCreator), args.Error(1)
}

func (m *mockInviteRepo) ListPendingByUser(ctx context.Context, userID string) ([]InviteWithGroup, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]InviteWithGroup), args.Error(1)
}

func (m *mockInviteRepo) Create(ctx context.Context, groupID, code, createdBy string, expiresAt time.Time) (Invite, error) {
	args := m.Called(ctx, groupID, code, createdBy, expiresAt)
	return args.Get(0).(Invite), args.Error(1)
}

func (m *mockInviteRepo) Use(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockInviteRepo) Revoke(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
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
func (m *mockGroupMemberRepo) Create(ctx context.Context, groupID, userID, role string) (group.GroupMember, error) {
	args := m.Called(ctx, groupID, userID, role)
	return args.Get(0).(group.GroupMember), args.Error(1)
}
func (m *mockGroupMemberRepo) UpdateRole(ctx context.Context, groupID, userID, role string) (group.GroupMember, error) {
	args := m.Called(ctx, groupID, userID, role)
	return args.Get(0).(group.GroupMember), args.Error(1)
}
func (m *mockGroupMemberRepo) Remove(ctx context.Context, groupID, userID string) error {
	args := m.Called(ctx, groupID, userID)
	return args.Error(0)
}
func (m *mockGroupMemberRepo) IsMember(ctx context.Context, groupID, userID string) (bool, error) {
	args := m.Called(ctx, groupID, userID)
	return args.Get(0).(bool), args.Error(1)
}
func (m *mockGroupMemberRepo) CountActive(ctx context.Context, groupID string) (int64, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).(int64), args.Error(1)
}

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
func (m *mockGroupRepo) GetUserIDsByGroupID(ctx context.Context, groupID string) ([]string, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).([]string), args.Error(1)
}

func TestInviteService_Use_Success(t *testing.T) {
	mockInvite := new(mockInviteRepo)
	mockMember := new(mockGroupMemberRepo)
	mockGroup := new(mockGroupRepo)
	svc := NewInviteService(mockInvite, mockMember, mockGroup)
	ctx := context.Background()

	future := time.Now().Add(24 * time.Hour)
	mockInvite.On("GetByID", ctx, "inv1").Return(Invite{ID: "inv1", GroupID: "g1", ExpiresAt: future}, nil)
	mockInvite.On("Use", ctx, "inv1").Return(nil)
	mockMember.On("Create", ctx, "g1", "u1", "member").Return(group.GroupMember{}, nil)

	err := svc.Use(ctx, "inv1", "u1")

	assert.NoError(t, err)
	mockInvite.AssertExpectations(t)
	mockMember.AssertExpectations(t)
}

func TestInviteService_Use_NotFound(t *testing.T) {
	mockInvite := new(mockInviteRepo)
	svc := NewInviteService(mockInvite, new(mockGroupMemberRepo), new(mockGroupRepo))
	ctx := context.Background()

	mockInvite.On("GetByID", ctx, "999").Return(Invite{}, assert.AnError)

	err := svc.Use(ctx, "999", "u1")

	assert.ErrorIs(t, err, ErrInviteNotFound)
}

func TestInviteService_Use_AlreadyUsed(t *testing.T) {
	mockInvite := new(mockInviteRepo)
	svc := NewInviteService(mockInvite, new(mockGroupMemberRepo), new(mockGroupRepo))
	ctx := context.Background()

	now := time.Now()
	mockInvite.On("GetByID", ctx, "inv1").Return(Invite{ID: "inv1", UsedAt: &now}, nil)

	err := svc.Use(ctx, "inv1", "u1")

	assert.ErrorIs(t, err, ErrInviteAlreadyUsed)
	mockInvite.AssertNotCalled(t, "Use")
}

func TestInviteService_Use_Revoked(t *testing.T) {
	mockInvite := new(mockInviteRepo)
	svc := NewInviteService(mockInvite, new(mockGroupMemberRepo), new(mockGroupRepo))
	ctx := context.Background()

	now := time.Now()
	mockInvite.On("GetByID", ctx, "inv1").Return(Invite{ID: "inv1", RevokedAt: &now}, nil)

	err := svc.Use(ctx, "inv1", "u1")

	assert.ErrorIs(t, err, ErrInviteRevoked)
	mockInvite.AssertNotCalled(t, "Use")
}

func TestInviteService_Use_Expired(t *testing.T) {
	mockInvite := new(mockInviteRepo)
	svc := NewInviteService(mockInvite, new(mockGroupMemberRepo), new(mockGroupRepo))
	ctx := context.Background()

	past := time.Now().Add(-24 * time.Hour)
	mockInvite.On("GetByID", ctx, "inv1").Return(Invite{ID: "inv1", ExpiresAt: past}, nil)

	err := svc.Use(ctx, "inv1", "u1")

	assert.ErrorIs(t, err, ErrInviteExpired)
	mockInvite.AssertNotCalled(t, "Use")
}

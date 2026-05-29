package domain

import (
	"context"
	"testing"
	"time"

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

func TestInviteService_Use_Success(t *testing.T) {
	mockInvite := new(mockInviteRepo)
	mockMember := new(mockGroupMemberRepo)
	mockGroup := new(mockGroupRepo)
	svc := NewInviteService(mockInvite, mockMember, mockGroup)
	ctx := context.Background()

	future := time.Now().Add(24 * time.Hour)
	mockInvite.On("GetByID", ctx, "inv1").Return(Invite{ID: "inv1", GroupID: "g1", ExpiresAt: future}, nil)
	mockInvite.On("Use", ctx, "inv1").Return(nil)
	mockMember.On("Create", ctx, "g1", "u1", "member").Return(GroupMember{}, nil)

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

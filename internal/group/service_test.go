package group

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGroupRepo struct {
	mock.Mock
}

func (m *mockGroupRepo) GetByID(ctx context.Context, id string) (Group, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Group), args.Error(1)
}

func (m *mockGroupRepo) ListByUserID(ctx context.Context, userID string) ([]Group, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]Group), args.Error(1)
}

func (m *mockGroupRepo) Create(ctx context.Context, params CreateGroupParams) (Group, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(Group), args.Error(1)
}

func (m *mockGroupRepo) Update(ctx context.Context, id string, params UpdateGroupParams) (Group, error) {
	args := m.Called(ctx, id, params)
	return args.Get(0).(Group), args.Error(1)
}

func (m *mockGroupRepo) RegenerateInviteToken(ctx context.Context, id, token string) error {
	args := m.Called(ctx, id, token)
	return args.Error(0)
}

type mockGroupMemberRepo struct {
	mock.Mock
}

func (m *mockGroupMemberRepo) GetByID(ctx context.Context, id string) (GroupMember, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(GroupMember), args.Error(1)
}

func (m *mockGroupMemberRepo) Get(ctx context.Context, groupID, userID string) (GroupMember, error) {
	args := m.Called(ctx, groupID, userID)
	return args.Get(0).(GroupMember), args.Error(1)
}

func (m *mockGroupMemberRepo) ListByGroup(ctx context.Context, groupID string) ([]GroupMemberWithUser, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).([]GroupMemberWithUser), args.Error(1)
}

func (m *mockGroupMemberRepo) Create(ctx context.Context, groupID, userID string, isAdmin bool) (GroupMember, error) {
	args := m.Called(ctx, groupID, userID, isAdmin)
	return args.Get(0).(GroupMember), args.Error(1)
}

func (m *mockGroupMemberRepo) UpdateRole(ctx context.Context, groupID, userID string, isAdmin bool) (GroupMember, error) {
	args := m.Called(ctx, groupID, userID, isAdmin)
	return args.Get(0).(GroupMember), args.Error(1)
}

func (m *mockGroupMemberRepo) Remove(ctx context.Context, groupID, userID string) error {
	args := m.Called(ctx, groupID, userID)
	return args.Error(0)
}

func (m *mockGroupMemberRepo) Count(ctx context.Context, groupID string) (int64, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).(int64), args.Error(1)
}

func TestGroupService_Create(t *testing.T) {
	mockGroup := new(mockGroupRepo)
	mockMember := new(mockGroupMemberRepo)
	svc := NewGroupService(mockGroup, mockMember, nil)
	ctx := context.Background()

	params := CreateGroupParams{Name: "Test Group"}
	mockGroup.On("Create", ctx, params).Return(Group{ID: "1", Name: "Test Group"}, nil)
	mockMember.On("Create", ctx, "1", "user1", true).Return(GroupMember{}, nil)

	group, err := svc.Create(ctx, params, "user1")

	assert.NoError(t, err)
	assert.Equal(t, "1", group.ID)
	mockGroup.AssertExpectations(t)
	mockMember.AssertExpectations(t)
}

func TestGroupService_GetByID_Success(t *testing.T) {
	mockGroup := new(mockGroupRepo)
	svc := NewGroupService(mockGroup, new(mockGroupMemberRepo), nil)
	ctx := context.Background()

	mockGroup.On("GetByID", ctx, "1").Return(Group{ID: "1", Name: "Test"}, nil)

	group, err := svc.GetByID(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, "1", group.ID)
}

func TestGroupService_GetByID_NotFound(t *testing.T) {
	mockGroup := new(mockGroupRepo)
	svc := NewGroupService(mockGroup, new(mockGroupMemberRepo), nil)
	ctx := context.Background()

	mockGroup.On("GetByID", ctx, "999").Return(Group{}, assert.AnError)

	_, err := svc.GetByID(ctx, "999")

	assert.ErrorIs(t, err, ErrGroupNotFound)
}

func TestGroupService_RemoveMember_Owner(t *testing.T) {
	mockGroup := new(mockGroupRepo)
	mockMember := new(mockGroupMemberRepo)
	svc := NewGroupService(mockGroup, mockMember, nil)
	ctx := context.Background()

	mockMember.On("Get", ctx, "g1", "u1").Return(GroupMember{IsAdmin: true}, nil)

	err := svc.RemoveMember(ctx, "g1", "u1")

	assert.ErrorIs(t, err, ErrCannotRemoveOwner)
	mockMember.AssertExpectations(t)
}

func TestGroupService_RemoveMember_Success(t *testing.T) {
	mockGroup := new(mockGroupRepo)
	mockMember := new(mockGroupMemberRepo)
	svc := NewGroupService(mockGroup, mockMember, nil)
	ctx := context.Background()

	mockMember.On("Get", ctx, "g1", "u1").Return(GroupMember{IsAdmin: false}, nil)
	mockMember.On("Remove", ctx, "g1", "u1").Return(nil)

	err := svc.RemoveMember(ctx, "g1", "u1")

	assert.NoError(t, err)
	mockMember.AssertExpectations(t)
}

func TestGroupService_GetMember_NotFound(t *testing.T) {
	mockGroup := new(mockGroupRepo)
	mockMember := new(mockGroupMemberRepo)
	svc := NewGroupService(mockGroup, mockMember, nil)
	ctx := context.Background()

	mockMember.On("Get", ctx, "g1", "u1").Return(GroupMember{}, assert.AnError)

	_, err := svc.GetMember(ctx, "g1", "u1")

	assert.ErrorIs(t, err, ErrMemberNotFound)
}

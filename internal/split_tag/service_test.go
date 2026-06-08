package split_tag

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSplitTagRepo struct {
	mock.Mock
}

func (m *mockSplitTagRepo) GetByID(ctx context.Context, id string) (SplitTag, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(SplitTag), args.Error(1)
}

func (m *mockSplitTagRepo) ListByGroup(ctx context.Context, groupID string) ([]SplitTag, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).([]SplitTag), args.Error(1)
}

func (m *mockSplitTagRepo) Create(ctx context.Context, groupID, name, createdBy string) (SplitTag, error) {
	args := m.Called(ctx, groupID, name, createdBy)
	return args.Get(0).(SplitTag), args.Error(1)
}

func (m *mockSplitTagRepo) Delete(ctx context.Context, id, groupID string) error {
	args := m.Called(ctx, id, groupID)
	return args.Error(0)
}

func (m *mockSplitTagRepo) ListMembers(ctx context.Context, splitTagID string) ([]SplitTagMemberWithUser, error) {
	args := m.Called(ctx, splitTagID)
	return args.Get(0).([]SplitTagMemberWithUser), args.Error(1)
}

func (m *mockSplitTagRepo) AddMember(ctx context.Context, splitTagID, userID string) (SplitTagMember, error) {
	args := m.Called(ctx, splitTagID, userID)
	return args.Get(0).(SplitTagMember), args.Error(1)
}

func (m *mockSplitTagRepo) RemoveMember(ctx context.Context, splitTagID, userID string) error {
	args := m.Called(ctx, splitTagID, userID)
	return args.Error(0)
}

func TestSplitTagService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(mockSplitTagRepo)
	svc := NewSplitTagService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetByID", ctx, "999").Return(SplitTag{}, assert.AnError)

	_, err := svc.GetByID(ctx, "999")

	assert.ErrorIs(t, err, ErrSplitTagNotFound)
}

func TestSplitTagService_GetByID_Success(t *testing.T) {
	mockRepo := new(mockSplitTagRepo)
	svc := NewSplitTagService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetByID", ctx, "1").Return(SplitTag{ID: "1", Name: "Test Tag"}, nil)

	tag, err := svc.GetByID(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, "1", tag.ID)
}

func TestSplitTagService_Create(t *testing.T) {
	mockRepo := new(mockSplitTagRepo)
	svc := NewSplitTagService(mockRepo)
	ctx := context.Background()

	mockRepo.On("Create", ctx, "g1", "Test Tag", "u1").Return(SplitTag{ID: "st1", Name: "Test Tag"}, nil)

	tag, err := svc.Create(ctx, "g1", "Test Tag", "u1")

	assert.NoError(t, err)
	assert.Equal(t, "st1", tag.ID)
}

func TestSplitTagService_Delete(t *testing.T) {
	mockRepo := new(mockSplitTagRepo)
	svc := NewSplitTagService(mockRepo)
	ctx := context.Background()

	mockRepo.On("Delete", ctx, "st1", "g1").Return(nil)

	err := svc.Delete(ctx, "st1", "g1")

	assert.NoError(t, err)
}

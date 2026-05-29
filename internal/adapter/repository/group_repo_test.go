package repository

import (
	"context"
	"testing"

	"doheem-backend/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroupRepo_CreateAndGetByID(t *testing.T) {
	q := newTxQueries(t)
	repo := NewGroupRepo(q)
	ctx := context.Background()

	group, err := repo.Create(ctx, domain.CreateGroupParams{
		Name:     "Test Group",
		Currency: "USD",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, group.ID)
	assert.Equal(t, "Test Group", group.Name)
	assert.Equal(t, "USD", group.Currency)
	assert.True(t, group.IsActive)

	got, err := repo.GetByID(ctx, group.ID)
	require.NoError(t, err)
	assert.Equal(t, group.ID, got.ID)
}

func TestGroupRepo_Update(t *testing.T) {
	q := newTxQueries(t)
	repo := NewGroupRepo(q)
	ctx := context.Background()

	group, err := repo.Create(ctx, domain.CreateGroupParams{
		Name:     "Original",
		Currency: "BRL",
	})
	require.NoError(t, err)

	newName := "Updated"
	updated, err := repo.Update(ctx, group.ID, domain.UpdateGroupParams{
		Name: &newName,
	})
	require.NoError(t, err)
	assert.Equal(t, "Updated", updated.Name)
}

func TestGroupRepo_SoftDelete(t *testing.T) {
	q := newTxQueries(t)
	repo := NewGroupRepo(q)
	ctx := context.Background()

	group, err := repo.Create(ctx, domain.CreateGroupParams{Name: "To Delete", Currency: "BRL"})
	require.NoError(t, err)

	err = repo.SoftDelete(ctx, group.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, group.ID)
	assert.Error(t, err)
}

func TestGroupRepo_ActivateDeactivate(t *testing.T) {
	q := newTxQueries(t)
	repo := NewGroupRepo(q)
	ctx := context.Background()

	group, err := repo.Create(ctx, domain.CreateGroupParams{Name: "Toggle", Currency: "BRL"})
	require.NoError(t, err)

	err = repo.Deactivate(ctx, group.ID)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, group.ID)
	require.NoError(t, err)
	assert.False(t, got.IsActive)
	assert.NotNil(t, got.InactiveSince)

	err = repo.Activate(ctx, group.ID)
	require.NoError(t, err)

	got, err = repo.GetByID(ctx, group.ID)
	require.NoError(t, err)
	assert.True(t, got.IsActive)
}

func TestGroupMemberRepo_Create(t *testing.T) {
	q := newTxQueries(t)
	groupRepo := NewGroupRepo(q)
	memberRepo := NewGroupMemberRepo(q)
	ctx := context.Background()

	user := createTestUser(t, q)
	group, err := groupRepo.Create(ctx, domain.CreateGroupParams{Name: "Members", Currency: "BRL"})
	require.NoError(t, err)

	member, err := memberRepo.Create(ctx, group.ID, user.ID, "admin")
	require.NoError(t, err)
	assert.NotEmpty(t, member.ID)
	assert.Equal(t, "admin", member.Role)
	assert.True(t, member.IsActive)
}

func TestGroupMemberRepo_ListByGroup(t *testing.T) {
	q := newTxQueries(t)
	groupRepo := NewGroupRepo(q)
	memberRepo := NewGroupMemberRepo(q)
	ctx := context.Background()

	user1 := createTestUser(t, q)
	user2 := createTestUser(t, q)
	group, err := groupRepo.Create(ctx, domain.CreateGroupParams{Name: "List Members", Currency: "BRL"})
	require.NoError(t, err)

	memberRepo.Create(ctx, group.ID, user1.ID, "admin")
	memberRepo.Create(ctx, group.ID, user2.ID, "member")

	members, err := memberRepo.ListByGroup(ctx, group.ID)
	require.NoError(t, err)
	assert.Len(t, members, 2)
}

func TestGroupMemberRepo_Get(t *testing.T) {
	q := newTxQueries(t)
	groupRepo := NewGroupRepo(q)
	memberRepo := NewGroupMemberRepo(q)
	ctx := context.Background()

	user := createTestUser(t, q)
	group, err := groupRepo.Create(ctx, domain.CreateGroupParams{Name: "Get Member", Currency: "BRL"})
	require.NoError(t, err)

	member, err := memberRepo.Create(ctx, group.ID, user.ID, "member")
	require.NoError(t, err)

	got, err := memberRepo.Get(ctx, group.ID, user.ID)
	require.NoError(t, err)
	assert.Equal(t, member.ID, got.ID)
}

func TestGroupMemberRepo_Remove(t *testing.T) {
	q := newTxQueries(t)
	groupRepo := NewGroupRepo(q)
	memberRepo := NewGroupMemberRepo(q)
	ctx := context.Background()

	user := createTestUser(t, q)
	group, err := groupRepo.Create(ctx, domain.CreateGroupParams{Name: "Remove Member", Currency: "BRL"})
	require.NoError(t, err)

	memberRepo.Create(ctx, group.ID, user.ID, "member")
	err = memberRepo.Remove(ctx, group.ID, user.ID)
	require.NoError(t, err)

	got, err := memberRepo.Get(ctx, group.ID, user.ID)
	require.NoError(t, err)
	assert.False(t, got.IsActive)
	assert.NotNil(t, got.LeftAt)
}

func TestGroupMemberRepo_UpdateRole(t *testing.T) {
	q := newTxQueries(t)
	groupRepo := NewGroupRepo(q)
	memberRepo := NewGroupMemberRepo(q)
	ctx := context.Background()

	user := createTestUser(t, q)
	group, err := groupRepo.Create(ctx, domain.CreateGroupParams{Name: "Role Update", Currency: "BRL"})
	require.NoError(t, err)

	memberRepo.Create(ctx, group.ID, user.ID, "member")
	updated, err := memberRepo.UpdateRole(ctx, group.ID, user.ID, "admin")
	require.NoError(t, err)
	assert.Equal(t, "admin", updated.Role)
}

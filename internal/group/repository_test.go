package group_test

import (
	"context"
	"testing"

	"doheem-backend/internal/dbtest"
	"doheem-backend/internal/group"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroupRepo_CreateAndGetByID(t *testing.T) {
	q := dbtest.NewTxQueries(t)
	repo := group.NewGroupRepo(q)
	ctx := context.Background()

	createdGroup, err := repo.Create(ctx, group.CreateGroupParams{
		Name: "Test Group",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, createdGroup.ID)
	assert.Equal(t, "Test Group", createdGroup.Name)

	got, err := repo.GetByID(ctx, createdGroup.ID)
	require.NoError(t, err)
	assert.Equal(t, createdGroup.ID, got.ID)
}

func TestGroupRepo_Update(t *testing.T) {
	q := dbtest.NewTxQueries(t)
	repo := group.NewGroupRepo(q)
	ctx := context.Background()

	createdGroup, err := repo.Create(ctx, group.CreateGroupParams{
		Name: "Original",
	})
	require.NoError(t, err)

	newName := "Updated"
	updated, err := repo.Update(ctx, createdGroup.ID, group.UpdateGroupParams{
		Name: &newName,
	})
	require.NoError(t, err)
	assert.Equal(t, "Updated", updated.Name)
}

func TestGroupMemberRepo_Create(t *testing.T) {
	q := dbtest.NewTxQueries(t)
	groupRepo := group.NewGroupRepo(q)
	memberRepo := group.NewGroupMemberRepo(q)
	ctx := context.Background()

	user := dbtest.CreateTestUser(t, q)
	createdGroup, err := groupRepo.Create(ctx, group.CreateGroupParams{Name: "Members"})
	require.NoError(t, err)

	member, err := memberRepo.Create(ctx, createdGroup.ID, user.ID, true)
	require.NoError(t, err)
	assert.NotEmpty(t, member.ID)
	assert.True(t, member.IsAdmin)
}

func TestGroupMemberRepo_ListByGroup(t *testing.T) {
	q := dbtest.NewTxQueries(t)
	groupRepo := group.NewGroupRepo(q)
	memberRepo := group.NewGroupMemberRepo(q)
	ctx := context.Background()

	user1 := dbtest.CreateTestUser(t, q)
	user2 := dbtest.CreateTestUser(t, q)
	createdGroup, err := groupRepo.Create(ctx, group.CreateGroupParams{Name: "List Members"})
	require.NoError(t, err)

	memberRepo.Create(ctx, createdGroup.ID, user1.ID, true)
	memberRepo.Create(ctx, createdGroup.ID, user2.ID, false)

	members, err := memberRepo.ListByGroup(ctx, createdGroup.ID)
	require.NoError(t, err)
	assert.Len(t, members, 2)
}

func TestGroupMemberRepo_Get(t *testing.T) {
	q := dbtest.NewTxQueries(t)
	groupRepo := group.NewGroupRepo(q)
	memberRepo := group.NewGroupMemberRepo(q)
	ctx := context.Background()

	user := dbtest.CreateTestUser(t, q)
	createdGroup, err := groupRepo.Create(ctx, group.CreateGroupParams{Name: "Get Member"})
	require.NoError(t, err)

	member, err := memberRepo.Create(ctx, createdGroup.ID, user.ID, false)
	require.NoError(t, err)

	got, err := memberRepo.Get(ctx, createdGroup.ID, user.ID)
	require.NoError(t, err)
	assert.Equal(t, member.ID, got.ID)
}

func TestGroupMemberRepo_Remove(t *testing.T) {
	q := dbtest.NewTxQueries(t)
	groupRepo := group.NewGroupRepo(q)
	memberRepo := group.NewGroupMemberRepo(q)
	ctx := context.Background()

	user := dbtest.CreateTestUser(t, q)
	createdGroup, err := groupRepo.Create(ctx, group.CreateGroupParams{Name: "Remove Member"})
	require.NoError(t, err)

	memberRepo.Create(ctx, createdGroup.ID, user.ID, false)
	err = memberRepo.Remove(ctx, createdGroup.ID, user.ID)
	require.NoError(t, err)

	_, err = memberRepo.Get(ctx, createdGroup.ID, user.ID)
	assert.Error(t, err)
}

func TestGroupMemberRepo_UpdateRole(t *testing.T) {
	q := dbtest.NewTxQueries(t)
	groupRepo := group.NewGroupRepo(q)
	memberRepo := group.NewGroupMemberRepo(q)
	ctx := context.Background()

	user := dbtest.CreateTestUser(t, q)
	createdGroup, err := groupRepo.Create(ctx, group.CreateGroupParams{Name: "Role Update"})
	require.NoError(t, err)

	memberRepo.Create(ctx, createdGroup.ID, user.ID, false)
	updated, err := memberRepo.UpdateRole(ctx, createdGroup.ID, user.ID, true)
	require.NoError(t, err)
	assert.True(t, updated.IsAdmin)
}

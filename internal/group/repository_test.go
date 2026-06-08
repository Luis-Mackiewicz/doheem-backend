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
		Name:     "Test Group",
		Currency: "USD",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, createdGroup.ID)
	assert.Equal(t, "Test Group", createdGroup.Name)
	assert.Equal(t, "USD", createdGroup.Currency)
	assert.True(t, createdGroup.IsActive)

	got, err := repo.GetByID(ctx, createdGroup.ID)
	require.NoError(t, err)
	assert.Equal(t, createdGroup.ID, got.ID)
}

func TestGroupRepo_Update(t *testing.T) {
	q := dbtest.NewTxQueries(t)
	repo := group.NewGroupRepo(q)
	ctx := context.Background()

	createdGroup, err := repo.Create(ctx, group.CreateGroupParams{
		Name:     "Original",
		Currency: "BRL",
	})
	require.NoError(t, err)

	newName := "Updated"
	updated, err := repo.Update(ctx, createdGroup.ID, group.UpdateGroupParams{
		Name: &newName,
	})
	require.NoError(t, err)
	assert.Equal(t, "Updated", updated.Name)
}

func TestGroupRepo_SoftDelete(t *testing.T) {
	q := dbtest.NewTxQueries(t)
	repo := group.NewGroupRepo(q)
	ctx := context.Background()

	createdGroup, err := repo.Create(ctx, group.CreateGroupParams{Name: "To Delete", Currency: "BRL"})
	require.NoError(t, err)

	err = repo.SoftDelete(ctx, createdGroup.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, createdGroup.ID)
	assert.Error(t, err)
}

func TestGroupRepo_ActivateDeactivate(t *testing.T) {
	q := dbtest.NewTxQueries(t)
	repo := group.NewGroupRepo(q)
	ctx := context.Background()

	createdGroup, err := repo.Create(ctx, group.CreateGroupParams{Name: "Toggle", Currency: "BRL"})
	require.NoError(t, err)

	err = repo.Deactivate(ctx, createdGroup.ID)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, createdGroup.ID)
	require.NoError(t, err)
	assert.False(t, got.IsActive)
	assert.NotNil(t, got.InactiveSince)

	err = repo.Activate(ctx, createdGroup.ID)
	require.NoError(t, err)

	got, err = repo.GetByID(ctx, createdGroup.ID)
	require.NoError(t, err)
	assert.True(t, got.IsActive)
}

func TestGroupMemberRepo_Create(t *testing.T) {
	q := dbtest.NewTxQueries(t)
	groupRepo := group.NewGroupRepo(q)
	memberRepo := group.NewGroupMemberRepo(q)
	ctx := context.Background()

	user := dbtest.CreateTestUser(t, q)
	createdGroup, err := groupRepo.Create(ctx, group.CreateGroupParams{Name: "Members", Currency: "BRL"})
	require.NoError(t, err)

	member, err := memberRepo.Create(ctx, createdGroup.ID, user.ID, "admin")
	require.NoError(t, err)
	assert.NotEmpty(t, member.ID)
	assert.Equal(t, "admin", member.Role)
	assert.True(t, member.IsActive)
}

func TestGroupMemberRepo_ListByGroup(t *testing.T) {
	q := dbtest.NewTxQueries(t)
	groupRepo := group.NewGroupRepo(q)
	memberRepo := group.NewGroupMemberRepo(q)
	ctx := context.Background()

	user1 := dbtest.CreateTestUser(t, q)
	user2 := dbtest.CreateTestUser(t, q)
	createdGroup, err := groupRepo.Create(ctx, group.CreateGroupParams{Name: "List Members", Currency: "BRL"})
	require.NoError(t, err)

	memberRepo.Create(ctx, createdGroup.ID, user1.ID, "admin")
	memberRepo.Create(ctx, createdGroup.ID, user2.ID, "member")

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
	createdGroup, err := groupRepo.Create(ctx, group.CreateGroupParams{Name: "Get Member", Currency: "BRL"})
	require.NoError(t, err)

	member, err := memberRepo.Create(ctx, createdGroup.ID, user.ID, "member")
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
	createdGroup, err := groupRepo.Create(ctx, group.CreateGroupParams{Name: "Remove Member", Currency: "BRL"})
	require.NoError(t, err)

	memberRepo.Create(ctx, createdGroup.ID, user.ID, "member")
	err = memberRepo.Remove(ctx, createdGroup.ID, user.ID)
	require.NoError(t, err)

	got, err := memberRepo.Get(ctx, createdGroup.ID, user.ID)
	require.NoError(t, err)
	assert.False(t, got.IsActive)
	assert.NotNil(t, got.LeftAt)
}

func TestGroupMemberRepo_UpdateRole(t *testing.T) {
	q := dbtest.NewTxQueries(t)
	groupRepo := group.NewGroupRepo(q)
	memberRepo := group.NewGroupMemberRepo(q)
	ctx := context.Background()

	user := dbtest.CreateTestUser(t, q)
	createdGroup, err := groupRepo.Create(ctx, group.CreateGroupParams{Name: "Role Update", Currency: "BRL"})
	require.NoError(t, err)

	memberRepo.Create(ctx, createdGroup.ID, user.ID, "member")
	updated, err := memberRepo.UpdateRole(ctx, createdGroup.ID, user.ID, "admin")
	require.NoError(t, err)
	assert.Equal(t, "admin", updated.Role)
}

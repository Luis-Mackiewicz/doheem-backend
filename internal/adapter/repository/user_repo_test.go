package repository

import (
	"context"
	"testing"

	"doheem-backend/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepo_CreateAndGetByID(t *testing.T) {
	q := newTxQueries(t)
	repo := NewUserRepo(q)
	ctx := context.Background()

	params := domain.CreateUserParams{
		Name:         "John Doe",
		Email:        "john@example.com",
		PasswordHash: "hashed-password",
	}

	user, err := repo.Create(ctx, params)
	require.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)

	got, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
	assert.Equal(t, user.Name, got.Name)
}

func TestUserRepo_GetByEmail(t *testing.T) {
	q := newTxQueries(t)
	repo := NewUserRepo(q)
	ctx := context.Background()

	params := domain.CreateUserParams{
		Name:         "Jane Doe",
		Email:        "jane@example.com",
		PasswordHash: "hashed-password",
	}

	user, err := repo.Create(ctx, params)
	require.NoError(t, err)

	got, err := repo.GetByEmail(ctx, "jane@example.com")
	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
}

func TestUserRepo_GetByEmail_NotFound(t *testing.T) {
	q := newTxQueries(t)
	repo := NewUserRepo(q)
	ctx := context.Background()

	_, err := repo.GetByEmail(ctx, "nonexistent@example.com")
	assert.Error(t, err)
}

func TestUserRepo_Update(t *testing.T) {
	q := newTxQueries(t)
	repo := NewUserRepo(q)
	ctx := context.Background()

	user, err := repo.Create(ctx, domain.CreateUserParams{
		Name:         "Old Name",
		Email:        "update@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	newName := "New Name"
	newEmail := "newemail@example.com"
	updated, err := repo.Update(ctx, user.ID, domain.UpdateUserParams{
		Name: &newName,
		Email: &newEmail,
	})
	require.NoError(t, err)
	assert.Equal(t, "New Name", updated.Name)
	assert.Equal(t, "newemail@example.com", updated.Email)
}

func TestUserRepo_UpdatePassword(t *testing.T) {
	q := newTxQueries(t)
	repo := NewUserRepo(q)
	ctx := context.Background()

	user, err := repo.Create(ctx, domain.CreateUserParams{
		Name:         "Password User",
		Email:        "password@example.com",
		PasswordHash: "old-hash",
	})
	require.NoError(t, err)

	err = repo.UpdatePassword(ctx, user.ID, "new-hash")
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "new-hash", got.PasswordHash)
}

func TestUserRepo_Delete(t *testing.T) {
	q := newTxQueries(t)
	repo := NewUserRepo(q)
	ctx := context.Background()

	user, err := repo.Create(ctx, domain.CreateUserParams{
		Name:         "Delete Me",
		Email:        "delete@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	err = repo.Delete(ctx, user.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, user.ID)
	assert.Error(t, err)
}

func TestUserRepo_List(t *testing.T) {
	q := newTxQueries(t)
	repo := NewUserRepo(q)
	ctx := context.Background()

	users, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Empty(t, users)
}

func TestUserRepo_EmailUnique(t *testing.T) {
	q := newTxQueries(t)
	repo := NewUserRepo(q)
	ctx := context.Background()

	_, err := repo.Create(ctx, domain.CreateUserParams{
		Name:         "First",
		Email:        "unique@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	_, err = repo.Create(ctx, domain.CreateUserParams{
		Name:         "Second",
		Email:        "unique@example.com",
		PasswordHash: "hash",
	})
	assert.Error(t, err)
}

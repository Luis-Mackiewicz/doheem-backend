package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"doheem-backend/internal/adapter/repository/db"
	"doheem-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:18.4",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(120*time.Second),
		),
	)
	if err != nil {
		panic("failed to start postgres container: " + err.Error())
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic("failed to get connection string: " + err.Error())
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	migration, err := os.ReadFile("migrations/001_init.up.sql")
	if err != nil {
		panic("failed to read migration file: " + err.Error())
	}

	if _, err := pool.Exec(ctx, string(migration)); err != nil {
		panic("failed to run migration: " + err.Error())
	}

	if err := pool.Ping(ctx); err != nil {
		panic("failed to ping database: " + err.Error())
	}

	testPool = pool

	code := m.Run()

	pool.Close()
	pgContainer.Terminate(ctx)
	os.Exit(code)
}

func newTxQueries(t *testing.T) *db.Queries {
	t.Helper()
	ctx := context.Background()
	tx, err := testPool.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	t.Cleanup(func() {
		tx.Rollback(ctx)
	})
	return db.New(tx)
}

var testUserCounter int

func createTestUser(t *testing.T, q *db.Queries) domain.User {
	t.Helper()
	testUserCounter++
	repo := NewUserRepo(q)
	user, err := repo.Create(context.Background(), domain.CreateUserParams{
		Name:         "Test User",
		Email:        "test-" + t.Name() + "-" + fmt.Sprint(testUserCounter) + "@example.com",
		PasswordHash: "hash",
	})
	require.NoError(t, err)
	return user
}

func createTestGroup(t *testing.T, q *db.Queries) domain.Group {
	t.Helper()
	repo := NewGroupRepo(q)
	group, err := repo.Create(context.Background(), domain.CreateGroupParams{
		Name:     "Test Group " + t.Name(),
		Currency: "BRL",
	})
	require.NoError(t, err)
	return group
}

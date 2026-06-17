package dbtest

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"doheem-backend/internal/db"
	"doheem-backend/internal/group"
	"doheem-backend/internal/user"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func strPtr(s string) *string {
	return &s
}

var TestPool *pgxpool.Pool
var testUserCounter int

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

	migrations := []string{
		"../db/migrations/001_init.up.sql",
		"../db/migrations/002_alter_users_not_null.up.sql",
		"../db/migrations/003_drop_group_cnpj_cep.up.sql",
		"../db/migrations/004_add_parent_expense_id_index.up.sql",
		"../db/migrations/005_add_expense_composite_index.up.sql",
		"../db/migrations/006_add_receipt_to_splits.up.sql",
		"../db/migrations/007_add_created_by_to_expenses.up.sql",
		"../db/migrations/008_add_fixed_source_id.up.sql",
	}
	for _, path := range migrations {
		migration, err := os.ReadFile(path)
		if err != nil {
			panic("failed to read migration file " + path + ": " + err.Error())
		}
		if _, err := pool.Exec(ctx, string(migration)); err != nil {
			panic("failed to run migration " + path + ": " + err.Error())
		}
	}

	if err := pool.Ping(ctx); err != nil {
		panic("failed to ping database: " + err.Error())
	}

	TestPool = pool

	code := m.Run()

	pool.Close()
	pgContainer.Terminate(ctx)
	os.Exit(code)
}

func NewTxQueries(t *testing.T) *db.Queries {
	t.Helper()
	ctx := context.Background()
	tx, err := TestPool.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	t.Cleanup(func() {
		tx.Rollback(ctx)
	})
	return db.New(tx)
}

func CreateTestUser(t *testing.T, q *db.Queries) user.User {
	t.Helper()
	testUserCounter++
	repo := user.NewUserRepo(q)
	now := time.Now()
	u, err := repo.Create(context.Background(), user.CreateUserParams{
		Name:         "Test User",
		Email:        "test-" + t.Name() + "-" + fmt.Sprint(testUserCounter) + "@example.com",
		PasswordHash: "hash",
		Phone:        strPtr("11999999999"),
		Document:     strPtr("123.456.789-00"),
		BirthDate:    &now,
		Cep:          strPtr("01001-000"),
	})
	require.NoError(t, err)
	return u
}

func CreateTestGroup(t *testing.T, q *db.Queries, owner user.User) group.Group {
	t.Helper()
	testUserCounter++
	repo := group.NewGroupRepo(q)
	memberRepo := group.NewGroupMemberRepo(q)
	g, err := repo.Create(context.Background(), group.CreateGroupParams{
		Name: "Test Group " + t.Name(),
	})
	require.NoError(t, err)
	_, err = memberRepo.Create(context.Background(), g.ID, owner.ID, true)
	require.NoError(t, err)
	return g
}

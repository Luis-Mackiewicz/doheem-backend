package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"doheem-backend/internal/config"
	"doheem-backend/internal/db"
	"doheem-backend/internal/expense"
	"doheem-backend/internal/group"
	adapterhttp "doheem-backend/internal/http"
	"doheem-backend/internal/notification"
	"doheem-backend/internal/task"
	"doheem-backend/internal/user"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := config.Load()
	initLogger(cfg)

	pool := initPostgres(ctx, cfg.DatabaseURL)
	defer pool.Close()

	rdb := initRedis(ctx, cfg.RedisURL)
	defer rdb.Close()

	q := db.New(pool)

	userRepo := user.NewUserRepo(q)
	groupRepo := group.NewGroupRepo(q)
	groupMemberRepo := group.NewGroupMemberRepo(q)
	expenseRepo := expense.NewExpenseRepo(q)
	expenseSplitRepo := expense.NewExpenseSplitRepo(q)
	categoryRepo := expense.NewExpenseCategoryRepo(q)
	taskRepo := task.NewTaskRepo(q)
	taskOccurrenceRepo := task.NewTaskOccurrenceRepo(q)
	notificationRepo := notification.NewNotificationRepo(q)
	refreshTokenRepo := user.NewRefreshTokenRepo(q)

	userSvc := user.NewUserService(userRepo, refreshTokenRepo)
	groupSvc := group.NewGroupService(groupRepo, groupMemberRepo)
	notificationSvc := notification.NewNotificationService(notificationRepo)
	expenseSvc := expense.NewExpenseService(expenseRepo, expenseSplitRepo, categoryRepo, groupMemberRepo, notificationRepo)
	taskSvc := task.NewTaskService(taskRepo, taskOccurrenceRepo, groupMemberRepo, notificationRepo)

	startFixedExpenseScheduler(expenseSvc)

	jwtSvc := adapterhttp.NewJWTService(cfg.JWTSecret, cfg.JWTExpiresIn, cfg.JWTRefreshExpiresIn)
	router := adapterhttp.NewRouter(jwtSvc, userSvc, groupSvc, expenseSvc, taskSvc, notificationSvc, rdb, pool)

	runServer(ctx, router.Handler(), cfg.Port)
}

func initLogger(cfg config.Config) {
	var slogHandler slog.Handler
	opts := &slog.HandlerOptions{AddSource: cfg.AppEnv == "development"}

	if cfg.LogFormat == "json" {
		slogHandler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		slogHandler = slog.NewTextHandler(os.Stdout, opts)
	}
	slog.SetDefault(slog.New(slogHandler))
}

func initPostgres(ctx context.Context, url string) *pgxpool.Pool {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		slog.Error("failed to parse database config", "error", err)
		os.Exit(1)
	}
	cfg.MaxConns = 25

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	if err := pool.Ping(ctx); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to database")
	return pool
}

func initRedis(ctx context.Context, url string) *redis.Client {
	rdb, err := redis.ParseURL(url)
	if err != nil {
		slog.Error("failed to parse redis URL", "error", err)
		os.Exit(1)
	}

	client := redis.NewClient(rdb)
	if err := client.Ping(ctx).Err(); err != nil {
		slog.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to redis")
	return client
}

func runServer(ctx context.Context, handler http.Handler, port string) {
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("Doheem server is running", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}
	slog.Info("server stopped")
}

func startFixedExpenseScheduler(svc *expense.ExpenseService) {
	go func() {
		ctx := context.Background()
		if err := svc.AutoRestoreFixedExpenses(ctx); err != nil {
			slog.Warn("fixed expense restore on startup failed", "error", err)
		} else {
			slog.Info("fixed expenses restored on startup")
		}
	}()

	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			ctx := context.Background()
			if err := svc.AutoRestoreFixedExpenses(ctx); err != nil {
				slog.Warn("fixed expense restore failed", "error", err)
			}
		}
	}()
}

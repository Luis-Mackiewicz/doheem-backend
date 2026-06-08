package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	adapterhttp "doheem-backend/internal/adapter/http"
	"doheem-backend/internal/adapter/repository"
	"doheem-backend/internal/adapter/repository/db"
	"doheem-backend/internal/config"
	"doheem-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := config.Load()

	var slogHandler slog.Handler
	opts := &slog.HandlerOptions{AddSource: cfg.AppEnv == "development"}
	if cfg.LogFormat == "json" {
		slogHandler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		slogHandler = slog.NewTextHandler(os.Stdout, opts)
	}
	slog.SetDefault(slog.New(slogHandler))

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to database")

	redisOpts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		slog.Error("failed to parse redis URL", "error", err)
		os.Exit(1)
	}
	rdb := redis.NewClient(redisOpts)
	if err := rdb.Ping(ctx).Err(); err != nil {
		slog.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}
	defer rdb.Close()
	slog.Info("connected to redis")

	queries := db.New(pool)

	userRepo := repository.NewUserRepo(queries)
	groupRepo := repository.NewGroupRepo(queries)
	groupMemberRepo := repository.NewGroupMemberRepo(queries)
	expenseRepo := repository.NewExpenseRepo(queries)
	expenseSplitRepo := repository.NewExpenseSplitRepo(queries)
	installmentRepo := repository.NewInstallmentRepo(queries)
	categoryRepo := repository.NewExpenseCategoryRepo(queries)
	paymentRepo := repository.NewPaymentRepo(queries)
	paymentAttachmentRepo := repository.NewPaymentAttachmentRepo(queries)
	taskRepo := repository.NewTaskRepo(queries)
	taskOccurrenceRepo := repository.NewTaskOccurrenceRepo(queries)
	inviteRepo := repository.NewInviteRepo(queries)
	notificationRepo := repository.NewNotificationRepo(queries)
	splitTagRepo := repository.NewSplitTagRepo(queries)
	refreshTokenRepo := repository.NewRefreshTokenRepo(queries)

	userSvc := domain.NewUserService(userRepo, refreshTokenRepo)
	groupSvc := domain.NewGroupService(groupRepo, groupMemberRepo)
	notificationSvc := domain.NewNotificationService(notificationRepo)

	expenseSvc := domain.NewExpenseService(expenseRepo, expenseSplitRepo, installmentRepo, categoryRepo, groupMemberRepo)
	paymentSvc := domain.NewPaymentService(paymentRepo, paymentAttachmentRepo)
	taskSvc := domain.NewTaskService(taskRepo, taskOccurrenceRepo)
	inviteSvc := domain.NewInviteService(inviteRepo, groupMemberRepo, groupRepo)
	splitTagSvc := domain.NewSplitTagService(splitTagRepo)

	jwtSvc := adapterhttp.NewJWTService(cfg.JWTSecret, cfg.JWTExpiresIn, cfg.JWTRefreshExpiresIn)

	router := adapterhttp.NewRouter(jwtSvc, userSvc, groupSvc, expenseSvc, paymentSvc, taskSvc, inviteSvc, notificationSvc, splitTagSvc, rdb, pool)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("Doheem server is running", "port", cfg.Port)
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

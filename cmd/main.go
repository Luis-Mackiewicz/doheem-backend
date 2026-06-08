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
	"doheem-backend/internal/invite"
	"doheem-backend/internal/notification"
	"doheem-backend/internal/payment"
	"doheem-backend/internal/split_tag"
	"doheem-backend/internal/task"
	"doheem-backend/internal/user"

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

	q := db.New(pool)

	userRepo := user.NewUserRepo(q)
	groupRepo := group.NewGroupRepo(q)
	groupMemberRepo := group.NewGroupMemberRepo(q)
	expenseRepo := expense.NewExpenseRepo(q)
	expenseSplitRepo := expense.NewExpenseSplitRepo(q)
	installmentRepo := expense.NewInstallmentRepo(q)
	categoryRepo := expense.NewExpenseCategoryRepo(q)
	paymentRepo := payment.NewPaymentRepo(q)
	paymentAttachmentRepo := payment.NewPaymentAttachmentRepo(q)
	taskRepo := task.NewTaskRepo(q)
	taskOccurrenceRepo := task.NewTaskOccurrenceRepo(q)
	inviteRepo := invite.NewInviteRepo(q)
	notificationRepo := notification.NewNotificationRepo(q)
	splitTagRepo := split_tag.NewSplitTagRepo(q)
	refreshTokenRepo := user.NewRefreshTokenRepo(q)

	userSvc := user.NewUserService(userRepo, refreshTokenRepo)
	groupSvc := group.NewGroupService(groupRepo, groupMemberRepo)
	notificationSvc := notification.NewNotificationService(notificationRepo)

	expenseSvc := expense.NewExpenseService(expenseRepo, expenseSplitRepo, installmentRepo, categoryRepo, groupMemberRepo)
	paymentSvc := payment.NewPaymentService(paymentRepo, paymentAttachmentRepo)
	taskSvc := task.NewTaskService(taskRepo, taskOccurrenceRepo)
	inviteSvc := invite.NewInviteService(inviteRepo, groupMemberRepo, groupRepo)
	splitTagSvc := split_tag.NewSplitTagService(splitTagRepo)

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

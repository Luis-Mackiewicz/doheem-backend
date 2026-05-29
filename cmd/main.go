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
	"doheem-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://doheem_dev_user:simple_pswd@localhost:5432/doheem_dev_db?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, dbURL)
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

	userSvc := domain.NewUserService(userRepo)
	groupSvc := domain.NewGroupService(groupRepo, groupMemberRepo)
	expenseSvc := domain.NewExpenseService(expenseRepo, expenseSplitRepo, installmentRepo, categoryRepo, groupMemberRepo)
	paymentSvc := domain.NewPaymentService(paymentRepo, paymentAttachmentRepo)
	taskSvc := domain.NewTaskService(taskRepo, taskOccurrenceRepo)
	inviteSvc := domain.NewInviteService(inviteRepo, groupMemberRepo, groupRepo)
	notificationSvc := domain.NewNotificationService(notificationRepo)
	splitTagSvc := domain.NewSplitTagService(splitTagRepo)

	router := adapterhttp.NewRouter(userSvc, groupSvc, expenseSvc, paymentSvc, taskSvc, inviteSvc, notificationSvc, splitTagSvc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router.Handler(),
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

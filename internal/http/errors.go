package http

import (
	"errors"
	"log/slog"
	"net/http"

	"doheem-backend/internal/expense"
	"doheem-backend/internal/group"
	"doheem-backend/internal/invite"
	"doheem-backend/internal/notification"
	"doheem-backend/internal/payment"
	"doheem-backend/internal/split_tag"
	"doheem-backend/internal/task"
	"doheem-backend/internal/user"
)

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	status := toHTTPStatus(err)
	if status >= http.StatusInternalServerError {
		slog.ErrorContext(r.Context(), "internal server error", "error", err, "request_id", GetRequestID(r.Context()))
	}
	respondError(w, status, err.Error())
}

func toHTTPStatus(err error) int {
	switch {
	case
		errors.Is(err, user.ErrUserNotFound),
		errors.Is(err, group.ErrGroupNotFound),
		errors.Is(err, expense.ErrExpenseNotFound),
		errors.Is(err, payment.ErrPaymentNotFound),
		errors.Is(err, task.ErrTaskNotFound),
		errors.Is(err, invite.ErrInviteNotFound),
		errors.Is(err, notification.ErrNotificationNotFound),
		errors.Is(err, split_tag.ErrSplitTagNotFound),
		errors.Is(err, group.ErrMemberNotFound),
		errors.Is(err, task.ErrTaskOccurrenceNotFound),
		errors.Is(err, expense.ErrCategoryNotFound):
		return http.StatusNotFound

	case
		errors.Is(err, user.ErrInvalidCredentials):
		return http.StatusUnauthorized

	case
		errors.Is(err, group.ErrCannotRemoveOwner):
		return http.StatusForbidden

	case
		errors.Is(err, user.ErrEmailAlreadyExists),
		errors.Is(err, group.ErrMemberAlreadyExists):
		return http.StatusConflict

	case
		errors.Is(err, expense.ErrInvalidSplitTotal),
		errors.Is(err, task.ErrInvalidDueDate),
		errors.Is(err, payment.ErrPaymentAlreadyConfirmed),
		errors.Is(err, payment.ErrPaymentAlreadyCancelled),
		errors.Is(err, invite.ErrInviteExpired),
		errors.Is(err, invite.ErrInviteAlreadyUsed),
		errors.Is(err, invite.ErrInviteRevoked):
		return http.StatusBadRequest

	default:
		return http.StatusInternalServerError
	}
}

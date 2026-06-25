package http

import (
	"errors"
	"log/slog"
	"net/http"

	"doheem-backend/internal/expense"
	"doheem-backend/internal/group"
	"doheem-backend/internal/notification"
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
		errors.Is(err, task.ErrTaskNotFound),
		errors.Is(err, notification.ErrNotificationNotFound),
		errors.Is(err, group.ErrMemberNotFound),
		errors.Is(err, task.ErrTaskOccurrenceNotFound),
		errors.Is(err, expense.ErrCategoryNotFound):
		return http.StatusNotFound

	case
		errors.Is(err, user.ErrInvalidCredentials):
		return http.StatusUnauthorized

	case
		errors.Is(err, group.ErrCannotRemoveOwner),
		errors.Is(err, group.ErrForbidden),
		errors.Is(err, expense.ErrForbidden),
		errors.Is(err, task.ErrForbidden):
		return http.StatusForbidden

	case
		errors.Is(err, user.ErrEmailAlreadyExists),
		errors.Is(err, user.ErrDocumentAlreadyExists),
		errors.Is(err, user.ErrPhoneAlreadyExists),
		errors.Is(err, group.ErrMemberAlreadyExists),
		errors.Is(err, group.ErrGroupFull),
		errors.Is(err, expense.ErrCannotDeleteWithPaidSplits),
		errors.Is(err, expense.ErrCannotEditWithPaidSplits),
		errors.Is(err, expense.ErrSplitAlreadyPaid):
		return http.StatusConflict

	case
		errors.Is(err, expense.ErrInvalidSplitTotal),
		errors.Is(err, expense.ErrCannotEditInstallmentChild),
		errors.Is(err, expense.ErrCannotEditInstallmentParent),
		errors.Is(err, expense.ErrFixedWithInstallments),
		errors.Is(err, task.ErrInvalidDueDate),
		errors.Is(err, notification.ErrReminderLimitExceeded),
		errors.Is(err, notification.ErrReminderTooSoon),
		errors.Is(err, task.ErrInvalidStatusTransition):
		return http.StatusBadRequest

	default:
		return http.StatusInternalServerError
	}
}

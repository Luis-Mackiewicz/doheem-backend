package http

import (
	"errors"
	"log/slog"
	"net/http"

	"doheem-backend/internal/domain"
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
		errors.Is(err, domain.ErrUserNotFound),
		errors.Is(err, domain.ErrGroupNotFound),
		errors.Is(err, domain.ErrExpenseNotFound),
		errors.Is(err, domain.ErrPaymentNotFound),
		errors.Is(err, domain.ErrTaskNotFound),
		errors.Is(err, domain.ErrInviteNotFound),
		errors.Is(err, domain.ErrNotificationNotFound),
		errors.Is(err, domain.ErrSplitTagNotFound),
		errors.Is(err, domain.ErrMemberNotFound),
		errors.Is(err, domain.ErrTaskOccurrenceNotFound),
		errors.Is(err, domain.ErrCategoryNotFound):
		return http.StatusNotFound

	case
		errors.Is(err, domain.ErrInvalidCredentials):
		return http.StatusUnauthorized

	case
		errors.Is(err, domain.ErrCannotRemoveOwner):
		return http.StatusForbidden

	case
		errors.Is(err, domain.ErrEmailAlreadyExists),
		errors.Is(err, domain.ErrMemberAlreadyExists):
		return http.StatusConflict

	case
		errors.Is(err, domain.ErrInvalidSplitTotal),
		errors.Is(err, domain.ErrInvalidDueDate),
		errors.Is(err, domain.ErrPaymentAlreadyConfirmed),
		errors.Is(err, domain.ErrPaymentAlreadyCancelled),
		errors.Is(err, domain.ErrInviteExpired),
		errors.Is(err, domain.ErrInviteAlreadyUsed),
		errors.Is(err, domain.ErrInviteRevoked):
		return http.StatusBadRequest

	default:
		return http.StatusInternalServerError
	}
}

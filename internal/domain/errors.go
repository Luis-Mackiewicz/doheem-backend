package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")

	ErrGroupNotFound       = errors.New("group not found")
	ErrMemberNotFound      = errors.New("member not found")
	ErrMemberAlreadyExists = errors.New("member already exists")
	ErrCannotRemoveOwner   = errors.New("cannot remove owner from group")

	ErrExpenseNotFound   = errors.New("expense not found")
	ErrInvalidSplitTotal = errors.New("split amounts must equal total amount")
	ErrCategoryNotFound  = errors.New("category not found")

	ErrPaymentNotFound           = errors.New("payment not found")
	ErrPaymentAlreadyConfirmed   = errors.New("payment already confirmed")
	ErrPaymentAlreadyCancelled   = errors.New("payment already cancelled")

	ErrTaskNotFound           = errors.New("task not found")
	ErrTaskOccurrenceNotFound = errors.New("task occurrence not found")
	ErrInvalidDueDate         = errors.New("due date must be in the future")

	ErrInviteNotFound    = errors.New("invite not found")
	ErrInviteExpired     = errors.New("invite has expired")
	ErrInviteAlreadyUsed = errors.New("invite already used")
	ErrInviteRevoked     = errors.New("invite has been revoked")

	ErrNotificationNotFound = errors.New("notification not found")

	ErrSplitTagNotFound = errors.New("split tag not found")
)

package invite

import (
	"context"
	"errors"
	"time"
)

type Invite struct {
	ID        string
	GroupID   string
	Code      string
	CreatedBy string
	ExpiresAt time.Time
	UsedAt    *time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

type InviteWithGroup struct {
	Invite
	GroupName string
}

type InviteWithCreator struct {
	Invite
	CreatedByName string
}

type InviteRepository interface {
	GetByID(ctx context.Context, id string) (Invite, error)
	GetByCode(ctx context.Context, code string) (InviteWithGroup, error)
	ListByGroup(ctx context.Context, groupID string) ([]InviteWithCreator, error)
	ListPendingByUser(ctx context.Context, userID string) ([]InviteWithGroup, error)
	Create(ctx context.Context, groupID, code, createdBy string, expiresAt time.Time) (Invite, error)
	Use(ctx context.Context, id string) error
	Revoke(ctx context.Context, id string) error
}

var (
	ErrInviteNotFound    = errors.New("invite not found")
	ErrInviteExpired     = errors.New("invite has expired")
	ErrInviteAlreadyUsed = errors.New("invite already used")
	ErrInviteRevoked     = errors.New("invite has been revoked")
)

package group

import (
	"context"
	"errors"
	"time"
)

type Group struct {
	ID            string
	Name          string
	Currency      string
	IsActive      bool
	InactiveSince *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

type CreateGroupParams struct {
	Name     string
	Currency string
}

type UpdateGroupParams struct {
	Name     *string
	Currency *string
}

type GroupMember struct {
	ID       string
	GroupID  string
	UserID   string
	Role     string
	JoinedAt time.Time
	LeftAt   *time.Time
	IsActive bool
}

type GroupMemberWithUser struct {
	GroupMember
	UserName  string
	UserEmail string
	AvatarURL *string
}

type GroupRepository interface {
	GetByID(ctx context.Context, id string) (Group, error)
	ListByUserID(ctx context.Context, userID string) ([]Group, error)
	Create(ctx context.Context, params CreateGroupParams) (Group, error)
	Update(ctx context.Context, id string, params UpdateGroupParams) (Group, error)
	SoftDelete(ctx context.Context, id string) error
	Deactivate(ctx context.Context, id string) error
	Activate(ctx context.Context, id string) error
}

type GroupMemberRepository interface {
	GetByID(ctx context.Context, id string) (GroupMember, error)
	Get(ctx context.Context, groupID, userID string) (GroupMember, error)
	ListByGroup(ctx context.Context, groupID string) ([]GroupMemberWithUser, error)
	Create(ctx context.Context, groupID, userID, role string) (GroupMember, error)
	UpdateRole(ctx context.Context, groupID, userID, role string) (GroupMember, error)
	Remove(ctx context.Context, groupID, userID string) error
	CountActive(ctx context.Context, groupID string) (int64, error)
}

var (
	ErrGroupNotFound       = errors.New("group not found")
	ErrMemberNotFound      = errors.New("member not found")
	ErrMemberAlreadyExists = errors.New("member already exists")
	ErrCannotRemoveOwner   = errors.New("cannot remove owner from group")
)

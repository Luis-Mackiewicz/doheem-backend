package group

import (
	"context"
	"errors"
	"time"
)

type Group struct {
	ID          string
	Name        string
	Description string
	MonthlyFee  float64
	Cnpj        string
	Cep         string
	PhotoURL    *string
	InviteToken *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateGroupParams struct {
	Name string
}

type UpdateGroupParams struct {
	Name        *string
	Description *string
	MonthlyFee  *float64
	Cnpj        *string
	Cep         *string
	PhotoURL    *string
}

type GroupMember struct {
	ID       string
	GroupID  string
	UserID   string
	IsAdmin  bool
	JoinedAt time.Time
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
	RegenerateInviteToken(ctx context.Context, id, token string) error
}

type GroupMemberRepository interface {
	GetByID(ctx context.Context, id string) (GroupMember, error)
	Get(ctx context.Context, groupID, userID string) (GroupMember, error)
	ListByGroup(ctx context.Context, groupID string) ([]GroupMemberWithUser, error)
	Create(ctx context.Context, groupID, userID string, isAdmin bool) (GroupMember, error)
	UpdateRole(ctx context.Context, groupID, userID string, isAdmin bool) (GroupMember, error)
	Remove(ctx context.Context, groupID, userID string) error
	Count(ctx context.Context, groupID string) (int64, error)
}

var (
	ErrGroupNotFound       = errors.New("group not found")
	ErrMemberNotFound      = errors.New("member not found")
	ErrMemberAlreadyExists = errors.New("member already exists")
	ErrCannotRemoveOwner   = errors.New("cannot remove the group owner")
)

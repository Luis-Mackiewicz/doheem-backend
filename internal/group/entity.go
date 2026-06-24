package group

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

type Group struct {
	ID          string
	Name        string
	Description string
	MonthlyFee  decimal.Decimal
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
	MonthlyFee  *decimal.Decimal
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
	UserPhone *string
	AvatarURL *string
}

type GroupRepository interface {
	GetByID(ctx context.Context, id string) (Group, error)
	ListByUserID(ctx context.Context, userID string) ([]Group, error)
	Create(ctx context.Context, params CreateGroupParams) (Group, error)
	Update(ctx context.Context, id string, params UpdateGroupParams) (Group, error)
	RegenerateInviteToken(ctx context.Context, id, token string) error
	Delete(ctx context.Context, id string) error
}

type GroupMemberRepository interface {
	GetByID(ctx context.Context, id string) (GroupMember, error)
	Get(ctx context.Context, groupID, userID string) (GroupMember, error)
	ListByGroup(ctx context.Context, groupID string) ([]GroupMemberWithUser, error)
	Create(ctx context.Context, groupID, userID string, isAdmin bool) (GroupMember, error)
	UpdateRole(ctx context.Context, groupID, userID string, isAdmin bool) (GroupMember, error)
	Remove(ctx context.Context, groupID, userID string) error
	Count(ctx context.Context, groupID string) (int64, error)
	CountAdmins(ctx context.Context, groupID string) (int64, error)
}

var (
	ErrGroupNotFound        = errors.New("grupo não encontrado")
	ErrMemberNotFound       = errors.New("membro não encontrado")
	ErrMemberAlreadyExists  = errors.New("membro já existe")
	ErrCannotRemoveOwner    = errors.New("não é possível remover o proprietário do grupo")
	ErrGroupFull            = errors.New("o grupo atingiu o máximo de 30 membros")
	ErrUserHasPendingDebts  = errors.New("usuário possui débitos pendentes no grupo")
	ErrUserHasPendingTasks  = errors.New("usuário possui tarefas pendentes no grupo")
	ErrNoOtherAdmin         = errors.New("não é possível sair: não há outros administradores no grupo")
)

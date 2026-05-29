package domain

import (
	"context"
	"time"
)

type SplitTag struct {
	ID        string
	GroupID   string
	Name      string
	CreatedBy string
	CreatedAt time.Time
}

type SplitTagMember struct {
	ID         string
	SplitTagID string
	UserID     string
}

type SplitTagMemberWithUser struct {
	SplitTagMember
	UserName string
}

type SplitTagRepository interface {
	GetByID(ctx context.Context, id string) (SplitTag, error)
	ListByGroup(ctx context.Context, groupID string) ([]SplitTag, error)
	Create(ctx context.Context, groupID, name, createdBy string) (SplitTag, error)
	Delete(ctx context.Context, id, groupID string) error
	ListMembers(ctx context.Context, splitTagID string) ([]SplitTagMemberWithUser, error)
	AddMember(ctx context.Context, splitTagID, userID string) (SplitTagMember, error)
	RemoveMember(ctx context.Context, splitTagID, userID string) error
}

package repository

import (
	"context"

	"doheem-backend/internal/adapter/repository/db"
	"doheem-backend/internal/domain"
)

type SplitTagRepo struct {
	q *db.Queries
}

func NewSplitTagRepo(q *db.Queries) *SplitTagRepo {
	return &SplitTagRepo{q: q}
}

func (r *SplitTagRepo) GetByID(ctx context.Context, id string) (domain.SplitTag, error) {
	st, err := r.q.GetSplitTagByID(ctx, uuidFromString(id))
	if err != nil {
		return domain.SplitTag{}, err
	}
	return domainSplitTag(st), nil
}

func (r *SplitTagRepo) ListByGroup(ctx context.Context, groupID string) ([]domain.SplitTag, error) {
	tags, err := r.q.ListSplitTagsByGroup(ctx, uuidFromString(groupID))
	if err != nil {
		return nil, err
	}
	return domainSplitTags(tags), nil
}

func (r *SplitTagRepo) Create(ctx context.Context, groupID, name, createdBy string) (domain.SplitTag, error) {
	st, err := r.q.CreateSplitTag(ctx, db.CreateSplitTagParams{
		GroupID:   uuidFromString(groupID),
		Name:      name,
		CreatedBy: uuidFromString(createdBy),
	})
	if err != nil {
		return domain.SplitTag{}, err
	}
	return domainSplitTag(st), nil
}

func (r *SplitTagRepo) Delete(ctx context.Context, id, groupID string) error {
	return r.q.DeleteSplitTag(ctx, db.DeleteSplitTagParams{
		ID:      uuidFromString(id),
		GroupID: uuidFromString(groupID),
	})
}

func (r *SplitTagRepo) ListMembers(ctx context.Context, splitTagID string) ([]domain.SplitTagMemberWithUser, error) {
	rows, err := r.q.ListSplitTagMembers(ctx, uuidFromString(splitTagID))
	if err != nil {
		return nil, err
	}
	return domainSplitTagMembersWithUser(rows), nil
}

func (r *SplitTagRepo) AddMember(ctx context.Context, splitTagID, userID string) (domain.SplitTagMember, error) {
	stm, err := r.q.AddSplitTagMember(ctx, db.AddSplitTagMemberParams{
		SplitTagID: uuidFromString(splitTagID),
		UserID:     uuidFromString(userID),
	})
	if err != nil {
		return domain.SplitTagMember{}, err
	}
	return domain.SplitTagMember{
		ID:         uuidToString(stm.ID),
		SplitTagID: uuidToString(stm.SplitTagID),
		UserID:     uuidToString(stm.UserID),
	}, nil
}

func (r *SplitTagRepo) RemoveMember(ctx context.Context, splitTagID, userID string) error {
	return r.q.RemoveSplitTagMember(ctx, db.RemoveSplitTagMemberParams{
		SplitTagID: uuidFromString(splitTagID),
		UserID:     uuidFromString(userID),
	})
}

package domain

import (
	"context"
)

type SplitTagService struct {
	repo SplitTagRepository
}

func NewSplitTagService(repo SplitTagRepository) *SplitTagService {
	return &SplitTagService{repo: repo}
}

func (s *SplitTagService) Create(ctx context.Context, groupID, name, createdBy string) (SplitTag, error) {
	return s.repo.Create(ctx, groupID, name, createdBy)
}

func (s *SplitTagService) GetByID(ctx context.Context, id string) (SplitTag, error) {
	tag, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return SplitTag{}, ErrSplitTagNotFound
	}
	return tag, nil
}

func (s *SplitTagService) ListByGroup(ctx context.Context, groupID string) ([]SplitTag, error) {
	return s.repo.ListByGroup(ctx, groupID)
}

func (s *SplitTagService) Delete(ctx context.Context, id, groupID string) error {
	return s.repo.Delete(ctx, id, groupID)
}

func (s *SplitTagService) ListMembers(ctx context.Context, splitTagID string) ([]SplitTagMemberWithUser, error) {
	return s.repo.ListMembers(ctx, splitTagID)
}

func (s *SplitTagService) AddMember(ctx context.Context, splitTagID, userID string) (SplitTagMember, error) {
	return s.repo.AddMember(ctx, splitTagID, userID)
}

func (s *SplitTagService) RemoveMember(ctx context.Context, splitTagID, userID string) error {
	return s.repo.RemoveMember(ctx, splitTagID, userID)
}

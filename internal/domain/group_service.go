package domain

import (
	"context"
	"errors"
)

var (
	ErrGroupNotFound    = errors.New("group not found")
	ErrMemberNotFound   = errors.New("member not found")
	ErrMemberAlreadyExists = errors.New("member already exists")
	ErrCannotRemoveOwner   = errors.New("cannot remove owner from group")
)

type GroupService struct {
	groupRepo       GroupRepository
	groupMemberRepo GroupMemberRepository
}

func NewGroupService(groupRepo GroupRepository, groupMemberRepo GroupMemberRepository) *GroupService {
	return &GroupService{
		groupRepo:       groupRepo,
		groupMemberRepo: groupMemberRepo,
	}
}

func (s *GroupService) Create(ctx context.Context, params CreateGroupParams, creatorID string) (Group, error) {
	group, err := s.groupRepo.Create(ctx, params)
	if err != nil {
		return Group{}, err
	}

	_, err = s.groupMemberRepo.Create(ctx, group.ID, creatorID, "owner")
	if err != nil {
		return Group{}, err
	}

	return group, nil
}

func (s *GroupService) GetByID(ctx context.Context, id string) (Group, error) {
	group, err := s.groupRepo.GetByID(ctx, id)
	if err != nil {
		return Group{}, ErrGroupNotFound
	}
	return group, nil
}

func (s *GroupService) ListByUser(ctx context.Context, userID string) ([]Group, error) {
	return s.groupRepo.ListByUserID(ctx, userID)
}

func (s *GroupService) Update(ctx context.Context, id string, params UpdateGroupParams) (Group, error) {
	return s.groupRepo.Update(ctx, id, params)
}

func (s *GroupService) SoftDelete(ctx context.Context, id string) error {
	return s.groupRepo.SoftDelete(ctx, id)
}

func (s *GroupService) Deactivate(ctx context.Context, id string) error {
	return s.groupRepo.Deactivate(ctx, id)
}

func (s *GroupService) Activate(ctx context.Context, id string) error {
	return s.groupRepo.Activate(ctx, id)
}

func (s *GroupService) AddMember(ctx context.Context, groupID, userID, role string) (GroupMember, error) {
	return s.groupMemberRepo.Create(ctx, groupID, userID, role)
}

func (s *GroupService) RemoveMember(ctx context.Context, groupID, userID string) error {
	member, err := s.groupMemberRepo.Get(ctx, groupID, userID)
	if err != nil {
		return ErrMemberNotFound
	}
	if member.Role == "owner" {
		return ErrCannotRemoveOwner
	}
	return s.groupMemberRepo.Remove(ctx, groupID, userID)
}

func (s *GroupService) UpdateMemberRole(ctx context.Context, groupID, userID, role string) (GroupMember, error) {
	return s.groupMemberRepo.UpdateRole(ctx, groupID, userID, role)
}

func (s *GroupService) ListMembers(ctx context.Context, groupID string) ([]GroupMemberWithUser, error) {
	return s.groupMemberRepo.ListByGroup(ctx, groupID)
}

func (s *GroupService) GetMember(ctx context.Context, groupID, userID string) (GroupMember, error) {
	member, err := s.groupMemberRepo.Get(ctx, groupID, userID)
	if err != nil {
		return GroupMember{}, ErrMemberNotFound
	}
	return member, nil
}

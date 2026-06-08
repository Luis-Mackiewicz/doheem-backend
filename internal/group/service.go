package group

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
)

func generateInviteToken() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 32)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

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

	_, err = s.groupMemberRepo.Create(ctx, group.ID, creatorID, true)
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

func (s *GroupService) AddMember(ctx context.Context, groupID, userID string, isAdmin bool) (GroupMember, error) {
	return s.groupMemberRepo.Create(ctx, groupID, userID, isAdmin)
}

func (s *GroupService) RemoveMember(ctx context.Context, groupID, userID string) error {
	member, err := s.groupMemberRepo.Get(ctx, groupID, userID)
	if err != nil {
		return ErrMemberNotFound
	}
	if member.IsAdmin {
		return ErrCannotRemoveOwner
	}
	return s.groupMemberRepo.Remove(ctx, groupID, userID)
}

func (s *GroupService) UpdateMemberRole(ctx context.Context, groupID, userID string, isAdmin bool) (GroupMember, error) {
	return s.groupMemberRepo.UpdateRole(ctx, groupID, userID, isAdmin)
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

func (s *GroupService) Join(ctx context.Context, groupID, userID string) error {
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return ErrGroupNotFound
	}

	if group.InviteToken == nil {
		return fmt.Errorf("group has no invite token")
	}

	_, err = s.groupMemberRepo.Create(ctx, groupID, userID, false)
	if err != nil {
		return err
	}

	return nil
}

func (s *GroupService) RegenerateInviteToken(ctx context.Context, groupID string) (*string, error) {
	token := generateInviteToken()
	err := s.groupRepo.RegenerateInviteToken(ctx, groupID, token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (s *GroupService) GetInviteToken(ctx context.Context, groupID string) (*string, error) {
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return nil, ErrGroupNotFound
	}
	return group.InviteToken, nil
}

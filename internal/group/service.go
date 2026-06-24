package group

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"doheem-backend/internal/audit_log"
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

type contextKey string

const actorIDKey contextKey = "user_id"

func actorIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(actorIDKey).(string)
	return id
}

type GroupService struct {
	groupRepo       GroupRepository
	groupMemberRepo GroupMemberRepository
	auditLogRepo    audit_log.AuditLogRepository
}

func NewGroupService(groupRepo GroupRepository, groupMemberRepo GroupMemberRepository, auditLogRepo audit_log.AuditLogRepository) *GroupService {
	return &GroupService{
		groupRepo:       groupRepo,
		groupMemberRepo: groupMemberRepo,
		auditLogRepo:    auditLogRepo,
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

	if s.auditLogRepo != nil {
		s.auditLogRepo.Create(ctx, group.ID, creatorID, "group", group.ID, "created", map[string]interface{}{
			"name": params.Name,
		})
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
	group, err := s.groupRepo.Update(ctx, id, params)
	if err != nil {
		return Group{}, err
	}

	if actorID := actorIDFromContext(ctx); actorID != "" && s.auditLogRepo != nil {
		changes := make(map[string]interface{})
		if params.Name != nil {
			changes["name"] = *params.Name
		}
		if params.Description != nil {
			changes["description"] = *params.Description
		}
		if params.MonthlyFee != nil {
			changes["monthly_fee"] = *params.MonthlyFee
		}
		if params.PhotoURL != nil {
			changes["photo_url"] = *params.PhotoURL
		}
		s.auditLogRepo.Create(ctx, id, actorID, "group", id, "updated", changes)
	}

	return group, nil
}

func (s *GroupService) AddMember(ctx context.Context, groupID, userID string, isAdmin bool) (GroupMember, error) {
	count, err := s.groupMemberRepo.Count(ctx, groupID)
	if err != nil {
		return GroupMember{}, err
	}
	if count >= 30 {
		return GroupMember{}, ErrGroupFull
	}
	member, err := s.groupMemberRepo.Create(ctx, groupID, userID, isAdmin)
	if err != nil {
		return GroupMember{}, err
	}

	if actorID := actorIDFromContext(ctx); actorID != "" && s.auditLogRepo != nil {
		s.auditLogRepo.Create(ctx, groupID, actorID, "group_member", member.ID, "member_added", map[string]interface{}{
			"user_id":  userID,
			"is_admin": isAdmin,
		})
	}

	return member, nil
}

func (s *GroupService) RemoveMember(ctx context.Context, groupID, userID string) error {
	member, err := s.groupMemberRepo.Get(ctx, groupID, userID)
	if err != nil {
		return ErrMemberNotFound
	}
	if member.IsAdmin {
		return ErrCannotRemoveOwner
	}
	if err := s.groupMemberRepo.Remove(ctx, groupID, userID); err != nil {
		return err
	}

	if actorID := actorIDFromContext(ctx); actorID != "" && s.auditLogRepo != nil {
		s.auditLogRepo.Create(ctx, groupID, actorID, "group_member", member.ID, "member_removed", map[string]interface{}{
			"user_id": userID,
		})
	}

	return nil
}

func (s *GroupService) UpdateMemberRole(ctx context.Context, groupID, userID string, isAdmin bool) (GroupMember, error) {
	member, err := s.groupMemberRepo.UpdateRole(ctx, groupID, userID, isAdmin)
	if err != nil {
		return GroupMember{}, err
	}

	if actorID := actorIDFromContext(ctx); actorID != "" && s.auditLogRepo != nil {
		s.auditLogRepo.Create(ctx, groupID, actorID, "group_member", member.ID, "member_role_updated", map[string]interface{}{
			"user_id":  userID,
			"is_admin": isAdmin,
		})
	}

	return member, nil
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

	count, err := s.groupMemberRepo.Count(ctx, groupID)
	if err != nil {
		return err
	}
	if count >= 30 {
		return ErrGroupFull
	}

	member, err := s.groupMemberRepo.Create(ctx, groupID, userID, false)
	if err != nil {
		return err
	}

	if s.auditLogRepo != nil {
		s.auditLogRepo.Create(ctx, groupID, userID, "group_member", member.ID, "member_joined", nil)
	}

	return nil
}

func (s *GroupService) RegenerateInviteToken(ctx context.Context, groupID string) (*string, error) {
	token := generateInviteToken()
	err := s.groupRepo.RegenerateInviteToken(ctx, groupID, token)
	if err != nil {
		return nil, err
	}

	if actorID := actorIDFromContext(ctx); actorID != "" && s.auditLogRepo != nil {
		s.auditLogRepo.Create(ctx, groupID, actorID, "group", groupID, "invite_token_regenerated", nil)
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

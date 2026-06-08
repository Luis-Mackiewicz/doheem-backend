package invite

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"doheem-backend/internal/group"
)

type InviteService struct {
	inviteRepo  InviteRepository
	memberRepo  group.GroupMemberRepository
	groupRepo   group.GroupRepository
}

func NewInviteService(inviteRepo InviteRepository, memberRepo group.GroupMemberRepository, groupRepo group.GroupRepository) *InviteService {
	return &InviteService{
		inviteRepo: inviteRepo,
		memberRepo: memberRepo,
		groupRepo:  groupRepo,
	}
}

func generateInviteCode() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	code := make([]byte, 8)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[n.Int64()]
	}
	return string(code), nil
}

func (s *InviteService) Create(ctx context.Context, groupID, createdBy string, expiresAt time.Time) (Invite, error) {
	code, err := generateInviteCode()
	if err != nil {
		return Invite{}, err
	}
	return s.inviteRepo.Create(ctx, groupID, code, createdBy, expiresAt)
}

func (s *InviteService) GetByCode(ctx context.Context, code string) (InviteWithGroup, error) {
	return s.inviteRepo.GetByCode(ctx, code)
}

func (s *InviteService) ListByGroup(ctx context.Context, groupID string) ([]InviteWithCreator, error) {
	return s.inviteRepo.ListByGroup(ctx, groupID)
}

func (s *InviteService) ListPendingByUser(ctx context.Context, userID string) ([]InviteWithGroup, error) {
	return s.inviteRepo.ListPendingByUser(ctx, userID)
}

func (s *InviteService) Use(ctx context.Context, id, userID string) error {
	invite, err := s.inviteRepo.GetByID(ctx, id)
	if err != nil {
		return ErrInviteNotFound
	}
	if invite.UsedAt != nil {
		return ErrInviteAlreadyUsed
	}
	if invite.RevokedAt != nil {
		return ErrInviteRevoked
	}
	if time.Now().After(invite.ExpiresAt) {
		return ErrInviteExpired
	}

	if err := s.inviteRepo.Use(ctx, id); err != nil {
		return err
	}

	_, err = s.memberRepo.Create(ctx, invite.GroupID, userID, "member")
	return err
}

func (s *InviteService) Revoke(ctx context.Context, id string) error {
	return s.inviteRepo.Revoke(ctx, id)
}

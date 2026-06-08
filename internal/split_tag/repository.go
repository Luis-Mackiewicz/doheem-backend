package split_tag

import (
	"context"

	"doheem-backend/internal/db"
)

type SplitTagRepo struct {
	q *db.Queries
}

func NewSplitTagRepo(q *db.Queries) *SplitTagRepo {
	return &SplitTagRepo{q: q}
}

func (r *SplitTagRepo) GetByID(ctx context.Context, id string) (SplitTag, error) {
	st, err := r.q.GetSplitTagByID(ctx, db.UUIDFromString(id))
	if err != nil {
		return SplitTag{}, err
	}
	return toSplitTag(st), nil
}

func (r *SplitTagRepo) ListByGroup(ctx context.Context, groupID string) ([]SplitTag, error) {
	tags, err := r.q.ListSplitTagsByGroup(ctx, db.UUIDFromString(groupID))
	if err != nil {
		return nil, err
	}
	return toSplitTags(tags), nil
}

func (r *SplitTagRepo) Create(ctx context.Context, groupID, name, createdBy string) (SplitTag, error) {
	st, err := r.q.CreateSplitTag(ctx, db.CreateSplitTagParams{
		GroupID:   db.UUIDFromString(groupID),
		Name:      name,
		CreatedBy: db.UUIDFromString(createdBy),
	})
	if err != nil {
		return SplitTag{}, err
	}
	return toSplitTag(st), nil
}

func (r *SplitTagRepo) Delete(ctx context.Context, id, groupID string) error {
	return r.q.DeleteSplitTag(ctx, db.DeleteSplitTagParams{
		ID:      db.UUIDFromString(id),
		GroupID: db.UUIDFromString(groupID),
	})
}

func (r *SplitTagRepo) ListMembers(ctx context.Context, splitTagID string) ([]SplitTagMemberWithUser, error) {
	rows, err := r.q.ListSplitTagMembers(ctx, db.UUIDFromString(splitTagID))
	if err != nil {
		return nil, err
	}
	return toSplitTagMembersWithUser(rows), nil
}

func (r *SplitTagRepo) AddMember(ctx context.Context, splitTagID, userID string) (SplitTagMember, error) {
	stm, err := r.q.AddSplitTagMember(ctx, db.AddSplitTagMemberParams{
		SplitTagID: db.UUIDFromString(splitTagID),
		UserID:     db.UUIDFromString(userID),
	})
	if err != nil {
		return SplitTagMember{}, err
	}
	return SplitTagMember{
		ID:         db.UUIDToString(stm.ID),
		SplitTagID: db.UUIDToString(stm.SplitTagID),
		UserID:     db.UUIDToString(stm.UserID),
	}, nil
}

func (r *SplitTagRepo) RemoveMember(ctx context.Context, splitTagID, userID string) error {
	return r.q.RemoveSplitTagMember(ctx, db.RemoveSplitTagMemberParams{
		SplitTagID: db.UUIDFromString(splitTagID),
		UserID:     db.UUIDFromString(userID),
	})
}

func toSplitTag(st db.SplitTag) SplitTag {
	return SplitTag{
		ID:        db.UUIDToString(st.ID),
		GroupID:   db.UUIDToString(st.GroupID),
		Name:      st.Name,
		CreatedBy: db.UUIDToString(st.CreatedBy),
		CreatedAt: st.CreatedAt.Time,
	}
}

func toSplitTags(tags []db.SplitTag) []SplitTag {
	result := make([]SplitTag, len(tags))
	for i, t := range tags {
		result[i] = toSplitTag(t)
	}
	return result
}

func toSplitTagMemberWithUser(row db.ListSplitTagMembersRow) SplitTagMemberWithUser {
	return SplitTagMemberWithUser{
		SplitTagMember: SplitTagMember{
			ID:         db.UUIDToString(row.ID),
			SplitTagID: db.UUIDToString(row.SplitTagID),
			UserID:     db.UUIDToString(row.UserID),
		},
		UserName: row.UserName,
	}
}

func toSplitTagMembersWithUser(rows []db.ListSplitTagMembersRow) []SplitTagMemberWithUser {
	result := make([]SplitTagMemberWithUser, len(rows))
	for i, r := range rows {
		result[i] = toSplitTagMemberWithUser(r)
	}
	return result
}

package repository

import (
	"encoding/json"
	"strconv"
	"time"

	"doheem-backend/internal/adapter/repository/db"
	"doheem-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgtype"
)

func uuidFromString(s string) pgtype.UUID {
	var u pgtype.UUID
	if s != "" {
		u.Scan(s)
	}
	return u
}

func uuidToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	v, _ := u.Value()
	return v.(string)
}

func numericFromFloat64(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	n.Scan(f)
	return n
}

func numericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	v, _ := n.Value()
	s, ok := v.(string)
	if !ok {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func textFromStringPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func textToStringPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func uuidFromStringPtrOrZero(s *string) pgtype.UUID {
	if s == nil {
		return pgtype.UUID{}
	}
	return uuidFromString(*s)
}

func uuidFromStringPtr(s *string) pgtype.UUID {
	if s == nil {
		return pgtype.UUID{}
	}
	return uuidFromString(*s)
}

func timeFromTimestamptz(t pgtype.Timestamptz) time.Time {
	return t.Time
}

func timestamptzFromTime(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func timestamptzToTimePtr(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func dateFromTime(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}

func dateToTimePtr(t pgtype.Date) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func int2FromInt16Ptr(i *int16) pgtype.Int2 {
	if i == nil {
		return pgtype.Int2{}
	}
	return pgtype.Int2{Int16: *i, Valid: true}
}

func int2ToInt16Ptr(i pgtype.Int2) *int16 {
	if !i.Valid {
		return nil
	}
	return &i.Int16
}

func jsonToMap(b []byte) map[string]interface{} {
	if b == nil {
		return nil
	}
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	return m
}

func mapToJSON(m map[string]interface{}) []byte {
	if m == nil {
		return nil
	}
	b, _ := json.Marshal(m)
	return b
}

// ── User ──

func domainUser(u db.User) domain.User {
	return domain.User{
		ID:           uuidToString(u.ID),
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		AvatarURL:    textToStringPtr(u.AvatarUrl),
		CreatedAt:    u.CreatedAt.Time,
		UpdatedAt:    u.UpdatedAt.Time,
	}
}

func domainUsers(users []db.User) []domain.User {
	result := make([]domain.User, len(users))
	for i, u := range users {
		result[i] = domainUser(u)
	}
	return result
}

// ── Group ──

func domainGroup(g db.Group) domain.Group {
	return domain.Group{
		ID:            uuidToString(g.ID),
		Name:          g.Name,
		Currency:      g.Currency,
		IsActive:      g.IsActive,
		InactiveSince: timestamptzToTimePtr(g.InactiveSince),
		CreatedAt:     g.CreatedAt.Time,
		UpdatedAt:     g.UpdatedAt.Time,
		DeletedAt:     timestamptzToTimePtr(g.DeletedAt),
	}
}

func domainGroups(groups []db.Group) []domain.Group {
	result := make([]domain.Group, len(groups))
	for i, g := range groups {
		result[i] = domainGroup(g)
	}
	return result
}

// ── GroupMember ──

func domainGroupMember(gm db.GroupMember) domain.GroupMember {
	return domain.GroupMember{
		ID:       uuidToString(gm.ID),
		GroupID:  uuidToString(gm.GroupID),
		UserID:   uuidToString(gm.UserID),
		Role:     gm.Role,
		JoinedAt: gm.JoinedAt.Time,
		LeftAt:   timestamptzToTimePtr(gm.LeftAt),
		IsActive: gm.IsActive,
	}
}

func domainGroupMemberWithUser(row db.ListGroupMembersRow) domain.GroupMemberWithUser {
	return domain.GroupMemberWithUser{
		GroupMember: domain.GroupMember{
			ID:       uuidToString(row.ID),
			GroupID:  uuidToString(row.GroupID),
			UserID:   uuidToString(row.UserID),
			Role:     row.Role,
			JoinedAt: row.JoinedAt.Time,
			LeftAt:   timestamptzToTimePtr(row.LeftAt),
			IsActive: row.IsActive,
		},
		UserName:  row.Name,
		UserEmail: row.Email,
		AvatarURL: textToStringPtr(row.AvatarUrl),
	}
}

// ── Expense ──

func domainExpense(e db.Expense) domain.Expense {
	return domain.Expense{
		ID:               uuidToString(e.ID),
		GroupID:          uuidToString(e.GroupID),
		CreatedBy:        uuidToString(e.CreatedBy),
		Description:      e.Description,
		TotalAmount:      numericToFloat64(e.TotalAmount),
		ExpenseDate:      e.ExpenseDate.Time,
		DueDate:          dateToTimePtr(e.DueDate),
		CategoryID:       uuidToStringPtr(e.CategoryID),
		SplitType:        e.SplitType,
		IsInstallment:    e.IsInstallment,
		InstallmentCount: int2ToInt16Ptr(e.InstallmentCount),
		CreatedAt:        e.CreatedAt.Time,
		UpdatedAt:        e.UpdatedAt.Time,
	}
}

func domainExpenses(expenses []db.Expense) []domain.Expense {
	result := make([]domain.Expense, len(expenses))
	for i, e := range expenses {
		result[i] = domainExpense(e)
	}
	return result
}

func uuidToStringPtr(u pgtype.UUID) *string {
	if !u.Valid {
		return nil
	}
	s := uuidToString(u)
	return &s
}

// ── ExpenseCategory ──

func domainExpenseCategory(ec db.ExpenseCategory) domain.ExpenseCategory {
	return domain.ExpenseCategory{
		ID:        uuidToString(ec.ID),
		GroupID:   uuidToString(ec.GroupID),
		Name:      ec.Name,
		CreatedAt: ec.CreatedAt.Time,
	}
}

func domainExpenseCategories(categories []db.ExpenseCategory) []domain.ExpenseCategory {
	result := make([]domain.ExpenseCategory, len(categories))
	for i, c := range categories {
		result[i] = domainExpenseCategory(c)
	}
	return result
}

// ── ExpenseSplit ──

func domainExpenseSplit(es db.ExpenseSplit) domain.ExpenseSplit {
	return domain.ExpenseSplit{
		ID:        uuidToString(es.ID),
		ExpenseID: uuidToString(es.ExpenseID),
		UserID:    uuidToString(es.UserID),
		Amount:    numericToFloat64(es.Amount),
		IsPaid:    es.IsPaid,
		PaidAt:    timestamptzToTimePtr(es.PaidAt),
		CreatedAt: es.CreatedAt.Time,
	}
}

func domainExpenseSplitWithUser(row db.ListExpenseSplitsByExpenseRow) domain.ExpenseSplitWithUser {
	return domain.ExpenseSplitWithUser{
		ExpenseSplit: domain.ExpenseSplit{
			ID:        uuidToString(row.ID),
			ExpenseID: uuidToString(row.ExpenseID),
			UserID:    uuidToString(row.UserID),
			Amount:    numericToFloat64(row.Amount),
			IsPaid:    row.IsPaid,
			PaidAt:    timestamptzToTimePtr(row.PaidAt),
			CreatedAt: row.CreatedAt.Time,
		},
		UserName:  row.UserName,
		UserEmail: row.UserEmail,
	}
}

func domainExpenseSplitsWithUser(rows []db.ListExpenseSplitsByExpenseRow) []domain.ExpenseSplitWithUser {
	result := make([]domain.ExpenseSplitWithUser, len(rows))
	for i, r := range rows {
		result[i] = domainExpenseSplitWithUser(r)
	}
	return result
}

// ── Installment ──

func domainInstallment(i db.Installment) domain.Installment {
	return domain.Installment{
		ID:                uuidToString(i.ID),
		ExpenseID:         uuidToString(i.ExpenseID),
		InstallmentNumber: i.InstallmentNumber,
		Amount:            numericToFloat64(i.Amount),
		DueDate:           i.DueDate.Time,
		IsPaid:            i.IsPaid,
		PaidAt:            timestamptzToTimePtr(i.PaidAt),
		CreatedAt:         i.CreatedAt.Time,
	}
}

func domainInstallments(installments []db.Installment) []domain.Installment {
	result := make([]domain.Installment, len(installments))
	for i, inst := range installments {
		result[i] = domainInstallment(inst)
	}
	return result
}

// ── Payment ──

func domainPayment(p db.Payment) domain.Payment {
	return domain.Payment{
		ID:          uuidToString(p.ID),
		GroupID:     uuidToString(p.GroupID),
		PayerID:     uuidToString(p.PayerID),
		ReceiverID:  uuidToString(p.ReceiverID),
		Amount:      numericToFloat64(p.Amount),
		PaymentDate: p.PaymentDate.Time,
		Status:      p.Status,
		Notes:       textToStringPtr(p.Notes),
		CreatedAt:   p.CreatedAt.Time,
		ConfirmedAt: timestamptzToTimePtr(p.ConfirmedAt),
		CancelledAt: timestamptzToTimePtr(p.CancelledAt),
	}
}

func domainPayments(payments []db.Payment) []domain.Payment {
	result := make([]domain.Payment, len(payments))
	for i, p := range payments {
		result[i] = domainPayment(p)
	}
	return result
}

func domainPaymentWithUsers(row db.ListPaymentsByGroupRow) domain.PaymentWithUsers {
	return domain.PaymentWithUsers{
		Payment: domain.Payment{
			ID:          uuidToString(row.ID),
			GroupID:     uuidToString(row.GroupID),
			PayerID:     uuidToString(row.PayerID),
			ReceiverID:  uuidToString(row.ReceiverID),
			Amount:      numericToFloat64(row.Amount),
			PaymentDate: row.PaymentDate.Time,
			Status:      row.Status,
			Notes:       textToStringPtr(row.Notes),
			CreatedAt:   row.CreatedAt.Time,
			ConfirmedAt: timestamptzToTimePtr(row.ConfirmedAt),
			CancelledAt: timestamptzToTimePtr(row.CancelledAt),
		},
		PayerName:   row.PayerName,
		ReceiverName: row.ReceiverName,
	}
}

func domainPaymentsWithUsers(rows []db.ListPaymentsByGroupRow) []domain.PaymentWithUsers {
	result := make([]domain.PaymentWithUsers, len(rows))
	for i, r := range rows {
		result[i] = domainPaymentWithUsers(r)
	}
	return result
}

func domainPaymentWithGroup(row db.ListPaymentsByUserRow) domain.PaymentWithGroup {
	return domain.PaymentWithGroup{
		Payment: domain.Payment{
			ID:          uuidToString(row.ID),
			GroupID:     uuidToString(row.GroupID),
			PayerID:     uuidToString(row.PayerID),
			ReceiverID:  uuidToString(row.ReceiverID),
			Amount:      numericToFloat64(row.Amount),
			PaymentDate: row.PaymentDate.Time,
			Status:      row.Status,
			Notes:       textToStringPtr(row.Notes),
			CreatedAt:   row.CreatedAt.Time,
			ConfirmedAt: timestamptzToTimePtr(row.ConfirmedAt),
			CancelledAt: timestamptzToTimePtr(row.CancelledAt),
		},
		GroupName: row.GroupName,
	}
}

func domainPaymentsWithGroup(rows []db.ListPaymentsByUserRow) []domain.PaymentWithGroup {
	result := make([]domain.PaymentWithGroup, len(rows))
	for i, r := range rows {
		result[i] = domainPaymentWithGroup(r)
	}
	return result
}

// ── PaymentAttachment ──

func domainPaymentAttachment(pa db.PaymentAttachment) domain.PaymentAttachment {
	return domain.PaymentAttachment{
		ID:         uuidToString(pa.ID),
		PaymentID:  uuidToString(pa.PaymentID),
		FilePath:   pa.FilePath,
		FileType:   pa.FileType,
		FileSize:   pa.FileSize,
		UploadedAt: pa.UploadedAt.Time,
	}
}

func domainPaymentAttachments(attachments []db.PaymentAttachment) []domain.PaymentAttachment {
	result := make([]domain.PaymentAttachment, len(attachments))
	for i, a := range attachments {
		result[i] = domainPaymentAttachment(a)
	}
	return result
}

// ── Task ──

func domainTask(t db.Task) domain.Task {
	return domain.Task{
		ID:               uuidToString(t.ID),
		GroupID:          uuidToString(t.GroupID),
		Title:            t.Title,
		Description:      textToStringPtr(t.Description),
		AssignedTo:       uuidToStringPtr(t.AssignedTo),
		Category:         textToStringPtr(t.Category),
		IsRecurring:      t.IsRecurring,
		RecurringPeriod:  textToStringPtr(t.RecurringPeriod),
		RecurringEndedAt: timestamptzToTimePtr(t.RecurringEndedAt),
		CreatedBy:        uuidToString(t.CreatedBy),
		CreatedAt:        t.CreatedAt.Time,
		UpdatedAt:        t.UpdatedAt.Time,
	}
}

func domainTasks(tasks []db.Task) []domain.Task {
	result := make([]domain.Task, len(tasks))
	for i, t := range tasks {
		result[i] = domainTask(t)
	}
	return result
}

// ── TaskOccurrence ──

func domainTaskOccurrence(to db.TaskOccurrence) domain.TaskOccurrence {
	return domain.TaskOccurrence{
		ID:          uuidToString(to.ID),
		TaskID:      uuidToString(to.TaskID),
		DueDate:     to.DueDate.Time,
		Status:      to.Status,
		CompletedBy: uuidToStringPtr(to.CompletedBy),
		CompletedAt: timestamptzToTimePtr(to.CompletedAt),
		DiscardedAt: timestamptzToTimePtr(to.DiscardedAt),
		CreatedAt:   to.CreatedAt.Time,
	}
}

func domainTaskOccurrenceWithTask(row db.ListPendingTaskOccurrencesByUserRow) domain.TaskOccurrenceWithTask {
	return domain.TaskOccurrenceWithTask{
		TaskOccurrence: domain.TaskOccurrence{
			ID:          uuidToString(row.ID),
			TaskID:      uuidToString(row.TaskID),
			DueDate:     row.DueDate.Time,
			Status:      row.Status,
			CompletedBy: uuidToStringPtr(row.CompletedBy),
			CompletedAt: timestamptzToTimePtr(row.CompletedAt),
			DiscardedAt: timestamptzToTimePtr(row.DiscardedAt),
			CreatedAt:   row.CreatedAt.Time,
		},
		TaskTitle: row.TaskTitle,
		GroupID:   uuidToString(row.GroupID),
	}
}

func domainTaskOccurrencesWithTask(rows []db.ListPendingTaskOccurrencesByUserRow) []domain.TaskOccurrenceWithTask {
	result := make([]domain.TaskOccurrenceWithTask, len(rows))
	for i, r := range rows {
		result[i] = domainTaskOccurrenceWithTask(r)
	}
	return result
}

// ── Invite ──

func domainInvite(i db.Invite) domain.Invite {
	return domain.Invite{
		ID:        uuidToString(i.ID),
		GroupID:   uuidToString(i.GroupID),
		Code:      i.Code,
		CreatedBy: uuidToString(i.CreatedBy),
		ExpiresAt: i.ExpiresAt.Time,
		UsedAt:    timestamptzToTimePtr(i.UsedAt),
		RevokedAt: timestamptzToTimePtr(i.RevokedAt),
		CreatedAt: i.CreatedAt.Time,
	}
}

func domainInviteWithGroup(row db.GetInviteByCodeRow) domain.InviteWithGroup {
	return domain.InviteWithGroup{
		Invite: domain.Invite{
			ID:        uuidToString(row.ID),
			GroupID:   uuidToString(row.GroupID),
			Code:      row.Code,
			CreatedBy: uuidToString(row.CreatedBy),
			ExpiresAt: row.ExpiresAt.Time,
			UsedAt:    timestamptzToTimePtr(row.UsedAt),
			RevokedAt: timestamptzToTimePtr(row.RevokedAt),
			CreatedAt: row.CreatedAt.Time,
		},
		GroupName: row.GroupName,
	}
}

// ── SplitTag ──

func domainSplitTag(st db.SplitTag) domain.SplitTag {
	return domain.SplitTag{
		ID:        uuidToString(st.ID),
		GroupID:   uuidToString(st.GroupID),
		Name:      st.Name,
		CreatedBy: uuidToString(st.CreatedBy),
		CreatedAt: st.CreatedAt.Time,
	}
}

func domainSplitTags(tags []db.SplitTag) []domain.SplitTag {
	result := make([]domain.SplitTag, len(tags))
	for i, t := range tags {
		result[i] = domainSplitTag(t)
	}
	return result
}

func domainSplitTagMemberWithUser(row db.ListSplitTagMembersRow) domain.SplitTagMemberWithUser {
	return domain.SplitTagMemberWithUser{
		SplitTagMember: domain.SplitTagMember{
			ID:         uuidToString(row.ID),
			SplitTagID: uuidToString(row.SplitTagID),
			UserID:     uuidToString(row.UserID),
		},
		UserName: row.UserName,
	}
}

func domainSplitTagMembersWithUser(rows []db.ListSplitTagMembersRow) []domain.SplitTagMemberWithUser {
	result := make([]domain.SplitTagMemberWithUser, len(rows))
	for i, r := range rows {
		result[i] = domainSplitTagMemberWithUser(r)
	}
	return result
}

// ── Notification ──

func domainNotification(n db.Notification) domain.Notification {
	return domain.Notification{
		ID:        uuidToString(n.ID),
		UserID:    uuidToString(n.UserID),
		GroupID:   uuidToStringPtr(n.GroupID),
		Type:      n.Type,
		Title:     n.Title,
		Message:   n.Message,
		IsRead:    n.IsRead,
		ReadAt:    timestamptzToTimePtr(n.ReadAt),
		CreatedAt: n.CreatedAt.Time,
	}
}

func domainNotifications(notifications []db.Notification) []domain.Notification {
	result := make([]domain.Notification, len(notifications))
	for i, n := range notifications {
		result[i] = domainNotification(n)
	}
	return result
}

// ── AuditLog ──

func domainAuditLog(al db.AuditLog) domain.AuditLog {
	return domain.AuditLog{
		ID:         uuidToString(al.ID),
		GroupID:    uuidToStringPtr(al.GroupID),
		UserID:     uuidToStringPtr(al.UserID),
		EntityType: al.EntityType,
		EntityID:   uuidToString(al.EntityID),
		Action:     al.Action,
		Changes:    jsonToMap(al.Changes),
		CreatedAt:  al.CreatedAt.Time,
	}
}

func domainAuditLogs(logs []db.AuditLog) []domain.AuditLog {
	result := make([]domain.AuditLog, len(logs))
	for i, l := range logs {
		result[i] = domainAuditLog(l)
	}
	return result
}

// ── UserBalance ──

func domainUserBalance(row db.GetUserBalanceInGroupRow) domain.UserBalance {
	return domain.UserBalance{
		TotalOwed: numericToFloat64(row.TotalOwed),
		TotalPaid: numericToFloat64(row.TotalPaid),
	}
}

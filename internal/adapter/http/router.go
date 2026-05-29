package http

import (
	"net/http"

	"doheem-backend/internal/domain"
)

type Router struct {
	user         *UserHandler
	group        *GroupHandler
	expense      *ExpenseHandler
	payment      *PaymentHandler
	task         *TaskHandler
	invite       *InviteHandler
	notification *NotificationHandler
	splitTag     *SplitTagHandler
}

func NewRouter(
	userSvc *domain.UserService,
	groupSvc *domain.GroupService,
	expenseSvc *domain.ExpenseService,
	paymentSvc *domain.PaymentService,
	taskSvc *domain.TaskService,
	inviteSvc *domain.InviteService,
	notificationSvc *domain.NotificationService,
	splitTagSvc *domain.SplitTagService,
) *Router {
	return &Router{
		user:         NewUserHandler(userSvc),
		group:        NewGroupHandler(groupSvc),
		expense:      NewExpenseHandler(expenseSvc),
		payment:      NewPaymentHandler(paymentSvc),
		task:         NewTaskHandler(taskSvc),
		invite:       NewInviteHandler(inviteSvc),
		notification: NewNotificationHandler(notificationSvc),
		splitTag:     NewSplitTagHandler(splitTagSvc),
	}
}

func (rt *Router) Handler() http.Handler {
	mux := http.NewServeMux()

	health := &HealthHandler{}
	mux.HandleFunc("GET /", health.HealthCheck)

	mux.HandleFunc("POST /api/auth/register", rt.user.Register)
	mux.HandleFunc("POST /api/auth/login", rt.user.Login)

	mux.Handle("GET /api/users/me", rt.auth(rt.user.GetProfile))
	mux.Handle("PUT /api/users/me", rt.auth(rt.user.UpdateProfile))
	mux.Handle("PUT /api/users/me/password", rt.auth(rt.user.ChangePassword))

	mux.Handle("POST /api/groups", rt.auth(rt.group.Create))
	mux.Handle("GET /api/groups", rt.auth(rt.group.List))
	mux.Handle("GET /api/groups/{id}", rt.auth(rt.group.GetByID))
	mux.Handle("PUT /api/groups/{id}", rt.auth(rt.group.Update))
	mux.Handle("DELETE /api/groups/{id}", rt.auth(rt.group.SoftDelete))
	mux.Handle("PATCH /api/groups/{id}/deactivate", rt.auth(rt.group.Deactivate))
	mux.Handle("PATCH /api/groups/{id}/activate", rt.auth(rt.group.Activate))
	mux.Handle("GET /api/groups/{id}/members", rt.auth(rt.group.ListMembers))
	mux.Handle("POST /api/groups/{id}/members", rt.auth(rt.group.AddMember))
	mux.Handle("PUT /api/groups/{id}/members/{userId}", rt.auth(rt.group.UpdateMemberRole))
	mux.Handle("DELETE /api/groups/{id}/members/{userId}", rt.auth(rt.group.RemoveMember))

	mux.Handle("POST /api/groups/{groupId}/expenses", rt.auth(rt.expense.Create))
	mux.Handle("GET /api/groups/{groupId}/expenses", rt.auth(rt.expense.ListByGroup))
	mux.Handle("GET /api/expenses/{id}", rt.auth(rt.expense.GetByID))
	mux.Handle("PUT /api/expenses/{id}", rt.auth(rt.expense.Update))
	mux.Handle("DELETE /api/expenses/{id}", rt.auth(rt.expense.Delete))
	mux.Handle("GET /api/expenses/{id}/splits", rt.auth(rt.expense.ListSplits))
	mux.Handle("PATCH /api/expenses/splits/{id}/pay", rt.auth(rt.expense.MarkSplitAsPaid))
	mux.Handle("GET /api/expenses/{id}/installments", rt.auth(rt.expense.ListInstallments))
	mux.Handle("PATCH /api/expenses/installments/{id}/pay", rt.auth(rt.expense.MarkInstallmentAsPaid))

	mux.Handle("POST /api/groups/{groupId}/categories", rt.auth(rt.expense.CreateCategory))
	mux.Handle("GET /api/groups/{groupId}/categories", rt.auth(rt.expense.ListCategories))
	mux.Handle("PUT /api/categories/{id}", rt.auth(rt.expense.UpdateCategory))
	mux.Handle("DELETE /api/categories/{id}", rt.auth(rt.expense.DeleteCategory))

	mux.Handle("POST /api/groups/{groupId}/payments", rt.auth(rt.payment.Create))
	mux.Handle("GET /api/groups/{groupId}/payments", rt.auth(rt.payment.ListByGroup))
	mux.Handle("GET /api/payments/{id}", rt.auth(rt.payment.GetByID))
	mux.Handle("PATCH /api/payments/{id}/confirm", rt.auth(rt.payment.Confirm))
	mux.Handle("PATCH /api/payments/{id}/cancel", rt.auth(rt.payment.Cancel))
	mux.Handle("DELETE /api/payments/{id}", rt.auth(rt.payment.Delete))

	mux.Handle("POST /api/groups/{groupId}/tasks", rt.auth(rt.task.Create))
	mux.Handle("GET /api/groups/{groupId}/tasks", rt.auth(rt.task.ListByGroup))
	mux.Handle("GET /api/tasks/{id}", rt.auth(rt.task.GetByID))
	mux.Handle("PUT /api/tasks/{id}", rt.auth(rt.task.Update))
	mux.Handle("DELETE /api/tasks/{id}", rt.auth(rt.task.Delete))
	mux.Handle("GET /api/tasks/{id}/occurrences", rt.auth(rt.task.ListOccurrences))
	mux.Handle("POST /api/tasks/{taskId}/occurrences", rt.auth(rt.task.CreateOccurrence))
	mux.Handle("PATCH /api/tasks/occurrences/{id}/complete", rt.auth(rt.task.CompleteOccurrence))
	mux.Handle("PATCH /api/tasks/occurrences/{id}/discard", rt.auth(rt.task.DiscardOccurrence))

	mux.Handle("POST /api/groups/{groupId}/invites", rt.auth(rt.invite.Create))
	mux.Handle("GET /api/groups/{groupId}/invites", rt.auth(rt.invite.ListByGroup))
	mux.Handle("GET /api/invites/pending", rt.auth(rt.invite.ListPending))
	mux.Handle("POST /api/invites/{id}/accept", rt.auth(rt.invite.Accept))
	mux.Handle("PATCH /api/invites/{id}/revoke", rt.auth(rt.invite.Revoke))

	mux.Handle("GET /api/notifications", rt.auth(rt.notification.List))
	mux.Handle("GET /api/notifications/unread", rt.auth(rt.notification.ListUnread))
	mux.Handle("PATCH /api/notifications/{id}/read", rt.auth(rt.notification.MarkAsRead))
	mux.Handle("PATCH /api/notifications/read-all", rt.auth(rt.notification.MarkAllAsRead))
	mux.Handle("DELETE /api/notifications/{id}", rt.auth(rt.notification.Delete))

	mux.Handle("POST /api/groups/{groupId}/split-tags", rt.auth(rt.splitTag.Create))
	mux.Handle("GET /api/groups/{groupId}/split-tags", rt.auth(rt.splitTag.ListByGroup))
	mux.Handle("DELETE /api/split-tags/{id}", rt.auth(rt.splitTag.Delete))
	mux.Handle("GET /api/split-tags/{id}/members", rt.auth(rt.splitTag.ListMembers))
	mux.Handle("POST /api/split-tags/{id}/members", rt.auth(rt.splitTag.AddMember))
	mux.Handle("DELETE /api/split-tags/{id}/members/{userId}", rt.auth(rt.splitTag.RemoveMember))

	var h http.Handler = mux
	h = LoggingMiddleware(h)
	h = CORSMiddleware(h)
	h = RecoveryMiddleware(h)

	return h
}

func (rt *Router) auth(handler http.HandlerFunc) http.Handler {
	return AuthMiddleware(http.HandlerFunc(handler))
}

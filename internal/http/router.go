package http

import (
	"net/http"
	"time"

	"doheem-backend/internal/expense"
	"doheem-backend/internal/group"
	"doheem-backend/internal/notification"
	"doheem-backend/internal/task"
	"doheem-backend/internal/user"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Router struct {
	jwt          *JWTService
	user         *UserHandler
	group        *GroupHandler
	expense      *ExpenseHandler
	task         *TaskHandler
	notification *NotificationHandler
	rdb          *redis.Client
	pool         *pgxpool.Pool
}

func NewRouter(
	jwt *JWTService,
	userSvc *user.UserService,
	groupSvc *group.GroupService,
	expenseSvc *expense.ExpenseService,
	taskSvc *task.TaskService,
	notificationSvc *notification.NotificationService,
	rdb *redis.Client,
	pool *pgxpool.Pool,
) *Router {
	return &Router{
		jwt:          jwt,
		user:         NewUserHandler(userSvc, jwt),
		group:        NewGroupHandler(groupSvc, rdb, expenseSvc, taskSvc),
		expense:      NewExpenseHandler(expenseSvc),
		task:         NewTaskHandler(taskSvc),
		notification: NewNotificationHandler(notificationSvc),
		rdb:          rdb,
		pool:         pool,
	}
}

func (rt *Router) Handler() http.Handler {
	mux := http.NewServeMux()

	healthHandler := NewHealthHandler(rt.pool, rt.rdb)
	mux.HandleFunc("GET /", healthHandler.HealthCheck)

	authLimiter := RateLimitMiddleware(NewRateLimiter(rt.rdb, 10, time.Minute))
	mux.Handle("POST /api/auth/register", authLimiter(http.HandlerFunc(rt.user.Register)))
	mux.Handle("POST /api/auth/login", authLimiter(http.HandlerFunc(rt.user.Login)))
	mux.Handle("POST /api/auth/refresh", authLimiter(http.HandlerFunc(rt.user.Refresh)))
	mux.Handle("POST /api/auth/logout", rt.auth(rt.user.Logout))

	mux.Handle("GET /api/users/me", rt.auth(rt.user.GetProfile))
	mux.Handle("PUT /api/users/me", rt.auth(rt.user.UpdateProfile))
	mux.Handle("PUT /api/users/me/password", rt.auth(rt.user.ChangePassword))

	mux.Handle("POST /api/groups", rt.auth(rt.group.Create))
	mux.Handle("GET /api/groups", rt.auth(rt.group.List))
	mux.Handle("GET /api/groups/{id}", rt.auth(rt.group.GetByID))
	mux.Handle("PUT /api/groups/{id}", rt.auth(rt.group.Update))
	// DELETE /api/groups/{id} — not implemented yet
	mux.Handle("GET /api/groups/{id}/members", rt.auth(rt.group.ListMembers))
	mux.Handle("POST /api/groups/{id}/members", rt.auth(rt.group.AddMember))
	mux.Handle("PUT /api/groups/{id}/members/{userId}", rt.auth(rt.group.UpdateMemberRole))
	mux.Handle("DELETE /api/groups/{id}/members/{userId}", rt.auth(rt.group.RemoveMember))

	mux.Handle("POST /api/groups/{id}/leave", rt.auth(rt.group.Leave))
	mux.Handle("POST /api/groups/{id}/join", rt.auth(rt.group.Join))
	mux.Handle("POST /api/groups/{id}/regenerate-invite", rt.auth(rt.group.RegenerateInvite))
	mux.Handle("GET /api/groups/{id}/invite-token", rt.auth(rt.group.GetInviteToken))

	mux.Handle("POST /api/groups/{groupId}/expenses", rt.auth(rt.expense.Create))
	mux.Handle("GET /api/groups/{groupId}/expenses", rt.auth(rt.expense.ListByGroup))
	mux.Handle("GET /api/expenses/{id}", rt.auth(rt.expense.GetByID))
	mux.Handle("PUT /api/expenses/{id}", rt.auth(rt.expense.Update))
	mux.Handle("DELETE /api/expenses/{id}", rt.auth(rt.expense.Delete))
	mux.Handle("GET /api/expenses/{id}/splits", rt.auth(rt.expense.ListSplits))
	mux.Handle("PATCH /api/expenses/splits/{id}/pay", rt.auth(rt.expense.MarkSplitAsPaid))
	mux.Handle("GET /api/expenses/{id}/installments", rt.auth(rt.expense.ListByParent))

	mux.Handle("POST /api/categories", rt.auth(rt.expense.CreateCategory))
	mux.Handle("GET /api/categories", rt.auth(rt.expense.ListCategories))
	mux.Handle("PUT /api/categories/{id}", rt.auth(rt.expense.UpdateCategory))
	mux.Handle("DELETE /api/categories/{id}", rt.auth(rt.expense.DeleteCategory))

	mux.Handle("POST /api/groups/{groupId}/tasks", rt.auth(rt.task.Create))
	mux.Handle("GET /api/groups/{groupId}/tasks", rt.auth(rt.task.ListByGroup))
	mux.Handle("GET /api/tasks/{id}", rt.auth(rt.task.GetByID))
	mux.Handle("PUT /api/tasks/{id}", rt.auth(rt.task.Update))
	mux.Handle("DELETE /api/tasks/{id}", rt.auth(rt.task.Delete))
	mux.Handle("GET /api/tasks/{id}/occurrences", rt.auth(rt.task.ListOccurrences))
	mux.Handle("POST /api/tasks/{taskId}/occurrences", rt.auth(rt.task.CreateOccurrence))
	mux.Handle("PATCH /api/tasks/occurrences/{id}/complete", rt.auth(rt.task.CompleteOccurrence))
	mux.Handle("PATCH /api/tasks/occurrences/{id}/discard", rt.auth(rt.task.DiscardOccurrence))

	mux.Handle("POST /api/notifications", rt.auth(rt.notification.Create))
	mux.Handle("GET /api/notifications", rt.auth(rt.notification.List))
	mux.Handle("GET /api/notifications/unread", rt.auth(rt.notification.ListUnread))
	mux.Handle("PATCH /api/notifications/{id}/read", rt.auth(rt.notification.MarkAsRead))
	mux.Handle("PATCH /api/notifications/read-all", rt.auth(rt.notification.MarkAllAsRead))
	mux.Handle("DELETE /api/notifications", rt.auth(rt.notification.DeleteAll))
	mux.Handle("DELETE /api/notifications/{id}", rt.auth(rt.notification.Delete))

	var h http.Handler = mux
	h = RateLimitMiddleware(NewRateLimiter(rt.rdb, 100, time.Minute))(h)
	h = LoggingMiddleware(h)
	h = CORSMiddleware(h)
	h = RecoveryMiddleware(h)

	return h
}

func (rt *Router) auth(handler http.HandlerFunc) http.Handler {
	return rt.jwt.AuthMiddleware(http.HandlerFunc(handler))
}

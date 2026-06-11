package http

import (
	"net/http"
	"time"

	"doheem-backend/internal/expense"
)

type ExpenseHandler struct {
	svc *expense.ExpenseService
}

func NewExpenseHandler(svc *expense.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{svc: svc}
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("groupId")

	var req struct {
		Description    string    `json:"description"      validate:"required"`
		Amount         float64   `json:"amount"            validate:"required,gt=0"`
		CategoryID     string    `json:"category_id"       validate:"required"`
		CompetenceDate string    `json:"competence_date"   validate:"required"`
		DueDate        string    `json:"due_date"          validate:"required"`
		PaidBy         string    `json:"paid_by"           validate:"required"`
		SplitMode      string    `json:"split_mode"        validate:"required,oneof=equal custom"`
		Installments   int32     `json:"installments"`
		FirstDueDate   *string   `json:"first_due_date,omitempty"`
		IsFixed        bool      `json:"is_fixed"`
		Splits         []struct {
			UserID string  `json:"user_id" validate:"required"`
			Amount float64 `json:"amount"   validate:"required,gt=0"`
		} `json:"splits,omitempty"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}

	competenceDate, err := time.Parse("2006-01-02", req.CompetenceDate)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid competence_date, use YYYY-MM-DD")
		return
	}
	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid due_date, use YYYY-MM-DD")
		return
	}
	var firstDueDate *time.Time
	if req.FirstDueDate != nil {
		t, err := time.Parse("2006-01-02", *req.FirstDueDate)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid first_due_date, use YYYY-MM-DD")
			return
		}
		firstDueDate = &t
	}

	splits := make([]expense.CreateExpenseSplitParams, len(req.Splits))
	for i, s := range req.Splits {
		splits[i] = expense.CreateExpenseSplitParams{
			UserID: s.UserID,
			Amount: s.Amount,
		}
	}

	e, err := h.svc.Create(r.Context(), expense.CreateExpenseWithSplitsParams{
		Expense: expense.CreateExpenseParams{
			GroupID:        groupID,
			Description:    req.Description,
			Amount:         req.Amount,
			CategoryID:     req.CategoryID,
			CompetenceDate: competenceDate,
			DueDate:        dueDate,
			PaidBy:         req.PaidBy,
			SplitMode:      req.SplitMode,
			Installments:   req.Installments,
			FirstDueDate:   firstDueDate,
			IsFixed:        req.IsFixed,
		},
		Splits: splits,
	})
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusCreated, toExpenseResponse(e))
}

func (h *ExpenseHandler) ListByGroup(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("groupId")
	limit, offset := parsePagination(r)

	expenses, err := h.svc.ListByGroup(r.Context(), groupID, int32(limit), int32(offset))
	if err != nil {
		handleError(w, r, err)
		return
	}

	total, err := h.svc.CountByGroup(r.Context(), groupID)
	if err != nil {
		handleError(w, r, err)
		return
	}

	respondJSON(w, http.StatusOK, paginatedResponse{Data: toExpenseResponses(expenses), Total: total})
}

func (h *ExpenseHandler) ListSplits(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	splits, err := h.svc.ListSplitsByExpense(r.Context(), id)
	if err != nil {
		handleError(w, r, err)
		return
	}
	limit, offset := parsePagination(r)
	items, total := paginate(toExpenseSplitResponses(splits), limit, offset)
	respondJSON(w, http.StatusOK, paginatedResponse{Data: items, Total: total})
}

func (h *ExpenseHandler) ListByParent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	expenses, err := h.svc.ListByParent(r.Context(), id)
	if err != nil {
		handleError(w, r, err)
		return
	}
	limit, offset := parsePagination(r)
	items, total := paginate(toExpenseResponses(expenses), limit, offset)
	respondJSON(w, http.StatusOK, paginatedResponse{Data: items, Total: total})
}

func (h *ExpenseHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.svc.ListAllCategories(r.Context())
	if err != nil {
		handleError(w, r, err)
		return
	}
	limit, offset := parsePagination(r)
	items, total := paginate(toCategoryResponses(categories), limit, offset)
	respondJSON(w, http.StatusOK, paginatedResponse{Data: items, Total: total})
}

func (h *ExpenseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	expense, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, toExpenseResponse(expense))
}

func (h *ExpenseHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := r.Context().Value(UserIDKey).(string)
	var req struct {
		Description    *string  `json:"description,omitempty"`
		Amount         *float64 `json:"amount,omitempty"         validate:"omitempty,gt=0"`
		CompetenceDate *string  `json:"competence_date,omitempty"`
		DueDate        *string  `json:"due_date,omitempty"`
		CategoryID     *string  `json:"category_id,omitempty"`
		SplitMode      *string  `json:"split_mode,omitempty"     validate:"omitempty,oneof=equal custom"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}

	params := expense.UpdateExpenseParams{
		Description: req.Description,
		Amount:      req.Amount,
		SplitMode:   req.SplitMode,
		CategoryID:  req.CategoryID,
	}
	if req.CompetenceDate != nil {
		t, err := time.Parse("2006-01-02", *req.CompetenceDate)
		if err == nil {
			params.CompetenceDate = &t
		}
	}
	if req.DueDate != nil {
		t, err := time.Parse("2006-01-02", *req.DueDate)
		if err == nil {
			params.DueDate = &t
		}
	}

	expense, err := h.svc.Update(r.Context(), id, params, userID)
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, toExpenseResponse(expense))
}

func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := r.Context().Value(UserIDKey).(string)
	if err := h.svc.Delete(r.Context(), id, userID); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ExpenseHandler) MarkSplitAsPaid(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.MarkSplitAsPaid(r.Context(), id); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ExpenseHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Slug  string `json:"slug"  validate:"required"`
		Label string `json:"label" validate:"required"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	category, err := h.svc.CreateCategory(r.Context(), req.Slug, req.Label)
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusCreated, toCategoryResponse(category))
}

func (h *ExpenseHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		Slug  string `json:"slug"  validate:"required"`
		Label string `json:"label" validate:"required"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	category, err := h.svc.UpdateCategory(r.Context(), id, req.Slug, req.Label)
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusOK, toCategoryResponse(category))
}

func (h *ExpenseHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.DeleteCategory(r.Context(), id); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type expenseResponse struct {
	ID               string   `json:"id"`
	GroupID          string   `json:"group_id"`
	Description      string   `json:"description"`
	Amount           float64  `json:"amount"`
	CategoryID       string   `json:"category_id"`
	CompetenceDate   string   `json:"competence_date"`
	DueDate          string   `json:"due_date"`
	PaidBy           string   `json:"paid_by"`
	SplitMode        string   `json:"split_mode"`
	Installments     int32    `json:"installments"`
	FirstDueDate     *string  `json:"first_due_date,omitempty"`
	IsFixed          bool     `json:"is_fixed"`
	ParentExpenseID  *string  `json:"parent_expense_id,omitempty"`
	InstallmentIndex *int32   `json:"installment_index,omitempty"`
	InstallmentTotal *int32   `json:"installment_total,omitempty"`
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
}

type expenseSplitResponse struct {
	ID        string  `json:"id"`
	ExpenseID string  `json:"expense_id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
	IsPaid    bool    `json:"is_paid"`
	CreatedAt string  `json:"created_at"`
}

type expenseSplitWithUserResponse struct {
	ID        string  `json:"id"`
	ExpenseID string  `json:"expense_id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
	IsPaid    bool    `json:"is_paid"`
	CreatedAt string  `json:"created_at"`
	UserName  string  `json:"user_name"`
	UserEmail string  `json:"user_email"`
}

type categoryResponse struct {
	ID        string `json:"id"`
	Slug      string `json:"slug"`
	Label     string `json:"label"`
	CreatedAt string `json:"created_at"`
}

func toExpenseResponse(e expense.Expense) expenseResponse {
	r := expenseResponse{
		ID:               e.ID,
		GroupID:          e.GroupID,
		Description:      e.Description,
		Amount:           e.Amount,
		CategoryID:       e.CategoryID,
		CompetenceDate:   e.CompetenceDate.Format("2006-01-02"),
		DueDate:          e.DueDate.Format("2006-01-02"),
		PaidBy:           e.PaidBy,
		SplitMode:        e.SplitMode,
		Installments:     e.Installments,
		IsFixed:          e.IsFixed,
		ParentExpenseID:  e.ParentExpenseID,
		InstallmentIndex: e.InstallmentIndex,
		InstallmentTotal: e.InstallmentTotal,
		CreatedAt:        e.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        e.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if e.FirstDueDate != nil {
		s := e.FirstDueDate.Format("2006-01-02")
		r.FirstDueDate = &s
	}
	return r
}

func toExpenseResponses(expenses []expense.Expense) []expenseResponse {
	res := make([]expenseResponse, len(expenses))
	for i, e := range expenses {
		res[i] = toExpenseResponse(e)
	}
	return res
}

func toExpenseSplitResponse(s expense.ExpenseSplit) expenseSplitResponse {
	return expenseSplitResponse{
		ID:        s.ID,
		ExpenseID: s.ExpenseID,
		UserID:    s.UserID,
		Amount:    s.Amount,
		IsPaid:    s.IsPaid,
		CreatedAt: s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toExpenseSplitResponses(splits []expense.ExpenseSplitWithUser) []expenseSplitWithUserResponse {
	res := make([]expenseSplitWithUserResponse, len(splits))
	for i, s := range splits {
		res[i] = expenseSplitWithUserResponse{
			ID:        s.ID,
			ExpenseID: s.ExpenseID,
			UserID:    s.UserID,
			Amount:    s.Amount,
			IsPaid:    s.IsPaid,
			CreatedAt: s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UserName:  s.UserName,
			UserEmail: s.UserEmail,
		}
	}
	return res
}

func toCategoryResponse(c expense.ExpenseCategory) categoryResponse {
	return categoryResponse{
		ID:        c.ID,
		Slug:      c.Slug,
		Label:     c.Label,
		CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toCategoryResponses(categories []expense.ExpenseCategory) []categoryResponse {
	res := make([]categoryResponse, len(categories))
	for i, c := range categories {
		res[i] = toCategoryResponse(c)
	}
	return res
}

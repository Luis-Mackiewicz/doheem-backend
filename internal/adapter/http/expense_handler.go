package http

import (
	"errors"
	"net/http"
	"time"

	"doheem-backend/internal/domain"
)

type ExpenseHandler struct {
	svc *domain.ExpenseService
}

func NewExpenseHandler(svc *domain.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{svc: svc}
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	groupID := r.PathValue("groupId")

	var req struct {
		Description      string   `json:"description"       validate:"required"`
		TotalAmount      float64  `json:"total_amount"       validate:"required,gt=0"`
		ExpenseDate      string   `json:"expense_date"       validate:"required"`
		DueDate          *string  `json:"due_date,omitempty"`
		CategoryID       *string  `json:"category_id,omitempty"`
		SplitType        string   `json:"split_type"         validate:"required,oneof=equal custom"`
		IsInstallment    bool     `json:"is_installment"`
		InstallmentCount *int16   `json:"installment_count,omitempty"`
		Splits           []struct {
			UserID string  `json:"user_id" validate:"required"`
			Amount float64 `json:"amount"   validate:"required,gt=0"`
		} `json:"splits,omitempty"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}

	expenseDate, err := time.Parse("2006-01-02", req.ExpenseDate)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid expense_date, use YYYY-MM-DD")
		return
	}
	var dueDate *time.Time
	if req.DueDate != nil {
		t, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid due_date, use YYYY-MM-DD")
			return
		}
		dueDate = &t
	}

	splits := make([]domain.CreateExpenseSplitParams, len(req.Splits))
	for i, s := range req.Splits {
		splits[i] = domain.CreateExpenseSplitParams{
			UserID: s.UserID,
			Amount: s.Amount,
		}
	}

	expense, err := h.svc.Create(r.Context(), domain.CreateExpenseWithSplitsParams{
		Expense: domain.CreateExpenseParams{
			GroupID:          groupID,
			CreatedBy:        userID,
			Description:      req.Description,
			TotalAmount:      req.TotalAmount,
			ExpenseDate:      expenseDate,
			DueDate:          dueDate,
			CategoryID:       req.CategoryID,
			SplitType:        req.SplitType,
			IsInstallment:    req.IsInstallment,
			InstallmentCount: req.InstallmentCount,
		},
		Splits: splits,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidSplitTotal) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, domain.ErrCategoryNotFound) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toExpenseResponse(expense))
}

func (h *ExpenseHandler) ListByGroup(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("groupId")
	expenses, err := h.svc.ListByGroup(r.Context(), groupID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toExpenseResponses(expenses))
}

func (h *ExpenseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	expense, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toExpenseResponse(expense))
}

func (h *ExpenseHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		Description *string  `json:"description,omitempty"`
		TotalAmount *float64 `json:"total_amount,omitempty" validate:"omitempty,gt=0"`
		ExpenseDate *string  `json:"expense_date,omitempty"`
		DueDate     *string  `json:"due_date,omitempty"`
		CategoryID  *string  `json:"category_id,omitempty"`
		SplitType   *string  `json:"split_type,omitempty"   validate:"omitempty,oneof=equal custom"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}

	params := domain.UpdateExpenseParams{
		Description: req.Description,
		TotalAmount: req.TotalAmount,
		SplitType:   req.SplitType,
		CategoryID:  req.CategoryID,
	}
	if req.ExpenseDate != nil {
		t, err := time.Parse("2006-01-02", *req.ExpenseDate)
		if err == nil {
			params.ExpenseDate = &t
		}
	}
	if req.DueDate != nil {
		t, err := time.Parse("2006-01-02", *req.DueDate)
		if err == nil {
			params.DueDate = &t
		}
	}

	expense, err := h.svc.Update(r.Context(), id, params)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toExpenseResponse(expense))
}

func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ExpenseHandler) ListSplits(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	splits, err := h.svc.ListSplitsByExpense(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toExpenseSplitResponses(splits))
}

func (h *ExpenseHandler) MarkSplitAsPaid(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.MarkSplitAsPaid(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ExpenseHandler) ListInstallments(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	installments, err := h.svc.ListInstallmentsByExpense(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toInstallmentResponses(installments))
}

func (h *ExpenseHandler) MarkInstallmentAsPaid(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.MarkInstallmentAsPaid(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ExpenseHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("groupId")
	var req struct {
		Name string `json:"name" validate:"required"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	category, err := h.svc.CreateCategory(r.Context(), groupID, req.Name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toCategoryResponse(category))
}

func (h *ExpenseHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("groupId")
	categories, err := h.svc.ListCategoriesByGroup(r.Context(), groupID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toCategoryResponses(categories))
}

func (h *ExpenseHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		GroupID string `json:"group_id" validate:"required"`
		Name    string `json:"name"     validate:"required"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	category, err := h.svc.UpdateCategory(r.Context(), id, req.GroupID, req.Name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toCategoryResponse(category))
}

func (h *ExpenseHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	groupID := r.URL.Query().Get("group_id")
	if groupID == "" {
		respondError(w, http.StatusBadRequest, "group_id query param required")
		return
	}
	if err := h.svc.DeleteCategory(r.Context(), id, groupID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type expenseResponse struct {
	ID               string   `json:"id"`
	GroupID          string   `json:"group_id"`
	CreatedBy        string   `json:"created_by"`
	Description      string   `json:"description"`
	TotalAmount      float64  `json:"total_amount"`
	ExpenseDate      string   `json:"expense_date"`
	DueDate          *string  `json:"due_date,omitempty"`
	CategoryID       *string  `json:"category_id,omitempty"`
	SplitType        string   `json:"split_type"`
	IsInstallment    bool     `json:"is_installment"`
	InstallmentCount *int16   `json:"installment_count,omitempty"`
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

type installmentResponse struct {
	ID                string  `json:"id"`
	ExpenseID         string  `json:"expense_id"`
	InstallmentNumber int16   `json:"installment_number"`
	Amount            float64 `json:"amount"`
	DueDate           string  `json:"due_date"`
	IsPaid            bool    `json:"is_paid"`
	CreatedAt         string  `json:"created_at"`
}

type categoryResponse struct {
	ID        string `json:"id"`
	GroupID   string `json:"group_id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

func toExpenseResponse(e domain.Expense) expenseResponse {
	r := expenseResponse{
		ID:               e.ID,
		GroupID:          e.GroupID,
		CreatedBy:        e.CreatedBy,
		Description:      e.Description,
		TotalAmount:      e.TotalAmount,
		ExpenseDate:      e.ExpenseDate.Format("2006-01-02"),
		SplitType:        e.SplitType,
		IsInstallment:    e.IsInstallment,
		InstallmentCount: e.InstallmentCount,
		CreatedAt:        e.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        e.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		CategoryID:       e.CategoryID,
	}
	if e.DueDate != nil {
		s := e.DueDate.Format("2006-01-02")
		r.DueDate = &s
	}
	return r
}

func toExpenseResponses(expenses []domain.Expense) []expenseResponse {
	res := make([]expenseResponse, len(expenses))
	for i, e := range expenses {
		res[i] = toExpenseResponse(e)
	}
	return res
}

func toExpenseSplitResponse(s domain.ExpenseSplit) expenseSplitResponse {
	return expenseSplitResponse{
		ID:        s.ID,
		ExpenseID: s.ExpenseID,
		UserID:    s.UserID,
		Amount:    s.Amount,
		IsPaid:    s.IsPaid,
		CreatedAt: s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toExpenseSplitResponses(splits []domain.ExpenseSplitWithUser) []expenseSplitWithUserResponse {
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

func toInstallmentResponse(i domain.Installment) installmentResponse {
	return installmentResponse{
		ID:                i.ID,
		ExpenseID:         i.ExpenseID,
		InstallmentNumber: i.InstallmentNumber,
		Amount:            i.Amount,
		DueDate:           i.DueDate.Format("2006-01-02"),
		IsPaid:            i.IsPaid,
		CreatedAt:         i.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toInstallmentResponses(installments []domain.Installment) []installmentResponse {
	res := make([]installmentResponse, len(installments))
	for i, inst := range installments {
		res[i] = toInstallmentResponse(inst)
	}
	return res
}

func toCategoryResponse(c domain.ExpenseCategory) categoryResponse {
	return categoryResponse{
		ID:        c.ID,
		GroupID:   c.GroupID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toCategoryResponses(categories []domain.ExpenseCategory) []categoryResponse {
	res := make([]categoryResponse, len(categories))
	for i, c := range categories {
		res[i] = toCategoryResponse(c)
	}
	return res
}

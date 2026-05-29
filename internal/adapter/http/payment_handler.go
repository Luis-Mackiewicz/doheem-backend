package http

import (
	"net/http"
	"time"

	"doheem-backend/internal/domain"
)

type PaymentHandler struct {
	svc *domain.PaymentService
}

func NewPaymentHandler(svc *domain.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

// Create creates a new payment in a group
// @Summary Create a payment
// @Description Create a new payment record in a group
// @Tags Payments
// @Accept json
// @Produce json
// @Param groupId path string true "Group ID"
// @Success 201 {object} paymentResponse
// @Failure 400 {object} map[string]any "Validation error"
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Group not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/groups/{groupId}/payments [post]
// @Security BearerAuth
func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("groupId")
	var req struct {
		PayerID     string  `json:"payer_id"     validate:"required"`
		ReceiverID  string  `json:"receiver_id"  validate:"required"`
		Amount      float64 `json:"amount"       validate:"required,gt=0"`
		PaymentDate string  `json:"payment_date" validate:"required"`
		Notes       *string `json:"notes,omitempty"`
	}
	if errs := decodeAndValidate(r, &req); errs != nil {
		respondValidationError(w, errs)
		return
	}
	paymentDate, err := time.Parse("2006-01-02", req.PaymentDate)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid payment_date, use YYYY-MM-DD")
		return
	}
	payment, err := h.svc.Create(r.Context(), domain.CreatePaymentParams{
		GroupID:     groupID,
		PayerID:     req.PayerID,
		ReceiverID:  req.ReceiverID,
		Amount:      req.Amount,
		PaymentDate: paymentDate,
		Notes:       req.Notes,
	})
	if err != nil {
		handleError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, toPaymentResponse(payment))
}

// ListByGroup lists payments in a group
// @Summary List payments by group
// @Description List all payments for a specific group
// @Tags Payments
// @Accept json
// @Produce json
// @Param groupId path string true "Group ID"
// @Success 200 {array} paymentWithUserResponse
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Group not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/groups/{groupId}/payments [get]
// @Security BearerAuth
func (h *PaymentHandler) ListByGroup(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("groupId")
	payments, err := h.svc.ListByGroup(r.Context(), groupID)
	if err != nil {
		handleError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, toPaymentWithUserResponses(payments))
}

// GetByID gets a payment by ID
// @Summary Get payment by ID
// @Description Get a single payment by its unique identifier
// @Tags Payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} paymentResponse
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Payment not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/payments/{id} [get]
// @Security BearerAuth
func (h *PaymentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	payment, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, toPaymentResponse(payment))
}

// Confirm confirms a payment
// @Summary Confirm a payment
// @Description Confirm a pending payment
// @Tags Payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Success 204 {object} nil
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Payment not found"
// @Failure 409 {object} map[string]any "Conflict"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/payments/{id}/confirm [patch]
// @Security BearerAuth
func (h *PaymentHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Confirm(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Cancel cancels a payment
// @Summary Cancel a payment
// @Description Cancel a pending payment
// @Tags Payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Success 204 {object} nil
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Payment not found"
// @Failure 409 {object} map[string]any "Conflict"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/payments/{id}/cancel [patch]
// @Security BearerAuth
func (h *PaymentHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Cancel(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Delete deletes a payment
// @Summary Delete a payment
// @Description Permanently delete a payment
// @Tags Payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Success 204 {object} nil
// @Failure 401 {object} map[string]any "Unauthorized"
// @Failure 404 {object} map[string]any "Payment not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/payments/{id} [delete]
// @Security BearerAuth
func (h *PaymentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type paymentResponse struct {
	ID          string   `json:"id"`
	GroupID     string   `json:"group_id"`
	PayerID     string   `json:"payer_id"`
	ReceiverID  string   `json:"receiver_id"`
	Amount      float64  `json:"amount"`
	PaymentDate string   `json:"payment_date"`
	Status      string   `json:"status"`
	Notes       *string  `json:"notes,omitempty"`
	CreatedAt   string   `json:"created_at"`
	ConfirmedAt *string  `json:"confirmed_at,omitempty"`
	CancelledAt *string  `json:"cancelled_at,omitempty"`
}

type paymentWithUserResponse struct {
	ID           string   `json:"id"`
	GroupID      string   `json:"group_id"`
	PayerID      string   `json:"payer_id"`
	ReceiverID   string   `json:"receiver_id"`
	Amount       float64  `json:"amount"`
	PaymentDate  string   `json:"payment_date"`
	Status       string   `json:"status"`
	Notes        *string  `json:"notes,omitempty"`
	CreatedAt    string   `json:"created_at"`
	ConfirmedAt  *string  `json:"confirmed_at,omitempty"`
	CancelledAt  *string  `json:"cancelled_at,omitempty"`
	PayerName    string   `json:"payer_name"`
	ReceiverName string   `json:"receiver_name"`
}

func toPaymentResponse(p domain.Payment) paymentResponse {
	r := paymentResponse{
		ID:          p.ID,
		GroupID:     p.GroupID,
		PayerID:     p.PayerID,
		ReceiverID:  p.ReceiverID,
		Amount:      p.Amount,
		PaymentDate: p.PaymentDate.Format("2006-01-02"),
		Status:      p.Status,
		Notes:       p.Notes,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if p.ConfirmedAt != nil {
		s := p.ConfirmedAt.Format("2006-01-02T15:04:05Z07:00")
		r.ConfirmedAt = &s
	}
	if p.CancelledAt != nil {
		s := p.CancelledAt.Format("2006-01-02T15:04:05Z07:00")
		r.CancelledAt = &s
	}
	return r
}

func toPaymentWithUserResponses(payments []domain.PaymentWithUsers) []paymentWithUserResponse {
	res := make([]paymentWithUserResponse, len(payments))
	for i, p := range payments {
		r := paymentWithUserResponse{
			ID:           p.ID,
			GroupID:      p.GroupID,
			PayerID:      p.PayerID,
			ReceiverID:   p.ReceiverID,
			Amount:       p.Amount,
			PaymentDate:  p.PaymentDate.Format("2006-01-02"),
			Status:       p.Status,
			Notes:        p.Notes,
			CreatedAt:    p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			PayerName:    p.PayerName,
			ReceiverName: p.ReceiverName,
		}
		if p.ConfirmedAt != nil {
			s := p.ConfirmedAt.Format("2006-01-02T15:04:05Z07:00")
			r.ConfirmedAt = &s
		}
		if p.CancelledAt != nil {
			s := p.CancelledAt.Format("2006-01-02T15:04:05Z07:00")
			r.CancelledAt = &s
		}
		res[i] = r
	}
	return res
}

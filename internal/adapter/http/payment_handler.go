package http

import (
	"encoding/json"
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

func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("groupId")
	var req struct {
		PayerID     string  `json:"payer_id"`
		ReceiverID  string  `json:"receiver_id"`
		Amount      float64 `json:"amount"`
		PaymentDate string  `json:"payment_date"`
		Notes       *string `json:"notes,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
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
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, toPaymentResponse(payment))
}

func (h *PaymentHandler) ListByGroup(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("groupId")
	payments, err := h.svc.ListByGroup(r.Context(), groupID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toPaymentWithUserResponses(payments))
}

func (h *PaymentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	payment, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, toPaymentResponse(payment))
}

func (h *PaymentHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Confirm(r.Context(), id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PaymentHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Cancel(r.Context(), id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PaymentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
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

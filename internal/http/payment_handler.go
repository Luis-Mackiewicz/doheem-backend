package http

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"doheem-backend/internal/payment"

	"github.com/google/uuid"
)

type PaymentHandler struct {
	svc *payment.PaymentService
}

func NewPaymentHandler(svc *payment.PaymentService) *PaymentHandler {
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
	p, err := h.svc.Create(r.Context(), payment.CreatePaymentParams{
		GroupID:     groupID,
		PayerID:     req.PayerID,
		ReceiverID:  req.ReceiverID,
		Amount:      req.Amount,
		PaymentDate: paymentDate,
		Notes:       req.Notes,
	})
	if err != nil {
		handleError(w, r, err)
		return
	}
	respondJSON(w, http.StatusCreated, toPaymentResponse(p))
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
		handleError(w, r, err)
		return
	}
	limit, offset := parsePagination(r)
	items, total := paginate(toPaymentWithUserResponses(payments), limit, offset)
	respondJSON(w, http.StatusOK, paginatedResponse{Data: items, Total: total})
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
		handleError(w, r, err)
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
		handleError(w, r, err)
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
		handleError(w, r, err)
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
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UploadAttachment uploads a file attachment for a payment
// @Summary Upload payment attachment
// @Description Upload a file (jpeg, png, gif, pdf) as an attachment for a payment
// @Tags Payments
// @Accept mpfd
// @Produce json
// @Param id path string true "Payment ID"
// @Param file formData file true "Attachment file"
// @Success 201 {object} paymentAttachmentResponse
// @Failure 400 {object} map[string]any "Invalid file or request"
// @Failure 404 {object} map[string]any "Payment not found"
// @Failure 500 {object} map[string]any "Internal server error"
// @Router /api/payments/{id}/attachments [post]
// @Security BearerAuth
func (h *PaymentHandler) UploadAttachment(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("id")

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "file too large or invalid multipart form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		respondError(w, http.StatusBadRequest, "missing file field")
		return
	}
	defer file.Close()

	allowedTypes := map[string]bool{
		"image/jpeg": true, "image/png": true, "image/gif": true,
		"application/pdf": true,
	}
	contentType := header.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		respondError(w, http.StatusBadRequest, "invalid file type, allowed: jpeg, png, gif, pdf")
		return
	}

	ext := filepath.Ext(header.Filename)
	filename := uuid.NewString() + ext
	uploadDir := "./uploads"
	os.MkdirAll(uploadDir, 0755)
	dst, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to save file")
		return
	}

	attachment, err := h.svc.AddAttachment(r.Context(), paymentID, filepath.Join(uploadDir, filename), contentType, int32(header.Size))
	if err != nil {
		handleError(w, r, err)
		return
	}

	respondJSON(w, http.StatusCreated, toPaymentAttachmentResponse(attachment))
}

type paymentAttachmentResponse struct {
	ID        string `json:"id"`
	PaymentID string `json:"payment_id"`
	FilePath  string `json:"file_path"`
	FileType  string `json:"file_type"`
	FileSize  int32  `json:"file_size"`
	UploadedAt string `json:"uploaded_at"`
}

func toPaymentAttachmentResponse(a payment.PaymentAttachment) paymentAttachmentResponse {
	return paymentAttachmentResponse{
		ID:         a.ID,
		PaymentID:  a.PaymentID,
		FilePath:   a.FilePath,
		FileType:   a.FileType,
		FileSize:   a.FileSize,
		UploadedAt: a.UploadedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
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

func toPaymentResponse(p payment.Payment) paymentResponse {
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

func toPaymentWithUserResponses(payments []payment.PaymentWithUsers) []paymentWithUserResponse {
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

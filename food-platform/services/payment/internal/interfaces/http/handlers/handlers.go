package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/food-platform/payment/internal/application"
	"github.com/food-platform/payment/internal/domain"
	"github.com/food-platform/shared/errors"
	"github.com/food-platform/shared/logging"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, err error) {
	statusCode, code, message := errors.ToHTTP(err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{"code": code, "message": message},
	})
}

// ============ Charge Payment ============

type ChargePaymentHandler struct {
	uc *application.ChargePaymentUseCase
}

func NewChargePaymentHandler(uc *application.ChargePaymentUseCase) *ChargePaymentHandler {
	return &ChargePaymentHandler{uc: uc}
}

type chargePaymentRequest struct {
	OrderID string `json:"order_id"`
	Method  string `json:"method"`
	Amount  float64 `json:"amount"`
}

func (h *ChargePaymentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userIDStr := logging.GetUserID(r.Context())
	customerID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, errors.ErrUnauthorized)
		return
	}

	var req chargePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.ErrInvalidJSON)
		return
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		writeError(w, errors.New("INVALID_ID", "Invalid order ID", 400))
		return
	}

	idempotencyKey := r.Header.Get("X-Idempotency-Key")
	if idempotencyKey == "" {
		idempotencyKey = uuid.New().String()
	}

	cmd := application.ChargePaymentCommand{
		OrderID:        orderID,
		CustomerID:     customerID,
		Method:         domain.PaymentMethod(req.Method),
		Amount:         req.Amount,
		IdempotencyKey: idempotencyKey,
	}

	result, err := h.uc.Execute(r.Context(), cmd)
	if err != nil {
		logging.FromContext(r.Context()).Error("charge_payment_failed", "error", err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// ============ Get Payment ============

type GetPaymentHandler struct {
	uc *application.GetPaymentUseCase
}

func NewGetPaymentHandler(uc *application.GetPaymentUseCase) *GetPaymentHandler {
	return &GetPaymentHandler{uc: uc}
}

func (h *GetPaymentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, errors.New("INVALID_ID", "Invalid payment ID", 400))
		return
	}

	result, err := h.uc.Execute(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// ============ Get Payment by Order ============

type GetPaymentByOrderHandler struct {
	uc *application.GetPaymentByOrderUseCase
}

func NewGetPaymentByOrderHandler(uc *application.GetPaymentByOrderUseCase) *GetPaymentByOrderHandler {
	return &GetPaymentByOrderHandler{uc: uc}
}

func (h *GetPaymentByOrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	orderIDStr := chi.URLParam(r, "orderId")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		writeError(w, errors.New("INVALID_ID", "Invalid order ID", 400))
		return
	}

	result, err := h.uc.Execute(r.Context(), orderID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// ============ Refund ============

type RefundHandler struct {
	uc *application.RefundPaymentUseCase
}

func NewRefundHandler(uc *application.RefundPaymentUseCase) *RefundHandler {
	return &RefundHandler{uc: uc}
}

type refundRequest struct {
	Amount float64 `json:"amount"`
	Reason string  `json:"reason"`
	Type   string  `json:"type"`
}

func (h *RefundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	paymentID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, errors.New("INVALID_ID", "Invalid payment ID", 400))
		return
	}

	var req refundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.ErrInvalidJSON)
		return
	}

	if req.Reason == "" {
		writeError(w, errors.New("VALIDATION_FAILED", "refund reason is required", 422))
		return
	}

	cmd := application.RefundCommand{
		PaymentID: paymentID,
		Amount:    req.Amount,
		Reason:    req.Reason,
		Type:      domain.RefundType(req.Type),
	}

	result, err := h.uc.Execute(r.Context(), cmd)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// ============ Health ============

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "payment",
		"version": "1.0.0",
	})
}

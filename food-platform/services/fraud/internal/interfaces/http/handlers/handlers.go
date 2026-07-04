package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/food-platform/fraud/internal/application"
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

type ScoreOrderHandler struct {
	uc *application.ScoreOrderUseCase
}

func NewScoreOrderHandler(uc *application.ScoreOrderUseCase) *ScoreOrderHandler {
	return &ScoreOrderHandler{uc: uc}
}

type scoreOrderRequest struct {
	OrderID     string  `json:"order_id"`
	CustomerID  string  `json:"customer_id"`
	OrderAmount float64 `json:"order_amount"`
	IsNewUser   bool    `json:"is_new_user"`
	OrderCount  int     `json:"order_count"`
	RefundCount int     `json:"refund_count"`
}

func (h *ScoreOrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req scoreOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.ErrInvalidJSON)
		return
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		writeError(w, errors.New("INVALID_ID", "Invalid order ID", 400))
		return
	}

	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		writeError(w, errors.New("INVALID_ID", "Invalid customer ID", 400))
		return
	}

	result, err := h.uc.Execute(r.Context(), application.ScoreOrderCommand{
		OrderID:     orderID,
		CustomerID:  customerID,
		OrderAmount: req.OrderAmount,
		IsNewUser:   req.IsNewUser,
		OrderCount:  req.OrderCount,
		RefundCount: req.RefundCount,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

type GetTrustScoreHandler struct {
	uc *application.GetTrustScoreUseCase
}

func NewGetTrustScoreHandler(uc *application.GetTrustScoreUseCase) *GetTrustScoreHandler {
	return &GetTrustScoreHandler{uc: uc}
}

func (h *GetTrustScoreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	customerIDStr := chi.URLParam(r, "customerId")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		writeError(w, errors.New("INVALID_ID", "Invalid customer ID", 400))
		return
	}

	result, err := h.uc.Execute(r.Context(), customerID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok", "service": "fraud", "version": "1.0.0",
	})
}

// Suppress unused import
var _ = logging.GetUserID

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/food-platform/delivery-matching/internal/application"
	"github.com/food-platform/shared/errors"
	"github.com/food-platform/shared/logging"
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

type MatchOrderHandler struct {
	uc *application.MatchOrderUseCase
}

func NewMatchOrderHandler(uc *application.MatchOrderUseCase) *MatchOrderHandler {
	return &MatchOrderHandler{uc: uc}
}

type matchOrderRequest struct {
	OrderID       string `json:"order_id"`
	RestaurantLat float64 `json:"restaurant_lat"`
	RestaurantLng float64 `json:"restaurant_lng"`
	CustomerLat   float64 `json:"customer_lat"`
	CustomerLng   float64 `json:"customer_lng"`
}

func (h *MatchOrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req matchOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.ErrInvalidJSON)
		return
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		writeError(w, errors.New("INVALID_ID", "Invalid order ID", 400))
		return
	}

	result, err := h.uc.Execute(r.Context(), application.MatchOrderCommand{
		OrderID:       orderID,
		RestaurantLat: req.RestaurantLat,
		RestaurantLng: req.RestaurantLng,
		CustomerLat:   req.CustomerLat,
		CustomerLng:   req.CustomerLng,
	})
	if err != nil {
		writeJSON(w, http.StatusOK, result) // Return result even on "no drivers"
		return
	}

	writeJSON(w, http.StatusOK, result)
}

type AcceptOrderHandler struct {
	uc *application.AcceptOrderUseCase
}

func NewAcceptOrderHandler(uc *application.AcceptOrderUseCase) *AcceptOrderHandler {
	return &AcceptOrderHandler{uc: uc}
}

type acceptOrderRequest struct {
	OrderID string `json:"order_id"`
}

func (h *AcceptOrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	driverIDStr := logging.GetUserID(r.Context())
	driverID, err := uuid.Parse(driverIDStr)
	if err != nil {
		writeError(w, errors.ErrUnauthorized)
		return
	}

	var req acceptOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.ErrInvalidJSON)
		return
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		writeError(w, errors.New("INVALID_ID", "Invalid order ID", 400))
		return
	}

	err = h.uc.Execute(r.Context(), application.AcceptOrderCommand{
		OrderID:  orderID,
		DriverID: driverID,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "assigned"})
}

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok", "service": "delivery-matching", "version": "1.0.0",
	})
}

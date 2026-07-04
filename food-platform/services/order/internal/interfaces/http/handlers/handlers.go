package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/food-platform/order/internal/application"
	"github.com/food-platform/order/internal/domain"
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

// ============ Create Order ============

type CreateOrderHandler struct {
	uc *application.CreateOrderUseCase
}

func NewCreateOrderHandler(uc *application.CreateOrderUseCase) *CreateOrderHandler {
	return &CreateOrderHandler{uc: uc}
}

type createOrderItemRequest struct {
	MenuItemID string `json:"menu_item_id"`
	Name       string `json:"name"`
	Quantity   int    `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	Notes      string `json:"notes"`
}

type createOrderRequest struct {
	RestaurantID    string                  `json:"restaurant_id"`
	Items           []createOrderItemRequest `json:"items"`
	DeliveryAddress string                  `json:"delivery_address"`
	Latitude        float64                 `json:"latitude"`
	Longitude       float64                 `json:"longitude"`
	PaymentMethod   string                  `json:"payment_method"`
	DeliveryFee     float64                 `json:"delivery_fee"`
	Notes           string                  `json:"notes"`
}

func (h *CreateOrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userIDStr := logging.GetUserID(r.Context())
	customerID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, errors.ErrUnauthorized)
		return
	}

	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.ErrInvalidJSON)
		return
	}

	restaurantID, err := uuid.Parse(req.RestaurantID)
	if err != nil {
		writeError(w, errors.New("INVALID_ID", "Invalid restaurant ID", 400))
		return
	}

	items := make([]application.CreateOrderItemCommand, len(req.Items))
	for i, item := range req.Items {
		menuItemID, err := uuid.Parse(item.MenuItemID)
		if err != nil {
			writeError(w, errors.New("INVALID_ID", "Invalid menu item ID", 400))
			return
		}
		items[i] = application.CreateOrderItemCommand{
			MenuItemID: menuItemID,
			Name:       item.Name,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			Notes:      item.Notes,
		}
	}

	// Idempotency key
	idempotencyKey := r.Header.Get("X-Idempotency-Key")
	if idempotencyKey == "" {
		idempotencyKey = uuid.New().String()
	}

	cmd := application.CreateOrderCommand{
		CustomerID:      customerID,
		RestaurantID:    restaurantID,
		Items:           items,
		DeliveryAddress: req.DeliveryAddress,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		PaymentMethod:   domain.PaymentMethod(req.PaymentMethod),
		DeliveryFee:     req.DeliveryFee,
		ServiceFeeRate:  0.05,
		VATRate:         0.14,
		Discount:        0,
		Notes:           req.Notes,
	}

	result, err := h.uc.Execute(r.Context(), cmd)
	if err != nil {
		logging.FromContext(r.Context()).Error("create_order_failed", "error", err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// ============ Get Order ============

type GetOrderHandler struct {
	uc *application.GetOrderUseCase
}

func NewGetOrderHandler(uc *application.GetOrderUseCase) *GetOrderHandler {
	return &GetOrderHandler{uc: uc}
}

func (h *GetOrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, errors.New("INVALID_ID", "Invalid order ID", 400))
		return
	}

	result, err := h.uc.Execute(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// ============ Get Active Orders ============

type GetActiveOrdersHandler struct {
	uc *application.GetActiveOrdersUseCase
}

func NewGetActiveOrdersHandler(uc *application.GetActiveOrdersUseCase) *GetActiveOrdersHandler {
	return &GetActiveOrdersHandler{uc: uc}
}

func (h *GetActiveOrdersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userIDStr := logging.GetUserID(r.Context())
	customerID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, errors.ErrUnauthorized)
		return
	}

	orders, err := h.uc.Execute(r.Context(), customerID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"orders": orders,
		"total":  len(orders),
	})
}

// ============ Get Order History ============

type GetOrderHistoryHandler struct {
	uc *application.GetOrderHistoryUseCase
}

func NewGetOrderHistoryHandler(uc *application.GetOrderHistoryUseCase) *GetOrderHistoryHandler {
	return &GetOrderHistoryHandler{uc: uc}
}

func (h *GetOrderHistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userIDStr := logging.GetUserID(r.Context())
	customerID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, errors.ErrUnauthorized)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	orders, err := h.uc.Execute(r.Context(), customerID, limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"orders": orders,
		"total":  len(orders),
	})
}

// ============ Cancel Order ============

type CancelOrderHandler struct {
	uc *application.CancelOrderUseCase
}

func NewCancelOrderHandler(uc *application.CancelOrderUseCase) *CancelOrderHandler {
	return &CancelOrderHandler{uc: uc}
}

type cancelOrderRequest struct {
	Reason string `json:"reason"`
}

func (h *CancelOrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, errors.New("INVALID_ID", "Invalid order ID", 400))
		return
	}

	var req cancelOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.ErrInvalidJSON)
		return
	}

	cmd := application.CancelOrderCommand{
		OrderID: orderID,
		Reason:  req.Reason,
	}

	if err := h.uc.Execute(r.Context(), cmd); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
}

// ============ Update Status ============

type UpdateStatusHandler struct {
	uc *application.UpdateOrderStatusUseCase
}

func NewUpdateStatusHandler(uc *application.UpdateOrderStatusUseCase) *UpdateStatusHandler {
	return &UpdateStatusHandler{uc: uc}
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

func (h *UpdateStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, errors.New("INVALID_ID", "Invalid order ID", 400))
		return
	}

	var req updateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.ErrInvalidJSON)
		return
	}

	cmd := application.UpdateStatusCommand{
		OrderID: orderID,
		Status:  domain.OrderStatus(req.Status),
	}

	if err := h.uc.Execute(r.Context(), cmd); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// ============ Health ============

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "order",
		"version": "1.0.0",
	})
}

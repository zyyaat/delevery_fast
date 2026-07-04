package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/food-platform/driver-management/internal/application"
	"github.com/food-platform/driver-management/internal/domain"
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

// ============ Register Driver ============

type RegisterDriverHandler struct {
	uc *application.RegisterDriverUseCase
}

func NewRegisterDriverHandler(uc *application.RegisterDriverUseCase) *RegisterDriverHandler {
	return &RegisterDriverHandler{uc: uc}
}

type registerDriverRequest struct {
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	VehicleType string `json:"vehicle_type"`
}

func (h *RegisterDriverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userIDStr := logging.GetUserID(r.Context())
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, errors.ErrUnauthorized)
		return
	}

	var req registerDriverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.ErrInvalidJSON)
		return
	}

	result, err := h.uc.Execute(r.Context(), application.RegisterDriverCommand{
		UserID:      userID,
		Name:        req.Name,
		Phone:       req.Phone,
		VehicleType: req.VehicleType,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// ============ Get Driver ============

type GetDriverHandler struct {
	uc *application.GetDriverUseCase
}

func NewGetDriverHandler(uc *application.GetDriverUseCase) *GetDriverHandler {
	return &GetDriverHandler{uc: uc}
}

func (h *GetDriverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userIDStr := logging.GetUserID(r.Context())
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, errors.ErrUnauthorized)
		return
	}

	// Get driver by user ID (the authenticated user)
	result, err := h.uc.Execute(r.Context(), userID)
	if err != nil {
		// Try by driver ID from URL
		idStr := chi.URLParam(r, "id")
		if idStr != "" {
			driverID, err := uuid.Parse(idStr)
			if err == nil {
				getByDriverID := application.NewGetDriverByUserUseCase(nil)
				_ = getByDriverID
				result, err = h.uc.Execute(r.Context(), driverID)
				if err != nil {
					writeError(w, err)
					return
				}
				writeJSON(w, http.StatusOK, result)
				return
			}
		}
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// ============ Update Status ============

type UpdateStatusHandler struct {
	uc *application.UpdateStatusUseCase
}

func NewUpdateStatusHandler(uc *application.UpdateStatusUseCase) *UpdateStatusHandler {
	return &UpdateStatusHandler{uc: uc}
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

func (h *UpdateStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userIDStr := logging.GetUserID(r.Context())
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, errors.ErrUnauthorized)
		return
	}

	// For now, use userID as driverID (in production, lookup driver by userID)
	var req updateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.ErrInvalidJSON)
		return
	}

	err = h.uc.Execute(r.Context(), application.UpdateStatusCommand{
		DriverID: userID,
		Status:   req.Status,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// ============ Update Location ============

type UpdateLocationHandler struct {
	uc *application.UpdateLocationUseCase
}

func NewUpdateLocationHandler(uc *application.UpdateLocationUseCase) *UpdateLocationHandler {
	return &UpdateLocationHandler{uc: uc}
}

type updateLocationRequest struct {
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Heading float64 `json:"heading"`
	Speed   float64 `json:"speed"`
}

func (h *UpdateLocationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userIDStr := logging.GetUserID(r.Context())
	driverID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, errors.ErrUnauthorized)
		return
	}

	var req updateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.ErrInvalidJSON)
		return
	}

	err = h.uc.Execute(r.Context(), application.UpdateLocationCommand{
		DriverID: driverID,
		Lat:      req.Lat,
		Lng:      req.Lng,
		Heading:  req.Heading,
		Speed:    req.Speed,
	})
	if err != nil {
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
		"service": "driver-management",
		"version": "1.0.0",
	})
}

// Suppress unused import
var _ = domain.DriverStatusOffline

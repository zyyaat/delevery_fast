package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/food-platform/geo/internal/application"
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
	driverIDStr := logging.GetUserID(r.Context())
	driverID, err := uuid.Parse(driverIDStr)
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

type GetLocationHandler struct {
	uc *application.GetLocationUseCase
}

func NewGetLocationHandler(uc *application.GetLocationUseCase) *GetLocationHandler {
	return &GetLocationHandler{uc: uc}
}

func (h *GetLocationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	driverIDStr := r.URL.Query().Get("driver_id")
	driverID, err := uuid.Parse(driverIDStr)
	if err != nil {
		writeError(w, errors.New("INVALID_ID", "Invalid driver ID", 400))
		return
	}

	result, err := h.uc.Execute(r.Context(), driverID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

type FindNearbyHandler struct {
	uc *application.FindNearbyUseCase
}

func NewFindNearbyHandler(uc *application.FindNearbyUseCase) *FindNearbyHandler {
	return &FindNearbyHandler{uc: uc}
}

func (h *FindNearbyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lng, _ := strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
	radius, _ := strconv.ParseFloat(r.URL.Query().Get("radius"), 64)
	count, _ := strconv.Atoi(r.URL.Query().Get("count"))

	if lat == 0 || lng == 0 {
		writeError(w, errors.New("MISSING_PARAMS", "lat and lng required", 400))
		return
	}

	drivers, err := h.uc.Execute(r.Context(), application.FindNearbyCommand{
		Lat: lat, Lng: lng, RadiusKm: radius, Count: count,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"drivers": drivers,
		"total":   len(drivers),
	})
}

type CalculateETAHandler struct {
	uc *application.CalculateETAUseCase
}

func NewCalculateETAHandler(uc *application.CalculateETAUseCase) *CalculateETAHandler {
	return &CalculateETAHandler{uc: uc}
}

func (h *CalculateETAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	originLat, _ := strconv.ParseFloat(r.URL.Query().Get("origin_lat"), 64)
	originLng, _ := strconv.ParseFloat(r.URL.Query().Get("origin_lng"), 64)
	destLat, _ := strconv.ParseFloat(r.URL.Query().Get("dest_lat"), 64)
	destLng, _ := strconv.ParseFloat(r.URL.Query().Get("dest_lng"), 64)

	result, err := h.uc.Execute(r.Context(), application.CalculateETACommand{
		OriginLat: originLat, OriginLng: originLng,
		DestLat: destLat, DestLng: destLng,
	})
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
		"status": "ok", "service": "geo", "version": "1.0.0",
	})
}

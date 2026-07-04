// Package handlers contains HTTP handlers for the Restaurant Catalog Service.
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/food-platform/services/restaurant-catalog/internal/application"
	"github.com/food-platform/services/restaurant-catalog/internal/domain"
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
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	})
}

func decodeJSON(r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return errors.ErrInvalidJSON
	}
	return json.NewDecoder(r.Body).Decode(dst)
}

// ============ Create Restaurant ============

type CreateRestaurantHandler struct {
	uc *application.CreateRestaurantUseCase
}

func NewCreateRestaurantHandler(uc *application.CreateRestaurantUseCase) *CreateRestaurantHandler {
	return &CreateRestaurantHandler{uc: uc}
}

type createRestaurantRequest struct {
	Name         string   `json:"name"`
	Latitude     float64  `json:"latitude"`
	Longitude    float64  `json:"longitude"`
	Address      string   `json:"address"`
	City         string   `json:"city"`
	CuisineTypes []string `json:"cuisine_types"`
}

func (h *CreateRestaurantHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req createRestaurantRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, errors.ErrInvalidJSON)
		return
	}

	cuisines := make([]domain.CuisineType, len(req.CuisineTypes))
	for i, c := range req.CuisineTypes {
		cuisines[i] = domain.CuisineType(c)
	}

	cmd := application.CreateRestaurantCommand{
		Name:         req.Name,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Address:      req.Address,
		City:         req.City,
		CuisineTypes: cuisines,
	}

	result, err := h.uc.Execute(r.Context(), cmd)
	if err != nil {
		logging.FromContext(r.Context()).Error("create_restaurant_failed", "error", err)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// ============ Get Restaurant ============

type GetRestaurantHandler struct {
	uc *application.GetRestaurantUseCase
}

func NewGetRestaurantHandler(uc *application.GetRestaurantUseCase) *GetRestaurantHandler {
	return &GetRestaurantHandler{uc: uc}
}

func (h *GetRestaurantHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, errors.New("INVALID_ID", "Invalid restaurant ID", 400))
		return
	}

	result, err := h.uc.Execute(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// ============ Find Nearby ============

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
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if lat == 0 || lng == 0 {
		writeError(w, errors.New("MISSING_PARAMS", "lat and lng query params required", 400))
		return
	}

	q := application.NearbyQuery{
		Latitude:  lat,
		Longitude: lng,
		RadiusKm:  radius,
		Limit:     limit,
		Offset:    offset,
	}

	restaurants, err := h.uc.Execute(r.Context(), q)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"restaurants": restaurants,
		"total":       len(restaurants),
		"has_more":    len(restaurants) >= (limit),
	})
}

// ============ Search ============

type SearchHandler struct {
	uc *application.SearchRestaurantsUseCase
}

func NewSearchHandler(uc *application.SearchRestaurantsUseCase) *SearchHandler {
	return &SearchHandler{uc: uc}
}

func (h *SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		writeError(w, errors.New("MISSING_QUERY", "q parameter required", 400))
		return
	}

	lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lng, _ := strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit == 0 {
		limit = 20
	}

	restaurants, err := h.uc.Execute(r.Context(), q, lat, lng, limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"restaurants": restaurants,
		"total":       len(restaurants),
	})
}

// ============ Health ============

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "restaurant-catalog",
		"version": "1.0.0",
	})
}

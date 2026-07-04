package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/food-platform/services/menu/internal/application"
	"github.com/food-platform/shared/errors"
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

// ============ Create Menu Item ============

type CreateMenuItemHandler struct {
	uc *application.CreateMenuItemUseCase
}

func NewCreateMenuItemHandler(uc *application.CreateMenuItemUseCase) *CreateMenuItemHandler {
	return &CreateMenuItemHandler{uc: uc}
}

type createMenuItemRequest struct {
	CategoryID   string                     `json:"category_id"`
	Name         string                     `json:"name"`
	Description  string                     `json:"description"`
	Price        float64                    `json:"price"`
	ImageURL     string                     `json:"image_url"`
	PrepTimeMin  int                        `json:"prep_time_minutes"`
	Modifiers    []createModifierRequest    `json:"modifiers"`
}

type createModifierRequest struct {
	Name           string                    `json:"name"`
	Required       bool                      `json:"required"`
	MultipleChoice bool                      `json:"multiple_choice"`
	Options        []createModifierOptionReq `json:"options"`
}

type createModifierOptionReq struct {
	Name       string  `json:"name"`
	PriceDelta float64 `json:"price_delta"`
}

func (h *CreateMenuItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	restaurantIDStr := chi.URLParam(r, "restaurantId")
	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		writeError(w, errors.New("INVALID_ID", "Invalid restaurant ID", 400))
		return
	}

	var req createMenuItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.ErrInvalidJSON)
		return
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		writeError(w, errors.New("INVALID_ID", "Invalid category ID", 400))
		return
	}

	mods := make([]application.CreateModifierCommand, len(req.Modifiers))
	for i, m := range req.Modifiers {
		opts := make([]application.CreateModifierOptionCommand, len(m.Options))
		for j, o := range m.Options {
			opts[j] = application.CreateModifierOptionCommand{
				Name:       o.Name,
				PriceDelta: o.PriceDelta,
			}
		}
		mods[i] = application.CreateModifierCommand{
			Name:           m.Name,
			Required:       m.Required,
			MultipleChoice: m.MultipleChoice,
			Options:        opts,
		}
	}

	cmd := application.CreateMenuItemCommand{
		RestaurantID: restaurantID,
		CategoryID:   categoryID,
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		ImageURL:     req.ImageURL,
		PrepTimeMin:  req.PrepTimeMin,
		Modifiers:    mods,
	}

	result, err := h.uc.Execute(r.Context(), cmd)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// ============ Get Menu ============

type GetMenuHandler struct {
	uc *application.GetMenuUseCase
}

func NewGetMenuHandler(uc *application.GetMenuUseCase) *GetMenuHandler {
	return &GetMenuHandler{uc: uc}
}

func (h *GetMenuHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	restaurantIDStr := chi.URLParam(r, "restaurantId")
	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		writeError(w, errors.New("INVALID_ID", "Invalid restaurant ID", 400))
		return
	}

	result, err := h.uc.Execute(r.Context(), restaurantID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"categories": result,
	})
}

// ============ Toggle Availability ============

type ToggleAvailabilityHandler struct {
	uc *application.ToggleAvailabilityUseCase
}

func NewToggleAvailabilityHandler(uc *application.ToggleAvailabilityUseCase) *ToggleAvailabilityHandler {
	return &ToggleAvailabilityHandler{uc: uc}
}

type toggleRequest struct {
	IsAvailable bool `json:"is_available"`
}

func (h *ToggleAvailabilityHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	itemIDStr := chi.URLParam(r, "itemId")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		writeError(w, errors.New("INVALID_ID", "Invalid item ID", 400))
		return
	}

	var req toggleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.ErrInvalidJSON)
		return
	}

	if err := h.uc.Execute(r.Context(), itemID, req.IsAvailable); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// ============ Create Category ============

type CreateCategoryHandler struct {
	uc *application.CreateCategoryUseCase
}

func NewCreateCategoryHandler(uc *application.CreateCategoryUseCase) *CreateCategoryHandler {
	return &CreateCategoryHandler{uc: uc}
}

type createCategoryRequest struct {
	Name         string `json:"name"`
	DisplayOrder int    `json:"display_order"`
}

func (h *CreateCategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	restaurantIDStr := chi.URLParam(r, "restaurantId")
	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		writeError(w, errors.New("INVALID_ID", "Invalid restaurant ID", 400))
		return
	}

	var req createCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.ErrInvalidJSON)
		return
	}

	cmd := application.CreateCategoryCommand{
		RestaurantID: restaurantID,
		Name:         req.Name,
		DisplayOrder: req.DisplayOrder,
	}

	result, err := h.uc.Execute(r.Context(), cmd)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// ============ Health ============

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "menu",
		"version": "1.0.0",
	})
}

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/food-platform/analytics/internal/application"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

type DashboardHandler struct {
	uc *application.GetDashboardStatsUseCase
}

func NewDashboardHandler(uc *application.GetDashboardStatsUseCase) *DashboardHandler {
	return &DashboardHandler{uc: uc}
}

func (h *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.uc.Execute(r.Context()))
}

type ZonesHandler struct {
	uc *application.GetZoneMetricsUseCase
}

func NewZonesHandler(uc *application.GetZoneMetricsUseCase) *ZonesHandler {
	return &ZonesHandler{uc: uc}
}

func (h *ZonesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.uc.Execute(r.Context()))
}

type IncidentsHandler struct {
	uc *application.GetIncidentsUseCase
}

func NewIncidentsHandler(uc *application.GetIncidentsUseCase) *IncidentsHandler {
	return &IncidentsHandler{uc: uc}
}

func (h *IncidentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.uc.Execute(r.Context()))
}

type ForecastHandler struct {
	uc *application.GetForecastUseCase
}

func NewForecastHandler(uc *application.GetForecastUseCase) *ForecastHandler {
	return &ForecastHandler{uc: uc}
}

func (h *ForecastHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.uc.Execute(r.Context()))
}

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok", "service": "analytics", "version": "1.0.0",
	})
}

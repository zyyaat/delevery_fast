// Package handlers contains the HTTP handlers for the Auth Service.
package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/food-platform/shared/errors"
	"github.com/food-platform/shared/logging"
	"github.com/food-platform/auth/internal/application"
	"github.com/food-platform/auth/internal/domain"
	"github.com/google/uuid"
)

// ============ Helpers ============

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeError writes an error response.
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

// decodeJSON decodes a JSON body into the given destination.
func decodeJSON(r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return errors.ErrInvalidJSON
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

// extractIPAddress extracts the client IP address from the request.
// It handles X-Forwarded-For and X-Real-IP headers (set by API Gateway/Proxy).
func extractIPAddress(r *http.Request) string {
	// X-Forwarded-For: client, proxy1, proxy2
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	// Fallback to RemoteAddr (strip port)
	idx := strings.LastIndex(r.RemoteAddr, ":")
	if idx > 0 {
		return r.RemoteAddr[:idx]
	}
	return r.RemoteAddr
}

// ============ Send OTP Handler ============

// SendOTPHandler handles POST /api/v1/auth/otp/send
type SendOTPHandler struct {
	useCase *application.SendOTPUseCase
}

// NewSendOTPHandler creates a new SendOTPHandler.
func NewSendOTPHandler(useCase *application.SendOTPUseCase) *SendOTPHandler {
	return &SendOTPHandler{useCase: useCase}
}

type sendOTPRequest struct {
	Phone string `json:"phone"`
	Role  string `json:"role"`
}

type sendOTPResponse struct {
	RequestID         string `json:"request_id"`
	ExpiresIn         int    `json:"expires_in"`
	AttemptsRemaining int    `json:"attempts_remaining"`
}

func (h *SendOTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req sendOTPRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, errors.ErrInvalidJSON)
		return
	}

	if req.Phone == "" {
		writeError(w, errors.ErrInvalidPhone.WithDetails(map[string]interface{}{"field": "phone"}))
		return
	}
	if req.Role == "" {
		writeError(w, errors.ErrValidation.WithDetails(map[string]interface{}{"field": "role", "message": "role is required"}))
		return
	}

	cmd := application.SendOTPCommand{
		Phone: req.Phone,
		Role:  domain.UserRole(req.Role),
	}

	result, err := h.useCase.Execute(r.Context(), cmd)
	if err != nil {
		logging.FromContext(r.Context()).Error("send_otp_failed",
			"phone", req.Phone,
			"role", req.Role,
			"error", err,
		)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, sendOTPResponse{
		RequestID:         result.RequestID,
		ExpiresIn:         result.ExpiresInSeconds,
		AttemptsRemaining: result.AttemptsRemaining,
	})
}

// ============ Verify OTP Handler ============

// VerifyOTPHandler handles POST /api/v1/auth/otp/verify
type VerifyOTPHandler struct {
	useCase *application.VerifyOTPUseCase
}

// NewVerifyOTPHandler creates a new VerifyOTPHandler.
func NewVerifyOTPHandler(useCase *application.VerifyOTPUseCase) *VerifyOTPHandler {
	return &VerifyOTPHandler{useCase: useCase}
}

type verifyOTPRequest struct {
	RequestID         string `json:"request_id"`
	Code              string `json:"code"`
	DeviceFingerprint string `json:"device_fingerprint"`
}

type userResponse struct {
	ID         string `json:"id"`
	Phone      string `json:"phone"`
	Email      string `json:"email,omitempty"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	TrustScore int    `json:"trust_score"`
}

type authResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"`
	TokenType    string       `json:"token_type"`
	User         userResponse `json:"user"`
}

func (h *VerifyOTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req verifyOTPRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, errors.ErrInvalidJSON)
		return
	}

	if req.RequestID == "" || req.Code == "" {
		writeError(w, errors.ErrValidation.WithDetails(map[string]interface{}{
			"message": "request_id and code are required",
		}))
		return
	}

	// Get device fingerprint from header or body
	deviceFingerprint := req.DeviceFingerprint
	if deviceFingerprint == "" {
		deviceFingerprint = r.Header.Get("X-Device-Fingerprint")
	}

	cmd := application.VerifyOTPCommand{
		RequestID:         req.RequestID,
		Code:              req.Code,
		DeviceFingerprint: deviceFingerprint,
		UserAgent:         r.UserAgent(),
		IPAddress:         extractIPAddress(r),
	}

	result, err := h.useCase.Execute(r.Context(), cmd)
	if err != nil {
		logging.FromContext(r.Context()).Error("verify_otp_failed",
			"request_id", req.RequestID,
			"error", err,
		)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, authResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		TokenType:    "Bearer",
		User: userResponse{
			ID:         result.User.ID().String(),
			Phone:      result.User.Phone(),
			Email:      result.User.Email(),
			Name:       result.User.Name(),
			Role:       string(result.User.Role()),
			TrustScore: result.User.TrustScore(),
		},
	})
}

// ============ Refresh Token Handler ============

// RefreshTokenHandler handles POST /api/v1/auth/refresh
type RefreshTokenHandler struct {
	useCase *application.RefreshTokenUseCase
}

// NewRefreshTokenHandler creates a new RefreshTokenHandler.
func NewRefreshTokenHandler(useCase *application.RefreshTokenUseCase) *RefreshTokenHandler {
	return &RefreshTokenHandler{useCase: useCase}
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *RefreshTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req refreshTokenRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, errors.ErrInvalidJSON)
		return
	}

	if req.RefreshToken == "" {
		writeError(w, errors.ErrRefreshTokenInvalid)
		return
	}

	cmd := application.RefreshTokenCommand{
		RefreshToken: req.RefreshToken,
	}

	result, err := h.useCase.Execute(r.Context(), cmd)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, authResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		TokenType:    "Bearer",
		User: userResponse{
			ID:         result.User.ID().String(),
			Phone:      result.User.Phone(),
			Email:      result.User.Email(),
			Name:       result.User.Name(),
			Role:       string(result.User.Role()),
			TrustScore: result.User.TrustScore(),
		},
	})
}

// ============ Logout Handler ============

// LogoutHandler handles POST /api/v1/auth/logout
type LogoutHandler struct {
	useCase *application.LogoutUseCase
}

// NewLogoutHandler creates a new LogoutHandler.
func NewLogoutHandler(useCase *application.LogoutUseCase) *LogoutHandler {
	return &LogoutHandler{useCase: useCase}
}

func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by Auth middleware)
	userIDStr := logging.GetUserID(r.Context())
	if userIDStr == "" {
		writeError(w, errors.ErrUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, errors.ErrUnauthorized)
		return
	}

	// Optional: revoke specific session via header
	var sessionID *uuid.UUID
	if sid := r.Header.Get("X-Session-ID"); sid != "" {
		if parsed, err := uuid.Parse(sid); err == nil {
			sessionID = &parsed
		}
	}

	cmd := application.LogoutCommand{
		UserID:    userID,
		SessionID: sessionID,
	}

	if err := h.useCase.Execute(r.Context(), cmd); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "logged_out",
	})
}

// ============ Health Handler ============

// HealthHandler returns the health status of the service.
type HealthHandler struct {
	version string
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{version: version}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "auth",
		"version": h.version,
	})
}

// ReadyHandler returns the readiness status (checks dependencies).
type ReadyHandler struct {
	dbCheck func(ctx interface{}) error
}

// NewReadyHandler creates a new ReadyHandler.
func NewReadyHandler() *ReadyHandler {
	return &ReadyHandler{}
}

func (h *ReadyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Check DB, Redis, Kafka
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ready",
	})
}

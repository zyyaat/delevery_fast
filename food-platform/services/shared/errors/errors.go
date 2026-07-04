// Package errors provides standardized error types used across all services.
// These errors map cleanly to HTTP status codes via the errors.ToHTTP function.
package errors

import (
        "errors"
        "fmt"
        "net/http"
)

// ============ Error Types ============

// Error is the base error type with a code, message, and HTTP status.
type Error struct {
        Code       string                 `json:"code"`
        Message    string                 `json:"message"`
        StatusCode int                    `json:"-"`
        Details    map[string]interface{} `json:"details,omitempty"`
}

func (e *Error) Error() string {
        if e.Details != nil {
                return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Details)
        }
        return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error { return nil }

// New creates a new Error with the given code, message, and status code.
func New(code, message string, statusCode int) *Error {
        return &Error{
                Code:       code,
                Message:    message,
                StatusCode: statusCode,
        }
}

// WithDetails attaches details to the error.
func (e *Error) WithDetails(details map[string]interface{}) *Error {
        e.Details = details
        return e
}

// ============ Common Errors ============

var (
        // 400 Bad Request
        ErrBadRequest         = New("BAD_REQUEST", "Invalid request", http.StatusBadRequest)
        ErrValidation         = New("VALIDATION_FAILED", "Validation failed", http.StatusBadRequest)
        ErrInvalidJSON        = New("INVALID_JSON", "Invalid JSON body", http.StatusBadRequest)
        ErrInvalidOTP         = New("AUTH_INVALID_OTP", "Invalid or expired OTP", http.StatusBadRequest)
        ErrInvalidPhone       = New("INVALID_PHONE", "Invalid phone number", http.StatusBadRequest)

        // 401 Unauthorized
        ErrUnauthorized       = New("UNAUTHORIZED", "Authentication required", http.StatusUnauthorized)
        ErrTokenExpired       = New("AUTH_TOKEN_EXPIRED", "Access token expired", http.StatusUnauthorized)
        ErrInvalidToken       = New("AUTH_INVALID_TOKEN", "Invalid token", http.StatusUnauthorized)
        ErrRefreshInvalid     = New("AUTH_REFRESH_INVALID", "Invalid or used refresh token", http.StatusUnauthorized)

        // 403 Forbidden
        ErrForbidden          = New("PERMISSION_DENIED", "Insufficient permissions", http.StatusForbidden)
        ErrBiometricRequired  = New("AUTH_BIOMETRIC_REQUIRED", "Biometric verification required", http.StatusForbidden)
        ErrDualApproval       = New("DUAL_APPROVAL_REQUIRED", "Dual approval required", http.StatusForbidden)

        // 404 Not Found
        ErrNotFound           = New("NOT_FOUND", "Resource not found", http.StatusNotFound)
        ErrUserNotFound       = New("USER_NOT_FOUND", "User not found", http.StatusNotFound)
        ErrSessionNotFound    = New("SESSION_NOT_FOUND", "Session not found", http.StatusNotFound)

        // Session errors
        ErrSessionExpired     = New("AUTH_SESSION_EXPIRED", "Session has expired", http.StatusUnauthorized)
        ErrRefreshTokenInvalid = New("AUTH_REFRESH_INVALID", "Invalid or used refresh token", http.StatusUnauthorized)

        // 409 Conflict
        ErrConflict           = New("CONFLICT", "Resource conflict", http.StatusConflict)
        ErrUserExists         = New("USER_EXISTS", "User already exists", http.StatusConflict)
        ErrDuplicateRequest   = New("DUPLICATE_REQUEST", "Idempotency key already used", http.StatusConflict)

        // 422 Unprocessable Entity
        ErrInvalidTransition  = New("INVALID_TRANSITION", "Invalid state transition", http.StatusUnprocessableEntity)

        // 429 Too Many Requests
        ErrRateLimited        = New("RATE_LIMIT_EXCEEDED", "Too many requests", http.StatusTooManyRequests)

        // 500 Internal Server Error
        ErrInternal           = New("INTERNAL_ERROR", "Internal server error", http.StatusInternalServerError)
        ErrDatabase           = New("DATABASE_ERROR", "Database error", http.StatusInternalServerError)
        ErrRedis              = New("REDIS_ERROR", "Cache error", http.StatusInternalServerError)
        ErrKafka              = New("KAFKA_ERROR", "Event bus error", http.StatusInternalServerError)

        // 503 Service Unavailable
        ErrUnavailable        = New("SERVICE_UNAVAILABLE", "Service temporarily unavailable", http.StatusServiceUnavailable)
)

// ============ Helpers ============

// Wrap wraps a standard error into our Error type with a new message.
func Wrap(err error, code, message string, statusCode int) *Error {
        return &Error{
                Code:       code,
                Message:    message,
                StatusCode: statusCode,
                Details:    map[string]interface{}{"cause": err.Error()},
        }
}

// As attempts to extract an *Error from any error chain.
func As(err error) (*Error, bool) {
        var e *Error
        if errors.As(err, &e) {
                return e, true
        }
        return nil, false
}

// ToHTTP converts any error to an HTTP status code.
// Returns 500 for unknown errors.
func ToHTTP(err error) (int, string, string) {
        if err == nil {
                return http.StatusOK, "", ""
        }
        if e, ok := As(err); ok {
                return e.StatusCode, e.Code, e.Message
        }
        return http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error"
}

// Is checks if an error matches a target error code.
func Is(err error, target *Error) bool {
        e, ok := As(err)
        if !ok {
                return false
        }
        return e.Code == target.Code
}

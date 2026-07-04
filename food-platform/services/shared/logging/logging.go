// Package logging provides structured logging utilities using log/slog.
package logging

import (
	"context"
	"log/slog"
	"os"
)

// contextKey is a private type for context keys in this package.
type contextKey string

const (
	// RequestIDKey is the context key for request IDs.
	RequestIDKey contextKey = "request_id"
	// UserIDKey is the context key for user IDs.
	UserIDKey contextKey = "user_id"
	// SessionIDKey is the context key for session IDs.
	SessionIDKey contextKey = "session_id"
)

// Setup configures the global slog logger with the given level.
// Level can be: "debug", "info", "warn", "error".
func Setup(level string) {
	var lvl slog.Level
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "info":
		lvl = slog.LevelInfo
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     lvl,
		AddSource: true,
	})
	slog.SetDefault(slog.New(handler))
}

// FromContext extracts request-scoped fields from context and returns a logger with them.
func FromContext(ctx context.Context) *slog.Logger {
	attrs := []slog.Attr{}

	if reqID, ok := ctx.Value(RequestIDKey).(string); ok && reqID != "" {
		attrs = append(attrs, slog.String("request_id", reqID))
	}
	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		attrs = append(attrs, slog.String("user_id", userID))
	}
	if sessionID, ok := ctx.Value(SessionIDKey).(string); ok && sessionID != "" {
		attrs = append(attrs, slog.String("session_id", sessionID))
	}

	if len(attrs) == 0 {
		return slog.Default()
	}
	return slog.Default().With(slog.Group("context", toAny(attrs)...))
}

// toAny converts []slog.Attr to []any for slog.Group.
func toAny(attrs []slog.Attr) []any {
	result := make([]any, len(attrs)*2)
	for i, attr := range attrs {
		result[i*2] = attr.Key
		result[i*2+1] = attr.Value.Any()
	}
	return result
}

// WithRequestID adds a request ID to the context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// WithUserID adds a user ID to the context.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// WithSessionID adds a session ID to the context.
func WithSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, SessionIDKey, sessionID)
}

// GetRequestID retrieves the request ID from context.
func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}

// GetUserID retrieves the user ID from context.
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

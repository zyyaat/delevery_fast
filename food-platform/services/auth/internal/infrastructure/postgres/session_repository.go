package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/food-platform/services/auth/internal/domain"
	"github.com/google/uuid"
)

// ============ Session Repository ============

// SessionRepository implements application.SessionRepository using PostgreSQL.
type SessionRepository struct {
	db *sql.DB
}

// NewSessionRepository creates a new SessionRepository.
func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Save inserts a new session into the database.
func (r *SessionRepository) Save(ctx context.Context, session *domain.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, refresh_token, device_fingerprint, user_agent, ip_address, expires_at, created_at, revoked_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	var revokedAt interface{}
	if session.RevokedAt() != nil {
		revokedAt = *session.RevokedAt()
	}

	_, err := r.db.ExecContext(ctx, query,
		session.ID(),
		session.UserID(),
		session.RefreshToken(),
		session.DeviceFingerprint(),
		session.UserAgent(),
		session.IPAddress(),
		session.ExpiresAt(),
		session.CreatedAt(),
		revokedAt,
	)
	if err != nil {
		return fmt.Errorf("session_repo.Save: %w", err)
	}
	return nil
}

// FindByID retrieves a session by ID.
func (r *SessionRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	query := `
		SELECT id, user_id, refresh_token, device_fingerprint, user_agent, ip_address, expires_at, created_at, revoked_at
		FROM sessions
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return scanSession(row)
}

// FindByRefreshToken retrieves a session by refresh token.
func (r *SessionRepository) FindByRefreshToken(ctx context.Context, token string) (*domain.Session, error) {
	query := `
		SELECT id, user_id, refresh_token, device_fingerprint, user_agent, ip_address, expires_at, created_at, revoked_at
		FROM sessions
		WHERE refresh_token = $1
	`
	row := r.db.QueryRowContext(ctx, query, token)
	return scanSession(row)
}

// Revoke marks a session as revoked.
func (r *SessionRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE sessions SET revoked_at = $2 WHERE id = $1 AND revoked_at IS NULL`
	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query, id, now)
	if err != nil {
		return fmt.Errorf("session_repo.Revoke: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrSessionNotFound
	}
	return nil
}

// RevokeAllForUser revokes all active sessions for a user.
func (r *SessionRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE sessions SET revoked_at = $2 WHERE user_id = $1 AND revoked_at IS NULL`
	now := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, query, userID, now)
	if err != nil {
		return fmt.Errorf("session_repo.RevokeAllForUser: %w", err)
	}
	return nil
}

func scanSession(s scanner) (*domain.Session, error) {
	var (
		id                uuid.UUID
		userID            uuid.UUID
		refreshToken      string
		deviceFingerprint string
		userAgent         string
		ipAddress         string
		expiresAt         time.Time
		createdAt         time.Time
		revokedAt         *time.Time
	)

	err := s.Scan(&id, &userID, &refreshToken, &deviceFingerprint, &userAgent, &ipAddress, &expiresAt, &createdAt, &revokedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrSessionNotFound
		}
		return nil, fmt.Errorf("scanSession: %w", err)
	}

	session := &domain.Session{
		id:                id,
		userID:            userID,
		refreshToken:      refreshToken,
		deviceFingerprint: deviceFingerprint,
		userAgent:         userAgent,
		ipAddress:         ipAddress,
		expiresAt:         expiresAt,
		createdAt:         createdAt,
		revokedAt:         revokedAt,
	}
	return session, nil
}

// ============ Refresh Token Repository ============

// RefreshTokenRepository implements application.RefreshTokenRepository using PostgreSQL.
type RefreshTokenRepository struct {
	db *sql.DB
}

// NewRefreshTokenRepository creates a new RefreshTokenRepository.
func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Save inserts a new refresh token.
func (r *RefreshTokenRepository) Save(ctx context.Context, rt *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, session_id, token, expires_at, created_at, used_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	var usedAt interface{}
	if rt.UsedAt() != nil {
		usedAt = *rt.UsedAt()
	}

	_, err := r.db.ExecContext(ctx, query,
		rt.ID(),
		rt.UserID(),
		rt.SessionID(),
		rt.Token(),
		rt.ExpiresAt(),
		rt.CreatedAt(),
		usedAt,
	)
	if err != nil {
		return fmt.Errorf("refresh_repo.Save: %w", err)
	}
	return nil
}

// FindByToken retrieves a refresh token by its string value.
func (r *RefreshTokenRepository) FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, session_id, token, expires_at, created_at, used_at
		FROM refresh_tokens
		WHERE token = $1
	`
	row := r.db.QueryRowContext(ctx, query, token)
	return scanRefreshToken(row)
}

// MarkUsed marks a refresh token as used.
func (r *RefreshTokenRepository) MarkUsed(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE refresh_tokens SET used_at = $2 WHERE id = $1 AND used_at IS NULL`
	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query, id, now)
	if err != nil {
		return fmt.Errorf("refresh_repo.MarkUsed: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrRefreshTokenUsed
	}
	return nil
}

func scanRefreshToken(s scanner) (*domain.RefreshToken, error) {
	var (
		id        uuid.UUID
		userID    uuid.UUID
		sessionID uuid.UUID
		token     string
		expiresAt time.Time
		createdAt time.Time
		usedAt    *time.Time
	)

	err := s.Scan(&id, &userID, &sessionID, &token, &expiresAt, &createdAt, &usedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrRefreshTokenInvalid
		}
		return nil, fmt.Errorf("scanRefreshToken: %w", err)
	}

	return &domain.RefreshToken{
		id:        id,
		userID:    userID,
		sessionID: sessionID,
		token:     token,
		expiresAt: expiresAt,
		createdAt: createdAt,
		usedAt:    usedAt,
	}, nil
}

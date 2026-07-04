package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/food-platform/services/auth/internal/domain"
	"github.com/google/uuid"
)

// ============ OTP Repository ============

// OTPRepository implements application.OTPRepository using PostgreSQL.
type OTPRepository struct {
	db *sql.DB
}

// NewOTPRepository creates a new OTPRepository.
func NewOTPRepository(db *sql.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

// Save inserts or updates an OTP.
func (r *OTPRepository) Save(ctx context.Context, otp *domain.OTP) error {
	query := `
		INSERT INTO otps (id, phone, code, status, attempts_used, max_attempts, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			attempts_used = EXCLUDED.attempts_used
	`
	_, err := r.db.ExecContext(ctx, query,
		otp.ID(),
		otp.Phone(),
		otp.Code(),
		string(otp.Status()),
		otp.AttemptsUsed(),
		otp.MaxAttempts(),
		otp.ExpiresAt(),
		otp.CreatedAt(),
	)
	if err != nil {
		return fmt.Errorf("otp_repo.Save: %w", err)
	}
	return nil
}

// FindByID retrieves an OTP by ID.
func (r *OTPRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.OTP, error) {
	query := `
		SELECT id, phone, code, status, attempts_used, max_attempts, expires_at, created_at
		FROM otps
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return scanOTP(row)
}

// FindLatestByPhone retrieves the most recent OTP for a phone number.
func (r *OTPRepository) FindLatestByPhone(ctx context.Context, phone string) (*domain.OTP, error) {
	query := `
		SELECT id, phone, code, status, attempts_used, max_attempts, expires_at, created_at
		FROM otps
		WHERE phone = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, domain.NormalizePhone(phone))
	return scanOTP(row)
}

func scanOTP(s scanner) (*domain.OTP, error) {
	var (
		id            uuid.UUID
		phone         string
		code          string
		status        string
		attemptsUsed  int
		maxAttempts   int
		expiresAt     time.Time
		createdAt     time.Time
	)

	err := s.Scan(&id, &phone, &code, &status, &attemptsUsed, &maxAttempts, &expiresAt, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("OTP not found")
		}
		return nil, fmt.Errorf("scanOTP: %w", err)
	}

	return reconstructOTP(id, phone, code, domain.OTPStatus(status), attemptsUsed, maxAttempts, expiresAt, createdAt), nil
}

func reconstructOTP(
	id uuid.UUID,
	phone, code string,
	status domain.OTPStatus,
	attemptsUsed, maxAttempts int,
	expiresAt, createdAt time.Time,
) *domain.OTP {
	return &domain.OTP{
		id:           id,
		phone:        phone,
		code:         code,
		status:       status,
		attemptsUsed: attemptsUsed,
		maxAttempts:  maxAttempts,
		expiresAt:    expiresAt,
		createdAt:    createdAt,
	}
}

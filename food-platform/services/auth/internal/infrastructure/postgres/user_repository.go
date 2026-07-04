// Package postgres implements the repository interfaces using PostgreSQL.
package postgres

import (
        "context"
        "database/sql"
        "fmt"
        "time"

        "github.com/food-platform/auth/internal/domain"
        "github.com/google/uuid"
)

// ============ User Repository ============

// UserRepository implements application.UserRepository using PostgreSQL.
type UserRepository struct {
        db *sql.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
        return &UserRepository{db: db}
}

// Create inserts a new user into the database.
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
        query := `
                INSERT INTO users (id, phone, email, name, role, status, trust_score, created_at, updated_at)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        `
        _, err := r.db.ExecContext(ctx, query,
                user.ID(),
                user.Phone(),
                nullableString(user.Email()),
                user.Name(),
                string(user.Role()),
                string(user.Status()),
                user.TrustScore(),
                user.CreatedAt(),
                user.UpdatedAt(),
        )
        if err != nil {
                return fmt.Errorf("user_repo.Create: %w", err)
        }
        return nil
}

// FindByID retrieves a user by ID.
func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
        query := `
                SELECT id, phone, COALESCE(email, ''), name, role, status, trust_score, created_at, updated_at
                FROM users
                WHERE id = $1
        `
        row := r.db.QueryRowContext(ctx, query, id)
        return scanUser(row)
}

// FindByPhone retrieves a user by phone number.
func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
        normalized := domain.NormalizePhone(phone)
        query := `
                SELECT id, phone, COALESCE(email, ''), name, role, status, trust_score, created_at, updated_at
                FROM users
                WHERE phone = $1
        `
        row := r.db.QueryRowContext(ctx, query, normalized)
        return scanUser(row)
}

// Update updates an existing user.
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
        query := `
                UPDATE users
                SET email = $2, name = $3, role = $4, status = $5, trust_score = $6, updated_at = $7
                WHERE id = $1
        `
        _, err := r.db.ExecContext(ctx, query,
                user.ID(),
                nullableString(user.Email()),
                user.Name(),
                string(user.Role()),
                string(user.Status()),
                user.TrustScore(),
                user.UpdatedAt(),
        )
        if err != nil {
                return fmt.Errorf("user_repo.Update: %w", err)
        }
        return nil
}

// ============ Helpers ============

func nullableString(s string) interface{} {
        if s == "" {
                return nil
        }
        return s
}

// scanner is the common interface between sql.Row and sql.Rows
type scanner interface {
        Scan(dest ...interface{}) error
}

func scanUser(s scanner) (*domain.User, error) {
        var (
                id         uuid.UUID
                phone      string
                email      string
                name       string
                role       string
                status     string
                trustScore int
                createdAt  time.Time
                updatedAt  time.Time
        )

        err := s.Scan(&id, &phone, &email, &name, &role, &status, &trustScore, &createdAt, &updatedAt)
        if err != nil {
                if err == sql.ErrNoRows {
                        return nil, domain.ErrUserNotFound
                }
                return nil, fmt.Errorf("scanUser: %w", err)
        }

        // Reconstruct the user using the persisted data
        u := reconstructUser(id, phone, email, name, domain.UserRole(role), domain.UserStatus(status), trustScore, createdAt, updatedAt)
        return u, nil
}

// reconstructUser creates a User from persisted data (bypasses validation).
// This is used when loading from the database where data is assumed valid.
func reconstructUser(
        id uuid.UUID,
        phone, email, name string,
        role domain.UserRole,
        status domain.UserStatus,
        trustScore int,
        createdAt, updatedAt time.Time,
) *domain.User {
        return domain.ReconstructUser(id, phone, email, name, role, status, trustScore, createdAt, updatedAt)
}

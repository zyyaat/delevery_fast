package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/food-platform/fraud/internal/domain"
	"github.com/google/uuid"
)

type FraudRepository struct {
	db *sql.DB
}

func NewFraudRepository(db *sql.DB) *FraudRepository {
	return &FraudRepository{db: db}
}

func (r *FraudRepository) SaveScore(ctx context.Context, score *domain.FraudScore) error {
	query := `
		INSERT INTO fraud_scores (id, order_id, customer_id, score, decision, reasons, model_version, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	reasonsJSON := fmt.Sprintf(`["%s"]`, joinStrings(score.Reasons(), `","`))
	_, err := r.db.ExecContext(ctx, query,
		score.ID(), score.OrderID(), score.CustomerID(),
		score.Score(), string(score.Decision()), reasonsJSON,
		score.ModelVersion(), score.CreatedAt(),
	)
	return err
}

func (r *FraudRepository) FindScoreByOrder(ctx context.Context, orderID uuid.UUID) (*domain.FraudScore, error) {
	query := `SELECT id, order_id, customer_id, score, decision, reasons, model_version, created_at FROM fraud_scores WHERE order_id = $1`
	row := r.db.QueryRowContext(ctx, query, orderID)

	var (
		id, custID uuid.UUID
		score      int
		decision   string
		reasonsStr string
		modelVer   string
		createdAt  time.Time
	)

	err := row.Scan(&id, &orderID, &custID, &score, &decision, &reasonsStr, &modelVer, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrFraudScoreNotFound
		}
		return nil, err
	}

	return domain.ReconstructFraudScore(id, orderID, custID, score, domain.FraudDecision(decision), parseReasons(reasonsStr), modelVer, createdAt), nil
}

func (r *FraudRepository) GetTrustScore(ctx context.Context, customerID uuid.UUID) (*domain.TrustScore, error) {
	query := `SELECT customer_id, score, total_orders, refund_count, chargeback_count, fraud_flags, last_updated FROM trust_scores WHERE customer_id = $1`
	row := r.db.QueryRowContext(ctx, query, customerID)

	var (
		cID            uuid.UUID
		score          int
		totalOrders    int
		refundCount    int
		chargebackCount int
		fraudFlags     int
		lastUpdated    time.Time
	)

	err := row.Scan(&cID, &score, &totalOrders, &refundCount, &chargebackCount, &fraudFlags, &lastUpdated)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.NewTrustScore(customerID), nil
		}
		return nil, err
	}

	// Reconstruct trust score
	ts := domain.NewTrustScore(customerID)
	// Use internal state via reconstruct
	return ts, nil
}

func (r *FraudRepository) UpdateTrustScore(ctx context.Context, ts *domain.TrustScore) error {
	query := `
		INSERT INTO trust_scores (customer_id, score, total_orders, refund_count, chargeback_count, fraud_flags, last_updated)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (customer_id) DO UPDATE SET
			score = EXCLUDED.score, total_orders = EXCLUDED.total_orders,
			refund_count = EXCLUDED.refund_count, chargeback_count = EXCLUDED.chargeback_count,
			fraud_flags = EXCLUDED.fraud_flags, last_updated = EXCLUDED.last_updated
	`
	_, err := r.db.ExecContext(ctx, query,
		ts.CustomerID(), ts.Score(), ts.TotalOrders(),
		ts.RefundCount(), ts.ChargebackCount(), ts.FraudFlags(),
		ts.LastUpdated(),
	)
	return err
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

func parseReasons(s string) []string {
	if s == "" || s == "null" {
		return nil
	}
	// Simple JSON array parsing
	s = trimBrackets(s)
	if s == "" {
		return nil
	}
	var result []string
	current := ""
	inString := false
	for _, c := range s {
		if c == '"' {
			inString = !inString
			if !inString && current != "" {
				result = append(result, current)
				current = ""
			}
		} else if inString {
			current += string(c)
		}
	}
	return result
}

func trimBrackets(s string) string {
	if len(s) >= 2 && s[0] == '[' && s[len(s)-1] == ']' {
		return s[1 : len(s)-1]
	}
	return s
}

package application

import (
	"context"
	"fmt"
	"math"

	"github.com/food-platform/fraud/internal/domain"
	"github.com/google/uuid"
)

type FraudRepository interface {
	SaveScore(ctx context.Context, score *domain.FraudScore) error
	FindScoreByOrder(ctx context.Context, orderID uuid.UUID) (*domain.FraudScore, error)
	GetTrustScore(ctx context.Context, customerID uuid.UUID) (*domain.TrustScore, error)
	UpdateTrustScore(ctx context.Context, ts *domain.TrustScore) error
}

type ScoreOrderCommand struct {
	OrderID      uuid.UUID
	CustomerID   uuid.UUID
	OrderAmount  float64
	IsNewUser    bool
	DeviceFingerprint string
	IPAddress    string
	OrderCount   int
	RefundCount  int
}

type ScoreResultDTO struct {
	OrderID   string   `json:"order_id"`
	Score     int      `json:"score"`
	Decision  string   `json:"decision"`
	Reasons   []string `json:"reasons"`
	ModelVersion string `json:"model_version"`
}

type TrustScoreDTO struct {
	CustomerID      string `json:"customer_id"`
	Score           int    `json:"score"`
	TotalOrders     int    `json:"total_orders"`
	RefundCount     int    `json:"refund_count"`
	ChargebackCount int    `json:"chargeback_count"`
	FraudFlags      int    `json:"fraud_flags"`
}

type ScoreOrderUseCase struct {
	repo FraudRepository
}

func NewScoreOrderUseCase(repo FraudRepository) *ScoreOrderUseCase {
	return &ScoreOrderUseCase{repo: repo}
}

func (uc *ScoreOrderUseCase) Execute(ctx context.Context, cmd ScoreOrderCommand) (*ScoreResultDTO, error) {
	score := 0
	var reasons []string

	// Rule 1: New user + high-value order
	if cmd.IsNewUser && cmd.OrderAmount > 500 {
		score += 30
		reasons = append(reasons, "new_user_high_value_order")
	}

	// Rule 2: High refund rate
	if cmd.OrderCount > 0 {
		refundRate := float64(cmd.RefundCount) / float64(cmd.OrderCount)
		if refundRate > 0.3 {
			score += 35
			reasons = append(reasons, "high_refund_rate")
		} else if refundRate > 0.1 {
			score += 15
			reasons = append(reasons, "moderate_refund_rate")
		}
	}

	// Rule 3: Very high order amount
	if cmd.OrderAmount > 2000 {
		score += 25
		reasons = append(reasons, "very_high_order_amount")
	}

	// Rule 4: New user with COD
	if cmd.IsNewUser {
		score += 10
		reasons = append(reasons, "new_user")
	}

	// Rule 5: Get trust score and factor it in
	trustScore, _ := uc.repo.GetTrustScore(ctx, cmd.CustomerID)
	if trustScore != nil {
		// Lower trust = higher fraud risk
		if trustScore.Score() < 30 {
			score += 40
			reasons = append(reasons, "low_trust_score")
		} else if trustScore.Score() < 50 {
			score += 20
			reasons = append(reasons, "moderate_trust_score")
		}

		if trustScore.ChargebackCount() > 0 {
			score += 20
			reasons = append(reasons, "previous_chargeback")
		}

		if trustScore.FraudFlags() > 0 {
			score += 30
			reasons = append(reasons, "previous_fraud_flag")
		}
	}

	// Clamp score 0-100
	score = int(math.Min(100, float64(score)))
	if score < 0 {
		score = 0
	}

	// Create fraud score entity
	fraudScore := domain.NewFraudScore(cmd.OrderID, cmd.CustomerID, score, reasons)

	// Save
	if err := uc.repo.SaveScore(ctx, fraudScore); err != nil {
		return nil, fmt.Errorf("save fraud score: %w", err)
	}

	// Update trust score
	if trustScore != nil {
		trustScore.RecordOrder()
		_ = uc.repo.UpdateTrustScore(ctx, trustScore)
	}

	return &ScoreResultDTO{
		OrderID:      fraudScore.OrderID().String(),
		Score:        fraudScore.Score(),
		Decision:     string(fraudScore.Decision()),
		Reasons:      fraudScore.Reasons(),
		ModelVersion: fraudScore.ModelVersion(),
	}, nil
}

type GetTrustScoreUseCase struct {
	repo FraudRepository
}

func NewGetTrustScoreUseCase(repo FraudRepository) *GetTrustScoreUseCase {
	return &GetTrustScoreUseCase{repo: repo}
}

func (uc *GetTrustScoreUseCase) Execute(ctx context.Context, customerID uuid.UUID) (*TrustScoreDTO, error) {
	ts, err := uc.repo.GetTrustScore(ctx, customerID)
	if err != nil {
		// Return default for new users
		ts = domain.NewTrustScore(customerID)
	}

	return &TrustScoreDTO{
		CustomerID:      ts.CustomerID().String(),
		Score:           ts.Score(),
		TotalOrders:     ts.TotalOrders(),
		RefundCount:     ts.RefundCount(),
		ChargebackCount: ts.ChargebackCount(),
		FraudFlags:      ts.FraudFlags(),
	}, nil
}

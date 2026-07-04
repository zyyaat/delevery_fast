package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrFraudScoreNotFound = errors.New("fraud score not found")
)

type FraudDecision string

const (
	DecisionApprove FraudDecision = "approve"
	DecisionReview  FraudDecision = "review"
	DecisionBlock   FraudDecision = "block"
)

type FraudScore struct {
	id          uuid.UUID
	orderID     uuid.UUID
	customerID  uuid.UUID
	score       int // 0-100 (higher = more risky)
	decision    FraudDecision
	reasons     []string
	modelVersion string
	createdAt   time.Time
}

func NewFraudScore(orderID, customerID uuid.UUID, score int, reasons []string) *FraudScore {
	decision := DecisionApprove
	if score >= 70 {
		decision = DecisionBlock
	} else if score >= 40 {
		decision = DecisionReview
	}

	return &FraudScore{
		id:          uuid.New(),
		orderID:     orderID,
		customerID:  customerID,
		score:       score,
		decision:    decision,
		reasons:     reasons,
		modelVersion: "v1.0",
		createdAt:   time.Now().UTC(),
	}
}

func (f *FraudScore) ID() uuid.UUID        { return f.id }
func (f *FraudScore) OrderID() uuid.UUID   { return f.orderID }
func (f *FraudScore) CustomerID() uuid.UUID { return f.customerID }
func (f *FraudScore) Score() int           { return f.score }
func (f *FraudScore) Decision() FraudDecision { return f.decision }
func (f *FraudScore) Reasons() []string    { return f.reasons }
func (f *FraudScore) ModelVersion() string { return f.modelVersion }
func (f *FraudScore) CreatedAt() time.Time { return f.createdAt }

// TrustScore represents a customer's overall trustworthiness.
type TrustScore struct {
	customerID    uuid.UUID
	score         int // 0-100 (higher = more trustworthy)
	totalOrders   int
	refundCount   int
	chargebackCount int
	fraudFlags    int
	lastUpdated   time.Time
}

func NewTrustScore(customerID uuid.UUID) *TrustScore {
	return &TrustScore{
		customerID:  customerID,
		score:       50, // Default for new users
		lastUpdated: time.Now().UTC(),
	}
}

func (t *TrustScore) CustomerID() uuid.UUID { return t.customerID }
func (t *TrustScore) Score() int            { return t.score }
func (t *TrustScore) TotalOrders() int      { return t.totalOrders }
func (t *TrustScore) RefundCount() int      { return t.refundCount }
func (t *TrustScore) ChargebackCount() int  { return t.chargebackCount }
func (t *TrustScore) FraudFlags() int       { return t.fraudFlags }
func (t *TrustScore) LastUpdated() time.Time { return t.lastUpdated }

func (t *TrustScore) RecordOrder() {
	t.totalOrders++
	t.recalculate()
}

func (t *TrustScore) RecordRefund() {
	t.refundCount++
	t.recalculate()
}

func (t *TrustScore) RecordChargeback() {
	t.chargebackCount++
	t.score -= 15
	t.recalculate()
}

func (t *TrustScore) RecordFraudFlag() {
	t.fraudFlags++
	t.score -= 20
	t.recalculate()
}

func (t *TrustScore) recalculate() {
	// Base score from order history (more orders = more trust)
	if t.totalOrders > 0 {
		refundRate := float64(t.refundCount) / float64(t.totalOrders)
		if refundRate > 0.3 {
			t.score -= 20
		} else if refundRate > 0.1 {
			t.score -= 10
		}
	}

	// Penalize chargebacks and fraud flags
	t.score -= t.chargebackCount * 15
	t.score -= t.fraudFlags * 20

	// Clamp 0-100
	if t.score < 0 {
		t.score = 0
	} else if t.score > 100 {
		t.score = 100
	}

	t.lastUpdated = time.Now().UTC()
}

// ReconstructFraudScore creates a FraudScore from DB data.
func ReconstructFraudScore(
	id, orderID, customerID uuid.UUID,
	score int, decision FraudDecision, reasons []string,
	modelVersion string, createdAt time.Time,
) *FraudScore {
	return &FraudScore{
		id: id, orderID: orderID, customerID: customerID,
		score: score, decision: decision, reasons: reasons,
		modelVersion: modelVersion, createdAt: createdAt,
	}
}

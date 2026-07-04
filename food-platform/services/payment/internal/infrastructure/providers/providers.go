// Package providers implements payment provider integrations.
package providers

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/food-platform/payment/internal/domain"
)

// ============ Vodafone Cash Provider ============

// VodafoneCashProvider implements the Vodafone Cash merchant API.
// TODO: Replace mock with actual API calls.
type VodafoneCashProvider struct {
	merchantID string
	apiKey     string
	sandbox    bool
}

func NewVodafoneCashProvider(merchantID, apiKey string, sandbox bool) *VodafoneCashProvider {
	return &VodafoneCashProvider{merchantID: merchantID, apiKey: apiKey, sandbox: sandbox}
}

func (p *VodafoneCashProvider) Charge(ctx context.Context, payment *domain.Payment) (*domain.ProviderResult, error) {
	if p.sandbox {
		// Simulate success in sandbox mode
		time.Sleep(100 * time.Millisecond) // Simulate API latency

		txnID := fmt.Sprintf("VFC-%d-%s", time.Now().Unix(), payment.ID().String()[:8])
		redirectURL := ""

		// In sandbox, Vodafone Cash requires a redirect to their payment page
		if p.sandbox {
			redirectURL = fmt.Sprintf("https://sandbox.vodafone.com.eg/pay/%s", payment.ID())
		}

		slog.Info("vodafone_cash_charge_mock",
			"payment_id", payment.ID(),
			"amount", payment.Amount(),
			"txn_id", txnID,
		)

		return &domain.ProviderResult{
			TransactionID: txnID,
			RedirectURL:   redirectURL,
			Status:        domain.PaymentStatusCaptured,
		}, nil
	}

	// TODO: Implement real API call
	// POST https://api.vodafone.com.eg/merchant/charge
	// Body: { merchant_id, amount, customer_phone, order_id, callback_url }
	return nil, fmt.Errorf("vodafone cash production API not yet implemented")
}

func (p *VodafoneCashProvider) Refund(ctx context.Context, payment *domain.Payment, amount float64) error {
	if p.sandbox {
		slog.Info("vodafone_cash_refund_mock",
			"payment_id", payment.ID(),
			"amount", amount,
		)
		return nil
	}
	// TODO: Implement real refund API
	return fmt.Errorf("vodafone cash refund not yet implemented")
}

// ============ InstaPay Provider ============

// InstaPayProvider implements the InstaPay merchant API.
type InstaPayProvider struct {
	apiKey  string
	sandbox bool
}

func NewInstaPayProvider(apiKey string, sandbox bool) *InstaPayProvider {
	return &InstaPayProvider{apiKey: apiKey, sandbox: sandbox}
}

func (p *InstaPayProvider) Charge(ctx context.Context, payment *domain.Payment) (*domain.ProviderResult, error) {
	if p.sandbox {
		time.Sleep(150 * time.Millisecond)

		txnID := fmt.Sprintf("INSTA-%d-%s", time.Now().Unix(), payment.ID().String()[:8])
		redirectURL := fmt.Sprintf("https://sandbox.instapay.com/pay/%s", payment.ID())

		slog.Info("instapay_charge_mock", "payment_id", payment.ID(), "txn_id", txnID)

		return &domain.ProviderResult{
			TransactionID: txnID,
			RedirectURL:   redirectURL,
			Status:        domain.PaymentStatusCaptured,
		}, nil
	}

	return nil, fmt.Errorf("instapay production API not yet implemented")
}

func (p *InstaPayProvider) Refund(ctx context.Context, payment *domain.Payment, amount float64) error {
	if p.sandbox {
		slog.Info("instapay_refund_mock", "payment_id", payment.ID(), "amount", amount)
		return nil
	}
	return fmt.Errorf("instapay refund not yet implemented")
}

// ============ Paymob (Card) Provider ============

// PaymobProvider implements the Paymob (card processing) API.
type PaymobProvider struct {
	apiKey     string
	merchantID string
	sandbox    bool
}

func NewPaymobProvider(apiKey, merchantID string, sandbox bool) *PaymobProvider {
	return &PaymobProvider{apiKey: apiKey, merchantID: merchantID, sandbox: sandbox}
}

func (p *PaymobProvider) Charge(ctx context.Context, payment *domain.Payment) (*domain.ProviderResult, error) {
	if p.sandbox {
		time.Sleep(200 * time.Millisecond)

		txnID := fmt.Sprintf("CARD-%d-%s", time.Now().Unix(), payment.ID().String()[:8])
		redirectURL := fmt.Sprintf("https://sandbox.paymob.com/iframe/%s", payment.ID())

		// Simulate occasional decline (5% chance)
		if rand.Float64() < 0.05 {
			return &domain.ProviderResult{
				Status:       domain.PaymentStatusFailed,
				ErrorMessage: "card_declined",
			}, nil
		}

		slog.Info("paymob_charge_mock", "payment_id", payment.ID(), "txn_id", txnID)

		return &domain.ProviderResult{
			TransactionID: txnID,
			RedirectURL:   redirectURL,
			Status:        domain.PaymentStatusCaptured,
		}, nil
	}

	// TODO: Implement real Paymob API:
	// 1. POST /auth/tokens (get auth token)
	// 2. POST /ecommerce/orders (create order)
	// 3. POST /acceptance/payment_keys (get payment key)
	// 4. Redirect to iframe with payment key
	return nil, fmt.Errorf("paymob production API not yet implemented")
}

func (p *PaymobProvider) Refund(ctx context.Context, payment *domain.Payment, amount float64) error {
	if p.sandbox {
		slog.Info("paymob_refund_mock", "payment_id", payment.ID(), "amount", amount)
		return nil
	}
	// TODO: POST /acceptance/refund
	return fmt.Errorf("paymob refund not yet implemented")
}

// ============ Provider Factory ============

// Factory creates providers based on payment method.
type Factory struct {
	vodafoneCash *VodafoneCashProvider
	instaPay     *InstaPayProvider
	paymob       *PaymobProvider
}

func NewFactory(vodafone *VodafoneCashProvider, instaPay *InstaPayProvider, paymob *PaymobProvider) *Factory {
	return &Factory{vodafoneCash: vodafone, instaPay: instaPay, paymob: paymob}
}

func (f *Factory) GetProvider(method domain.PaymentMethod) (domain.Provider, error) {
	switch method {
	case domain.PaymentVodafoneCash:
		return f.vodafoneCash, nil
	case domain.PaymentInstaPay:
		return f.instaPay, nil
	case domain.PaymentCard:
		return f.paymob, nil
	case domain.PaymentCOD:
		return nil, nil // COD doesn't use a provider
	default:
		return nil, fmt.Errorf("unsupported payment method: %s", method)
	}
}

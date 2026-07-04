// Package jwt provides JWT token generation and validation.
package jwt

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/food-platform/services/auth/internal/domain"
	"github.com/google/uuid"
)

// ============ Generator ============

// Generator implements application.JWTGenerator using HMAC-SHA256.
// For production, use RS256 with a public/private key pair.
type Generator struct {
	secretKey   []byte
	issuer      string
	audience    string
	accessTTL   time.Duration
}

// Config holds the JWT generator configuration.
type Config struct {
	SecretKey string
	Issuer    string
	Audience  string
	AccessTTL time.Duration
}

// NewGenerator creates a new JWT Generator.
func NewGenerator(cfg Config) *Generator {
	if cfg.AccessTTL == 0 {
		cfg.AccessTTL = 15 * time.Minute
	}
	return &Generator{
		secretKey:   []byte(cfg.SecretKey),
		issuer:      cfg.Issuer,
		audience:    cfg.Audience,
		accessTTL:   cfg.AccessTTL,
	}
}

// Generate creates a new JWT access token for the given user.
// Returns the token string and its lifetime in seconds.
func (g *Generator) Generate(ctx context.Context, userID uuid.UUID, role domain.UserRole, sessionID uuid.UUID) (string, int, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(g.accessTTL)

	claims := Claims{
		Subject:   userID.String(),
		Issuer:    g.issuer,
		Audience:  g.audience,
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt.Unix(),
		Role:      string(role),
		SessionID: sessionID.String(),
		JWTID:     uuid.New().String(),
	}

	token, err := g.sign(claims)
	if err != nil {
		return "", 0, fmt.Errorf("jwt.Generate: %w", err)
	}

	return token, int(g.accessTTL.Seconds()), nil
}

// Claims represents the JWT claims.
type Claims struct {
	Subject   string `json:"sub"`
	Issuer    string `json:"iss"`
	Audience  string `json:"aud"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
	Role      string `json:"role"`
	SessionID string `json:"session_id"`
	JWTID     string `json:"jti"`
}

// sign creates a signed JWT token from claims.
func (g *Generator) sign(claims Claims) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsJSON)

	signingInput := encodedHeader + "." + encodedClaims
	signature := g.hmacSHA256(signingInput)

	return signingInput + "." + signature, nil
}

// hmacSHA256 creates an HMAC-SHA256 signature.
func (g *Generator) hmacSHA256(data string) string {
	h := hmac.New(sha256.New, g.secretKey)
	h.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// ============ Validator ============

// Validator validates JWT tokens.
type Validator struct {
	secretKey []byte
	issuer    string
	audience  string
}

// NewValidator creates a new JWT Validator.
func NewValidator(cfg Config) *Validator {
	return &Validator{
		secretKey: []byte(cfg.SecretKey),
		issuer:    cfg.Issuer,
		audience:  cfg.Audience,
	}
}

// Validate parses and validates a JWT token, returning the claims if valid.
func (v *Validator) Validate(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	signingInput := parts[0] + "." + parts[1]
	expectedSignature := v.hmacSHA256(signingInput)

	if !hmac.Equal([]byte(parts[2]), []byte(expectedSignature)) {
		return nil, fmt.Errorf("invalid signature")
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid claims: %w", err)
	}

	var claims Claims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, fmt.Errorf("invalid claims JSON: %w", err)
	}

	// Check expiration
	if time.Now().UTC().Unix() >= claims.ExpiresAt {
		return nil, fmt.Errorf("token expired")
	}

	// Check issuer
	if v.issuer != "" && claims.Issuer != v.issuer {
		return nil, fmt.Errorf("invalid issuer")
	}

	return &claims, nil
}

func (v *Validator) hmacSHA256(data string) string {
	h := hmac.New(sha256.New, v.secretKey)
	h.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

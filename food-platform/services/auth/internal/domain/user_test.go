package domain

import (
        "testing"
        "time"

        "github.com/google/uuid"
)

// ============ Phone Validation Tests ============

func TestValidatePhone_ValidNumbers(t *testing.T) {
        validPhones := []string{
                "01012345678",
                "01112345678",
                "01212345678",
                "01512345678",
        }

        for _, phone := range validPhones {
                t.Run("valid_"+phone, func(t *testing.T) {
                        err := ValidatePhone(phone)
                        if err != nil {
                                t.Errorf("expected nil, got %v for phone %s", err, phone)
                        }
                })
        }
}

func TestValidatePhone_InvalidNumbers(t *testing.T) {
        invalidPhones := []string{
                "",
                "123",
                "0101234567",    // too short
                "010123456789",  // too long
                "02012345678",   // invalid prefix
                "03012345678",   // invalid prefix
                "abc12345678",   // non-numeric
                "+2010123456",   // wrong length after normalization handled elsewhere
        }

        for _, phone := range invalidPhones {
                t.Run("invalid_"+phone, func(t *testing.T) {
                        err := ValidatePhone(phone)
                        if err == nil {
                                t.Errorf("expected error, got nil for phone %s", phone)
                        }
                })
        }
}

func TestNormalizePhone(t *testing.T) {
        tests := []struct {
                input    string
                expected string
        }{
                {"01012345678", "01012345678"},
                {"+201012345678", "01012345678"},
                {"201012345678", "01012345678"},
                {"010-1234-5678", "01012345678"},
                {" 01012345678 ", "01012345678"},
        }

        for _, tt := range tests {
                t.Run("normalize_"+tt.input, func(t *testing.T) {
                        result := NormalizePhone(tt.input)
                        if result != tt.expected {
                                t.Errorf("NormalizePhone(%q) = %q, want %q", tt.input, result, tt.expected)
                        }
                })
        }
}

// ============ User Tests ============

func TestNewUser_Valid(t *testing.T) {
        user, err := NewUser("01012345678", "Ahmed Mohamed", RoleCustomer)
        if err != nil {
                t.Fatalf("expected nil error, got %v", err)
        }
        if user == nil {
                t.Fatal("expected user, got nil")
        }
        if user.Phone() != "01012345678" {
                t.Errorf("expected phone 01012345678, got %s", user.Phone())
        }
        if user.Name() != "Ahmed Mohamed" {
                t.Errorf("expected name Ahmed Mohamed, got %s", user.Name())
        }
        if user.Role() != RoleCustomer {
                t.Errorf("expected role customer, got %s", user.Role())
        }
        if user.Status() != UserStatusActive {
                t.Errorf("expected status active, got %s", user.Status())
        }
        if user.TrustScore() != 50 {
                t.Errorf("expected trust score 50, got %d", user.TrustScore())
        }
        if user.ID().String() == "" {
                t.Error("expected non-empty ID")
        }
}

func TestNewUser_InvalidPhone(t *testing.T) {
        _, err := NewUser("invalid", "Ahmed", RoleCustomer)
        if err != ErrInvalidPhone {
                t.Errorf("expected ErrInvalidPhone, got %v", err)
        }
}

func TestNewUser_EmptyName(t *testing.T) {
        _, err := NewUser("01012345678", "", RoleCustomer)
        if err == nil {
                t.Error("expected error for empty name")
        }
}

func TestNewUser_EmptyRole(t *testing.T) {
        _, err := NewUser("01012345678", "Ahmed", "")
        if err == nil {
                t.Error("expected error for empty role")
        }
}

func TestNewUser_NormalizesPhone(t *testing.T) {
        user, err := NewUser("+201012345678", "Ahmed", RoleCustomer)
        if err != nil {
                t.Fatalf("expected nil error, got %v", err)
        }
        if user.Phone() != "01012345678" {
                t.Errorf("expected normalized phone, got %s", user.Phone())
        }
}

func TestUser_Suspend(t *testing.T) {
        user, _ := NewUser("01012345678", "Ahmed", RoleCustomer)

        if err := user.Suspend(); err != nil {
                t.Errorf("expected nil error, got %v", err)
        }
        if user.Status() != UserStatusSuspended {
                t.Errorf("expected suspended, got %s", user.Status())
        }
        if user.IsActive() {
                t.Error("expected IsActive() to be false")
        }
}

func TestUser_Reactivate(t *testing.T) {
        user, _ := NewUser("01012345678", "Ahmed", RoleCustomer)
        user.Suspend()

        if err := user.Reactivate(); err != nil {
                t.Errorf("expected nil error, got %v", err)
        }
        if user.Status() != UserStatusActive {
                t.Errorf("expected active, got %s", user.Status())
        }
        if !user.IsActive() {
                t.Error("expected IsActive() to be true")
        }
}

func TestUser_Delete(t *testing.T) {
        user, _ := NewUser("01012345678", "Ahmed", RoleCustomer)

        if err := user.Delete(); err != nil {
                t.Errorf("expected nil error, got %v", err)
        }
        if user.Status() != UserStatusDeleted {
                t.Errorf("expected deleted, got %s", user.Status())
        }

        // Cannot suspend or reactivate a deleted user
        if err := user.Suspend(); err != ErrInvalidTransition {
                t.Errorf("expected ErrInvalidTransition, got %v", err)
        }
        if err := user.Reactivate(); err != ErrInvalidTransition {
                t.Errorf("expected ErrInvalidTransition, got %v", err)
        }
}

func TestUser_SetTrustScore(t *testing.T) {
        user, _ := NewUser("01012345678", "Ahmed", RoleCustomer)

        tests := []struct {
                input    int
                expected int
        }{
                {0, 0},
                {50, 50},
                {100, 100},
                {-10, 0},   // clamp to 0
                {150, 100}, // clamp to 100
        }

        for _, tt := range tests {
                user.SetTrustScore(tt.input)
                if user.TrustScore() != tt.expected {
                        t.Errorf("SetTrustScore(%d) = %d, want %d", tt.input, user.TrustScore(), tt.expected)
                }
        }
}

func TestUser_SetEmail(t *testing.T) {
        user, _ := NewUser("01012345678", "Ahmed", RoleCustomer)

        user.SetEmail("ahmed@example.com")
        if user.Email() != "ahmed@example.com" {
                t.Errorf("expected email ahmed@example.com, got %s", user.Email())
        }
}

// ============ OTP Tests ============

func TestNewOTP(t *testing.T) {
        otp := NewOTP("01012345678", "123456")

        if otp.Phone() != "01012345678" {
                t.Errorf("expected phone 01012345678, got %s", otp.Phone())
        }
        if otp.Code() != "123456" {
                t.Errorf("expected code 123456, got %s", otp.Code())
        }
        if otp.Status() != OTPStatusPending {
                t.Errorf("expected status pending, got %s", otp.Status())
        }
        if otp.AttemptsUsed() != 0 {
                t.Errorf("expected 0 attempts used, got %d", otp.AttemptsUsed())
        }
        if otp.MaxAttempts() != 3 {
                t.Errorf("expected max 3 attempts, got %d", otp.MaxAttempts())
        }
        if !otp.CanVerify() {
                t.Error("expected CanVerify() to be true for new OTP")
        }
}

func TestOTP_Verify_CorrectCode(t *testing.T) {
        otp := NewOTP("01012345678", "123456")

        err := otp.Verify("123456")
        if err != nil {
                t.Errorf("expected nil error, got %v", err)
        }
        if otp.Status() != OTPStatusVerified {
                t.Errorf("expected status verified, got %s", otp.Status())
        }
}

func TestOTP_Verify_IncorrectCode(t *testing.T) {
        otp := NewOTP("01012345678", "123456")

        err := otp.Verify("000000")
        if err != ErrInvalidOTP {
                t.Errorf("expected ErrInvalidOTP, got %v", err)
        }
        if otp.AttemptsUsed() != 1 {
                t.Errorf("expected 1 attempt used, got %d", otp.AttemptsUsed())
        }
}

func TestOTP_Verify_AttemptsExceeded(t *testing.T) {
        otp := NewOTP("01012345678", "123456")

        // Use all 3 attempts
        otp.Verify("000000")
        otp.Verify("000000")
        otp.Verify("000000")

        // 4th attempt should fail with attempts exceeded
        err := otp.Verify("123456")
        if err != ErrOTPAttemptsExceeded && err != ErrInvalidOTP {
                t.Errorf("expected ErrOTPAttemptsExceeded or ErrInvalidOTP, got %v", err)
        }
}

func TestOTP_Verify_AfterVerified(t *testing.T) {
        otp := NewOTP("01012345678", "123456")

        // Verify successfully
        if err := otp.Verify("123456"); err != nil {
                t.Fatalf("expected nil error, got %v", err)
        }

        // Second verification should fail
        err := otp.Verify("123456")
        if err == nil {
                t.Error("expected error when verifying already-verified OTP")
        }
}

func TestOTP_IsExpired(t *testing.T) {
        otp := NewOTP("01012345678", "123456")
        if otp.IsExpired() {
                t.Error("expected new OTP to not be expired")
        }

        // Manually set expiry in the past (we can't do this directly, so we test the logic indirectly)
        // In practice, we'd need to wait 2 minutes or use a constructor that accepts a TTL
}

// ============ Session Tests ============

func TestNewSession(t *testing.T) {
        userID := uuid.New()
        session := NewSession(userID, "refresh-token-123", "device-fp", "Mozilla/5.0", "192.168.1.1", 24*time.Hour)

        if session.UserID() != userID {
                t.Error("expected user ID to match")
        }
        if session.RefreshToken() != "refresh-token-123" {
                t.Errorf("expected refresh token, got %s", session.RefreshToken())
        }
        if session.DeviceFingerprint() != "device-fp" {
                t.Errorf("expected device fingerprint, got %s", session.DeviceFingerprint())
        }
        if !session.IsActive() {
                t.Error("expected new session to be active")
        }
        if session.IsExpired() {
                t.Error("expected new session to not be expired")
        }
        if session.IsRevoked() {
                t.Error("expected new session to not be revoked")
        }
}

func TestSession_Revoke(t *testing.T) {
        session := NewSession(uuid.New(), "token", "fp", "ua", "ip", 24*time.Hour)

        session.Revoke()

        if !session.IsRevoked() {
                t.Error("expected session to be revoked")
        }
        if session.IsActive() {
                t.Error("expected session to not be active after revoke")
        }
        if session.RevokedAt() == nil {
                t.Error("expected RevokedAt to be set")
        }
}

func TestSession_IsExpired(t *testing.T) {
        // Create a session that's already expired
        session := NewSession(uuid.New(), "token", "fp", "ua", "ip", -1*time.Hour)

        if !session.IsExpired() {
                t.Error("expected session to be expired")
        }
        if session.IsActive() {
                t.Error("expected expired session to not be active")
        }
}

// ============ Refresh Token Tests ============

func TestNewRefreshToken(t *testing.T) {
        userID := uuid.New()
        sessionID := uuid.New()
        rt := NewRefreshToken(userID, sessionID, "token-123", 24*time.Hour)

        if rt.UserID() != userID {
                t.Error("expected user ID to match")
        }
        if rt.SessionID() != sessionID {
                t.Error("expected session ID to match")
        }
        if rt.Token() != "token-123" {
                t.Errorf("expected token, got %s", rt.Token())
        }
        if !rt.IsValid() {
                t.Error("expected new token to be valid")
        }
        if rt.IsUsed() {
                t.Error("expected new token to not be used")
        }
}

func TestRefreshToken_Use(t *testing.T) {
        rt := NewRefreshToken(uuid.New(), uuid.New(), "token", 24*time.Hour)

        if err := rt.Use(); err != nil {
                t.Errorf("expected nil error, got %v", err)
        }
        if !rt.IsUsed() {
                t.Error("expected token to be used")
        }
        if rt.IsValid() {
                t.Error("expected used token to not be valid")
        }

        // Using again should fail
        if err := rt.Use(); err != ErrRefreshTokenUsed {
                t.Errorf("expected ErrRefreshTokenUsed, got %v", err)
        }
}

func TestRefreshToken_IsExpired(t *testing.T) {
        rt := NewRefreshToken(uuid.New(), uuid.New(), "token", -1*time.Hour)

        if !rt.IsExpired() {
                t.Error("expected token to be expired")
        }
        if rt.IsValid() {
                t.Error("expected expired token to not be valid")
        }
}

// ============ Access Token Tests ============

func TestNewAccessToken(t *testing.T) {
        at := NewAccessToken("jwt-token", 900)

        if at.Token() != "jwt-token" {
                t.Errorf("expected token, got %s", at.Token())
        }
        if at.ExpiresIn() != 900 {
                t.Errorf("expected 900, got %d", at.ExpiresIn())
        }
}

// ============ Auth Result Tests ============

func TestNewAuthResult(t *testing.T) {
        user, _ := NewUser("01012345678", "Ahmed", RoleCustomer)
        at := NewAccessToken("token", 900)
        result := NewAuthResult(at, "refresh-token", user)

        if result.AccessToken().Token() != "token" {
                t.Error("expected access token to match")
        }
        if result.RefreshToken() != "refresh-token" {
                t.Errorf("expected refresh token, got %s", result.RefreshToken())
        }
        if result.User().Name() != "Ahmed" {
                t.Error("expected user to match")
        }
}

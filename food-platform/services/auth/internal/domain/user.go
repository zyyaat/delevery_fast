// Package domain contains the core business logic of the Auth Service.
// It has no external dependencies — pure Go types and functions.
package domain

import (
        "errors"
        "regexp"
        "time"

        "github.com/google/uuid"
)

// ============ Errors ============

var (
        ErrUserNotFound       = errors.New("user not found")
        ErrUserExists         = errors.New("user already exists")
        ErrUserInactive       = errors.New("user is inactive")
        ErrInvalidPhone       = errors.New("invalid phone number")
        ErrInvalidOTP         = errors.New("invalid or expired OTP")
        ErrOTPExpired         = errors.New("OTP has expired")
        ErrOTPAttemptsExceeded = errors.New("OTP attempts exceeded")
        ErrSessionNotFound    = errors.New("session not found")
        ErrSessionExpired     = errors.New("session has expired")
        ErrRefreshTokenInvalid = errors.New("refresh token is invalid")
        ErrRefreshTokenUsed   = errors.New("refresh token has already been used")
        ErrInvalidTransition  = errors.New("invalid state transition")
)

// ============ Enums ============

// UserRole represents the role of a user in the system.
type UserRole string

const (
        RoleCustomer         UserRole = "customer"
        RoleDriver           UserRole = "driver"
        RoleRestaurant       UserRole = "restaurant"
        RoleSupportL1        UserRole = "support_l1"
        RoleSupportL2        UserRole = "support_l2"
        RoleOpsManager       UserRole = "ops_manager"
        RoleFinance          UserRole = "finance"
        RoleSuperAdmin       UserRole = "super_admin"
        RoleFieldSupervisor  UserRole = "field_supervisor"
        RoleHR               UserRole = "hr"
        RoleReadOnlyAnalyst  UserRole = "read_only_analyst"
)

// UserStatus represents the account status of a user.
type UserStatus string

const (
        UserStatusActive    UserStatus = "active"
        UserStatusSuspended UserStatus = "suspended"
        UserStatusDeleted   UserStatus = "deleted"
)

// OTPStatus represents the state of an OTP code.
type OTPStatus string

const (
        OTPStatusPending  OTPStatus = "pending"
        OTPStatusVerified OTPStatus = "verified"
        OTPStatusExpired  OTPStatus = "expired"
        OTPStatusUsed     OTPStatus = "used"
)

// ============ Validation ============

// egyptianPhoneRegex matches valid Egyptian mobile numbers.
// Examples: 01012345678, 01112345678, 01212345678, 01512345678
var egyptianPhoneRegex = regexp.MustCompile(`^01[0-2,5][0-9]{8}$`)

// ValidatePhone validates an Egyptian phone number.
func ValidatePhone(phone string) error {
        if !egyptianPhoneRegex.MatchString(phone) {
                return ErrInvalidPhone
        }
        return nil
}

// NormalizePhone removes the country code and returns the 11-digit number.
// Examples:
//   "+201012345678" → "01012345678"
//   "201012345678"  → "01012345678"
//   "01012345678"   → "01012345678"
func NormalizePhone(phone string) string {
        // Remove all non-digits
        digits := make([]byte, 0, len(phone))
        for _, c := range phone {
                if c >= '0' && c <= '9' {
                        digits = append(digits, byte(c))
                }
        }

        s := string(digits)

        // Handle "+20" or "20" prefix
        // "+201012345678" → digits "201012345678" (12) → "0" + "1012345678" = "01012345678"
        if len(s) == 12 && s[:2] == "20" {
                return "0" + s[2:]
        }
        // "+2010123456789" (extra digit) → 13 digits, also handle
        if len(s) == 13 && s[:3] == "201" {
                return "0" + s[3:]
        }

        return s
}

// ============ Entities ============

// User represents an authenticated user in the system.
type User struct {
        id          uuid.UUID
        phone       string
        email       string
        name        string
        role        UserRole
        status      UserStatus
        trustScore  int
        createdAt   time.Time
        updatedAt   time.Time
}

// NewUser creates a new User with the given parameters.
// Validates the phone number and applies defaults.
func NewUser(phone, name string, role UserRole) (*User, error) {
        normalized := NormalizePhone(phone)
        if err := ValidatePhone(normalized); err != nil {
                return nil, err
        }

        if name == "" {
                return nil, errors.New("name is required")
        }

        if role == "" {
                return nil, errors.New("role is required")
        }

        now := time.Now().UTC()
        return &User{
                id:         uuid.New(),
                phone:      normalized,
                name:       name,
                role:       role,
                status:     UserStatusActive,
                trustScore: 50, // Default trust score for new users
                createdAt:  now,
                updatedAt:  now,
        }, nil
}

// ID returns the user's ID.
func (u *User) ID() uuid.UUID { return u.id }

// Phone returns the user's phone number (normalized).
func (u *User) Phone() string { return u.phone }

// Email returns the user's email address.
func (u *User) Email() string { return u.email }

// SetEmail sets the user's email address.
func (u *User) SetEmail(email string) {
        u.email = email
        u.updatedAt = time.Now().UTC()
}

// Name returns the user's name.
func (u *User) Name() string { return u.name }

// SetName sets the user's name.
func (u *User) SetName(name string) {
        u.name = name
        u.updatedAt = time.Now().UTC()
}

// Role returns the user's role.
func (u *User) Role() UserRole { return u.role }

// SetRole sets the user's role.
func (u *User) SetRole(role UserRole) {
        u.role = role
        u.updatedAt = time.Now().UTC()
}

// Status returns the user's account status.
func (u *User) Status() UserStatus { return u.status }

// IsActive returns true if the user is active.
func (u *User) IsActive() bool { return u.status == UserStatusActive }

// Suspend sets the user's status to suspended.
func (u *User) Suspend() error {
        if u.status == UserStatusDeleted {
                return ErrInvalidTransition
        }
        u.status = UserStatusSuspended
        u.updatedAt = time.Now().UTC()
        return nil
}

// Reactivate sets the user's status to active.
func (u *User) Reactivate() error {
        if u.status == UserStatusDeleted {
                return ErrInvalidTransition
        }
        u.status = UserStatusActive
        u.updatedAt = time.Now().UTC()
        return nil
}

// Delete sets the user's status to deleted (soft delete).
func (u *User) Delete() error {
        if u.status == UserStatusDeleted {
                return ErrInvalidTransition
        }
        u.status = UserStatusDeleted
        u.updatedAt = time.Now().UTC()
        return nil
}

// TrustScore returns the user's trust score (0-100).
func (u *User) TrustScore() int { return u.trustScore }

// SetTrustScore sets the user's trust score (0-100).
func (u *User) SetTrustScore(score int) {
        if score < 0 {
                score = 0
        } else if score > 100 {
                score = 100
        }
        u.trustScore = score
        u.updatedAt = time.Now().UTC()
}

// CreatedAt returns the user's creation timestamp.
func (u *User) CreatedAt() time.Time { return u.createdAt }

// UpdatedAt returns the user's last update timestamp.
func (u *User) UpdatedAt() time.Time { return u.updatedAt }

// ============ OTP ============

// OTP represents a one-time password sent to a user's phone.
type OTP struct {
        id              uuid.UUID
        phone           string
        code            string
        status          OTPStatus
        attemptsUsed    int
        maxAttempts     int
        expiresAt       time.Time
        createdAt       time.Time
}

// NewOTP creates a new OTP with the given phone number and code.
// The OTP expires after 2 minutes and allows 3 attempts.
func NewOTP(phone, code string) *OTP {
        now := time.Now().UTC()
        return &OTP{
                id:          uuid.New(),
                phone:       phone,
                code:        code,
                status:      OTPStatusPending,
                attemptsUsed: 0,
                maxAttempts: 3,
                expiresAt:   now.Add(2 * time.Minute),
                createdAt:   now,
        }
}

// ID returns the OTP's ID.
func (o *OTP) ID() uuid.UUID { return o.id }

// Phone returns the phone number the OTP was sent to.
func (o *OTP) Phone() string { return o.phone }

// Code returns the OTP code.
func (o *OTP) Code() string { return o.code }

// Status returns the OTP's current status.
func (o *OTP) Status() OTPStatus { return o.status }

// AttemptsUsed returns the number of failed verification attempts.
func (o *OTP) AttemptsUsed() int { return o.attemptsUsed }

// MaxAttempts returns the maximum number of attempts allowed.
func (o *OTP) MaxAttempts() int { return o.maxAttempts }

// ExpiresAt returns the OTP's expiration time.
func (o *OTP) ExpiresAt() time.Time { return o.expiresAt }

// CreatedAt returns the OTP's creation time.
func (o *OTP) CreatedAt() time.Time { return o.createdAt }

// IsExpired returns true if the OTP has expired.
func (o *OTP) IsExpired() bool {
        return time.Now().UTC().After(o.expiresAt)
}

// CanVerify returns true if the OTP can still be verified (not expired, not used, attempts remaining).
func (o *OTP) CanVerify() bool {
        if o.status != OTPStatusPending {
                return false
        }
        if o.IsExpired() {
                o.status = OTPStatusExpired
                return false
        }
        if o.attemptsUsed >= o.maxAttempts {
                return false
        }
        return true
}

// Verify attempts to verify the OTP with the given code.
// Returns nil on success, or an error if verification fails.
func (o *OTP) Verify(code string) error {
        if !o.CanVerify() {
                if o.IsExpired() {
                        return ErrOTPExpired
                }
                if o.attemptsUsed >= o.maxAttempts {
                        return ErrOTPAttemptsExceeded
                }
                return ErrInvalidOTP
        }

        o.attemptsUsed++

        if o.code != code {
                if o.attemptsUsed >= o.maxAttempts {
                        o.status = OTPStatusExpired
                }
                return ErrInvalidOTP
        }

        o.status = OTPStatusVerified
        return nil
}

// ============ Session ============

// Session represents an authenticated user session.
type Session struct {
        id           uuid.UUID
        userID       uuid.UUID
        refreshToken string
        deviceFingerprint string
        userAgent    string
        ipAddress    string
        expiresAt    time.Time
        createdAt    time.Time
        revokedAt    *time.Time
}

// NewSession creates a new Session for the given user.
func NewSession(
        userID uuid.UUID,
        refreshToken string,
        deviceFingerprint string,
        userAgent string,
        ipAddress string,
        ttl time.Duration,
) *Session {
        now := time.Now().UTC()
        return &Session{
                id:           uuid.New(),
                userID:       userID,
                refreshToken: refreshToken,
                deviceFingerprint: deviceFingerprint,
                userAgent:    userAgent,
                ipAddress:    ipAddress,
                expiresAt:    now.Add(ttl),
                createdAt:    now,
        }
}

// ID returns the session's ID.
func (s *Session) ID() uuid.UUID { return s.id }

// UserID returns the user ID associated with the session.
func (s *Session) UserID() uuid.UUID { return s.userID }

// RefreshToken returns the session's refresh token.
func (s *Session) RefreshToken() string { return s.refreshToken }

// DeviceFingerprint returns the device fingerprint associated with the session.
func (s *Session) DeviceFingerprint() string { return s.deviceFingerprint }

// UserAgent returns the user agent string from the original request.
func (s *Session) UserAgent() string { return s.userAgent }

// IPAddress returns the IP address from the original request.
func (s *Session) IPAddress() string { return s.ipAddress }

// ExpiresAt returns the session's expiration time.
func (s *Session) ExpiresAt() time.Time { return s.expiresAt }

// CreatedAt returns the session's creation time.
func (s *Session) CreatedAt() time.Time { return s.createdAt }

// RevokedAt returns the time the session was revoked, or nil if active.
func (s *Session) RevokedAt() *time.Time { return s.revokedAt }

// IsExpired returns true if the session has expired.
func (s *Session) IsExpired() bool {
        return time.Now().UTC().After(s.expiresAt)
}

// IsRevoked returns true if the session has been revoked.
func (s *Session) IsRevoked() bool {
        return s.revokedAt != nil
}

// IsActive returns true if the session is active (not expired, not revoked).
func (s *Session) IsActive() bool {
        return !s.IsExpired() && !s.IsRevoked()
}

// Revoke marks the session as revoked.
func (s *Session) Revoke() {
        now := time.Now().UTC()
        s.revokedAt = &now
}

// ============ Refresh Token ============

// RefreshToken represents a refresh token used to obtain new access tokens.
type RefreshToken struct {
        id           uuid.UUID
        userID       uuid.UUID
        sessionID    uuid.UUID
        token        string
        expiresAt    time.Time
        createdAt    time.Time
        usedAt       *time.Time
}

// NewRefreshToken creates a new RefreshToken for the given user and session.
func NewRefreshToken(userID, sessionID uuid.UUID, token string, ttl time.Duration) *RefreshToken {
        now := time.Now().UTC()
        return &RefreshToken{
                id:        uuid.New(),
                userID:    userID,
                sessionID: sessionID,
                token:     token,
                expiresAt: now.Add(ttl),
                createdAt: now,
        }
}

// ID returns the refresh token's ID.
func (rt *RefreshToken) ID() uuid.UUID { return rt.id }

// UserID returns the user ID associated with the refresh token.
func (rt *RefreshToken) UserID() uuid.UUID { return rt.userID }

// SessionID returns the session ID associated with the refresh token.
func (rt *RefreshToken) SessionID() uuid.UUID { return rt.sessionID }

// Token returns the refresh token string.
func (rt *RefreshToken) Token() string { return rt.token }

// ExpiresAt returns the refresh token's expiration time.
func (rt *RefreshToken) ExpiresAt() time.Time { return rt.expiresAt }

// CreatedAt returns the refresh token's creation time.
func (rt *RefreshToken) CreatedAt() time.Time { return rt.createdAt }

// UsedAt returns the time the token was used, or nil if unused.
func (rt *RefreshToken) UsedAt() *time.Time { return rt.usedAt }

// IsExpired returns true if the refresh token has expired.
func (rt *RefreshToken) IsExpired() bool {
        return time.Now().UTC().After(rt.expiresAt)
}

// IsUsed returns true if the refresh token has been used.
func (rt *RefreshToken) IsUsed() bool {
        return rt.usedAt != nil
}

// IsValid returns true if the refresh token is valid (not expired, not used).
func (rt *RefreshToken) IsValid() bool {
        return !rt.IsExpired() && !rt.IsUsed()
}

// Use marks the refresh token as used.
// Returns an error if the token has already been used or is expired.
func (rt *RefreshToken) Use() error {
        if rt.IsUsed() {
                return ErrRefreshTokenUsed
        }
        if rt.IsExpired() {
                return ErrSessionExpired
        }
        now := time.Now().UTC()
        rt.usedAt = &now
        return nil
}

// ============ Access Token ============

// AccessToken represents a JWT access token.
type AccessToken struct {
        token     string
        expiresIn int // seconds
}

// NewAccessToken creates a new AccessToken.
func NewAccessToken(token string, expiresIn int) *AccessToken {
        return &AccessToken{
                token:     token,
                expiresIn: expiresIn,
        }
}

// Token returns the access token string.
func (at *AccessToken) Token() string { return at.token }

// ExpiresIn returns the token's lifetime in seconds.
func (at *AccessToken) ExpiresIn() int { return at.expiresIn }

// ============ Auth Result ============

// AuthResult represents the result of a successful authentication.
type AuthResult struct {
        accessToken  *AccessToken
        refreshToken string
        user         *User
}

// NewAuthResult creates a new AuthResult.
func NewAuthResult(accessToken *AccessToken, refreshToken string, user *User) *AuthResult {
        return &AuthResult{
                accessToken:  accessToken,
                refreshToken: refreshToken,
                user:         user,
        }
}

// AccessToken returns the access token.
func (a *AuthResult) AccessToken() *AccessToken { return a.accessToken }

// RefreshToken returns the refresh token string.
func (a *AuthResult) RefreshToken() string { return a.refreshToken }

// User returns the authenticated user.
func (a *AuthResult) User() *User { return a.user }

// ============ Reconstructors (for DB hydration) ============

// ReconstructOTP creates an OTP from persisted data (bypasses validation).
func ReconstructOTP(
        id uuid.UUID,
        phone, code string,
        status OTPStatus,
        attemptsUsed, maxAttempts int,
        expiresAt, createdAt time.Time,
) *OTP {
        return &OTP{
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

// ReconstructSession creates a Session from persisted data.
func ReconstructSession(
        id, userID uuid.UUID,
        refreshToken, deviceFingerprint, userAgent, ipAddress string,
        expiresAt, createdAt time.Time,
        revokedAt *time.Time,
) *Session {
        return &Session{
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
}

// ReconstructRefreshToken creates a RefreshToken from persisted data.
func ReconstructRefreshToken(
        id, userID, sessionID uuid.UUID,
        token string,
        expiresAt, createdAt time.Time,
        usedAt *time.Time,
) *RefreshToken {
        return &RefreshToken{
                id:        id,
                userID:    userID,
                sessionID: sessionID,
                token:     token,
                expiresAt: expiresAt,
                createdAt: createdAt,
                usedAt:    usedAt,
        }
}

// ReconstructUser creates a User from persisted data (bypasses validation).
func ReconstructUser(
        id uuid.UUID,
        phone, email, name string,
        role UserRole,
        status UserStatus,
        trustScore int,
        createdAt, updatedAt time.Time,
) *User {
        return &User{
                id:         id,
                phone:      phone,
                email:      email,
                name:       name,
                role:       role,
                status:     status,
                trustScore: trustScore,
                createdAt:  createdAt,
                updatedAt:  updatedAt,
        }
}

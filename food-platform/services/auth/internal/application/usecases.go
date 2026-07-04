// Package application contains the use cases (application services) for the Auth Service.
// Use cases orchestrate the domain logic and infrastructure.
package application

import (
        "context"
        "crypto/rand"
        "errors"
        "fmt"
        "math/big"
        "time"

        stderrors "github.com/food-platform/shared/errors"
        "github.com/food-platform/shared/logging"
        "github.com/food-platform/auth/internal/domain"
        "github.com/google/uuid"
)

// UserRepository is the interface for persisting users.
type UserRepository interface {
        Create(ctx context.Context, user *domain.User) error
        FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
        FindByPhone(ctx context.Context, phone string) (*domain.User, error)
        Update(ctx context.Context, user *domain.User) error
}

// OTPRepository is the interface for persisting OTP codes.
type OTPRepository interface {
        Save(ctx context.Context, otp *domain.OTP) error
        FindByID(ctx context.Context, id uuid.UUID) (*domain.OTP, error)
        FindLatestByPhone(ctx context.Context, phone string) (*domain.OTP, error)
}

// SessionRepository is the interface for persisting sessions.
type SessionRepository interface {
        Save(ctx context.Context, session *domain.Session) error
        FindByID(ctx context.Context, id uuid.UUID) (*domain.Session, error)
        FindByRefreshToken(ctx context.Context, token string) (*domain.Session, error)
        Revoke(ctx context.Context, id uuid.UUID) error
        RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
}

// RefreshTokenRepository is the interface for persisting refresh tokens.
type RefreshTokenRepository interface {
        Save(ctx context.Context, token *domain.RefreshToken) error
        FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error)
        MarkUsed(ctx context.Context, id uuid.UUID) error
}

// SMSSender is the interface for sending SMS messages (OTP codes).
type SMSSender interface {
        SendOTP(ctx context.Context, phone, code string) error
}

// JWTGenerator is the interface for generating JWT access tokens.
type JWTGenerator interface {
        Generate(ctx context.Context, userID uuid.UUID, role domain.UserRole, sessionID uuid.UUID) (string, int, error)
}

// ============ DTOs ============

// SendOTPCommand is the input for sending an OTP.
type SendOTPCommand struct {
        Phone string
        Role  domain.UserRole
}

// SendOTPResult is the output of sending an OTP.
type SendOTPResult struct {
        RequestID        string
        ExpiresInSeconds int
        AttemptsRemaining int
}

// VerifyOTPCommand is the input for verifying an OTP.
type VerifyOTPCommand struct {
        RequestID         string
        Code              string
        DeviceFingerprint string
        UserAgent         string
        IPAddress         string
}

// AuthResult is the output of a successful authentication.
type AuthResult struct {
        AccessToken  string
        RefreshToken string
        ExpiresIn    int
        User         *domain.User
}

// RefreshTokenCommand is the input for refreshing an access token.
type RefreshTokenCommand struct {
        RefreshToken string
}

// LogoutCommand is the input for logging out.
type LogoutCommand struct {
        UserID   uuid.UUID
        SessionID *uuid.UUID // Optional; if nil, revoke all sessions
}

// ============ Use Cases ============

// SendOTPUseCase handles sending an OTP to a phone number.
type SendOTPUseCase struct {
        userRepo  UserRepository
        otpRepo   OTPRepository
        smsSender SMSSender
}

// NewSendOTPUseCase creates a new SendOTPUseCase.
func NewSendOTPUseCase(
        userRepo UserRepository,
        otpRepo OTPRepository,
        smsSender SMSSender,
) *SendOTPUseCase {
        return &SendOTPUseCase{
                userRepo:  userRepo,
                otpRepo:   otpRepo,
                smsSender: smsSender,
        }
}

// Execute sends an OTP to the given phone number.
// If the user doesn't exist, it creates a new user (for customer/driver roles).
func (uc *SendOTPUseCase) Execute(ctx context.Context, cmd SendOTPCommand) (*SendOTPResult, error) {
        // Normalize and validate phone
        phone := domain.NormalizePhone(cmd.Phone)
        if err := domain.ValidatePhone(phone); err != nil {
                return nil, stderrors.ErrInvalidPhone.WithDetails(map[string]interface{}{
                        "phone": cmd.Phone,
                })
        }

        // Check if user exists; create if not (for customer/driver)
        user, err := uc.userRepo.FindByPhone(ctx, phone)
        if err != nil {
                if errors.Is(err, domain.ErrUserNotFound) {
                        // Create new user for customer/driver roles
                        if cmd.Role == domain.RoleCustomer || cmd.Role == domain.RoleDriver {
                                name := "مستخدم" // Default name; will be updated in profile setup
                                if cmd.Role == domain.RoleDriver {
                                        name = "مندوب"
                                }
                                user, err = domain.NewUser(phone, name, cmd.Role)
                                if err != nil {
                                        return nil, stderrors.Wrap(err, "USER_CREATE_FAILED", "Failed to create user", 500)
                                }
                                if err := uc.userRepo.Create(ctx, user); err != nil {
                                        return nil, stderrors.Wrap(err, "USER_CREATE_FAILED", "Failed to create user", 500)
                                }
                                logging.FromContext(ctx).Info("user_created",
                                        "user_id", user.ID(),
                                        "phone", phone,
                                        "role", cmd.Role,
                                )
                        } else {
                                // For employee roles, user must be pre-created by HR
                                return nil, stderrors.ErrUserNotFound.WithDetails(map[string]interface{}{
                                        "phone": phone,
                                        "role":  cmd.Role,
                                        "hint":  "Employee accounts must be created by HR first",
                                })
                        }
                } else {
                        return nil, stderrors.Wrap(err, "DATABASE_ERROR", "Failed to find user", 500)
                }
        }

        // Check if user is active
        if !user.IsActive() {
                return nil, stderrors.New("USER_INACTIVE", "User account is not active", 403)
        }

        // Generate OTP code
        code, err := generateOTPCode()
        if err != nil {
                return nil, stderrors.Wrap(err, "INTERNAL_ERROR", "Failed to generate OTP", 500)
        }

        // Create OTP entity
        otp := domain.NewOTP(phone, code)
        if err := uc.otpRepo.Save(ctx, otp); err != nil {
                return nil, stderrors.Wrap(err, "DATABASE_ERROR", "Failed to save OTP", 500)
        }

        // Send OTP via SMS
        if err := uc.smsSender.SendOTP(ctx, phone, code); err != nil {
                logging.FromContext(ctx).Error("sms_send_failed",
                        "phone", phone,
                        "error", err,
                )
                // Don't fail the request; user can retry
                // In production, you might want to return an error here
        }

        logging.FromContext(ctx).Info("otp_sent",
                "phone", phone,
                "otp_id", otp.ID(),
        )

        return &SendOTPResult{
                RequestID:        otp.ID().String(),
                ExpiresInSeconds: 120,
                AttemptsRemaining: otp.MaxAttempts(),
        }, nil
}

// VerifyOTPUseCase handles verifying an OTP and authenticating the user.
type VerifyOTPUseCase struct {
        userRepo     UserRepository
        otpRepo      OTPRepository
        sessionRepo  SessionRepository
        refreshRepo  RefreshTokenRepository
        jwtGenerator JWTGenerator
}

// NewVerifyOTPUseCase creates a new VerifyOTPUseCase.
func NewVerifyOTPUseCase(
        userRepo UserRepository,
        otpRepo OTPRepository,
        sessionRepo SessionRepository,
        refreshRepo RefreshTokenRepository,
        jwtGenerator JWTGenerator,
) *VerifyOTPUseCase {
        return &VerifyOTPUseCase{
                userRepo:     userRepo,
                otpRepo:      otpRepo,
                sessionRepo:  sessionRepo,
                refreshRepo:  refreshRepo,
                jwtGenerator: jwtGenerator,
        }
}

// Execute verifies the OTP and returns auth tokens.
func (uc *VerifyOTPUseCase) Execute(ctx context.Context, cmd VerifyOTPCommand) (*AuthResult, error) {
        // Parse request ID
        requestID, err := uuid.Parse(cmd.RequestID)
        if err != nil {
                return nil, stderrors.New("INVALID_REQUEST", "Invalid request ID", 400)
        }

        // Find OTP
        otp, err := uc.otpRepo.FindByID(ctx, requestID)
        if err != nil {
                return nil, stderrors.ErrInvalidOTP.WithDetails(map[string]interface{}{
                        "request_id": cmd.RequestID,
                })
        }

        // Verify OTP code
        if err := otp.Verify(cmd.Code); err != nil {
                if errors.Is(err, domain.ErrOTPExpired) || errors.Is(err, domain.ErrOTPAttemptsExceeded) {
                        return nil, stderrors.ErrInvalidOTP.WithDetails(map[string]interface{}{
                                "reason":         err.Error(),
                                "attempts_used":  otp.AttemptsUsed(),
                                "max_attempts":   otp.MaxAttempts(),
                        })
                }
                return nil, stderrors.ErrInvalidOTP
        }

        // Find user by phone
        user, err := uc.userRepo.FindByPhone(ctx, otp.Phone())
        if err != nil {
                return nil, stderrors.ErrUserNotFound
        }

        // Check if user is active
        if !user.IsActive() {
                return nil, stderrors.New("USER_INACTIVE", "User account is not active", 403)
        }

        // Generate access token (JWT)
        accessToken, expiresIn, err := uc.jwtGenerator.Generate(ctx, user.ID(), user.Role(), uuid.New())
        if err != nil {
                return nil, stderrors.Wrap(err, "INTERNAL_ERROR", "Failed to generate access token", 500)
        }

        // Generate refresh token (random string)
        refreshTokenStr, err := generateRefreshToken()
        if err != nil {
                return nil, stderrors.Wrap(err, "INTERNAL_ERROR", "Failed to generate refresh token", 500)
        }

        // Create session
        session := domain.NewSession(
                user.ID(),
                refreshTokenStr,
                cmd.DeviceFingerprint,
                cmd.UserAgent,
                cmd.IPAddress,
                30*24*time.Hour, // 30 days
        )

        // Save session
        if err := uc.sessionRepo.Save(ctx, session); err != nil {
                return nil, stderrors.Wrap(err, "DATABASE_ERROR", "Failed to save session", 500)
        }

        // Save refresh token
        refreshToken := domain.NewRefreshToken(user.ID(), session.ID(), refreshTokenStr, 30*24*time.Hour)
        if err := uc.refreshRepo.Save(ctx, refreshToken); err != nil {
                return nil, stderrors.Wrap(err, "DATABASE_ERROR", "Failed to save refresh token", 500)
        }

        logging.FromContext(ctx).Info("user_authenticated",
                "user_id", user.ID(),
                "session_id", session.ID(),
                "role", user.Role(),
        )

        return &AuthResult{
                AccessToken:  accessToken,
                RefreshToken: refreshTokenStr,
                ExpiresIn:    expiresIn,
                User:         user,
        }, nil
}

// RefreshTokenUseCase handles refreshing an access token.
type RefreshTokenUseCase struct {
        userRepo     UserRepository
        sessionRepo  SessionRepository
        refreshRepo  RefreshTokenRepository
        jwtGenerator JWTGenerator
}

// NewRefreshTokenUseCase creates a new RefreshTokenUseCase.
func NewRefreshTokenUseCase(
        userRepo UserRepository,
        sessionRepo SessionRepository,
        refreshRepo RefreshTokenRepository,
        jwtGenerator JWTGenerator,
) *RefreshTokenUseCase {
        return &RefreshTokenUseCase{
                userRepo:     userRepo,
                sessionRepo:  sessionRepo,
                refreshRepo:  refreshRepo,
                jwtGenerator: jwtGenerator,
        }
}

// Execute validates the refresh token and issues a new access token.
// The old refresh token is invalidated (rotation).
func (uc *RefreshTokenUseCase) Execute(ctx context.Context, cmd RefreshTokenCommand) (*AuthResult, error) {
        // Find refresh token
        rt, err := uc.refreshRepo.FindByToken(ctx, cmd.RefreshToken)
        if err != nil {
                return nil, stderrors.ErrRefreshTokenInvalid
        }

        // Validate refresh token
        if !rt.IsValid() {
                return nil, stderrors.ErrRefreshTokenInvalid.WithDetails(map[string]interface{}{
                        "reason": "token expired or already used",
                })
        }

        // Find session
        session, err := uc.sessionRepo.FindByID(ctx, rt.SessionID())
        if err != nil {
                return nil, stderrors.ErrSessionNotFound
        }

        // Check session is active
        if !session.IsActive() {
                return nil, stderrors.ErrSessionExpired
        }

        // Find user
        user, err := uc.userRepo.FindByID(ctx, rt.UserID())
        if err != nil {
                return nil, stderrors.ErrUserNotFound
        }

        // Check user is active
        if !user.IsActive() {
                return nil, stderrors.New("USER_INACTIVE", "User account is not active", 403)
        }

        // Mark old refresh token as used (rotation)
        if err := uc.refreshRepo.MarkUsed(ctx, rt.ID()); err != nil {
                return nil, stderrors.Wrap(err, "DATABASE_ERROR", "Failed to mark refresh token as used", 500)
        }

        // Generate new access token
        accessToken, expiresIn, err := uc.jwtGenerator.Generate(ctx, user.ID(), user.Role(), session.ID())
        if err != nil {
                return nil, stderrors.Wrap(err, "INTERNAL_ERROR", "Failed to generate access token", 500)
        }

        // Generate new refresh token
        newRefreshTokenStr, err := generateRefreshToken()
        if err != nil {
                return nil, stderrors.Wrap(err, "INTERNAL_ERROR", "Failed to generate refresh token", 500)
        }

        // Save new refresh token
        newRT := domain.NewRefreshToken(user.ID(), session.ID(), newRefreshTokenStr, 30*24*time.Hour)
        if err := uc.refreshRepo.Save(ctx, newRT); err != nil {
                return nil, stderrors.Wrap(err, "DATABASE_ERROR", "Failed to save refresh token", 500)
        }

        logging.FromContext(ctx).Info("token_refreshed",
                "user_id", user.ID(),
                "session_id", session.ID(),
        )

        return &AuthResult{
                AccessToken:  accessToken,
                RefreshToken: newRefreshTokenStr,
                ExpiresIn:    expiresIn,
                User:         user,
        }, nil
}

// LogoutUseCase handles logging out a user.
type LogoutUseCase struct {
        sessionRepo SessionRepository
        refreshRepo RefreshTokenRepository
}

// NewLogoutUseCase creates a new LogoutUseCase.
func NewLogoutUseCase(
        sessionRepo SessionRepository,
        refreshRepo RefreshTokenRepository,
) *LogoutUseCase {
        return &LogoutUseCase{
                sessionRepo: sessionRepo,
                refreshRepo: refreshRepo,
        }
}

// Execute logs out the user by revoking their session(s).
func (uc *LogoutUseCase) Execute(ctx context.Context, cmd LogoutCommand) error {
        if cmd.SessionID != nil {
                // Revoke specific session
                if err := uc.sessionRepo.Revoke(ctx, *cmd.SessionID); err != nil {
                        return stderrors.Wrap(err, "DATABASE_ERROR", "Failed to revoke session", 500)
                }
        } else {
                // Revoke all sessions for the user
                if err := uc.sessionRepo.RevokeAllForUser(ctx, cmd.UserID); err != nil {
                        return stderrors.Wrap(err, "DATABASE_ERROR", "Failed to revoke sessions", 500)
                }
        }

        logging.FromContext(ctx).Info("user_logged_out",
                "user_id", cmd.UserID,
                "session_id", cmd.SessionID,
        )
        return nil
}

// ============ Helpers ============

// generateOTPCode generates a random 6-digit OTP code.
func generateOTPCode() (string, error) {
        max := big.NewInt(1000000)
        n, err := rand.Int(rand.Reader, max)
        if err != nil {
                return "", err
        }
        return fmt.Sprintf("%06d", n.Int64()), nil
}

// generateRefreshToken generates a random 32-byte refresh token (64 hex chars).
func generateRefreshToken() (string, error) {
        b := make([]byte, 32)
        if _, err := rand.Read(b); err != nil {
                return "", err
        }
        return fmt.Sprintf("%x", b), nil
}

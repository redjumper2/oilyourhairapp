package services

import (
	"context"
	"fmt"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/sparque/auth_module/config"
	"github.com/sparque/auth_module/internal/database"
	"github.com/sparque/auth_module/pkg/utils"
)

// OAuthService handles OAuth authentication
type OAuthService struct {
	db          *database.DB
	cfg         *config.Config
	authService *AuthService
}

// NewOAuthService creates a new OAuth service
func NewOAuthService(db *database.DB, cfg *config.Config) *OAuthService {
	return &OAuthService{
		db:          db,
		cfg:         cfg,
		authService: NewAuthService(db, cfg),
	}
}

// InitializeProviders sets up OAuth providers (call once at startup)
func InitializeProviders(cfg *config.Config) {
	goth.UseProviders(
		google.New(
			cfg.Google.ClientID,
			cfg.Google.ClientSecret,
			cfg.Google.CallbackURL,
			"email", "profile",
		),
	)
}

// HandleGoogleCallback processes the Google OAuth callback
func (s *OAuthService) HandleGoogleCallback(ctx context.Context, domain string, user goth.User) (string, error) {
	// Find or create user with Google auth
	dbUser, err := s.authService.findOrCreateUser(ctx, user.Email, domain, "google", user.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to find/create user: %w", err)
	}

	// Update user info from Google if needed
	// Could store name, avatar, etc. here

	// Generate JWT
	jwt, err := utils.GenerateJWT(
		dbUser.ID.Hex(),
		dbUser.Email,
		dbUser.Domain,
		dbUser.Role,
		dbUser.Permissions,
		s.cfg.JWT.Secret,
		s.cfg.JWT.ExpiryHours,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return jwt, nil
}

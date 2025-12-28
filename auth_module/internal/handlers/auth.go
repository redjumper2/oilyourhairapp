package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sparque/auth_module/config"
	"github.com/sparque/auth_module/internal/database"
	"github.com/sparque/auth_module/internal/services"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *services.AuthService
	cfg         *config.Config
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *database.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(db, cfg),
		cfg:         cfg,
	}
}

// RequestMagicLinkRequest represents the request body
type RequestMagicLinkRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// RequestMagicLink handles POST /auth/magic-link/request
func (h *AuthHandler) RequestMagicLink(c echo.Context) error {
	var req RequestMagicLinkRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Extract domain from Host header
	domain := c.Request().Host
	// Remove port if present
	if colonIndex := len(domain) - 1; colonIndex > 0 {
		for i := len(domain) - 1; i >= 0; i-- {
			if domain[i] == ':' {
				domain = domain[:i]
				break
			}
		}
	}

	// Request magic link
	err := h.authService.RequestMagicLink(c.Request().Context(), req.Email, domain)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to send magic link",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Magic link sent to your email",
	})
}

// VerifyMagicLink handles GET /auth/magic-link/verify?token=xxx
func (h *AuthHandler) VerifyMagicLink(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Token is required",
		})
	}

	// Verify token and get/create user
	jwt, user, err := h.authService.VerifyMagicLink(c.Request().Context(), token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": jwt,
		"user": map[string]interface{}{
			"id":          user.ID.Hex(),
			"email":       user.Email,
			"domain":      user.Domain,
			"role":        user.Role,
			"permissions": user.Permissions,
		},
	})
}

// VerifyInvitation handles GET /auth/invitation/verify?token=xxx
func (h *AuthHandler) VerifyInvitation(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Token is required",
		})
	}

	// Verify invitation
	invitation, domain, err := h.authService.VerifyInvitation(c.Request().Context(), token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	response := map[string]interface{}{
		"invitation_id": invitation.ID.Hex(),
		"role":          invitation.Role,
		"domain":        invitation.Domain,
		"expires_at":    invitation.ExpiresAt,
		"time_remaining": invitation.TimeRemaining().String(),
		"branding": map[string]interface{}{
			"company_name":  domain.Branding.CompanyName,
			"primary_color": domain.Branding.PrimaryColor,
			"logo_url":      domain.Branding.LogoURL,
		},
	}

	// Include email if specified
	if invitation.Email != nil {
		response["email"] = *invitation.Email
	}

	// Include promo info if present
	if invitation.Metadata != nil {
		response["promo_code"] = invitation.Metadata.PromoCode
		response["discount_percent"] = invitation.Metadata.DiscountPercent
	}

	return c.JSON(http.StatusOK, response)
}

// AcceptInvitationRequest represents the request body
type AcceptInvitationRequest struct {
	Token        string `json:"token" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	AuthProvider string `json:"auth_provider" validate:"required"` // "magic_link" or "google"
	ProviderID   string `json:"provider_id,omitempty"`            // Google user ID if applicable
}

// AcceptInvitation handles POST /auth/invitation/accept
func (h *AuthHandler) AcceptInvitation(c echo.Context) error {
	var req AcceptInvitationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Accept invitation and create user
	jwt, user, err := h.authService.AcceptInvitation(
		c.Request().Context(),
		req.Token,
		req.Email,
		req.AuthProvider,
		req.ProviderID,
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": jwt,
		"user": map[string]interface{}{
			"id":          user.ID.Hex(),
			"email":       user.Email,
			"domain":      user.Domain,
			"role":        user.Role,
			"permissions": user.Permissions,
		},
	})
}

// GetMe handles GET /auth/me (requires authentication)
func (h *AuthHandler) GetMe(c echo.Context) error {
	// User info injected by middleware
	user := c.Get("user")
	if user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	return c.JSON(http.StatusOK, user)
}

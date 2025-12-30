package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
	"github.com/sparque/auth_module/config"
	"github.com/sparque/auth_module/internal/database"
	"github.com/sparque/auth_module/internal/services"
	"go.mongodb.org/mongo-driver/bson"
)

// OAuthHandler handles OAuth endpoints
type OAuthHandler struct {
	oauthService *services.OAuthService
	authService  *services.AuthService
	db           *database.DB
	cfg          *config.Config
	store        *sessions.CookieStore
}

// NewOAuthHandler creates a new OAuth handler
func NewOAuthHandler(db *database.DB, cfg *config.Config) *OAuthHandler {
	// Initialize session store for Goth
	store := sessions.NewCookieStore([]byte(cfg.JWT.Secret))
	store.Options.HttpOnly = true
	store.Options.Secure = cfg.Server.Env == "production"
	gothic.Store = store

	return &OAuthHandler{
		oauthService: services.NewOAuthService(db, cfg),
		authService:  services.NewAuthService(db, cfg),
		db:           db,
		cfg:          cfg,
		store:        store,
	}
}

// GoogleLogin initiates Google OAuth flow
func (h *OAuthHandler) GoogleLogin(c echo.Context) error {
	// Extract domain from query param first, fallback to Host header
	domain := c.QueryParam("domain")
	if domain == "" {
		domain = extractDomain(c.Request().Host)
	}

	// Verify domain exists
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	var domainDoc interface{}
	err := h.db.Domains.FindOne(ctx, bson.M{"domain": domain, "status": "active"}).Decode(&domainDoc)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Domain not found or inactive",
		})
	}

	// Store domain and redirect URL in session for callback
	redirectURL := c.QueryParam("redirect")
	if redirectURL == "" {
		redirectURL = fmt.Sprintf("https://%s", domain)
	}

	session, _ := h.store.Get(c.Request(), "auth-session")
	session.Values["domain"] = domain
	session.Values["redirect"] = redirectURL
	session.Save(c.Request(), c.Response())

	// Use gothic to handle the OAuth flow
	// We need to wrap Echo context for gothic
	req := c.Request()
	res := c.Response().Writer

	// Set provider to google by modifying URL query parameters
	q := req.URL.Query()
	q.Set("provider", "google")
	req.URL.RawQuery = q.Encode()

	// Initiate OAuth
	gothic.BeginAuthHandler(res, req)

	return nil
}

// GoogleCallback handles Google OAuth callback
func (h *OAuthHandler) GoogleCallback(c echo.Context) error {
	// Get domain and redirect URL from session
	session, _ := h.store.Get(c.Request(), "auth-session")
	domain, ok := session.Values["domain"].(string)
	if !ok || domain == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Domain not found in session",
		})
	}

	redirectURL, _ := session.Values["redirect"].(string)
	if redirectURL == "" {
		redirectURL = fmt.Sprintf("https://%s", domain)
	}

	// Complete OAuth with gothic
	req := c.Request()
	res := c.Response().Writer

	gothUser, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to complete OAuth: " + err.Error(),
		})
	}

	// Process callback and create/login user
	jwt, err := h.oauthService.HandleGoogleCallback(c.Request().Context(), domain, gothUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to process OAuth callback: " + err.Error(),
		})
	}

	// Clear session
	session.Values["domain"] = nil
	session.Values["redirect"] = nil
	session.Save(req, res)

	// Redirect to customer domain with token in hash (for SPA compatibility)
	finalRedirect := fmt.Sprintf("%s#token=%s", redirectURL, jwt)

	return c.Redirect(http.StatusTemporaryRedirect, finalRedirect)
}

// GoogleCallbackJSON handles Google OAuth callback and returns JSON (for API clients)
func (h *OAuthHandler) GoogleCallbackJSON(c echo.Context) error {
	// Get domain from query param (for API clients that can't use sessions)
	domain := c.QueryParam("domain")
	if domain == "" {
		// Try to get from Host header
		domain = extractDomain(c.Request().Host)
	}

	// Complete OAuth with gothic
	req := c.Request()
	res := c.Response().Writer

	gothUser, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to complete OAuth: " + err.Error(),
		})
	}

	// Process callback and create/login user
	jwt, err := h.oauthService.HandleGoogleCallback(c.Request().Context(), domain, gothUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to process OAuth callback: " + err.Error(),
		})
	}

	// Get user info
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	var user interface{}
	err = h.db.Users.FindOne(ctx, bson.M{
		"email":      gothUser.Email,
		"domain":     domain,
		"deleted_at": nil,
	}).Decode(&user)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": jwt,
		"user":  user,
	})
}

// extractDomain removes port from host if present
func extractDomain(host string) string {
	for i := len(host) - 1; i >= 0; i-- {
		if host[i] == ':' {
			return host[:i]
		}
	}
	return host
}

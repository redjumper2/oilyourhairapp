package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sparque/auth_module/internal/services"
)

// APIKeyHandler handles API key management endpoints
type APIKeyHandler struct {
	apiKeyService *services.APIKeyService
}

// NewAPIKeyHandler creates a new API key handler
func NewAPIKeyHandler(apiKeyService *services.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{
		apiKeyService: apiKeyService,
	}
}

// CreateAPIKeyRequest represents the request to create an API key
type CreateAPIKeyRequest struct {
	Domain      string   `json:"domain" validate:"required"`
	Service     string   `json:"service" validate:"required"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
	ExpiresIn   int      `json:"expires_in"` // Days from now
}

// CreateAPIKey creates a new service-scoped API key
func (h *APIKeyHandler) CreateAPIKey(c echo.Context) error {
	var req CreateAPIKeyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Get authenticated user from context
	userID, _ := c.Get("user_id").(string)
	domain, _ := c.Get("domain").(string)
	role, _ := c.Get("role").(string)

	// Verify user is admin for this domain
	if domain != req.Domain || role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Only domain admins can create API keys",
		})
	}

	// Default to 365 days if not specified
	expiresIn := req.ExpiresIn
	if expiresIn == 0 {
		expiresIn = 365
	}

	// Create API key
	createReq := services.CreateAPIKeyRequest{
		Domain:      req.Domain,
		Service:     req.Service,
		Description: req.Description,
		Permissions: req.Permissions,
		ExpiresIn:   time.Duration(expiresIn) * 24 * time.Hour,
	}

	token, apiKey, err := h.apiKeyService.CreateAPIKey(c.Request().Context(), createReq, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// Return both the token (only shown once) and the key metadata
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"api_key": token,
		"key_id":  apiKey.KeyID,
		"metadata": map[string]interface{}{
			"id":          apiKey.ID.Hex(),
			"domain":      apiKey.Domain,
			"service":     apiKey.Service,
			"description": apiKey.Description,
			"permissions": apiKey.Permissions,
			"expires_at":  apiKey.ExpiresAt,
			"created_at":  apiKey.CreatedAt,
		},
		"warning": "Save this API key securely. It cannot be retrieved again.",
	})
}

// ListAPIKeys lists all API keys for a domain
func (h *APIKeyHandler) ListAPIKeys(c echo.Context) error {
	domain := c.Param("domain")
	service := c.QueryParam("service")

	// Get authenticated user from context
	userDomain, _ := c.Get("domain").(string)
	role, _ := c.Get("role").(string)

	// Verify user is admin for this domain
	if userDomain != domain || role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Only domain admins can list API keys",
		})
	}

	keys, err := h.apiKeyService.ListAPIKeys(c.Request().Context(), domain, service)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// Return metadata only (not the actual key)
	var response []map[string]interface{}
	for _, key := range keys {
		response = append(response, map[string]interface{}{
			"id":           key.ID.Hex(),
			"key_id":       key.KeyID,
			"domain":       key.Domain,
			"service":      key.Service,
			"description":  key.Description,
			"permissions":  key.Permissions,
			"expires_at":   key.ExpiresAt,
			"revoked":      key.Revoked,
			"created_by":   key.CreatedBy,
			"created_at":   key.CreatedAt,
			"last_used_at": key.LastUsedAt,
			"is_valid":     key.IsValid(),
			"days_until_expiry": key.DaysUntilExpiry(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"keys": response,
	})
}

// GetAPIKey retrieves a specific API key by ID
func (h *APIKeyHandler) GetAPIKey(c echo.Context) error {
	keyID := c.Param("keyId")

	// Get authenticated user from context
	userDomain, _ := c.Get("domain").(string)
	role, _ := c.Get("role").(string)

	key, err := h.apiKeyService.GetAPIKey(c.Request().Context(), keyID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "API key not found",
		})
	}

	// Verify user is admin for this domain
	if userDomain != key.Domain || role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Access denied",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":           key.ID.Hex(),
		"key_id":       key.KeyID,
		"domain":       key.Domain,
		"service":      key.Service,
		"description":  key.Description,
		"permissions":  key.Permissions,
		"expires_at":   key.ExpiresAt,
		"revoked":      key.Revoked,
		"created_by":   key.CreatedBy,
		"created_at":   key.CreatedAt,
		"last_used_at": key.LastUsedAt,
		"is_valid":     key.IsValid(),
		"days_until_expiry": key.DaysUntilExpiry(),
	})
}

// RevokeAPIKey revokes an API key
func (h *APIKeyHandler) RevokeAPIKey(c echo.Context) error {
	keyID := c.Param("keyId")

	// Get authenticated user from context
	userDomain, _ := c.Get("domain").(string)
	role, _ := c.Get("role").(string)

	// Get the key to verify ownership
	key, err := h.apiKeyService.GetAPIKey(c.Request().Context(), keyID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "API key not found",
		})
	}

	// Verify user is admin for this domain
	if userDomain != key.Domain || role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Access denied",
		})
	}

	err = h.apiKeyService.RevokeAPIKey(c.Request().Context(), keyID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "API key revoked successfully",
	})
}

// GetExpiringKeys returns API keys expiring soon
func (h *APIKeyHandler) GetExpiringKeys(c echo.Context) error {
	domain := c.Param("domain")

	// Default to 30 days if not specified
	days := 30
	if daysParam := c.QueryParam("days"); daysParam != "" {
		if _, err := time.ParseDuration(daysParam + "d"); err == nil {
			// Parse days parameter if provided
		}
	}

	// Get authenticated user from context
	userDomain, _ := c.Get("domain").(string)
	role, _ := c.Get("role").(string)

	// Verify user is admin for this domain
	if userDomain != domain || role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Only domain admins can view expiring keys",
		})
	}

	keys, err := h.apiKeyService.GetExpiringKeys(c.Request().Context(), domain, time.Duration(days)*24*time.Hour)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// Return metadata only
	var response []map[string]interface{}
	for _, key := range keys {
		response = append(response, map[string]interface{}{
			"id":                key.ID.Hex(),
			"key_id":            key.KeyID,
			"service":           key.Service,
			"description":       key.Description,
			"expires_at":        key.ExpiresAt,
			"days_until_expiry": key.DaysUntilExpiry(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"expiring_keys": response,
		"threshold_days": days,
	})
}

// RegisterRoutes registers API key management routes
func (h *APIKeyHandler) RegisterRoutes(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	api := e.Group("/api/v1")

	// All routes require JWT authentication
	api.POST("/api-keys", h.CreateAPIKey, authMiddleware)
	api.GET("/domains/:domain/api-keys", h.ListAPIKeys, authMiddleware)
	api.GET("/api-keys/:keyId", h.GetAPIKey, authMiddleware)
	api.DELETE("/api-keys/:keyId", h.RevokeAPIKey, authMiddleware)
	api.GET("/domains/:domain/api-keys/expiring", h.GetExpiringKeys, authMiddleware)
}

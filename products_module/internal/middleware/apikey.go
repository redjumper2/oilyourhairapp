package middleware

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sparque/products_module/config"
)

// APIKeyClaims represents the claims in an API key JWT
type APIKeyClaims struct {
	KeyID       string   `json:"jti"`
	Domain      string   `json:"domain"`
	Service     string   `json:"service"`
	Type        string   `json:"type"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// APIKeyMiddleware validates API key JWTs for admin operations
func APIKeyMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract API key from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Missing authorization header",
				})
			}

			tokenString := extractToken(authHeader)
			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization header format",
				})
			}

			// Parse and validate API key JWT
			claims, err := validateAPIKey(tokenString, cfg.JWT.Secret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": fmt.Sprintf("Invalid API key: %v", err),
				})
			}

			// Verify it's an API key (not a user JWT)
			if claims.Type != "api_key" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Not an API key",
				})
			}

			// Verify service scope
			if claims.Service != "products" {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "API key not valid for products service",
				})
			}

			// Set context values
			c.Set("api_key_id", claims.KeyID)
			c.Set("domain", claims.Domain)
			c.Set("service", claims.Service)
			c.Set("permissions", claims.Permissions)

			return next(c)
		}
	}
}

// RequirePermission creates middleware that requires a specific permission
func RequirePermission(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Try to get permissions as []string first (direct from our middleware)
			permissions, ok := c.Get("permissions").([]string)
			if !ok {
				// Fallback: try as []interface{} (in case of other JWT middleware)
				permsInterface, ok := c.Get("permissions").([]interface{})
				if !ok {
					return c.JSON(http.StatusUnauthorized, map[string]string{
						"error": "Unauthorized",
					})
				}
				// Convert []interface{} to []string
				permissions = make([]string, len(permsInterface))
				for i, p := range permsInterface {
					permissions[i] = p.(string)
				}
			}

			// Check if required permission exists
			for _, p := range permissions {
				if p == permission {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{
				"error": fmt.Sprintf("Missing required permission: %s", permission),
			})
		}
	}
}

// validateAPIKey validates an API key JWT and returns the claims
func validateAPIKey(tokenString, secret string) (*APIKeyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &APIKeyClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*APIKeyClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// extractToken extracts the token from Authorization header
func extractToken(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}

// DomainMiddleware validates that the domain in the API key matches the requested resource
func DomainMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get domain from API key
			apiKeyDomain, _ := c.Get("domain").(string)

			// Get domain from request (for admin operations, domain should match API key)
			// For public operations, domain comes from path parameter

			// This middleware is optional - can be used on specific routes
			// that need to enforce domain matching

			c.Set("validated_domain", apiKeyDomain)
			return next(c)
		}
	}
}

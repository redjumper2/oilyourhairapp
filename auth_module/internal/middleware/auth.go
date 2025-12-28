package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sparque/auth_module/config"
	"github.com/sparque/auth_module/pkg/utils"
)

// AuthMiddleware creates authentication middleware
func AuthMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Missing authorization header",
				})
			}

			tokenString := utils.ExtractToken(authHeader)
			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization header format",
				})
			}

			// Validate JWT
			claims, err := utils.ValidateJWT(tokenString, cfg.JWT.Secret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid or expired token",
				})
			}

			// Verify domain matches Host header (prevent cross-domain token reuse)
			requestDomain := extractDomain(c.Request().Host)
			if claims.Domain != requestDomain {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "Token domain mismatch",
				})
			}

			// Set user info in context
			c.Set("user", map[string]interface{}{
				"id":          claims.UserID,
				"email":       claims.Email,
				"domain":      claims.Domain,
				"role":        claims.Role,
				"permissions": claims.Permissions,
			})
			c.Set("user_id", claims.UserID)
			c.Set("domain", claims.Domain)
			c.Set("role", claims.Role)
			c.Set("permissions", claims.Permissions)

			return next(c)
		}
	}
}

// RequireRole creates middleware that requires specific role
func RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole := c.Get("role")
			if userRole == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized",
				})
			}

			role := userRole.(string)
			for _, r := range roles {
				if role == r {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Insufficient permissions",
			})
		}
	}
}

// RequirePermission creates middleware that requires specific permission
func RequirePermission(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userPermissions := c.Get("permissions")
			if userPermissions == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized",
				})
			}

			permissions := userPermissions.([]string)
			for _, p := range permissions {
				if p == permission {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Insufficient permissions",
			})
		}
	}
}

// extractDomain removes port from host if present
func extractDomain(host string) string {
	if idx := strings.Index(host, ":"); idx != -1 {
		return host[:idx]
	}
	return host
}

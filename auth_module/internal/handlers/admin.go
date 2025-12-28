package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sparque/auth_module/config"
	"github.com/sparque/auth_module/internal/database"
	"github.com/sparque/auth_module/internal/models"
	"github.com/sparque/auth_module/internal/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdminHandler handles admin endpoints
type AdminHandler struct {
	db                *database.DB
	cfg               *config.Config
	invitationService *services.InvitationService
	emailService      *services.EmailService
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(db *database.DB, cfg *config.Config) *AdminHandler {
	return &AdminHandler{
		db:                db,
		cfg:               cfg,
		invitationService: services.NewInvitationService(db, cfg),
		emailService:      services.NewEmailService(cfg),
	}
}

// GetDomainSettings handles GET /admin/domain/settings
func (h *AdminHandler) GetDomainSettings(c echo.Context) error {
	domain := c.Get("domain").(string)

	var domainDoc models.Domain
	err := h.db.Domains.FindOne(c.Request().Context(), bson.M{"domain": domain}).Decode(&domainDoc)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Domain not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"domain":   domainDoc.Domain,
		"name":     domainDoc.Name,
		"status":   domainDoc.Status,
		"settings": domainDoc.Settings,
		"branding": domainDoc.Branding,
	})
}

// UpdateDomainSettingsRequest represents the request body
type UpdateDomainSettingsRequest struct {
	Settings *models.DomainSettings `json:"settings,omitempty"`
	Branding *models.DomainBranding `json:"branding,omitempty"`
}

// UpdateDomainSettings handles PUT /admin/domain/settings
func (h *AdminHandler) UpdateDomainSettings(c echo.Context) error {
	domain := c.Get("domain").(string)

	var req UpdateDomainSettingsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Build update document
	update := bson.M{}
	if req.Settings != nil {
		update["settings"] = req.Settings
	}
	if req.Branding != nil {
		update["branding"] = req.Branding
	}

	if len(update) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No fields to update",
		})
	}

	// Update domain
	result, err := h.db.Domains.UpdateOne(
		c.Request().Context(),
		bson.M{"domain": domain},
		bson.M{"$set": update},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update domain settings",
		})
	}

	if result.MatchedCount == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Domain not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Domain settings updated successfully",
	})
}

// ListUsers handles GET /admin/users
func (h *AdminHandler) ListUsers(c echo.Context) error {
	domain := c.Get("domain").(string)

	// Optional filters
	role := c.QueryParam("role")

	filter := bson.M{
		"domain":     domain,
		"deleted_at": nil, // Exclude soft-deleted users
	}

	if role != "" {
		filter["role"] = role
	}

	cursor, err := h.db.Users.Find(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to list users",
		})
	}
	defer cursor.Close(c.Request().Context())

	var users []models.User
	if err = cursor.All(c.Request().Context(), &users); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to decode users",
		})
	}

	// Format response
	response := make([]map[string]interface{}, len(users))
	for i, u := range users {
		response[i] = map[string]interface{}{
			"id":            u.ID.Hex(),
			"email":         u.Email,
			"role":          u.Role,
			"auth_provider": u.AuthProvider,
			"created_at":    u.CreatedAt,
			"last_login":    u.LastLogin,
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"users": response,
		"count": len(users),
	})
}

// CreateInvitationRequest represents the request body
type CreateInvitationRequest struct {
	Email           string  `json:"email"`
	Role            string  `json:"role" validate:"required"`
	Type            string  `json:"type" validate:"required"` // "email" | "qr_code" | "email_with_qr"
	SingleUse       *bool   `json:"single_use,omitempty"`
	MaxUses         *int    `json:"max_uses,omitempty"`
	ExpiresInHours  *int    `json:"expires_in_hours,omitempty"`
	PromoCode       string  `json:"promo_code,omitempty"`
	Source          string  `json:"source,omitempty"`
	Ref             string  `json:"ref,omitempty"`
	DiscountPercent float64 `json:"discount_percent,omitempty"`
}

// InviteUser handles POST /admin/users/invite
func (h *AdminHandler) InviteUser(c echo.Context) error {
	domain := c.Get("domain").(string)
	userID := c.Get("user_id").(string)

	var req CreateInvitationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Default single_use based on email
	singleUse := req.Email != ""
	if req.SingleUse != nil {
		singleUse = *req.SingleUse
	}

	// Prepare metadata
	var metadata *models.InvitationMetadata
	if req.PromoCode != "" || req.Source != "" || req.Ref != "" {
		metadata = &models.InvitationMetadata{
			PromoCode:       req.PromoCode,
			Source:          req.Source,
			Ref:             req.Ref,
			DiscountPercent: req.DiscountPercent,
		}
	}

	// Determine expiry hours
	expiryHours := 0
	if req.ExpiresInHours != nil {
		expiryHours = *req.ExpiresInHours
	}

	// Create invitation
	invitation, qrCodeDataURL, err := h.invitationService.CreateInvitation(
		c.Request().Context(),
		&services.CreateInvitationRequest{
			Domain:         domain,
			Email:          req.Email,
			Role:           req.Role,
			Type:           req.Type,
			SingleUse:      singleUse,
			MaxUses:        req.MaxUses,
			Metadata:       metadata,
			ExpiresInHours: expiryHours,
			CreatedBy:      userID,
		},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create invitation: " + err.Error(),
		})
	}

	inviteURL := h.cfg.App.FrontendURL + "/invite?token=" + invitation.Token

	// Send email if type includes email
	if req.Type == "email" || req.Type == "email_with_qr" {
		if req.Email == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Email is required for email invitations",
			})
		}

		// Get domain branding
		var domainDoc models.Domain
		err = h.db.Domains.FindOne(c.Request().Context(), bson.M{"domain": domain}).Decode(&domainDoc)
		if err == nil {
			_ = h.emailService.SendInvitation(
				req.Email,
				domain,
				invitation,
				inviteURL,
				qrCodeDataURL,
				&domainDoc.Branding,
			)
		}
	}

	response := map[string]interface{}{
		"invitation_id": invitation.ID.Hex(),
		"url":           inviteURL,
		"token":         invitation.Token,
		"expires_at":    invitation.ExpiresAt,
	}

	if qrCodeDataURL != "" {
		response["qr_code"] = qrCodeDataURL
	}

	return c.JSON(http.StatusOK, response)
}

// UpdateUserRequest represents the request body
type UpdateUserRequest struct {
	Role        *string   `json:"role,omitempty"`
	Permissions *[]string `json:"permissions,omitempty"`
}

// UpdateUser handles PUT /admin/users/:id
func (h *AdminHandler) UpdateUser(c echo.Context) error {
	domain := c.Get("domain").(string)
	userID := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Build update document
	update := bson.M{}
	if req.Role != nil {
		update["role"] = *req.Role
		// Update permissions based on role if not explicitly provided
		if req.Permissions == nil {
			update["permissions"] = models.GetPermissionsForRole(*req.Role)
		}
	}
	if req.Permissions != nil {
		update["permissions"] = *req.Permissions
	}

	if len(update) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No fields to update",
		})
	}

	// Update user (only within same domain)
	result, err := h.db.Users.UpdateOne(
		c.Request().Context(),
		bson.M{"_id": objectID, "domain": domain, "deleted_at": nil},
		bson.M{"$set": update},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update user",
		})
	}

	if result.MatchedCount == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "User updated successfully",
	})
}

// DeleteUser handles DELETE /admin/users/:id (soft delete)
func (h *AdminHandler) DeleteUser(c echo.Context) error {
	domain := c.Get("domain").(string)
	adminUserID := c.Get("user_id").(string)
	userID := c.Param("id")

	// Prevent deleting yourself
	if userID == adminUserID {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Cannot delete yourself",
		})
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	// Check if this is the last admin
	var targetUser models.User
	err = h.db.Users.FindOne(
		c.Request().Context(),
		bson.M{"_id": objectID, "domain": domain, "deleted_at": nil},
	).Decode(&targetUser)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	if targetUser.Role == "admin" {
		// Count active admins
		count, err := h.db.Users.CountDocuments(
			c.Request().Context(),
			bson.M{"domain": domain, "role": "admin", "deleted_at": nil},
		)
		if err == nil && count <= 1 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Cannot delete the last admin user",
			})
		}
	}

	// Soft delete
	now := time.Now()
	result, err := h.db.Users.UpdateOne(
		c.Request().Context(),
		bson.M{"_id": objectID, "domain": domain},
		bson.M{"$set": bson.M{
			"deleted_at": now,
			"deleted_by": adminUserID,
		}},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete user",
		})
	}

	if result.MatchedCount == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "User deleted successfully",
	})
}

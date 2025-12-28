package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/skip2/go-qrcode"
	"github.com/sparque/auth_module/config"
	"github.com/sparque/auth_module/internal/database"
	"github.com/sparque/auth_module/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InvitationService handles invitation creation and management
type InvitationService struct {
	db  *database.DB
	cfg *config.Config
}

// NewInvitationService creates a new invitation service
func NewInvitationService(db *database.DB, cfg *config.Config) *InvitationService {
	return &InvitationService{
		db:  db,
		cfg: cfg,
	}
}

// CreateInvitationRequest represents a request to create an invitation
type CreateInvitationRequest struct {
	Domain          string
	Email           string                  // Empty for anonymous/promotional
	Role            string
	Permissions     []string               // Optional, defaults based on role
	Type            string                 // "email" | "qr_code" | "email_with_qr"
	SingleUse       bool
	MaxUses         *int
	Metadata        *models.InvitationMetadata
	ExpiresInHours  int                    // If 0, uses defaults
	ExpiresAt       *time.Time             // Absolute expiry, takes precedence over ExpiresInHours
	CreatedBy       string                 // Admin user ID or "system"
}

// CreateInvitation creates a new invitation and optionally generates QR code
func (s *InvitationService) CreateInvitation(ctx context.Context, req *CreateInvitationRequest) (*models.Invitation, string, error) {
	// Generate unique token
	token, err := generateSecureToken()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Set permissions based on role if not provided
	permissions := req.Permissions
	if len(permissions) == 0 {
		permissions = models.GetPermissionsForRole(req.Role)
	}

	// Calculate expiry
	expiresAt := s.calculateExpiry(req)

	// Prepare email (nil for anonymous)
	var email *string
	if req.Email != "" {
		email = &req.Email
	}

	// Set created by
	createdBy := req.CreatedBy
	if createdBy == "" {
		createdBy = "system"
	}

	// Create invitation
	invitation := &models.Invitation{
		ID:          primitive.NewObjectID(),
		Domain:      req.Domain,
		Token:       token,
		Email:       email,
		Role:        req.Role,
		Permissions: permissions,
		Type:        req.Type,
		SingleUse:   req.SingleUse,
		MaxUses:     req.MaxUses,
		UsesCount:   0,
		Metadata:    req.Metadata,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		ExpiresAt:   expiresAt,
		Status:      "pending",
	}

	// Insert into database
	_, err = s.db.Invitations.InsertOne(ctx, invitation)
	if err != nil {
		return nil, "", fmt.Errorf("failed to insert invitation: %w", err)
	}

	// Generate QR code if needed
	var qrCodeDataURL string
	if req.Type == "qr_code" || req.Type == "email_with_qr" {
		inviteURL := fmt.Sprintf("%s/invite?token=%s", s.cfg.App.FrontendURL, token)
		qrCodeDataURL, err = s.generateQRCode(inviteURL)
		if err != nil {
			return nil, "", fmt.Errorf("failed to generate QR code: %w", err)
		}
	}

	return invitation, qrCodeDataURL, nil
}

// calculateExpiry determines when invitation expires
func (s *InvitationService) calculateExpiry(req *CreateInvitationRequest) time.Time {
	// Absolute expiry takes precedence
	if req.ExpiresAt != nil {
		return *req.ExpiresAt
	}

	// Use provided hours
	if req.ExpiresInHours > 0 {
		return time.Now().Add(time.Duration(req.ExpiresInHours) * time.Hour)
	}

	// Use defaults based on type
	switch req.Type {
	case "email", "email_with_qr":
		return time.Now().Add(time.Duration(s.cfg.Invitation.Defaults.EmailExpiryHours) * time.Hour)
	case "qr_code":
		if req.SingleUse {
			// User-specific QR code
			return time.Now().Add(time.Duration(s.cfg.Invitation.Defaults.QRCodeExpiryHours) * time.Hour)
		}
		// Promotional QR code
		return time.Now().Add(time.Duration(s.cfg.Invitation.Defaults.PromotionalExpiryHours) * time.Hour)
	default:
		// Default to email expiry
		return time.Now().Add(time.Duration(s.cfg.Invitation.Defaults.EmailExpiryHours) * time.Hour)
	}
}

// generateQRCode generates a QR code image and returns it as a data URL
func (s *InvitationService) generateQRCode(url string) (string, error) {
	// Generate QR code PNG
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	// Convert to base64 data URL
	dataURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
	return dataURL, nil
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

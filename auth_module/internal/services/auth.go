package services

import (
	"context"
	"fmt"
	"time"

	"github.com/sparque/auth_module/config"
	"github.com/sparque/auth_module/internal/database"
	"github.com/sparque/auth_module/internal/models"
	"github.com/sparque/auth_module/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthService handles authentication logic
type AuthService struct {
	db           *database.DB
	cfg          *config.Config
	emailService *EmailService
}

// NewAuthService creates a new auth service
func NewAuthService(db *database.DB, cfg *config.Config) *AuthService {
	return &AuthService{
		db:           db,
		cfg:          cfg,
		emailService: NewEmailService(cfg),
	}
}

// RequestMagicLink sends a magic link to the user's email
func (s *AuthService) RequestMagicLink(ctx context.Context, email, domain string) error {
	// Verify domain exists
	var domainDoc models.Domain
	err := s.db.Domains.FindOne(ctx, bson.M{"domain": domain, "status": "active"}).Decode(&domainDoc)
	if err != nil {
		return fmt.Errorf("domain not found or inactive: %w", err)
	}

	// Generate token
	token, err := generateSecureToken()
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	// Create magic link token
	magicLink := models.MagicLinkToken{
		ID:        primitive.NewObjectID(),
		Email:     email,
		Domain:    domain,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(s.cfg.MagicLink.ExpiryMinutes) * time.Minute),
		CreatedAt: time.Now(),
	}

	_, err = s.db.MagicLinkTokens.InsertOne(ctx, magicLink)
	if err != nil {
		return fmt.Errorf("failed to create magic link token: %w", err)
	}

	// Send email
	magicLinkURL := fmt.Sprintf("%s/auth/verify?token=%s", s.cfg.App.FrontendURL, token)
	err = s.emailService.SendMagicLink(email, domain, magicLinkURL)
	if err != nil {
		return fmt.Errorf("failed to send magic link email: %w", err)
	}

	return nil
}

// VerifyMagicLink verifies the magic link token and creates/logs in user
func (s *AuthService) VerifyMagicLink(ctx context.Context, token string) (string, *models.User, error) {
	// Find and verify token
	var magicLink models.MagicLinkToken
	err := s.db.MagicLinkTokens.FindOne(ctx, bson.M{"token": token}).Decode(&magicLink)
	if err != nil {
		return "", nil, fmt.Errorf("invalid or expired token")
	}

	// Check expiry
	if time.Now().After(magicLink.ExpiresAt) {
		return "", nil, fmt.Errorf("token expired")
	}

	// Delete token (single use)
	_, _ = s.db.MagicLinkTokens.DeleteOne(ctx, bson.M{"token": token})

	// Find or create user
	user, err := s.findOrCreateUser(ctx, magicLink.Email, magicLink.Domain, "magic_link", "")
	if err != nil {
		return "", nil, fmt.Errorf("failed to find/create user: %w", err)
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	_, err = s.db.Users.UpdateOne(ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"last_login": now}},
	)

	// Generate JWT
	jwt, err := utils.GenerateJWT(
		user.ID.Hex(),
		user.Email,
		user.Domain,
		user.Role,
		user.Permissions,
		s.cfg.JWT.Secret,
		s.cfg.JWT.ExpiryHours,
	)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	return jwt, user, nil
}

// VerifyInvitation verifies invitation token and returns invitation details
func (s *AuthService) VerifyInvitation(ctx context.Context, token string) (*models.Invitation, *models.Domain, error) {
	// Find invitation
	var invitation models.Invitation
	err := s.db.Invitations.FindOne(ctx, bson.M{"token": token}).Decode(&invitation)
	if err != nil {
		return nil, nil, fmt.Errorf("invitation not found")
	}

	// Check if can be claimed
	if !invitation.CanBeClaimed() {
		return nil, nil, fmt.Errorf("invitation cannot be claimed (expired or already used)")
	}

	// Get domain info
	var domain models.Domain
	err = s.db.Domains.FindOne(ctx, bson.M{"domain": invitation.Domain}).Decode(&domain)
	if err != nil {
		return nil, nil, fmt.Errorf("domain not found: %w", err)
	}

	return &invitation, &domain, nil
}

// AcceptInvitation accepts an invitation and creates user
func (s *AuthService) AcceptInvitation(ctx context.Context, token, email, authProvider, providerID string) (string, *models.User, error) {
	// Verify invitation
	invitation, domain, err := s.VerifyInvitation(ctx, token)
	if err != nil {
		return "", nil, err
	}

	// If invitation has email specified, must match
	if invitation.Email != nil && *invitation.Email != email {
		return "", nil, fmt.Errorf("email does not match invitation")
	}

	// Check if user already claimed this invitation (for multi-use)
	if !invitation.SingleUse {
		var existingLog models.InvitationLog
		err = s.db.InvitationLogs.FindOne(ctx, bson.M{
			"invitation_id": invitation.ID,
			"user_email":    email,
		}).Decode(&existingLog)
		if err == nil {
			return "", nil, fmt.Errorf("you have already claimed this invitation")
		}
	}

	// Create user
	user, err := s.findOrCreateUser(ctx, email, invitation.Domain, authProvider, providerID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Update user role/permissions from invitation
	user.Role = invitation.Role
	user.Permissions = invitation.Permissions
	_, err = s.db.Users.UpdateOne(ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{
			"role":        user.Role,
			"permissions": user.Permissions,
		}},
	)

	// Update invitation
	now := time.Now()
	userIDStr := user.ID.Hex()
	update := bson.M{
		"$inc": bson.M{"uses_count": 1},
		"$set": bson.M{"claimed_at": now},
	}

	if invitation.SingleUse {
		update["$set"].(bson.M)["status"] = "claimed"
		update["$set"].(bson.M)["claimed_by"] = userIDStr
	} else if invitation.MaxUses != nil && invitation.UsesCount+1 >= *invitation.MaxUses {
		update["$set"].(bson.M)["status"] = "exhausted"
	}

	_, err = s.db.Invitations.UpdateOne(ctx, bson.M{"_id": invitation.ID}, update)
	if err != nil {
		return "", nil, fmt.Errorf("failed to update invitation: %w", err)
	}

	// Log invitation claim
	invitationLog := models.InvitationLog{
		ID:           primitive.NewObjectID(),
		InvitationID: invitation.ID,
		Domain:       invitation.Domain,
		Email:        email,
		ClaimedBy:    userIDStr,
		ClaimedAt:    now,
		UserEmail:    email,
	}
	if invitation.Metadata != nil {
		invitationLog.PromoCode = invitation.Metadata.PromoCode
		invitationLog.Source = invitation.Metadata.Source
		invitationLog.Ref = invitation.Metadata.Ref
	}
	_, _ = s.db.InvitationLogs.InsertOne(ctx, invitationLog)

	// Send invitation email if configured
	if domain.Branding.SupportEmail != "" {
		// Email sent during invitation creation, not here
	}

	// Generate JWT
	jwt, err := utils.GenerateJWT(
		user.ID.Hex(),
		user.Email,
		user.Domain,
		user.Role,
		user.Permissions,
		s.cfg.JWT.Secret,
		s.cfg.JWT.ExpiryHours,
	)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	return jwt, user, nil
}

// findOrCreateUser finds existing user or creates new one
func (s *AuthService) findOrCreateUser(ctx context.Context, email, domain, authProvider, providerID string) (*models.User, error) {
	// Try to find existing user
	var user models.User
	err := s.db.Users.FindOne(ctx, bson.M{
		"email":      email,
		"domain":     domain,
		"deleted_at": nil,
	}).Decode(&user)

	if err == nil {
		// User exists
		return &user, nil
	}

	// Create new user with default role
	var domainDoc models.Domain
	err = s.db.Domains.FindOne(ctx, bson.M{"domain": domain}).Decode(&domainDoc)
	if err != nil {
		return nil, fmt.Errorf("domain not found: %w", err)
	}

	defaultRole := domainDoc.Settings.DefaultRole
	user = models.User{
		ID:           primitive.NewObjectID(),
		Email:        email,
		Domain:       domain,
		AuthProvider: authProvider,
		ProviderID:   providerID,
		Role:         defaultRole,
		Permissions:  models.GetPermissionsForRole(defaultRole),
		CreatedAt:    time.Now(),
	}

	_, err = s.db.Users.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

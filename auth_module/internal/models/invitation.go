package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Invitation represents a user invitation (email or QR code)
type Invitation struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Domain      string              `bson:"domain" json:"domain"`
	Token       string              `bson:"token" json:"token"`                                     // Unique token for invitation URL
	Email       *string             `bson:"email,omitempty" json:"email,omitempty"`                 // null for anonymous/promotional
	Role        string              `bson:"role" json:"role"`
	Permissions []string            `bson:"permissions" json:"permissions"`
	Type        string              `bson:"type" json:"type"`                                       // "email" | "qr_code" | "email_with_qr"
	SingleUse   bool                `bson:"single_use" json:"single_use"`                           // true for user-specific, false for promotional
	MaxUses     *int                `bson:"max_uses,omitempty" json:"max_uses,omitempty"`           // Optional limit for multi-use
	UsesCount   int                 `bson:"uses_count" json:"uses_count"`                           // Track claims
	Metadata    *InvitationMetadata `bson:"metadata,omitempty" json:"metadata,omitempty"`           // Promotional tracking
	CreatedBy   string              `bson:"created_by" json:"created_by"`                           // Admin user_id
	CreatedAt   time.Time           `bson:"created_at" json:"created_at"`
	ExpiresAt   time.Time           `bson:"expires_at" json:"expires_at"`
	Status      string              `bson:"status" json:"status"`                                   // "pending" | "claimed" | "expired" | "exhausted"
	ClaimedAt   *time.Time          `bson:"claimed_at,omitempty" json:"claimed_at,omitempty"`
	ClaimedBy   *string             `bson:"claimed_by,omitempty" json:"claimed_by,omitempty"`       // User ID who claimed
}

// InvitationMetadata contains promotional/tracking info
type InvitationMetadata struct {
	PromoCode       string  `bson:"promo_code,omitempty" json:"promo_code,omitempty"`
	Source          string  `bson:"source,omitempty" json:"source,omitempty"`                       // "booth", "instagram", etc.
	Ref             string  `bson:"ref,omitempty" json:"ref,omitempty"`                             // Referrer (sales rep, affiliate)
	DiscountPercent float64 `bson:"discount_percent,omitempty" json:"discount_percent,omitempty"`
	CustomData      map[string]interface{} `bson:"custom_data,omitempty" json:"custom_data,omitempty"`
}

// InvitationLog tracks invitation history for analytics
type InvitationLog struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	InvitationID primitive.ObjectID `bson:"invitation_id" json:"invitation_id"`
	Domain       string             `bson:"domain" json:"domain"`
	Email        string             `bson:"email" json:"email"`
	PromoCode    string             `bson:"promo_code,omitempty" json:"promo_code,omitempty"`
	Source       string             `bson:"source,omitempty" json:"source,omitempty"`
	Ref          string             `bson:"ref,omitempty" json:"ref,omitempty"`
	ClaimedBy    string             `bson:"claimed_by" json:"claimed_by"`                           // User ID
	ClaimedAt    time.Time          `bson:"claimed_at" json:"claimed_at"`
	UserEmail    string             `bson:"user_email" json:"user_email"`
}

// MagicLinkToken represents a temporary token for magic link authentication
type MagicLinkToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email"`
	Domain    string             `bson:"domain" json:"domain"`
	Token     string             `bson:"token" json:"token"`
	ExpiresAt time.Time          `bson:"expires_at" json:"expires_at"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// IsExpired checks if invitation has expired
func (i *Invitation) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}

// CanBeClaimed checks if invitation can be claimed
func (i *Invitation) CanBeClaimed() bool {
	if i.Status != "pending" {
		return false
	}
	if i.IsExpired() {
		return false
	}
	if i.SingleUse && i.UsesCount > 0 {
		return false
	}
	if i.MaxUses != nil && i.UsesCount >= *i.MaxUses {
		return false
	}
	return true
}

// TimeRemaining returns remaining time before expiry
func (i *Invitation) TimeRemaining() time.Duration {
	return time.Until(i.ExpiresAt)
}

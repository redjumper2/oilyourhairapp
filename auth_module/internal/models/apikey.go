package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// APIKey represents a service-scoped API key for accessing modules
type APIKey struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Domain      string             `bson:"domain" json:"domain"`
	Service     string             `bson:"service" json:"service"` // e.g., "products", "orders", "shipping"
	KeyID       string             `bson:"key_id" json:"key_id"`   // Unique identifier (jti claim)
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Permissions []string           `bson:"permissions" json:"permissions"` // e.g., ["products.read", "products.write"]
	ExpiresAt   time.Time          `bson:"expires_at" json:"expires_at"`
	Revoked     bool               `bson:"revoked" json:"revoked"`
	CreatedBy   string             `bson:"created_by" json:"created_by"` // User ID who created the key
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	LastUsedAt  *time.Time         `bson:"last_used_at,omitempty" json:"last_used_at,omitempty"`
}

// IsExpired checks if the API key has expired
func (k *APIKey) IsExpired() bool {
	return time.Now().After(k.ExpiresAt)
}

// IsValid checks if the API key is valid (not revoked and not expired)
func (k *APIKey) IsValid() bool {
	return !k.Revoked && !k.IsExpired()
}

// TimeRemaining returns the duration until expiration
func (k *APIKey) TimeRemaining() time.Duration {
	return time.Until(k.ExpiresAt)
}

// DaysUntilExpiry returns the number of days until expiration
func (k *APIKey) DaysUntilExpiry() int {
	return int(k.TimeRemaining().Hours() / 24)
}

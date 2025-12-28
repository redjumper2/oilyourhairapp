package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in a specific domain
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email" json:"email"`
	Domain       string             `bson:"domain" json:"domain"`                                   // Domain isolation
	AuthProvider string             `bson:"auth_provider" json:"auth_provider"`                     // "google" | "magic_link"
	ProviderID   string             `bson:"provider_id,omitempty" json:"provider_id,omitempty"`     // Google user ID if applicable
	Role         string             `bson:"role" json:"role"`                                       // "admin" | "editor" | "viewer" | "customer"
	Permissions  []string           `bson:"permissions" json:"permissions"`                         // ["products.read", "products.write", ...]
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	LastLogin    *time.Time         `bson:"last_login,omitempty" json:"last_login,omitempty"`
	DeletedAt    *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`       // Soft delete
	DeletedBy    string             `bson:"deleted_by,omitempty" json:"deleted_by,omitempty"`       // Admin user_id who deleted
}

// IsDeleted checks if user is soft deleted
func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

// GetPermissionsForRole returns default permissions for a role
func GetPermissionsForRole(role string) []string {
	switch role {
	case "admin":
		return []string{
			"domain.settings.read", "domain.settings.write",
			"users.read", "users.write", "users.delete", "users.invite",
			"products.read", "products.write",
			"orders.read", "orders.write",
			"inventory.read", "inventory.write",
		}
	case "editor":
		return []string{
			"products.read", "products.write",
			"orders.read",
			"inventory.read", "inventory.write",
		}
	case "viewer":
		return []string{
			"products.read",
			"orders.read",
			"inventory.read",
		}
	case "customer":
		return []string{
			"products.read",
			"cart.read", "cart.write",
			"orders.read",
		}
	default:
		return []string{}
	}
}

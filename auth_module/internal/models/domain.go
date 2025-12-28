package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Domain represents a registered domain/tenant
type Domain struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Domain    string             `bson:"domain" json:"domain"`                       // e.g., "oilyourhair.com"
	Name      string             `bson:"name" json:"name"`                           // e.g., "Oil Your Hair"
	Status    string             `bson:"status" json:"status"`                       // "active" | "suspended"
	Settings  DomainSettings     `bson:"settings" json:"settings"`
	Branding  DomainBranding     `bson:"branding" json:"branding"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	CreatedBy string             `bson:"created_by" json:"created_by"`               // "system" or user_id
}

// DomainSettings contains domain-specific settings
type DomainSettings struct {
	AllowedAuthProviders []string `bson:"allowed_auth_providers" json:"allowed_auth_providers"` // ["google", "magic_link"]
	DefaultRole          string   `bson:"default_role" json:"default_role"`                     // "customer"
	RequireEmailVerification bool `bson:"require_email_verification" json:"require_email_verification"`
}

// DomainBranding contains branding/white-label settings
type DomainBranding struct {
	CompanyName  string `bson:"company_name" json:"company_name"`
	PrimaryColor string `bson:"primary_color" json:"primary_color"`     // hex color
	LogoURL      string `bson:"logo_url,omitempty" json:"logo_url,omitempty"`
	SupportEmail string `bson:"support_email,omitempty" json:"support_email,omitempty"`
}

// DefaultDomainSettings returns default settings for new domains
func DefaultDomainSettings() DomainSettings {
	return DomainSettings{
		AllowedAuthProviders:     []string{"google", "magic_link"},
		DefaultRole:              "customer",
		RequireEmailVerification: true,
	}
}

// DefaultDomainBranding returns default branding
func DefaultDomainBranding(companyName string) DomainBranding {
	return DomainBranding{
		CompanyName:  companyName,
		PrimaryColor: "#000000", // default black
		LogoURL:      "",
		SupportEmail: "",
	}
}

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product represents a product in the catalog
type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Domain      string             `bson:"domain" json:"domain"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	BasePrice   float64            `bson:"base_price" json:"base_price"`
	Images      []string           `bson:"images,omitempty" json:"images,omitempty"`
	Attributes  map[string]string  `bson:"attributes,omitempty" json:"attributes,omitempty"` // Flexible key/value (category, type, etc.)
	Variants    []ProductVariant   `bson:"variants,omitempty" json:"variants,omitempty"`
	Active      bool               `bson:"active" json:"active"`
	CreatedBy   string             `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// ProductVariant represents a product variant with specific attributes
type ProductVariant struct {
	ID         string            `bson:"id" json:"id"`                                   // Variant ID
	Attributes map[string]string `bson:"attributes" json:"attributes"`                   // Variant-specific attributes (size, color, etc.)
	Price      float64           `bson:"price" json:"price"`                             // Price override for this variant
	Stock      int               `bson:"stock" json:"stock"`                             // Available inventory
	SKU        string            `bson:"sku,omitempty" json:"sku,omitempty"`            // Stock Keeping Unit
	ImageIndex int               `bson:"image_index,omitempty" json:"image_index,omitempty"` // Index into product images array
}

// GetEffectivePrice returns the variant price if available, otherwise base price
func (v *ProductVariant) GetEffectivePrice(basePrice float64) float64 {
	if v.Price > 0 {
		return v.Price
	}
	return basePrice
}

// IsInStock checks if the variant has available stock
func (v *ProductVariant) IsInStock() bool {
	return v.Stock > 0
}

// CreateProductRequest represents the request to create a product
type CreateProductRequest struct {
	Name        string                   `json:"name" validate:"required"`
	Description string                   `json:"description"`
	BasePrice   float64                  `json:"base_price" validate:"required,gt=0"`
	Images      []string                 `json:"images"`
	Attributes  map[string]string        `json:"attributes"`
	Variants    []CreateVariantRequest   `json:"variants"`
}

// CreateVariantRequest represents a variant in the create request
type CreateVariantRequest struct {
	Attributes map[string]string `json:"attributes" validate:"required"`
	Price      float64           `json:"price"`
	Stock      int               `json:"stock" validate:"gte=0"`
	SKU        string            `json:"sku"`
	ImageIndex int               `json:"image_index"`
}

// UpdateProductRequest represents the request to update a product
type UpdateProductRequest struct {
	Name        *string                  `json:"name"`
	Description *string                  `json:"description"`
	BasePrice   *float64                 `json:"base_price"`
	Images      *[]string                `json:"images"`
	Attributes  *map[string]string       `json:"attributes"`
	Variants    *[]CreateVariantRequest  `json:"variants"`
	Active      *bool                    `json:"active"`
}

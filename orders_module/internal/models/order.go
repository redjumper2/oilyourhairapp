package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Order represents a customer order
type Order struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderNumber string             `bson:"order_number" json:"order_number"`
	Domain      string             `bson:"domain" json:"domain"`

	// Customer
	Customer Customer `bson:"customer" json:"customer"`

	// Order Items
	Items []OrderItem `bson:"items" json:"items"`

	// Pricing (MVP: USD only)
	Currency string  `bson:"currency" json:"currency"`
	Subtotal float64 `bson:"subtotal" json:"subtotal"`
	Tax      float64 `bson:"tax" json:"tax"`
	Shipping float64 `bson:"shipping" json:"shipping"`
	Total    float64 `bson:"total" json:"total"`

	// Payment
	Payment Payment `bson:"payment" json:"payment"`

	// Addresses
	ShippingAddress Address `bson:"shipping_address" json:"shipping_address"`
	BillingAddress  Address `bson:"billing_address" json:"billing_address"`

	// Status
	Status string `bson:"status" json:"status"` // pending, paid, processing, shipped, delivered, cancelled, refunded

	// Timestamps
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	PaidAt    *time.Time `bson:"paid_at,omitempty" json:"paid_at,omitempty"`

	// Metadata
	Notes      string `bson:"notes,omitempty" json:"notes,omitempty"`
	AdminNotes string `bson:"admin_notes,omitempty" json:"admin_notes,omitempty"`
}

// Customer information
type Customer struct {
	UserID string `bson:"user_id,omitempty" json:"user_id,omitempty"` // From auth_module (null for guest)
	Email  string `bson:"email" json:"email"`
	Name   string `bson:"name" json:"name"`
}

// OrderItem represents a line item in an order
type OrderItem struct {
	ProductID         string                 `bson:"product_id" json:"product_id"`
	ProductName       string                 `bson:"product_name" json:"product_name"` // Snapshot
	ProductImage      string                 `bson:"product_image" json:"product_image"` // Snapshot
	VariantID         string                 `bson:"variant_id" json:"variant_id"`
	VariantSKU        string                 `bson:"variant_sku" json:"variant_sku"`
	VariantAttributes map[string]interface{} `bson:"variant_attributes" json:"variant_attributes"`
	Quantity          int                    `bson:"quantity" json:"quantity"`
	UnitPrice         float64                `bson:"unit_price" json:"unit_price"`
	Total             float64                `bson:"total" json:"total"`
}

// Payment information
type Payment struct {
	Provider        string  `bson:"provider" json:"provider"`                   // stripe
	PaymentIntentID string  `bson:"payment_intent_id" json:"payment_intent_id"` // Stripe payment intent ID
	Status          string  `bson:"status" json:"status"`                       // pending, succeeded, failed, refunded
	Amount          int64   `bson:"amount" json:"amount"`                       // Stripe uses cents
	Currency        string  `bson:"currency" json:"currency"`
	ClientSecret    string  `bson:"client_secret,omitempty" json:"client_secret,omitempty"` // For frontend
}

// Address for shipping/billing
type Address struct {
	Name         string `bson:"name" json:"name"`
	AddressLine1 string `bson:"address_line1" json:"address_line1"`
	AddressLine2 string `bson:"address_line2,omitempty" json:"address_line2,omitempty"`
	City         string `bson:"city" json:"city"`
	State        string `bson:"state" json:"state"`
	PostalCode   string `bson:"postal_code" json:"postal_code"`
	Country      string `bson:"country" json:"country"`
	Phone        string `bson:"phone,omitempty" json:"phone,omitempty"`
}

// OrderCounter for generating sequential order numbers
type OrderCounter struct {
	ID       string `bson:"_id"` // Domain name
	Sequence int    `bson:"sequence"`
	Year     int    `bson:"year"`
}

// CreateOrderRequest is the request body for creating an order
type CreateOrderRequest struct {
	Customer        Customer    `json:"customer"`
	Items           []OrderItem `json:"items"`
	ShippingAddress Address     `json:"shipping_address"`
	BillingAddress  Address     `json:"billing_address"`
	Notes           string      `json:"notes,omitempty"`
}

// UpdateOrderDetailsRequest is the request body for updating order details
type UpdateOrderDetailsRequest struct {
	Customer        *Customer `json:"customer,omitempty"`
	ShippingAddress *Address  `json:"shipping_address,omitempty"`
	BillingAddress  *Address  `json:"billing_address,omitempty"`
}

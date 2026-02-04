package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StockTransaction represents a stock change record
type StockTransaction struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	Domain     string `bson:"domain" json:"domain"`
	ProductID  string `bson:"product_id" json:"product_id"`
	VariantID  string `bson:"variant_id" json:"variant_id"`

	// Transaction details
	Type     string `bson:"type" json:"type"` // sale, restock, adjustment, return
	Quantity int    `bson:"quantity" json:"quantity"` // Negative for decrease, positive for increase

	// Stock levels
	StockBefore int `bson:"stock_before" json:"stock_before"`
	StockAfter  int `bson:"stock_after" json:"stock_after"`

	// Reference
	OrderID     string `bson:"order_id,omitempty" json:"order_id,omitempty"`
	OrderNumber string `bson:"order_number,omitempty" json:"order_number,omitempty"`

	// Metadata
	CreatedBy string    `bson:"created_by" json:"created_by"` // system | user_id
	Reason    string    `bson:"reason" json:"reason"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

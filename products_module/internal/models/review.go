package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Review represents a product review
type Review struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Domain       string             `bson:"domain" json:"domain"`
	ProductID    string             `bson:"product_id,omitempty" json:"product_id,omitempty"`
	ProductName  string             `bson:"product_name" json:"product"`
	Name         string             `bson:"name" json:"name"`
	Rating       int                `bson:"rating" json:"rating"`
	Text         string             `bson:"text" json:"text"`
	Highlight    string             `bson:"highlight,omitempty" json:"highlight,omitempty"`
	HelpfulCount int                `bson:"helpful_count" json:"helpfulCount"`
	CreatedAt    time.Time          `bson:"created_at" json:"date"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"-"`
}

// CreateReviewRequest represents the request to create a review
type CreateReviewRequest struct {
	ProductID string `json:"product_id,omitempty"`
	Product   string `json:"product" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Rating    int    `json:"rating" binding:"required,min=1,max=5"`
	Text      string `json:"text" binding:"required"`
	Highlight string `json:"highlight,omitempty"`
}

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Contact represents a contact form submission
type Contact struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Domain    string             `bson:"domain" json:"domain"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Subject   string             `bson:"subject" json:"subject"`
	Message   string             `bson:"message" json:"message"`
	Status    string             `bson:"status" json:"status"` // "new", "read", "replied"
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// CreateContactRequest represents the request to create a contact submission
type CreateContactRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Message string `json:"message" binding:"required"`
}

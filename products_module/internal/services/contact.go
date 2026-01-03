package services

import (
	"context"
	"time"

	"github.com/sparque/products_module/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ContactService struct {
	collection *mongo.Collection
}

func NewContactService(db *mongo.Database) *ContactService {
	return &ContactService{
		collection: db.Collection("contacts"),
	}
}

// CreateContact creates a new contact submission
func (s *ContactService) CreateContact(ctx context.Context, domain string, req *models.CreateContactRequest) (*models.Contact, error) {
	contact := &models.Contact{
		ID:        primitive.NewObjectID(),
		Domain:    domain,
		Name:      req.Name,
		Email:     req.Email,
		Subject:   req.Subject,
		Message:   req.Message,
		Status:    "new",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := s.collection.InsertOne(ctx, contact)
	if err != nil {
		return nil, err
	}

	return contact, nil
}

// GetContacts retrieves all contact submissions for a domain
func (s *ContactService) GetContacts(ctx context.Context, domain string, status string) ([]models.Contact, error) {
	filter := bson.M{"domain": domain}

	// If status is provided, filter by it
	if status != "" {
		filter["status"] = status
	}

	// Sort by created_at descending (newest first)
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var contacts []models.Contact
	if err := cursor.All(ctx, &contacts); err != nil {
		return nil, err
	}

	return contacts, nil
}

// GetContactByID retrieves a single contact by ID
func (s *ContactService) GetContactByID(ctx context.Context, domain string, id string) (*models.Contact, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var contact models.Contact
	err = s.collection.FindOne(ctx, bson.M{
		"_id":    objID,
		"domain": domain,
	}).Decode(&contact)

	if err != nil {
		return nil, err
	}

	return &contact, nil
}

// UpdateContactStatus updates the status of a contact submission
func (s *ContactService) UpdateContactStatus(ctx context.Context, domain string, id string, status string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID, "domain": domain},
		bson.M{
			"$set": bson.M{
				"status":     status,
				"updated_at": time.Now(),
			},
		},
	)

	return err
}

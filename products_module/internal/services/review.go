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

type ReviewService struct {
	collection *mongo.Collection
}

func NewReviewService(db *mongo.Database) *ReviewService {
	return &ReviewService{
		collection: db.Collection("reviews"),
	}
}

// CreateReview creates a new review
func (s *ReviewService) CreateReview(ctx context.Context, domain string, req *models.CreateReviewRequest) (*models.Review, error) {
	review := &models.Review{
		ID:           primitive.NewObjectID(),
		Domain:       domain,
		ProductID:    req.ProductID,
		ProductName:  req.Product,
		Name:         req.Name,
		Rating:       req.Rating,
		Text:         req.Text,
		Highlight:    req.Highlight,
		HelpfulCount: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err := s.collection.InsertOne(ctx, review)
	if err != nil {
		return nil, err
	}

	return review, nil
}

// GetReviews retrieves all reviews for a domain
func (s *ReviewService) GetReviews(ctx context.Context, domain string, productID string) ([]models.Review, error) {
	filter := bson.M{"domain": domain}

	// If productID is provided, filter by it
	if productID != "" {
		filter["product_id"] = productID
	}

	// Sort by created_at descending (newest first)
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reviews []models.Review
	if err := cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}

	return reviews, nil
}

// GetReviewByID retrieves a single review by ID
func (s *ReviewService) GetReviewByID(ctx context.Context, domain string, id string) (*models.Review, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var review models.Review
	err = s.collection.FindOne(ctx, bson.M{
		"_id":    objID,
		"domain": domain,
	}).Decode(&review)

	if err != nil {
		return nil, err
	}

	return &review, nil
}

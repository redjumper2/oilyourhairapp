package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sparque/products_module/internal/database"
	"github.com/sparque/products_module/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProductService handles product operations
type ProductService struct {
	db *database.DB
}

// NewProductService creates a new product service
func NewProductService(db *database.DB) *ProductService {
	return &ProductService{db: db}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, domain string, req models.CreateProductRequest, createdBy string) (*models.Product, error) {
	// Verify domain exists and is active
	var domainDoc bson.M
	err := s.db.Domains.FindOne(ctx, bson.M{
		"domain": domain,
		"status": "active",
	}).Decode(&domainDoc)
	if err != nil {
		return nil, fmt.Errorf("domain not found or inactive: %w", err)
	}

	// Create variants with generated IDs
	variants := make([]models.ProductVariant, len(req.Variants))
	for i, v := range req.Variants {
		variants[i] = models.ProductVariant{
			ID:         uuid.New().String(),
			Attributes: v.Attributes,
			Price:      v.Price,
			Stock:      v.Stock,
			SKU:        v.SKU,
			ImageIndex: v.ImageIndex,
		}
	}

	// Create product
	product := models.Product{
		ID:          primitive.NewObjectID(),
		Domain:      domain,
		Name:        req.Name,
		Description: req.Description,
		BasePrice:   req.BasePrice,
		Images:      req.Images,
		Attributes:  req.Attributes,
		Variants:    variants,
		Active:      true,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = s.db.Products.InsertOne(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return &product, nil
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(ctx context.Context, domain, productID string) (*models.Product, error) {
	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	var product models.Product
	err = s.db.Products.FindOne(ctx, bson.M{
		"_id":    objID,
		"domain": domain,
	}).Decode(&product)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	return &product, nil
}

// ListProducts lists all products for a domain with optional filtering
func (s *ProductService) ListProducts(ctx context.Context, domain string, activeOnly bool, attributes map[string]string) ([]models.Product, error) {
	filter := bson.M{"domain": domain}

	if activeOnly {
		filter["active"] = true
	}

	// Add attribute filters
	for key, value := range attributes {
		filter[fmt.Sprintf("attributes.%s", key)] = value
	}

	cursor, err := s.db.Products.Find(ctx, filter, options.Find().SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	return products, nil
}

// SearchProducts searches products by text
func (s *ProductService) SearchProducts(ctx context.Context, domain, query string, activeOnly bool) ([]models.Product, error) {
	filter := bson.M{
		"domain": domain,
		"$text":  bson.M{"$search": query},
	}

	if activeOnly {
		filter["active"] = true
	}

	cursor, err := s.db.Products.Find(ctx, filter, options.Find().SetSort(bson.M{"score": bson.M{"$meta": "textScore"}}))
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	return products, nil
}

// UpdateProduct updates a product
func (s *ProductService) UpdateProduct(ctx context.Context, domain, productID string, req models.UpdateProductRequest) (*models.Product, error) {
	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	// Build update document
	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	if req.Name != nil {
		update["$set"].(bson.M)["name"] = *req.Name
	}
	if req.Description != nil {
		update["$set"].(bson.M)["description"] = *req.Description
	}
	if req.BasePrice != nil {
		update["$set"].(bson.M)["base_price"] = *req.BasePrice
	}
	if req.Images != nil {
		update["$set"].(bson.M)["images"] = *req.Images
	}
	if req.Attributes != nil {
		update["$set"].(bson.M)["attributes"] = *req.Attributes
	}
	if req.Active != nil {
		update["$set"].(bson.M)["active"] = *req.Active
	}
	if req.Variants != nil {
		// Regenerate variant IDs
		variants := make([]models.ProductVariant, len(*req.Variants))
		for i, v := range *req.Variants {
			variants[i] = models.ProductVariant{
				ID:         uuid.New().String(),
				Attributes: v.Attributes,
				Price:      v.Price,
				Stock:      v.Stock,
				SKU:        v.SKU,
				ImageIndex: v.ImageIndex,
			}
		}
		update["$set"].(bson.M)["variants"] = variants
	}

	result, err := s.db.Products.UpdateOne(ctx, bson.M{
		"_id":    objID,
		"domain": domain,
	}, update)

	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("product not found")
	}

	return s.GetProduct(ctx, domain, productID)
}

// DeleteProduct deletes a product (soft delete by setting active=false)
func (s *ProductService) DeleteProduct(ctx context.Context, domain, productID string, hardDelete bool) error {
	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}

	if hardDelete {
		result, err := s.db.Products.DeleteOne(ctx, bson.M{
			"_id":    objID,
			"domain": domain,
		})
		if err != nil {
			return fmt.Errorf("failed to delete product: %w", err)
		}
		if result.DeletedCount == 0 {
			return fmt.Errorf("product not found")
		}
	} else {
		// Soft delete
		result, err := s.db.Products.UpdateOne(ctx, bson.M{
			"_id":    objID,
			"domain": domain,
		}, bson.M{
			"$set": bson.M{
				"active":     false,
				"updated_at": time.Now(),
			},
		})
		if err != nil {
			return fmt.Errorf("failed to deactivate product: %w", err)
		}
		if result.MatchedCount == 0 {
			return fmt.Errorf("product not found")
		}
	}

	return nil
}

// UpdateStock updates stock for a specific variant
func (s *ProductService) UpdateStock(ctx context.Context, domain, productID, variantID string, quantity int) error {
	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}

	result, err := s.db.Products.UpdateOne(ctx,
		bson.M{
			"_id":          objID,
			"domain":       domain,
			"variants.id":  variantID,
		},
		bson.M{
			"$set": bson.M{
				"variants.$.stock": quantity,
				"updated_at":       time.Now(),
			},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("product or variant not found")
	}

	return nil
}

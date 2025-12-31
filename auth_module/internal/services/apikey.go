package services

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sparque/auth_module/config"
	"github.com/sparque/auth_module/internal/database"
	"github.com/sparque/auth_module/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// APIKeyService handles API key operations
type APIKeyService struct {
	db  *database.DB
	cfg *config.Config
}

// NewAPIKeyService creates a new API key service
func NewAPIKeyService(db *database.DB, cfg *config.Config) *APIKeyService {
	return &APIKeyService{
		db:  db,
		cfg: cfg,
	}
}

// CreateAPIKeyRequest represents the request to create an API key
type CreateAPIKeyRequest struct {
	Domain      string        `json:"domain" validate:"required"`
	Service     string        `json:"service" validate:"required"` // e.g., "products"
	Description string        `json:"description"`
	Permissions []string      `json:"permissions"` // e.g., ["products.read", "products.write"]
	ExpiresIn   time.Duration `json:"expires_in"`  // Duration from now (e.g., 365 days)
}

// CreateAPIKey creates a new service-scoped API key
func (s *APIKeyService) CreateAPIKey(ctx context.Context, req CreateAPIKeyRequest, createdBy string) (string, *models.APIKey, error) {
	// Verify domain exists and is active
	var domain models.Domain
	err := s.db.Domains.FindOne(ctx, bson.M{
		"domain": req.Domain,
		"status": "active",
	}).Decode(&domain)
	if err != nil {
		return "", nil, fmt.Errorf("domain not found or inactive: %w", err)
	}

	// Generate unique key ID (jti claim)
	keyID := uuid.New().String()

	// Calculate expiration
	expiresAt := time.Now().Add(req.ExpiresIn)

	// Create API key record
	apiKey := models.APIKey{
		ID:          primitive.NewObjectID(),
		Domain:      req.Domain,
		Service:     req.Service,
		KeyID:       keyID,
		Description: req.Description,
		Permissions: req.Permissions,
		ExpiresAt:   expiresAt,
		Revoked:     false,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
	}

	// Insert into database
	_, err = s.db.APIKeys.InsertOne(ctx, apiKey)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create API key: %w", err)
	}

	// Generate JWT token (this IS the API key)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti":         keyID,
		"domain":      req.Domain,
		"service":     req.Service,
		"type":        "api_key",
		"permissions": req.Permissions,
		"exp":         expiresAt.Unix(),
		"iat":         time.Now().Unix(),
	})

	// Sign the token
	tokenString, err := token.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return "", nil, fmt.Errorf("failed to sign API key: %w", err)
	}

	return tokenString, &apiKey, nil
}

// ValidateAPIKey validates an API key JWT and checks if it's revoked
func (s *APIKeyService) ValidateAPIKey(ctx context.Context, tokenString string) (*models.APIKey, error) {
	// Parse JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid API key: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid API key claims")
	}

	// Check if it's an API key (not a regular user JWT)
	keyType, _ := claims["type"].(string)
	if keyType != "api_key" {
		return nil, fmt.Errorf("not an API key")
	}

	// Get key ID
	keyID, ok := claims["jti"].(string)
	if !ok {
		return nil, fmt.Errorf("missing key ID")
	}

	// Check if key exists and is valid in database
	var apiKey models.APIKey
	err = s.db.APIKeys.FindOne(ctx, bson.M{"key_id": keyID}).Decode(&apiKey)
	if err != nil {
		return nil, fmt.Errorf("API key not found: %w", err)
	}

	// Check if key is valid
	if !apiKey.IsValid() {
		if apiKey.Revoked {
			return nil, fmt.Errorf("API key has been revoked")
		}
		return nil, fmt.Errorf("API key has expired")
	}

	// Update last used timestamp
	s.db.APIKeys.UpdateOne(ctx,
		bson.M{"key_id": keyID},
		bson.M{"$set": bson.M{"last_used_at": time.Now()}},
	)

	return &apiKey, nil
}

// ListAPIKeys lists all API keys for a domain
func (s *APIKeyService) ListAPIKeys(ctx context.Context, domain string, service string) ([]models.APIKey, error) {
	filter := bson.M{"domain": domain}
	if service != "" {
		filter["service"] = service
	}

	cursor, err := s.db.APIKeys.Find(ctx, filter, options.Find().SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys: %w", err)
	}
	defer cursor.Close(ctx)

	var keys []models.APIKey
	if err := cursor.All(ctx, &keys); err != nil {
		return nil, fmt.Errorf("failed to decode API keys: %w", err)
	}

	return keys, nil
}

// RevokeAPIKey revokes an API key
func (s *APIKeyService) RevokeAPIKey(ctx context.Context, keyID string) error {
	result, err := s.db.APIKeys.UpdateOne(ctx,
		bson.M{"key_id": keyID},
		bson.M{"$set": bson.M{"revoked": true}},
	)

	if err != nil {
		return fmt.Errorf("failed to revoke API key: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("API key not found")
	}

	return nil
}

// GetAPIKey retrieves an API key by ID
func (s *APIKeyService) GetAPIKey(ctx context.Context, keyID string) (*models.APIKey, error) {
	var apiKey models.APIKey
	err := s.db.APIKeys.FindOne(ctx, bson.M{"key_id": keyID}).Decode(&apiKey)
	if err != nil {
		return nil, fmt.Errorf("API key not found: %w", err)
	}

	return &apiKey, nil
}

// GetExpiringKeys returns API keys expiring within the specified duration
func (s *APIKeyService) GetExpiringKeys(ctx context.Context, domain string, within time.Duration) ([]models.APIKey, error) {
	expiryThreshold := time.Now().Add(within)

	cursor, err := s.db.APIKeys.Find(ctx, bson.M{
		"domain":     domain,
		"revoked":    false,
		"expires_at": bson.M{"$lte": expiryThreshold},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get expiring keys: %w", err)
	}
	defer cursor.Close(ctx)

	var keys []models.APIKey
	if err := cursor.All(ctx, &keys); err != nil {
		return nil, fmt.Errorf("failed to decode expiring keys: %w", err)
	}

	return keys, nil
}

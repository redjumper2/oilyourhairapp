package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB holds MongoDB collections
type DB struct {
	Client      *mongo.Client
	Database    *mongo.Database
	Products    *mongo.Collection
	AuthDB      *mongo.Database // Reference to auth_module database
	Domains     *mongo.Collection // Read-only access to domains
}

// Connect establishes connection to MongoDB and returns DB instance
func Connect(uri, dbName, authDBName string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	database := client.Database(dbName)
	authDB := client.Database(authDBName)

	db := &DB{
		Client:   client,
		Database: database,
		Products: database.Collection("products"),
		AuthDB:   authDB,
		Domains:  authDB.Collection("domains"),
	}

	// Create indexes
	if err := db.createIndexes(ctx); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	return db, nil
}

// Disconnect closes the MongoDB connection
func (db *DB) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.Client.Disconnect(ctx)
}

// createIndexes creates required database indexes
func (db *DB) createIndexes(ctx context.Context) error {
	// Products: index on domain for multi-tenancy filtering
	_, err := db.Products.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "domain", Value: 1}},
	})
	if err != nil {
		return fmt.Errorf("failed to create products domain index: %w", err)
	}

	// Products: compound index on {domain, active} for active product queries
	_, err = db.Products.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "domain", Value: 1}, {Key: "active", Value: 1}},
	})
	if err != nil {
		return fmt.Errorf("failed to create products domain/active index: %w", err)
	}

	// Products: text index on name and description for search
	_, err = db.Products.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "name", Value: "text"}, {Key: "description", Value: "text"}},
	})
	if err != nil {
		return fmt.Errorf("failed to create products text search index: %w", err)
	}

	// Products: index on attributes for filtering
	_, err = db.Products.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "attributes", Value: 1}},
	})
	if err != nil {
		return fmt.Errorf("failed to create products attributes index: %w", err)
	}

	return nil
}

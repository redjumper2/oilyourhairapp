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
	Client         *mongo.Client
	Database       *mongo.Database
	Domains        *mongo.Collection
	Users          *mongo.Collection
	Invitations    *mongo.Collection
	InvitationLogs *mongo.Collection
	MagicLinkTokens *mongo.Collection
}

// Connect establishes connection to MongoDB and returns DB instance
func Connect(uri, dbName string) (*DB, error) {
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

	db := &DB{
		Client:          client,
		Database:        database,
		Domains:         database.Collection("domains"),
		Users:           database.Collection("users"),
		Invitations:     database.Collection("invitations"),
		InvitationLogs:  database.Collection("invitation_logs"),
		MagicLinkTokens: database.Collection("magic_link_tokens"),
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
	// Domains: unique index on domain
	_, err := db.Domains.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "domain", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create domains index: %w", err)
	}

	// Users: compound unique index on {email, domain}
	_, err = db.Users.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}, {Key: "domain", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create users index: %w", err)
	}

	// Users: index on domain for queries
	_, err = db.Users.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "domain", Value: 1}},
	})
	if err != nil {
		return fmt.Errorf("failed to create users domain index: %w", err)
	}

	// Invitations: unique index on token
	_, err = db.Invitations.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "token", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create invitations token index: %w", err)
	}

	// Invitations: TTL index on expires_at for auto-deletion
	_, err = db.Invitations.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "expires_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	})
	if err != nil {
		return fmt.Errorf("failed to create invitations TTL index: %w", err)
	}

	// InvitationLogs: compound index on {invitation_id, user_email} to prevent duplicate claims
	_, err = db.InvitationLogs.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "invitation_id", Value: 1}, {Key: "user_email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create invitation_logs index: %w", err)
	}

	// MagicLinkTokens: unique index on token
	_, err = db.MagicLinkTokens.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "token", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create magic_link_tokens index: %w", err)
	}

	// MagicLinkTokens: TTL index on expires_at for auto-deletion
	_, err = db.MagicLinkTokens.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "expires_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	})
	if err != nil {
		return fmt.Errorf("failed to create magic_link_tokens TTL index: %w", err)
	}

	return nil
}

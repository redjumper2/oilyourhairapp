package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sparque/orders_module/internal/database"
	"github.com/sparque/orders_module/internal/models"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StripeService struct {
	db            *database.MongoDB
	webhookSecret string
}

func NewStripeService(db *database.MongoDB, webhookSecret string) *StripeService {
	return &StripeService{
		db:            db,
		webhookSecret: webhookSecret,
	}
}

// HandleWebhook processes Stripe webhook events
func (s *StripeService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	// Verify webhook signature
	event, err := webhook.ConstructEvent(payload, signature, s.webhookSecret)
	if err != nil {
		return fmt.Errorf("webhook signature verification failed: %w", err)
	}

	// Handle different event types
	switch event.Type {
	case "payment_intent.succeeded":
		return s.handlePaymentSucceeded(ctx, event.Data.Raw)

	case "payment_intent.payment_failed":
		return s.handlePaymentFailed(ctx, event.Data.Raw)

	default:
		// Unhandled event type
		return nil
	}
}

// handlePaymentSucceeded handles successful payment
func (s *StripeService) handlePaymentSucceeded(ctx context.Context, rawData json.RawMessage) error {
	var pi stripe.PaymentIntent
	if err := json.Unmarshal(rawData, &pi); err != nil {
		return fmt.Errorf("failed to unmarshal payment intent: %w", err)
	}

	// Find order by payment intent ID
	collection := s.db.GetCollection("orders")
	var order models.Order
	err := collection.FindOne(ctx, bson.M{
		"payment.payment_intent_id": pi.ID,
	}).Decode(&order)

	if err != nil {
		return fmt.Errorf("order not found for payment intent %s: %w", pi.ID, err)
	}

	// Update order status
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":           "paid",
			"payment.status":   "succeeded",
			"paid_at":          now,
			"updated_at":       now,
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": order.ID}, update)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	// Deduct stock (MVP: simple approach)
	if err := s.deductStock(ctx, &order); err != nil {
		// Log error but don't fail the webhook
		// In production, you'd want retry logic or alerting
		fmt.Printf("ERROR: Failed to deduct stock for order %s: %v\n", order.OrderNumber, err)
	}

	return nil
}

// handlePaymentFailed handles failed payment
func (s *StripeService) handlePaymentFailed(ctx context.Context, rawData json.RawMessage) error {
	var pi stripe.PaymentIntent
	if err := json.Unmarshal(rawData, &pi); err != nil {
		return fmt.Errorf("failed to unmarshal payment intent: %w", err)
	}

	// Update order payment status
	collection := s.db.GetCollection("orders")
	update := bson.M{
		"$set": bson.M{
			"payment.status": "failed",
			"updated_at":     time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{
		"payment.payment_intent_id": pi.ID,
	}, update)

	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}

// deductStock deducts stock for all items in an order
func (s *StripeService) deductStock(ctx context.Context, order *models.Order) error {
	stockTxCollection := s.db.GetCollection("stock_transactions")

	for _, item := range order.Items {
		// Create stock transaction record
		tx := &models.StockTransaction{
			ID:          primitive.NewObjectID(),
			Domain:      order.Domain,
			ProductID:   item.ProductID,
			VariantID:   item.VariantID,
			Type:        "sale",
			Quantity:    -item.Quantity, // Negative for decrease
			OrderID:     order.ID.Hex(),
			OrderNumber: order.OrderNumber,
			CreatedBy:   "system",
			Reason:      "Order payment confirmed",
			CreatedAt:   time.Now(),
		}

		// Note: StockBefore and StockAfter would need to be fetched from products API
		// For MVP, we'll store them as 0 and implement full integration later
		tx.StockBefore = 0 // TODO: Fetch from products API
		tx.StockAfter = 0  // TODO: Calculate after deduction

		_, err := stockTxCollection.InsertOne(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to create stock transaction: %w", err)
		}

		// TODO: Call products API to actually deduct stock
		// This will be implemented in Task 7
	}

	return nil
}

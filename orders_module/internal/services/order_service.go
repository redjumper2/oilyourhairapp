package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sparque/orders_module/internal/database"
	"github.com/sparque/orders_module/internal/models"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderService struct {
	db         *database.MongoDB
	stripeKey  string
}

func NewOrderService(db *database.MongoDB, stripeKey string) *OrderService {
	// Initialize Stripe with secret key
	stripe.Key = stripeKey

	return &OrderService{
		db:        db,
		stripeKey: stripeKey,
	}
}

// CreateOrder creates a new order and Stripe payment intent
func (s *OrderService) CreateOrder(ctx context.Context, req *models.CreateOrderRequest, domain string) (*models.Order, error) {
	// Calculate totals
	var subtotal float64
	for i := range req.Items {
		req.Items[i].Total = req.Items[i].UnitPrice * float64(req.Items[i].Quantity)
		subtotal += req.Items[i].Total
	}

	// MVP: Simple pricing (no tax, no shipping)
	tax := 0.0
	shipping := 0.0
	total := subtotal + tax + shipping

	// Generate order number
	orderNumber, err := s.generateOrderNumber(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("failed to generate order number: %w", err)
	}

	// Create Stripe Payment Intent
	amountCents := int64(total * 100) // Convert to cents
	pi, err := s.createPaymentIntent(amountCents, "usd", orderNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	// Create order
	now := time.Now()
	order := &models.Order{
		OrderNumber: orderNumber,
		Domain:      domain,
		Customer:    req.Customer,
		Items:       req.Items,
		Currency:    "USD",
		Subtotal:    subtotal,
		Tax:         tax,
		Shipping:    shipping,
		Total:       total,
		Payment: models.Payment{
			Provider:        "stripe",
			PaymentIntentID: pi.ID,
			Status:          "pending",
			Amount:          amountCents,
			Currency:        "usd",
			ClientSecret:    pi.ClientSecret,
		},
		ShippingAddress: req.ShippingAddress,
		BillingAddress:  req.BillingAddress,
		Status:          "pending",
		CreatedAt:       now,
		UpdatedAt:       now,
		Notes:           req.Notes,
	}

	// Save to database
	collection := s.db.GetCollection("orders")
	result, err := collection.InsertOne(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to insert order: %w", err)
	}

	order.ID = result.InsertedID.(primitive.ObjectID)

	return order, nil
}

// createPaymentIntent creates a Stripe payment intent
func (s *OrderService) createPaymentIntent(amount int64, currency, orderNumber string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
		Metadata: map[string]string{
			"order_number": orderNumber,
		},
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	return pi, nil
}

// generateOrderNumber generates a sequential order number
func (s *OrderService) generateOrderNumber(ctx context.Context, domain string) (string, error) {
	collection := s.db.GetCollection("order_counters")

	currentYear := time.Now().Year()

	// Find and update counter atomically
	filter := bson.M{"_id": domain}
	update := bson.M{
		"$inc": bson.M{"sequence": 1},
		"$setOnInsert": bson.M{"year": currentYear},
	}

	var counter models.OrderCounter
	err := collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&counter)

	if err != nil {
		return "", err
	}

	// If year changed, reset counter
	if counter.Year != currentYear {
		_, err := collection.UpdateOne(
			ctx,
			filter,
			bson.M{
				"$set": bson.M{
					"sequence": 1,
					"year":     currentYear,
				},
			},
		)
		if err != nil {
			return "", err
		}
		counter.Sequence = 1
		counter.Year = currentYear
	}

	// Format: ORD-2026-00001
	orderNumber := fmt.Sprintf("ORD-%d-%05d", counter.Year, counter.Sequence)
	return orderNumber, nil
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(ctx context.Context, orderID string, domain string) (*models.Order, error) {
	collection := s.db.GetCollection("orders")

	objectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	var order models.Order
	err = collection.FindOne(ctx, bson.M{
		"_id":    objectID,
		"domain": domain,
	}).Decode(&order)

	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("order not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order, nil
}

// ListOrders lists orders for a user
func (s *OrderService) ListOrders(ctx context.Context, userID, domain string, limit int) ([]*models.Order, error) {
	collection := s.db.GetCollection("orders")

	// For now, show all orders for the domain (TODO: Add proper user auth filtering)
	filter := bson.M{
		"domain": domain,
	}

	log.Printf("ListOrders - domain: %s, filter: %+v", domain, filter)

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	defer cursor.Close(ctx)

	var orders []*models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, fmt.Errorf("failed to decode orders: %w", err)
	}

	return orders, nil
}

// UpdateOrderStatus updates an order's status
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID, status, domain string) error {
	collection := s.db.GetCollection("orders")

	objectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return fmt.Errorf("invalid order ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{
		"_id":    objectID,
		"domain": domain,
	}, update)

	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

// ListAllOrders lists all orders (admin) with optional domain filter
func (s *OrderService) ListAllOrders(ctx context.Context, domain string, limit int) ([]*models.Order, error) {
	collection := s.db.GetCollection("orders")

	// Filter by domain if provided (for multi-tenant), otherwise get all
	filter := bson.M{}
	if domain != "" {
		filter["domain"] = domain
	}

	log.Printf("ListAllOrders - domain filter: %s, filter: %+v", domain, filter)

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list all orders: %w", err)
	}
	defer cursor.Close(ctx)

	var orders []*models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, fmt.Errorf("failed to decode orders: %w", err)
	}

	return orders, nil
}

// UpdateOrderDetails updates customer and address info for an order
func (s *OrderService) UpdateOrderDetails(ctx context.Context, orderID string, req *models.UpdateOrderDetailsRequest, domain string) error {
	collection := s.db.GetCollection("orders")

	objectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return fmt.Errorf("invalid order ID: %w", err)
	}

	// Build update document with only provided fields
	setFields := bson.M{
		"updated_at": time.Now(),
	}

	if req.Customer != nil {
		setFields["customer"] = req.Customer
	}
	if req.ShippingAddress != nil {
		setFields["shipping_address"] = req.ShippingAddress
	}
	if req.BillingAddress != nil {
		setFields["billing_address"] = req.BillingAddress
	}

	update := bson.M{"$set": setFields}

	result, err := collection.UpdateOne(ctx, bson.M{
		"_id":    objectID,
		"domain": domain,
	}, update)

	if err != nil {
		return fmt.Errorf("failed to update order details: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("order not found")
	}

	log.Printf("UpdateOrderDetails - Updated order %s with fields: %+v", orderID, setFields)

	return nil
}

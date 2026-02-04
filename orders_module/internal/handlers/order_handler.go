package handlers

import (
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sparque/orders_module/internal/database"
	"github.com/sparque/orders_module/internal/models"
	"github.com/sparque/orders_module/internal/services"
)

type OrderHandler struct {
	db            *database.MongoDB
	jwtSecret     string
	orderService  *services.OrderService
	stripeService *services.StripeService
}

func NewOrderHandler(db *database.MongoDB, jwtSecret, stripeKey, webhookSecret string) *OrderHandler {
	return &OrderHandler{
		db:            db,
		jwtSecret:     jwtSecret,
		orderService:  services.NewOrderService(db, stripeKey),
		stripeService: services.NewStripeService(db, webhookSecret),
	}
}

// CreateOrder creates a new order and Stripe payment intent
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	// Parse request
	var req models.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Debug: log received items
	for i, item := range req.Items {
		log.Printf("CreateOrder - Item %d: product_id=%s, product_name=%s, product_image=%s",
			i, item.ProductID, item.ProductName, item.ProductImage)
	}

	// Get domain from header (multi-tenant)
	domain := c.Request().Header.Get("Host")
	if domain == "" {
		domain = "oilyourhair.com" // Default for development
	}

	// Create order
	order, err := h.orderService.CreateOrder(c.Request().Context(), &req, domain)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, order)
}

// GetOrder retrieves an order by ID
func (h *OrderHandler) GetOrder(c echo.Context) error {
	orderID := c.Param("id")
	domain := c.Request().Header.Get("Host")
	if domain == "" {
		domain = "oilyourhair.com"
	}

	order, err := h.orderService.GetOrder(c.Request().Context(), orderID, domain)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	// TODO: Verify user owns this order (JWT validation)

	return c.JSON(http.StatusOK, order)
}

// ListOrders lists orders for the authenticated user
func (h *OrderHandler) ListOrders(c echo.Context) error {
	// TODO: Get userID from JWT token
	userID := "guest" // Placeholder
	domain := c.Request().Header.Get("Host")
	if domain == "" {
		domain = "oilyourhair.com"
	}

	log.Printf("ListOrders handler - Host header: %s, domain: %s", c.Request().Header.Get("Host"), domain)

	orders, err := h.orderService.ListOrders(c.Request().Context(), userID, domain, 50)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	log.Printf("ListOrders handler - returned %d orders", len(orders))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"orders": orders,
		"count":  len(orders),
	})
}

// ListAllOrders lists all orders (admin only)
func (h *OrderHandler) ListAllOrders(c echo.Context) error {
	// TODO: Verify admin role from JWT token

	// Get optional domain filter from query param
	domain := c.QueryParam("domain")

	// Get limit from query param (default 100)
	limit := 100

	orders, err := h.orderService.ListAllOrders(c.Request().Context(), domain, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	log.Printf("ListAllOrders handler - returned %d orders", len(orders))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"orders": orders,
		"count":  len(orders),
	})
}

// UpdateOrderDetails updates an order's customer and address information
func (h *OrderHandler) UpdateOrderDetails(c echo.Context) error {
	orderID := c.Param("id")
	domain := c.Request().Header.Get("Host")
	if domain == "" {
		domain = "oilyourhair.com"
	}

	var req models.UpdateOrderDetailsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := h.orderService.UpdateOrderDetails(c.Request().Context(), orderID, &req, domain); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Order details updated successfully",
	})
}

// UpdateOrderStatus updates an order's status (admin only)
func (h *OrderHandler) UpdateOrderStatus(c echo.Context) error {
	// TODO: Verify admin role from JWT

	orderID := c.Param("id")
	domain := c.Request().Header.Get("Host")
	if domain == "" {
		domain = "oilyourhair.com"
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := h.orderService.UpdateOrderStatus(c.Request().Context(), orderID, req.Status, domain); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Order status updated successfully",
	})
}

// StripeWebhook handles Stripe payment webhooks
func (h *OrderHandler) StripeWebhook(c echo.Context) error {
	// Read raw body
	payload, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to read request body",
		})
	}

	// Get Stripe signature header
	signature := c.Request().Header.Get("Stripe-Signature")
	if signature == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing Stripe signature",
		})
	}

	// Handle webhook
	if err := h.stripeService.HandleWebhook(c.Request().Context(), payload, signature); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}

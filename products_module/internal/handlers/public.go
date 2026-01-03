package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sparque/products_module/internal/models"
	"github.com/sparque/products_module/internal/services"
)

// PublicHandler handles public product operations (no auth required)
type PublicHandler struct {
	productService *services.ProductService
	reviewService  *services.ReviewService
	contactService *services.ContactService
}

// NewPublicHandler creates a new public handler
func NewPublicHandler(productService *services.ProductService, reviewService *services.ReviewService, contactService *services.ContactService) *PublicHandler {
	return &PublicHandler{
		productService: productService,
		reviewService:  reviewService,
		contactService: contactService,
	}
}

// ListProducts lists active products for a domain (public access)
func (h *PublicHandler) ListProducts(c echo.Context) error {
	domain := c.Param("domain")

	// Parse attribute filters from query params
	attributes := make(map[string]string)
	for key, values := range c.QueryParams() {
		if len(values) > 0 {
			attributes[key] = values[0]
		}
	}

	// Only show active products on public API
	products, err := h.productService.ListProducts(c.Request().Context(), domain, true, attributes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch products",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"products": products,
		"count":    len(products),
	})
}

// GetProduct retrieves a single product by ID (public access)
func (h *PublicHandler) GetProduct(c echo.Context) error {
	domain := c.Param("domain")
	productID := c.Param("id")

	product, err := h.productService.GetProduct(c.Request().Context(), domain, productID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Product not found",
		})
	}

	// Only return if product is active
	if !product.Active {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Product not found",
		})
	}

	return c.JSON(http.StatusOK, product)
}

// SearchProducts searches products by text (public access)
func (h *PublicHandler) SearchProducts(c echo.Context) error {
	domain := c.Param("domain")
	query := c.QueryParam("q")

	if query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Search query is required (use ?q=search-term)",
		})
	}

	// Only show active products on public API
	products, err := h.productService.SearchProducts(c.Request().Context(), domain, query, true)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to search products",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"products": products,
		"count":    len(products),
		"query":    query,
	})
}

// GetPromotions returns active promotions/sales for a domain
func (h *PublicHandler) GetPromotions(c echo.Context) error {
	domain := c.Param("domain")

	// Get all active products to check for discounts
	products, err := h.productService.ListProducts(c.Request().Context(), domain, true, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch promotions",
		})
	}

	// Check if any products have active discounts
	hasActiveDiscounts := false
	maxDiscount := 0

	for _, product := range products {
		if product.Discount != nil && product.Discount.IsDiscountActive() {
			hasActiveDiscounts = true
			discountPercent := product.GetDiscountPercentage()
			if discountPercent > maxDiscount {
				maxDiscount = discountPercent
			}
		}
	}

	if !hasActiveDiscounts {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"active":      false,
			"promotions":  []interface{}{},
		})
	}

	// Return active promotion info
	return c.JSON(http.StatusOK, map[string]interface{}{
		"active": true,
		"promotions": []map[string]interface{}{
			{
				"type":         "sale",
				"message":      fmt.Sprintf("Sale - Up to %d%% Off Select Styles", maxDiscount),
				"max_discount": maxDiscount,
			},
		},
	})
}

// CreateReview creates a new review (public access)
func (h *PublicHandler) CreateReview(c echo.Context) error {
	domain := c.Param("domain")

	var req models.CreateReviewRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate rating
	if req.Rating < 1 || req.Rating > 5 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Rating must be between 1 and 5",
		})
	}

	// Validate required fields
	if req.Product == "" || req.Name == "" || req.Text == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Product name, reviewer name, and review text are required",
		})
	}

	review, err := h.reviewService.CreateReview(c.Request().Context(), domain, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create review",
		})
	}

	return c.JSON(http.StatusCreated, review)
}

// ListReviews lists all reviews for a domain (public access)
func (h *PublicHandler) ListReviews(c echo.Context) error {
	domain := c.Param("domain")
	productID := c.QueryParam("product_id")

	reviews, err := h.reviewService.GetReviews(c.Request().Context(), domain, productID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch reviews",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"reviews": reviews,
		"count":   len(reviews),
	})
}

// GetReview retrieves a single review by ID (public access)
func (h *PublicHandler) GetReview(c echo.Context) error {
	domain := c.Param("domain")
	reviewID := c.Param("id")

	review, err := h.reviewService.GetReviewByID(c.Request().Context(), domain, reviewID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Review not found",
		})
	}

	return c.JSON(http.StatusOK, review)
}

// CreateContact creates a new contact submission (public access)
func (h *PublicHandler) CreateContact(c echo.Context) error {
	domain := c.Param("domain")

	var req models.CreateContactRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.Name == "" || req.Email == "" || req.Subject == "" || req.Message == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Name, email, subject, and message are required",
		})
	}

	contact, err := h.contactService.CreateContact(c.Request().Context(), domain, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to submit contact form",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Thank you for contacting us! We'll get back to you soon.",
		"id":      contact.ID.Hex(),
	})
}

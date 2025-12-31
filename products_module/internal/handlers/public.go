package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sparque/products_module/internal/services"
)

// PublicHandler handles public product operations (no auth required)
type PublicHandler struct {
	productService *services.ProductService
}

// NewPublicHandler creates a new public handler
func NewPublicHandler(productService *services.ProductService) *PublicHandler {
	return &PublicHandler{
		productService: productService,
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

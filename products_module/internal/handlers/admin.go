package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sparque/products_module/internal/models"
	"github.com/sparque/products_module/internal/services"
)

// AdminHandler handles admin product operations (requires API key)
type AdminHandler struct {
	productService *services.ProductService
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(productService *services.ProductService) *AdminHandler {
	return &AdminHandler{
		productService: productService,
	}
}

// CreateProduct creates a new product
func (h *AdminHandler) CreateProduct(c echo.Context) error {
	var req models.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Get domain and API key ID from context (set by middleware)
	domain, _ := c.Get("domain").(string)
	apiKeyID, _ := c.Get("api_key_id").(string)

	product, err := h.productService.CreateProduct(c.Request().Context(), domain, req, apiKeyID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, product)
}

// ListProducts lists all products for the authenticated domain
func (h *AdminHandler) ListProducts(c echo.Context) error {
	domain, _ := c.Get("domain").(string)

	// Query parameters
	activeOnly := c.QueryParam("active") == "true"

	// Parse attribute filters from query params
	attributes := make(map[string]string)
	for key, values := range c.QueryParams() {
		if key != "active" && len(values) > 0 {
			attributes[key] = values[0]
		}
	}

	products, err := h.productService.ListProducts(c.Request().Context(), domain, activeOnly, attributes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"products": products,
		"count":    len(products),
	})
}

// GetProduct retrieves a product by ID
func (h *AdminHandler) GetProduct(c echo.Context) error {
	domain, _ := c.Get("domain").(string)
	productID := c.Param("id")

	product, err := h.productService.GetProduct(c.Request().Context(), domain, productID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Product not found",
		})
	}

	return c.JSON(http.StatusOK, product)
}

// UpdateProduct updates a product
func (h *AdminHandler) UpdateProduct(c echo.Context) error {
	domain, _ := c.Get("domain").(string)
	productID := c.Param("id")

	var req models.UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	product, err := h.productService.UpdateProduct(c.Request().Context(), domain, productID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, product)
}

// DeleteProduct deletes a product
func (h *AdminHandler) DeleteProduct(c echo.Context) error {
	domain, _ := c.Get("domain").(string)
	productID := c.Param("id")

	// Check for hard delete flag
	hardDelete := c.QueryParam("hard") == "true"

	err := h.productService.DeleteProduct(c.Request().Context(), domain, productID, hardDelete)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	message := "Product deactivated successfully"
	if hardDelete {
		message = "Product deleted successfully"
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": message,
	})
}

// UpdateStock updates stock for a specific variant
func (h *AdminHandler) UpdateStock(c echo.Context) error {
	domain, _ := c.Get("domain").(string)
	productID := c.Param("id")
	variantID := c.Param("variantId")

	var req struct {
		Stock int `json:"stock" validate:"gte=0"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	err := h.productService.UpdateStock(c.Request().Context(), domain, productID, variantID, req.Stock)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Stock updated successfully",
	})
}

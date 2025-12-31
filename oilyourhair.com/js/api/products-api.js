/**
 * Products API Client for oilyourhair.com
 * Handles all communication with the products module API
 */

class ProductsAPI {
    constructor(baseUrl, domain, apiKey = null) {
        this.baseUrl = baseUrl || 'http://localhost:9091';
        this.domain = domain || 'oilyourhair.com';
        this.apiKey = apiKey;
    }

    // Public API endpoints (no auth required)

    async getProducts(filters = {}) {
        const params = new URLSearchParams(filters);
        const url = `${this.baseUrl}/api/v1/public/${this.domain}/products?${params}`;

        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`Failed to fetch products: ${response.statusText}`);
        }

        return await response.json();
    }

    async getProduct(productId) {
        const url = `${this.baseUrl}/api/v1/public/${this.domain}/products/${productId}`;

        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`Failed to fetch product: ${response.statusText}`);
        }

        return await response.json();
    }

    async searchProducts(query) {
        const url = `${this.baseUrl}/api/v1/public/${this.domain}/products/search?q=${encodeURIComponent(query)}`;

        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`Failed to search products: ${response.statusText}`);
        }

        return await response.json();
    }

    // Admin API endpoints (requires API key)

    async createProduct(productData) {
        if (!this.apiKey) {
            throw new Error('API key required for admin operations');
        }

        const response = await fetch(`${this.baseUrl}/api/v1/products`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${this.apiKey}`
            },
            body: JSON.stringify(productData)
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || `Failed to create product: ${response.statusText}`);
        }

        return await response.json();
    }

    async updateProduct(productId, updates) {
        if (!this.apiKey) {
            throw new Error('API key required for admin operations');
        }

        const response = await fetch(`${this.baseUrl}/api/v1/products/${productId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${this.apiKey}`
            },
            body: JSON.stringify(updates)
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || `Failed to update product: ${response.statusText}`);
        }

        return await response.json();
    }

    async deleteProduct(productId, hardDelete = false) {
        if (!this.apiKey) {
            throw new Error('API key required for admin operations');
        }

        const url = `${this.baseUrl}/api/v1/products/${productId}${hardDelete ? '?hard=true' : ''}`;

        const response = await fetch(url, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${this.apiKey}`
            }
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || `Failed to delete product: ${response.statusText}`);
        }

        return await response.json();
    }

    async listAllProducts(filters = {}) {
        if (!this.apiKey) {
            throw new Error('API key required for admin operations');
        }

        const params = new URLSearchParams(filters);
        const url = `${this.baseUrl}/api/v1/products?${params}`;

        const response = await fetch(url, {
            headers: {
                'Authorization': `Bearer ${this.apiKey}`
            }
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || `Failed to list products: ${response.statusText}`);
        }

        return await response.json();
    }

    async updateStock(productId, variantId, stock) {
        if (!this.apiKey) {
            throw new Error('API key required for admin operations');
        }

        const response = await fetch(
            `${this.baseUrl}/api/v1/products/${productId}/variants/${variantId}/stock`,
            {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${this.apiKey}`
                },
                body: JSON.stringify({ stock })
            }
        );

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || `Failed to update stock: ${response.statusText}`);
        }

        return await response.json();
    }
}

// Export for use in other files
window.ProductsAPI = ProductsAPI;

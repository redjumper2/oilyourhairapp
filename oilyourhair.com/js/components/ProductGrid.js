/**
 * ProductGrid Component
 * Displays products in a responsive grid
 */

class ProductGrid {
    constructor(container, api) {
        this.container = container;
        this.api = api;
        this.products = [];
        this.filters = {};
    }

    async load(filters = {}) {
        this.filters = filters;
        this.showLoading();

        try {
            const data = await this.api.getProducts(filters);
            this.products = data.products || [];
            this.render();
        } catch (error) {
            this.showError(error.message);
        }
    }

    showLoading() {
        this.container.innerHTML = `
            <div class="loading-state">
                <div class="spinner"></div>
                <p>Loading products...</p>
            </div>
        `;
    }

    showError(message) {
        this.container.innerHTML = `
            <div class="error-state">
                <p>⚠️ ${message}</p>
                <button onclick="location.reload()">Try Again</button>
            </div>
        `;
    }

    render() {
        if (this.products.length === 0) {
            this.container.innerHTML = `
                <div class="empty-state">
                    <p>No products found</p>
                </div>
            `;
            return;
        }

        const productsHTML = this.products.map(product => this.renderProductCard(product)).join('');

        this.container.innerHTML = `
            <div class="products-grid">
                ${productsHTML}
            </div>
        `;
    }

    renderProductCard(product) {
        const imageUrl = product.images && product.images[0]
            ? product.images[0]
            : 'https://via.placeholder.com/300x300?text=No+Image';

        const price = this.getPrice(product);
        const stockStatus = this.getStockStatus(product);

        return `
            <div class="product-card" data-product-id="${product.id}">
                <div class="product-image">
                    <img src="${imageUrl}" alt="${product.name}" loading="lazy">
                    ${stockStatus.badge}
                </div>
                <div class="product-info">
                    <h3 class="product-name">${product.name}</h3>
                    <p class="product-description">${this.truncate(product.description, 80)}</p>
                    <div class="product-meta">
                        ${this.renderAttributes(product.attributes)}
                    </div>
                    <div class="product-footer">
                        <span class="product-price">${price}</span>
                        <a href="product.html?id=${product.id}" class="btn btn-primary">View Details</a>
                    </div>
                </div>
            </div>
        `;
    }

    getPrice(product) {
        if (product.variants && product.variants.length > 0) {
            const prices = product.variants.map(v => v.price || product.base_price);
            const minPrice = Math.min(...prices);
            const maxPrice = Math.max(...prices);

            if (minPrice === maxPrice) {
                return `$${minPrice.toFixed(2)}`;
            }
            return `$${minPrice.toFixed(2)} - $${maxPrice.toFixed(2)}`;
        }
        return `$${product.base_price.toFixed(2)}`;
    }

    getStockStatus(product) {
        if (!product.variants || product.variants.length === 0) {
            return { inStock: true, badge: '' };
        }

        const totalStock = product.variants.reduce((sum, v) => sum + (v.stock || 0), 0);

        if (totalStock === 0) {
            return {
                inStock: false,
                badge: '<span class="stock-badge out-of-stock">Out of Stock</span>'
            };
        }

        if (totalStock < 10) {
            return {
                inStock: true,
                badge: '<span class="stock-badge low-stock">Low Stock</span>'
            };
        }

        return { inStock: true, badge: '' };
    }

    renderAttributes(attributes) {
        if (!attributes || Object.keys(attributes).length === 0) {
            return '';
        }

        const tags = Object.entries(attributes)
            .filter(([key]) => key !== 'category') // Category shown separately
            .slice(0, 3) // Show max 3 attributes
            .map(([key, value]) => `<span class="attribute-tag">${value}</span>`)
            .join('');

        return `<div class="product-attributes">${tags}</div>`;
    }

    truncate(text, maxLength) {
        if (!text) return '';
        if (text.length <= maxLength) return text;
        return text.substr(0, maxLength) + '...';
    }
}

window.ProductGrid = ProductGrid;

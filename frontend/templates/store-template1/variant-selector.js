/**
 * Generic Variant Selector Component
 * Works with any product type (shoes, oils, clothing, etc.)
 * Handles size, color, volume, and custom attributes
 */

class VariantSelector {
    constructor(product, onVariantChange) {
        this.product = product;
        this.onVariantChange = onVariantChange;
        this.selectedVariant = null;

        // Extract variant attributes (e.g., size, color, volume)
        this.variantAttributes = this.extractVariantAttributes();
        this.selectedAttributes = {};
    }

    /**
     * Extract unique attributes from all variants
     * Returns object like: { size: ['S', 'M', 'L'], color: ['Black', 'White'] }
     */
    extractVariantAttributes() {
        if (!this.product.variants || this.product.variants.length === 0) {
            return {};
        }

        const attributeMap = {};

        this.product.variants.forEach(variant => {
            if (!variant.attributes) return;

            Object.keys(variant.attributes).forEach(key => {
                if (!attributeMap[key]) {
                    attributeMap[key] = new Set();
                }
                attributeMap[key].add(variant.attributes[key]);
            });
        });

        // Convert Sets to sorted arrays
        const result = {};
        Object.keys(attributeMap).forEach(key => {
            result[key] = Array.from(attributeMap[key]).sort();
        });

        return result;
    }

    /**
     * Find matching variant based on selected attributes
     */
    findMatchingVariant() {
        if (!this.product.variants) return null;

        return this.product.variants.find(variant => {
            if (!variant.attributes) return false;

            return Object.keys(this.selectedAttributes).every(key => {
                return variant.attributes[key] === this.selectedAttributes[key];
            });
        });
    }

    /**
     * Check if a specific attribute value is available (in stock)
     */
    isAttributeAvailable(attributeKey, attributeValue) {
        if (!this.product.variants) return false;

        // Create temporary selection with this value
        const tempSelection = {
            ...this.selectedAttributes,
            [attributeKey]: attributeValue
        };

        // Find variants that match this selection
        const matchingVariants = this.product.variants.filter(variant => {
            if (!variant.attributes) return false;

            return Object.keys(tempSelection).every(key => {
                return variant.attributes[key] === tempSelection[key];
            });
        });

        // Check if any matching variant has stock
        return matchingVariants.some(v => v.stock && v.stock > 0);
    }

    /**
     * Handle attribute selection
     */
    selectAttribute(attributeKey, attributeValue) {
        this.selectedAttributes[attributeKey] = attributeValue;

        // Find matching variant
        this.selectedVariant = this.findMatchingVariant();

        // Notify callback
        if (this.onVariantChange) {
            this.onVariantChange(this.selectedVariant);
        }

        return this.selectedVariant;
    }

    /**
     * Render variant selector UI
     * Returns HTML string
     */
    render() {
        if (Object.keys(this.variantAttributes).length === 0) {
            return ''; // No variants to display
        }

        let html = '<div class="variant-selector">';

        Object.keys(this.variantAttributes).forEach(attributeKey => {
            const values = this.variantAttributes[attributeKey];
            const displayName = this.formatAttributeName(attributeKey);

            html += `
                <div class="variant-group">
                    <label class="variant-label">${displayName}</label>
                    <div class="variant-options">
                        ${values.map(value => {
                            const isAvailable = this.isAttributeAvailable(attributeKey, value);
                            const isSelected = this.selectedAttributes[attributeKey] === value;
                            const className = `variant-option ${isSelected ? 'selected' : ''} ${!isAvailable ? 'unavailable' : ''}`;

                            return `
                                <button
                                    class="${className}"
                                    data-attribute="${attributeKey}"
                                    data-value="${value}"
                                    ${!isAvailable ? 'disabled' : ''}
                                >
                                    ${value}
                                </button>
                            `;
                        }).join('')}
                    </div>
                </div>
            `;
        });

        html += '</div>';

        // Add selected variant info
        if (this.selectedVariant) {
            const stockStatus = this.selectedVariant.stock > 0
                ? `<span class="stock-available">In Stock (${this.selectedVariant.stock} available)</span>`
                : `<span class="stock-unavailable">Out of Stock</span>`;

            html += `
                <div class="variant-info">
                    ${this.selectedVariant.sku ? `<div class="variant-sku">SKU: ${this.selectedVariant.sku}</div>` : ''}
                    <div class="variant-stock">${stockStatus}</div>
                    ${this.selectedVariant.price !== this.product.base_price ?
                        `<div class="variant-price">Price: $${this.selectedVariant.price.toFixed(2)}</div>` : ''}
                </div>
            `;
        }

        return html;
    }

    /**
     * Format attribute name for display
     * e.g., "bottle_size" -> "Bottle Size"
     */
    formatAttributeName(key) {
        return key
            .split('_')
            .map(word => word.charAt(0).toUpperCase() + word.slice(1))
            .join(' ');
    }

    /**
     * Get currently selected variant
     */
    getSelectedVariant() {
        return this.selectedVariant;
    }

    /**
     * Check if all required attributes are selected
     */
    isComplete() {
        const requiredAttributes = Object.keys(this.variantAttributes);
        return requiredAttributes.every(key => this.selectedAttributes[key]);
    }
}

// Export for use in other scripts
if (typeof module !== 'undefined' && module.exports) {
    module.exports = VariantSelector;
}

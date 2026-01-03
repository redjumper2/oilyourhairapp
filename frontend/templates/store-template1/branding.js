/**
 * Branding System - Loads and applies brand configuration
 * Makes the e-commerce template white-label ready
 */

class BrandingManager {
    constructor() {
        this.config = null;
        this.brandingPath = '/branding.json';
    }

    async loadBranding() {
        try {
            const response = await fetch(this.brandingPath);
            if (!response.ok) {
                console.warn('Branding config not found, using defaults');
                return;
            }
            this.config = await response.json();
            this.applyBranding();
        } catch (error) {
            console.error('Failed to load branding:', error);
        }
    }

    applyBranding() {
        if (!this.config) return;

        // Apply color scheme via CSS variables
        this.applyColors();

        // Apply typography
        this.applyTypography();

        // Apply layout settings
        this.applyLayout();

        // Update logo and brand name
        this.updateBrandElements();

        console.log('✅ Branding applied:', this.config.brand_name);
    }

    applyColors() {
        const colors = this.config.colors;
        const root = document.documentElement;

        // Set CSS custom properties for colors
        root.style.setProperty('--brand-primary', colors.primary);
        root.style.setProperty('--brand-primary-dark', colors.primary_dark);
        root.style.setProperty('--brand-primary-light', colors.primary_light);
        root.style.setProperty('--brand-secondary', colors.secondary);
        root.style.setProperty('--brand-accent', colors.accent);
        root.style.setProperty('--brand-background', colors.background);
        root.style.setProperty('--brand-text', colors.text);
        root.style.setProperty('--brand-text-light', colors.text_light);
        root.style.setProperty('--brand-border', colors.border);
        root.style.setProperty('--brand-success', colors.success);
        root.style.setProperty('--brand-error', colors.error);
        root.style.setProperty('--brand-warning', colors.warning);
    }

    applyTypography() {
        const typography = this.config.typography;
        const root = document.documentElement;

        root.style.setProperty('--brand-font-family', typography.font_family);
        root.style.setProperty('--brand-heading-font', typography.heading_font);
        root.style.setProperty('--brand-h1-size', typography.h1_size);
        root.style.setProperty('--brand-h2-size', typography.h2_size);
        root.style.setProperty('--brand-h3-size', typography.h3_size);
        root.style.setProperty('--brand-body-size', typography.body_size);
        root.style.setProperty('--brand-small-size', typography.small_size);
    }

    applyLayout() {
        const layout = this.config.layout;
        const root = document.documentElement;

        root.style.setProperty('--brand-max-width', layout.max_width);
        root.style.setProperty('--brand-border-radius', layout.border_radius);
        root.style.setProperty('--brand-card-shadow', layout.card_shadow);

        // Spacing system
        const spacing = layout.spacing;
        root.style.setProperty('--brand-spacing-xs', spacing.xs);
        root.style.setProperty('--brand-spacing-sm', spacing.sm);
        root.style.setProperty('--brand-spacing-md', spacing.md);
        root.style.setProperty('--brand-spacing-lg', spacing.lg);
        root.style.setProperty('--brand-spacing-xl', spacing.xl);
    }

    updateBrandElements() {
        // Update page title
        const titleElement = document.querySelector('title');
        if (titleElement && !titleElement.textContent.includes(this.config.brand_name)) {
            titleElement.textContent = `${this.config.brand_name} - ${titleElement.textContent}`;
        }

        // Update logo elements
        const logoElements = document.querySelectorAll('[data-brand-logo]');
        logoElements.forEach(el => {
            if (this.config.logo.url) {
                el.innerHTML = `<img src="${this.config.logo.url}" alt="${this.config.brand_name}" style="height: ${this.config.logo.height}">`;
            } else {
                el.textContent = this.config.logo.text;
            }
        });

        // Update brand name elements
        const brandNameElements = document.querySelectorAll('[data-brand-name]');
        brandNameElements.forEach(el => {
            el.textContent = this.config.brand_name;
        });

        // Update tagline elements
        const taglineElements = document.querySelectorAll('[data-brand-tagline]');
        taglineElements.forEach(el => {
            el.textContent = this.config.tagline;
        });
    }

    getBrandConfig() {
        return this.config;
    }
}

// Create global branding instance
const branding = new BrandingManager();

// Auto-load branding when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => branding.loadBranding());
} else {
    branding.loadBranding();
}

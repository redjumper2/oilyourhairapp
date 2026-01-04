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
        if (!this.config) {
            console.warn('No config found, using defaults');
            return;
        }

        const root = document.documentElement;

        // Support both nested (colors.primary) and flat (primaryColor) structure
        const primaryColor = this.config.colors?.primary || this.config.primaryColor;
        const primaryLight = this.config.colors?.primary_light || this.config.primaryColorLight;
        const primaryDark = this.config.colors?.primary_dark || this.config.primaryColorDark;

        // Set CSS custom properties for colors
        if (primaryColor) root.style.setProperty('--brand-primary', primaryColor);
        if (primaryDark) root.style.setProperty('--brand-primary-dark', primaryDark);
        if (primaryLight) root.style.setProperty('--brand-primary-light', primaryLight);

        // Optional colors with defaults
        const secondary = this.config.colors?.secondary || this.config.secondaryColor || primaryColor;
        const accent = this.config.colors?.accent || this.config.accentColor || primaryLight;

        if (secondary) root.style.setProperty('--brand-secondary', secondary);
        if (accent) root.style.setProperty('--brand-accent', accent);
        if (this.config.colors?.background) root.style.setProperty('--brand-background', this.config.colors.background);
        if (this.config.colors?.text) root.style.setProperty('--brand-text', this.config.colors.text);
        if (this.config.colors?.text_light) root.style.setProperty('--brand-text-light', this.config.colors.text_light);
        if (this.config.colors?.border) root.style.setProperty('--brand-border', this.config.colors.border);
        if (this.config.colors?.success) root.style.setProperty('--brand-success', this.config.colors.success);
        if (this.config.colors?.error) root.style.setProperty('--brand-error', this.config.colors.error);
        if (this.config.colors?.warning) root.style.setProperty('--brand-warning', this.config.colors.warning);
    }

    applyTypography() {
        if (!this.config) {
            return;
        }

        // Only apply if typography config exists
        if (!this.config.typography) {
            return;
        }

        const typography = this.config.typography;
        const root = document.documentElement;

        if (typography.font_family) root.style.setProperty('--brand-font-family', typography.font_family);
        if (typography.heading_font) root.style.setProperty('--brand-heading-font', typography.heading_font);
        if (typography.h1_size) root.style.setProperty('--brand-h1-size', typography.h1_size);
        if (typography.h2_size) root.style.setProperty('--brand-h2-size', typography.h2_size);
        if (typography.h3_size) root.style.setProperty('--brand-h3-size', typography.h3_size);
        if (typography.body_size) root.style.setProperty('--brand-body-size', typography.body_size);
        if (typography.small_size) root.style.setProperty('--brand-small-size', typography.small_size);
    }

    applyLayout() {
        if (!this.config) {
            return;
        }

        // Only apply if layout config exists
        if (!this.config.layout) {
            return;
        }

        const layout = this.config.layout;
        const root = document.documentElement;

        if (layout.max_width) root.style.setProperty('--brand-max-width', layout.max_width);
        if (layout.border_radius) root.style.setProperty('--brand-border-radius', layout.border_radius);
        if (layout.card_shadow) root.style.setProperty('--brand-card-shadow', layout.card_shadow);

        // Spacing system
        if (layout.spacing) {
            const spacing = layout.spacing;
            if (spacing.xs) root.style.setProperty('--brand-spacing-xs', spacing.xs);
            if (spacing.sm) root.style.setProperty('--brand-spacing-sm', spacing.sm);
            if (spacing.md) root.style.setProperty('--brand-spacing-md', spacing.md);
            if (spacing.lg) root.style.setProperty('--brand-spacing-lg', spacing.lg);
            if (spacing.xl) root.style.setProperty('--brand-spacing-xl', spacing.xl);
        }
    }

    updateBrandElements() {
        if (!this.config) {
            return;
        }

        // Support both brand_name and brandName
        const brandName = this.config.brand_name || this.config.brandName;
        const siteName = this.config.site_name || this.config.siteName || brandName;

        // Update page title
        const titleElement = document.querySelector('title');
        if (titleElement && siteName && !titleElement.textContent.includes(siteName)) {
            titleElement.textContent = `${siteName} - ${titleElement.textContent}`;
        }

        // Update logo elements - support both nested and flat structure
        const logoUrl = this.config.logo?.url || this.config.logo;
        const logoText = this.config.logo?.text || brandName;
        const logoHeight = this.config.logo?.height || '40px';

        if (logoUrl || logoText) {
            const logoElements = document.querySelectorAll('[data-brand-logo]');
            logoElements.forEach(el => {
                if (typeof logoUrl === 'string' && logoUrl.startsWith('http')) {
                    el.innerHTML = `<img src="${logoUrl}" alt="${brandName || 'Logo'}" style="height: ${logoHeight}">`;
                } else if (logoText) {
                    el.textContent = logoText;
                }
            });
        }

        // Update brand name elements
        if (brandName) {
            const brandNameElements = document.querySelectorAll('[data-brand-name]');
            brandNameElements.forEach(el => {
                el.textContent = brandName;
            });
        }

        // Update tagline elements
        const tagline = this.config.tagline;
        if (tagline) {
            const taglineElements = document.querySelectorAll('[data-brand-tagline]');
            taglineElements.forEach(el => {
                el.textContent = tagline;
            });
        }
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

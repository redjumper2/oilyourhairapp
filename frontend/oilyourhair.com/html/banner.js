/**
 * Promotion Banner - Displays active sales and promotions
 * Fetches from promotions API and shows/hides dynamically
 */

class PromotionBanner {
    constructor() {
        this.bannerElement = null;
        this.apiEndpoint = null;
        this.checkInterval = 5 * 60 * 1000; // Check every 5 minutes
        this.intervalId = null;
    }

    async init(domain) {
        this.apiEndpoint = `/api/v1/public/${domain}/promotions`;
        this.createBannerElement();
        await this.fetchAndDisplay();

        // Periodically check for promotion updates
        this.intervalId = setInterval(() => this.fetchAndDisplay(), this.checkInterval);
    }

    createBannerElement() {
        // Create banner container
        this.bannerElement = document.createElement('div');
        this.bannerElement.id = 'promotion-banner';
        this.bannerElement.className = 'promotion-banner hidden';

        // Insert at the very top of body
        document.body.insertBefore(this.bannerElement, document.body.firstChild);
    }

    async fetchAndDisplay() {
        try {
            const response = await fetch(this.apiEndpoint);
            if (!response.ok) {
                console.warn('Failed to fetch promotions');
                this.hideBanner();
                return;
            }

            const data = await response.json();

            if (data.active && data.promotions && data.promotions.length > 0) {
                this.showBanner(data.promotions[0]);
            } else {
                this.hideBanner();
            }
        } catch (error) {
            console.error('Error fetching promotions:', error);
            this.hideBanner();
        }
    }

    showBanner(promotion) {
        const message = promotion.message || 'Special Offer Available';

        this.bannerElement.innerHTML = `
            <div class="promotion-banner-content">
                <span class="promotion-message">${message}</span>
                <button class="banner-close" aria-label="Close banner">&times;</button>
            </div>
        `;

        // Add close button handler
        const closeBtn = this.bannerElement.querySelector('.banner-close');
        closeBtn.addEventListener('click', () => this.hideBanner());

        // Show banner
        this.bannerElement.classList.remove('hidden');
    }

    hideBanner() {
        if (this.bannerElement) {
            this.bannerElement.classList.add('hidden');
        }
    }

    destroy() {
        if (this.intervalId) {
            clearInterval(this.intervalId);
        }
        if (this.bannerElement) {
            this.bannerElement.remove();
        }
    }
}

// Auto-initialize banner when branding is loaded
let promotionBanner = null;

function initPromotionBanner() {
    if (!branding || !branding.getBrandConfig()) {
        // Branding not loaded yet, try again in 100ms
        setTimeout(initPromotionBanner, 100);
        return;
    }

    const config = branding.getBrandConfig();
    if (config && config.domain) {
        promotionBanner = new PromotionBanner();
        promotionBanner.init(config.domain);
        console.log('✅ Promotion banner initialized');
    }
}

// Start initialization when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initPromotionBanner);
} else {
    initPromotionBanner();
}

# OilYourHair App - TODO List

## Integration Testing - Reviews (Priority 1)

**Status:** In Progress

- [x] **Add review creation tests to test-integration.sh**
  - Create test product for reviews
  - Test review creation via admin API
  - Test listing reviews via public API
  - Test filtering by product_id
  - Test filtering by min_rating
  - Create multiple reviews
  - Verify review counts

- [ ] **Run integration tests and verify**
  - Execute `./test-integration.sh`
  - Verify all 24 tests pass
  - Check JSON responses in logs

---

## Multi-Currency Support Implementation

### Backend Tasks

- [ ] **Add currency fields to product model**
  - Add `base_price_usd` field
  - Add `currency_overrides` object for manual pricing per currency
  - Update product creation/update endpoints
  - Migration script for existing products

- [ ] **Create currency conversion service**
  - Integrate exchange rate API (Open Exchange Rates or Fixer.io)
  - Implement rate caching (24-hour cache)
  - Add manual exchange rate override functionality
  - Create API endpoint to fetch current rates
  - Handle currency conversion logic

### Frontend Tasks

- [ ] **Add currency selector dropdown to header**
  - Design currency selector UI (🌐 USD | INR | EUR...)
  - Position in header navigation
  - Style dropdown menu
  - Handle currency change events

- [ ] **Update product display for multi-currency**
  - Convert prices to selected currency
  - Update all product listings (shop page, index page)
  - Update product detail views
  - Update "Add to Cart" price displays

- [ ] **Implement currency formatting per locale**
  - Format USD: $99.99
  - Format INR: ₹8,299
  - Format EUR: €99,99
  - Handle decimal places per currency
  - Add proper thousand separators

- [ ] **Update cart and checkout for multi-currency**
  - Convert cart totals to selected currency
  - Show subtotal, tax, shipping in selected currency
  - Update checkout display
  - Handle currency in order summary

- [ ] **Store user's currency preference**
  - Save selected currency to localStorage
  - Load currency preference on page load
  - Persist across page navigation
  - Handle fallback if no preference set

- [ ] **Add auto-detection of user's currency**
  - Detect user's country via browser locale
  - Map country to default currency
  - Auto-select currency on first visit
  - Allow user to override auto-detected currency

---

## Supported Currencies (Initial Launch)

- **USD** - United States Dollar ($)
- **INR** - Indian Rupee (₹)
- **EUR** - Euro (€) - Optional
- **GBP** - British Pound (£) - Optional

---

## Notes

- Base prices stored in USD
- Exchange rates cached for 24 hours
- Manual price overrides allowed for strategic pricing
- All conversions happen client-side for performance

---

## Embeddable Widgets (Optional - Future Enhancement)

**Goal:** Make features embeddable on multiple domains (e.g., WordPress, Shopify, any website)

### Approach: Web Components + Existing Go API

- [ ] **Create Web Components architecture**
  - Build custom elements using Lit or vanilla Web Components
  - Encapsulate styles using Shadow DOM
  - Single JS bundle for all widgets
  - CDN hosting for widgets.js

- [ ] **Build embeddable components**
  - `<oilyourhair-products domain="example.com">`
  - `<oilyourhair-reviews domain="example.com" product-id="123">`
  - `<oilyourhair-cart domain="example.com">`
  - `<oilyourhair-contact domain="example.com">`
  - `<oilyourhair-auth domain="example.com">`

- [ ] **Update backend for CORS**
  - Configure CORS to allow embedding domains
  - Add domain whitelist for security
  - Handle cross-origin authentication
  - API key per domain for usage tracking

- [ ] **Create widget configuration**
  - Admin panel to generate embed codes
  - Customization options (colors, sizes, features)
  - Domain registration/whitelisting
  - Usage analytics per domain

- [ ] **Documentation & Examples**
  - Embed code generator
  - Integration guides (WordPress, Shopify, HTML)
  - Customization documentation
  - Demo page showing all widgets

### Benefits
- One-line embed: `<script src="https://widgets.oilyourhair.com/v1/widgets.js"></script>`
- Works on any website
- Centralized updates (fix bugs once, all sites benefit)
- Analytics across all embedded instances
- Revenue opportunity (charge for widget usage)

---

**Created:** 2026-01-03
**Priority:** High (Multi-Currency), Medium (Embeddable Widgets)
**Target:** Next development cycle

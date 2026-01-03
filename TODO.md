# OilYourHair App - TODO List

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

**Created:** 2026-01-03
**Priority:** High
**Target:** Next development cycle

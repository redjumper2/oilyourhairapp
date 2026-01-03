# Store Template 1 - Modern E-Commerce

An Allbirds-inspired clean and modern e-commerce template with full white-label capabilities.

## Features

### 🎨 Dynamic Branding System
- CSS custom properties for theming
- Automatic color application from `branding.json`
- Customizable typography and layout
- Logo and brand name injection

### 🛍️ Product Management
- **Product Variants**: Support for any attribute type (size, color, volume, etc.)
- **Variant Selector**: Smart component that works with any product type
- **Stock Management**: Visual indicators for availability
- **Product Images**: Multi-image support with hover zoom

### 💰 Discount & Sales
- **Percentage & Fixed Discounts**: Flexible discount types
- **Time-Based Sales**: Start/end date support
- **Sale Badges**: Automatic "50% OFF" badges
- **Dynamic Banner**: Top banner shows/hides based on active promotions

### 🛒 Shopping Cart
- Add to cart functionality
- Quantity controls
- Discounted pricing in cart
- Total calculation

### 👤 Authentication
- Login/logout integration
- Protected routes
- User status display

### 📱 Responsive Design
- Mobile-first approach
- Tablet and desktop optimized
- Touch-friendly interactions

### ⚡ Performance
- Optimized images
- Minimal dependencies
- Fast page loads
- Smooth animations

## Template Files

### Pages
- **index.html** - Homepage with hero, mission, and featured products
- **shop.html** - Full product catalog with filters and cart
- **about.html** - About page
- **contact.html** - Contact form
- **reviews.html** - Customer reviews

### Components
- **branding.js** - Loads and applies brand configuration
- **banner.js** - Dynamic promotion banner
- **variant-selector.js** - Generic product variant selector
- **auth.js** - Authentication integration

### Styles
- **styles.css** - Global styles with CSS custom properties

## Configuration

### Required: branding.json

```json
{
  "domain": "yourstore.com",
  "brand_name": "Your Store",
  "tagline": "Your Tagline",
  "colors": {
    "primary": "#2E7D32",
    "primary_dark": "#1B5E20",
    "primary_light": "#4CAF50",
    ...
  },
  "typography": {
    "font_family": "-apple-system, sans-serif",
    ...
  },
  "layout": {
    "max_width": "1400px",
    ...
  },
  "logo": {
    "text": "Your Store",
    "url": "/logo.png",
    "height": "40px"
  }
}
```

### CSS Custom Properties

The template uses these CSS variables (auto-applied from branding.json):

```css
--brand-primary          /* Primary brand color */
--brand-primary-dark     /* Darker shade for hover states */
--brand-primary-light    /* Lighter shade for backgrounds */
--brand-secondary        /* Secondary color */
--brand-accent           /* Accent color */
--brand-background       /* Page background */
--brand-text             /* Primary text color */
--brand-text-light       /* Secondary text color */
--brand-border           /* Border color */
--brand-success          /* Success state color */
--brand-error            /* Error state color */
--brand-warning          /* Warning state color */
```

## Product Requirements

### Product Data Structure

Products must have this structure (via Products API):

```json
{
  "id": "...",
  "name": "Product Name",
  "description": "Product description",
  "base_price": 39.99,
  "images": ["https://..."],
  "attributes": {
    "type": "oil",
    "badge": "Bestseller",
    "features": "Organic,Vegan,Cruelty-Free"
  },
  "variants": [
    {
      "id": "variant-1",
      "attributes": {
        "size": "50ml",
        "scent": "Lavender"
      },
      "price": 39.99,
      "stock": 10,
      "sku": "PROD-50ML-LAV",
      "image_index": 0
    }
  ],
  "discount": {
    "active": true,
    "type": "percentage",
    "value": 50,
    "start_date": "2026-01-01T00:00:00Z",
    "end_date": "2026-12-31T23:59:59Z"
  },
  "active": true
}
```

### Required API Endpoints

- `GET /api/v1/public/{domain}/products` - List products
- `GET /api/v1/public/{domain}/products/{id}` - Get product
- `GET /api/v1/public/{domain}/promotions` - Get active promotions

## Design System

### Typography Scale
- H1: 2.5rem (40px)
- H2: 2rem (32px)
- H3: 1.5rem (24px)
- Body: 1rem (16px)
- Small: 0.875rem (14px)

### Spacing Scale
- XS: 0.5rem (8px)
- SM: 1rem (16px)
- MD: 1.5rem (24px)
- LG: 2rem (32px)
- XL: 3rem (48px)

### Border Radius
- Default: 10px
- Buttons: 4px
- Badges: 4px
- Input fields: 10px

### Shadows
- Cards: `0 2px 4px rgba(0,0,0,0.08)`
- Hover: `0 8px 16px rgba(0,0,0,0.12)`

## Browser Support

- Chrome (last 2 versions)
- Firefox (last 2 versions)
- Safari (last 2 versions)
- Edge (last 2 versions)
- Mobile browsers (iOS Safari, Chrome Mobile)

## Accessibility

- Semantic HTML5
- ARIA labels where needed
- Keyboard navigation support
- Color contrast compliance (WCAG AA)
- Focus indicators
- Alt text for images

## Version History

### 1.0.0 (2026-01-02)
- Initial release
- Product variants support
- Discount system
- Sale banners
- Dynamic branding
- Shopping cart
- Authentication integration
- Responsive design

## License

Proprietary - OilYourHair Platform

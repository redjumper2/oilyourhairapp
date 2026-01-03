# Frontend Template System

A multi-tenant e-commerce platform with template-based site generation.

## Directory Structure

```
frontend/
├── templates/              # Reusable site templates
│   └── store-template1/    # Allbirds-inspired e-commerce template
│       ├── *.html          # Template HTML files
│       ├── *.js            # Template JavaScript
│       ├── *.css           # Template styles
│       └── template.json   # Template metadata
│
├── sites/                  # Individual site configurations
│   └── oilyourhair.com/
│       ├── config.json     # Site configuration
│       ├── branding.json   # Brand colors, fonts, logo
│       ├── assets/         # Site-specific assets
│       └── public/         # Generated site (auto-generated)
│
└── build-site.sh           # Build script
```

## Quick Start

### 1. Create a New Site

Create a new directory under `sites/`:

```bash
mkdir -p sites/mynewsite.com/assets
```

### 2. Create Configuration Files

**sites/mynewsite.com/config.json:**
```json
{
  "domain": "mynewsite.com",
  "template": "store-template1",
  "template_version": "1.0.0",
  "enabled": true,
  "site_name": "My Store",
  "site_description": "Amazing products",
  "api": {
    "products_endpoint": "http://localhost:9091",
    "auth_endpoint": "http://localhost:9090"
  }
}
```

**sites/mynewsite.com/branding.json:**
```json
{
  "domain": "mynewsite.com",
  "brand_name": "My Store",
  "tagline": "Amazing Products Here",
  "colors": {
    "primary": "#1976d2",
    "primary_dark": "#115293",
    "primary_light": "#4791db",
    "secondary": "#64b5f6",
    "accent": "#42a5f5",
    "background": "#FFFFFF",
    "text": "#1a1a1a",
    "text_light": "#666",
    "border": "#e0e0e0",
    "success": "#4CAF50",
    "error": "#F44336",
    "warning": "#FF9800"
  },
  "typography": {
    "font_family": "-apple-system, sans-serif",
    "heading_font": "inherit",
    "h1_size": "2.5rem",
    "h2_size": "2rem",
    "h3_size": "1.5rem",
    "body_size": "1rem",
    "small_size": "0.875rem"
  },
  "layout": {
    "max_width": "1400px",
    "border_radius": "10px",
    "card_shadow": "0 4px 6px rgba(0,0,0,0.1)",
    "spacing": {
      "xs": "0.5rem",
      "sm": "1rem",
      "md": "1.5rem",
      "lg": "2rem",
      "xl": "3rem"
    }
  },
  "logo": {
    "text": "My Store",
    "url": "/logo.png",
    "height": "40px"
  }
}
```

### 3. Build the Site

```bash
./build-site.sh mynewsite.com
```

Or build all sites:

```bash
./build-site.sh
```

### 4. Add to Docker Compose

Add your site to `docker-compose.yml`:

```yaml
mynewsite-frontend:
  build:
    context: ./nginx
    dockerfile: Dockerfile
  container_name: mynewsite-frontend
  restart: unless-stopped
  ports:
    - "8081:80"
  volumes:
    - ./frontend/sites/mynewsite.com/public:/var/www/mynewsite.com/html:ro
    - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
    - ./nginx/sites-available:/etc/nginx/conf.d:ro
  networks:
    - auth-network
```

### 5. Start the Container

```bash
docker compose up -d mynewsite-frontend
```

## Templates

### store-template1

Modern e-commerce template inspired by Allbirds.

**Features:**
- Product variants (sizes, colors, attributes)
- Discount system with sale badges
- Dynamic sale banners
- Shopping cart
- Authentication integration
- Responsive design
- White-label branding

**Files:**
- `index.html` - Homepage
- `shop.html` - Product listing page
- `about.html` - About page
- `contact.html` - Contact page
- `reviews.html` - Reviews page
- `branding.js` - Branding system
- `banner.js` - Sale banner component
- `variant-selector.js` - Product variant selector
- `auth.js` - Authentication integration
- `styles.css` - Global styles

## Customization

### Colors

Edit `branding.json` to change the color scheme. All colors use CSS custom properties for easy theming.

### Logo

Place your logo in `sites/yoursite.com/assets/logo.png` and reference it in `branding.json`:

```json
{
  "logo": {
    "url": "/logo.png",
    "height": "40px"
  }
}
```

### Typography

Customize fonts and sizes in `branding.json`:

```json
{
  "typography": {
    "font_family": "'Inter', sans-serif",
    "h1_size": "3rem"
  }
}
```

## Development Workflow

1. **Modify Template:** Edit files in `templates/store-template1/`
2. **Rebuild Sites:** Run `./build-site.sh` to regenerate all sites
3. **Test Changes:** Changes are immediately reflected (mounted as read-only volumes)

## Creating New Templates

1. Create a new directory under `templates/`:
   ```bash
   mkdir templates/store-template2
   ```

2. Add your template files (HTML, CSS, JS)

3. Create `template.json` with metadata

4. Sites can now use your new template by setting `"template": "store-template2"` in their config.json

## API Integration

Each site connects to the Products API and Auth API via endpoints configured in `config.json`.

**Products API:** `/api/v1/public/{domain}/products`
**Promotions API:** `/api/v1/public/{domain}/promotions`

## Troubleshooting

### Site not loading

1. Check if the public directory exists:
   ```bash
   ls sites/yoursite.com/public/
   ```

2. Rebuild the site:
   ```bash
   ./build-site.sh yoursite.com
   ```

3. Restart the container:
   ```bash
   docker compose restart yoursite-frontend
   ```

### Files not updating

The public directory is mounted as read-only. After making changes to templates or branding:

1. Rebuild the site
2. Container will automatically pick up changes (no restart needed)

## License

Proprietary - OilYourHair Platform

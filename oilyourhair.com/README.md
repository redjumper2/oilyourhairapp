# Oil Your Hair - Product Catalog Website

Demo website showcasing integration with the products module API.

## Features

### Public Pages
- **Shop** (`/public/shop.html`) - Browse all products with filtering and search
- **Product Details** (`/public/product.html`) - View individual product details with variants

### Admin Pages
- **Product Management** (`/public/admin/products.html`) - Create, edit, and delete products (requires API key)

## Setup

### 1. Start Products Module

```bash
cd ../products_module
./products-module serve --config=config.dev.yaml
```

### 2. Create API Key (for admin access)

```bash
cd ../auth_module
./auth-module apikey create \
  --config=config.dev.yaml \
  --domain=oilyourhair.com \
  --service=products \
  --permissions=products.read,products.write
```

Save the API key - you'll need it for the admin page.

### 3. Serve the Website

You can use any static file server. For example:

```bash
# Using Python
cd public
python3 -m http.server 3000

# Using Node.js (npx http-server)
cd public
npx http-server -p 3000

# Using PHP
cd public
php -S localhost:3000
```

### 4. Access the Website

- **Shop**: http://localhost:3000/shop.html
- **Admin**: http://localhost:3000/admin/products.html

## API Configuration

The API is configured in the JavaScript files to point to:
- **Products API**: http://localhost:9091 (products_module)
- **Domain**: oilyourhair.com

To change these, edit:
- `js/api/products-api.js` - Update constructor defaults
- Or pass different values when creating the API instance

## Components

### JavaScript Components

**ProductsAPI** (`js/api/products-api.js`)
- API client for products module
- Handles both public and admin endpoints
- Supports authentication with API keys

**ProductGrid** (`js/components/ProductGrid.js`)
- Displays products in a responsive grid
- Supports filtering and pagination
- Shows stock status and variants

### Pages

**Shop Page** (`public/shop.html`)
- Product grid with search
- Filter by category, type, organic
- Responsive design

**Product Details** (`public/product.html`)
- Full product information
- Variant selector
- Image gallery
- Add to cart (placeholder)

**Admin Dashboard** (`public/admin/products.html`)
- Product CRUD operations
- Requires API key authentication
- Edit product details
- Manage inventory

## Styling

All styles in `css/products.css`:
- Responsive grid layout
- Mobile-friendly
- Brand colors from oilyourhair.com
- CSS variables for easy customization

## Usage Examples

### Public API (No Authentication)

```javascript
const api = new ProductsAPI('http://localhost:9091', 'oilyourhair.com');

// Get all products
const data = await api.getProducts();
console.log(data.products);

// Filter products
const oils = await api.getProducts({ category: 'oils' });

// Search products
const results = await api.searchProducts('coconut');

// Get single product
const product = await api.getProduct('product-id');
```

### Admin API (Requires API Key)

```javascript
const api = new ProductsAPI(
    'http://localhost:9091',
    'oilyourhair.com',
    'your-api-key-here'
);

// Create product
await api.createProduct({
    name: "Coconut Oil",
    description: "Pure organic coconut oil",
    base_price: 29.99,
    attributes: { category: "oils", type: "coconut" },
    variants: [
        { attributes: { size: "500ml" }, price: 29.99, stock: 50 }
    ]
});

// Update product
await api.updateProduct('product-id', {
    base_price: 27.99
});

// Delete product
await api.deleteProduct('product-id'); // Soft delete
await api.deleteProduct('product-id', true); // Hard delete
```

## Next Steps

- [ ] Add shopping cart functionality
- [ ] Integrate with payment gateway
- [ ] Add product reviews
- [ ] Implement wishlist
- [ ] Add order management
- [ ] Integrate with auth_module for user login

## Troubleshooting

**Products not loading?**
- Check that products_module is running on port 9091
- Verify domain is registered in auth_module
- Check browser console for errors

**Admin page not working?**
- Verify API key is correct
- Check API key has required permissions (products.read, products.write)
- Ensure API key is not expired

**CORS errors?**
- Products module should allow requests from your domain
- Check products_module CORS configuration in config.dev.yaml

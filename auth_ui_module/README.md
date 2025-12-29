# Auth UI Module

Customer-facing authentication UI for the multi-tenant auth system. Built with SvelteKit and Tailwind CSS.

## Features

- **Invitation Acceptance** - Branded invitation flow with QR code support
- **Magic Link Authentication** - Passwordless email authentication
- **Google OAuth** - Social login integration
- **Dynamic Branding** - Shows domain-specific logos, colors, company name
- **Minimal & Fast** - Small bundle size, quick load times
- **Static Deployment** - Builds to static files, deploy anywhere

## Project Structure

```
auth_ui_module/
├── src/
│   ├── routes/
│   │   ├── invite/          # Invitation acceptance page
│   │   ├── login/           # Magic link + Google OAuth
│   │   ├── verify/          # Magic link verification
│   │   └── oauth/google/    # Google OAuth redirect
│   ├── lib/
│   │   ├── api.js           # API client
│   │   ├── stores/
│   │   │   └── branding.js  # Branding state
│   │   └── components/      # Reusable components
│   └── app.css              # Tailwind styles
├── static/                  # Static assets
├── package.json
└── svelte.config.js
```

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Running auth_module backend API

### Installation

```bash
# Install dependencies
npm install

# Copy environment variables
cp .env.example .env

# Edit .env to point to your API
nano .env
```

### Development

```bash
# Start dev server
npm run dev

# Runs on http://localhost:5173
```

Visit:
- http://localhost:5173 - Home page
- http://localhost:5173/login?domain=testdomain.com - Login page
- http://localhost:5173/invite?token=... - Invitation page

### Build for Production

```bash
# Build static files
npm run build

# Preview production build
npm run preview
```

Output: `build/` directory contains static HTML/JS/CSS

## Environment Variables

```bash
# .env
VITE_API_URL=http://localhost:8080/api/v1  # Development
# VITE_API_URL=https://api.yoursaas.com/api/v1  # Production
```

## Usage Flows

### 1. Invitation Flow

```
Admin creates invite via CLI:
  └─> Returns URL: https://auth.yoursaas.com/invite?token=abc123&redirect=https://oilyourhair.com

User clicks link:
  1. UI verifies invitation token
  2. Fetches domain branding
  3. Shows invitation details (role, promo, etc.)
  4. User clicks "Accept"
  5. Creates user account
  6. Redirects to: https://oilyourhair.com#token=JWT_HERE
```

### 2. Magic Link Flow

```
User visits: https://auth.yoursaas.com/login?domain=oilyourhair.com&redirect=https://oilyourhair.com

Steps:
  1. UI shows branded login page
  2. User enters email
  3. Receives magic link email
  4. Clicks link → /verify?token=xxx&redirect=...
  5. UI verifies token
  6. Redirects to: https://oilyourhair.com#token=JWT_HERE
```

### 3. Google OAuth Flow

```
User clicks "Sign in with Google":
  1. Redirects to /oauth/google?domain=...&redirect=...
  2. UI redirects to backend API OAuth endpoint
  3. Standard Google OAuth dance
  4. Backend redirects to: https://oilyourhair.com#token=JWT_HERE
```

## Integration with Domain Apps

Domains need this minimal JavaScript to handle auth:

```html
<script>
  // Extract token from URL hash
  const hash = window.location.hash.substring(1);
  const params = new URLSearchParams(hash);
  const token = params.get('token');

  if (token) {
    // Store token
    localStorage.setItem('auth_token', token);

    // Clean URL
    window.history.replaceState({}, '', window.location.pathname);

    // Redirect to dashboard
    window.location = '/dashboard.html';
  }

  // Use token in API calls
  function apiCall(url) {
    const token = localStorage.getItem('auth_token');
    return fetch(url, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Host': 'oilyourhair.com'
      }
    });
  }
</script>
```

## Deployment

### Option 1: Vercel

```bash
npm install -g vercel
vercel --prod
```

### Option 2: Netlify

```bash
npm run build
# Drag build/ folder to Netlify
```

### Option 3: Cloudflare Pages

```bash
# Connect GitHub repo
# Build command: npm run build
# Output directory: build
```

### Option 4: Static Hosting (nginx, S3, etc.)

```bash
npm run build
# Upload build/ directory contents
```

## Customization

### Branding

Branding is fetched from the API based on domain. CSS variables are set dynamically:

```css
:root {
  --brand-primary: #000000;  /* From API */
  --brand-secondary: #666666; /* From API */
}
```

### Styling

Edit `src/app.css` or component classes. Uses Tailwind CSS.

## Development Notes

- **Hash Fragment**: JWT is passed in URL hash (#token=...) not query (?token=...) to avoid server logs
- **Branding Store**: Svelte store manages branding state across pages
- **API Client**: Centralized in `src/lib/api.js`
- **No Backend**: UI is pure frontend, calls auth_module API

## Troubleshooting

### API Connection Failed

- Check `VITE_API_URL` in `.env`
- Ensure auth_module backend is running
- Check CORS settings in backend

### Branding Not Loading

- Verify domain parameter in URL
- Check API branding endpoint
- Check browser console for errors

### Redirect Loop

- Verify redirect URL format
- Check domain's token extraction code
- Ensure localStorage is working

## Next Steps

- Add loading states/animations
- Add error boundaries
- Add analytics tracking
- Customize branding per domain
- Add remember me functionality

## License

Proprietary

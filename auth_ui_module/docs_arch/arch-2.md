Yes, that's a really solid pattern! You're basically building a centralized auth gateway that multiple domains use. This is how Auth0, Supabase, Clerk, etc. work.

  The Flow

  1. Admin (CLI): Creates invite for oilyourhair.com
     └─> Returns: https://auth.yoursaas.com/invite?token=abc123&redirect=https://oilyourhair.com

  2. Customer: Clicks invitation link
     └─> Lands on auth portal (your Svelte UI)

  3. Auth Portal:
     - Verifies invitation token
     - Fetches domain branding (Oil Your Hair logo, colors)
     - Shows: "You're invited to join Oil Your Hair as Editor"
     - Shows: Promo code, discount, etc.
     - Button: "Accept Invitation"

  4. User accepts:
     - Creates user record in MongoDB
     - Issues JWT token
     - Redirects to: https://oilyourhair.com/auth/callback?token=JWT_HERE

  5. oilyourhair.com:
     - Extracts JWT from URL
     - Stores in localStorage/cookie
     - Validates token (optional check with your API)
     - Redirects to /dashboard
     - User is now logged in

  Key Design Decisions

  1. Where is the auth portal hosted?

  Option A: Single shared domain (simpler)
  - auth.yoursaas.com for all tenants
  - Uses invitation token to determine which branding to show
  - Cleaner, one deployment

  Option B: Subdomain per tenant (more white-label)
  - auth.oilyourhair.com → CNAME to your auth portal
  - auth.fitnessclub.com → same portal, different branding
  - Feels more integrated but requires DNS setup per tenant

  I'd go with Option A - simpler to manage, still shows correct branding.

  2. Redirect URL strategy

  Option A: Store in domain settings
  {
    "domain": "oilyourhair.com",
    "redirect_url": "https://oilyourhair.com/auth/callback"
  }

  Option B: Pass in invitation metadata
  # CLI command
  ./auth-module invite create \
    --email=user@example.com \
    --redirect-url=https://oilyourhair.com/dashboard

  Option C: Embed in invite URL
  https://auth.yoursaas.com/invite?token=abc&redirect_to=https://oilyourhair.com

  I'd combine A + C: Default redirect in domain settings, but allow override in URL for flexibility.

  3. Token passing

  Most common pattern:
  - Redirect to: https://oilyourhair.com/auth/callback?token=JWT_HERE
  - Domain's frontend JavaScript:
  const params = new URLSearchParams(window.location.search);
  const token = params.get('token');
  localStorage.setItem('auth_token', token);
  window.location.href = '/dashboard'; // clean URL

  Alternative (more secure):
  - POST the token to domain's backend endpoint
  - Backend sets httpOnly cookie
  - Redirects to dashboard
  - Prevents token exposure in browser history

  4. Magic link flow (not invitation)

  Similar flow but initiated by user:
  1. User on oilyourhair.com clicks "Login"
     └─> Redirects to: https://auth.yoursaas.com/login?domain=oilyourhair.com

  2. Auth portal shows Oil Your Hair branding
     User enters email → receives magic link

  3. User clicks magic link
     └─> https://auth.yoursaas.com/verify?token=xyz&redirect=https://oilyourhair.com

  4. Auth portal verifies → redirects with JWT

  What the Domain UI Needs

  Each domain (oilyourhair.com, fitnessclub.com) needs minimal auth code:

  JavaScript:
  // Check if logged in
  const token = localStorage.getItem('auth_token');

  // Include in API requests
  fetch('https://api.yoursaas.com/products', {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Host': 'oilyourhair.com'  // Important!
    }
  });

  // Login button
  <button onclick="window.location='https://auth.yoursaas.com/login?domain=oilyourhair.com&redirect=https://oilyourhair.com'">
    Login
  </button>

  That's it! The domain UI just:
  1. Redirects to auth portal for login
  2. Receives JWT token back
  3. Stores it
  4. Uses it in API calls

  Benefits of This Approach

  ✅ Single auth UI - Maintain one Svelte app for all domains
  ✅ Consistent UX - All tenants get same smooth auth experience
  ✅ Branding still works - Shows correct logo/colors during auth
  ✅ Secure - Auth logic is centralized, harder to mess up
  ✅ Domain UI stays simple - They just handle token storage
  ✅ Works with any tech stack - Domain UI can be React, Vue, vanilla JS, whatever

  Potential Issues

  ⚠️ CORS - Domain UI needs to call your API from different origin
  - Solution: Proper CORS headers on your API

  ⚠️ Redirect validation - Need to whitelist redirect URLs to prevent open redirect attacks
  - Solution: Validate redirect URL matches domain's registered redirect_url

  ⚠️ Token in URL - JWT briefly visible in browser history/logs
  - Solution: Use short-lived exchange token, or POST to domain backend

  ⚠️ Multiple domains per customer - What if they want staging + production?
  - Solution: Allow multiple redirect URLs per domain

  CLI Extension Needed

  Add redirect URL management:
  # Set default redirect for domain
  ./auth-module domain settings update \
    --domain=oilyourhair.com \
    --redirect-url=https://oilyourhair.com/auth/callback

  # Create invite with custom redirect
  ./auth-module invite create \
    --email=user@example.com \
    --redirect-url=https://oilyourhair.com/dashboard

  My Take

  This is a really clean architecture. You're building a proper SaaS auth service. The auth portal becomes a lightweight, focused UI that does one thing well: authenticate users and hand them back to their apps.

  The Svelte UI scope stays tiny (just those 3-4 pages), and each domain gets a seamless auth experience without building it themselves.

  One question: Do you envision the domain UIs being separate applications your customers build, or are you building those too as part of the platform?

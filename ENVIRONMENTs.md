Great question! The environment is automatically detected based on the hostname you're accessing the website from. Here's how it works:

  Environment Detection

  The auth.js file determines the environment like this:

  const isProduction = window.location.hostname === 'oilyourhair.com' ||
                       window.location.hostname === 'www.oilyourhair.com';

  Development Environment

  You're in development when accessing via:
  - http://localhost:8080 (main website)
  - http://localhost:5173 (auth UI)
  - http://localhost:9090 (auth API)

  Behavior:
  - Auth portal: http://localhost:5173
  - API endpoint: http://localhost:9090/api/v1

  Production Environment

  You're in production when accessing via:
  - https://oilyourhair.com (main website)
  - https://www.oilyourhair.com (main website)
  - https://auth.oilyourhair.com (auth UI)
  - https://api.oilyourhair.com (auth API)

  Behavior:
  - Auth portal: https://auth.oilyourhair.com
  - API endpoint: https://api.oilyourhair.com/api/v1

  How to Tell Which Environment You're In

  1. Look at the URL bar in your browser:
    - localhost = Development
    - oilyourhair.com = Production
  2. Open browser console (F12) and run:
  window.location.hostname
    - Returns "localhost" = Development
    - Returns "oilyourhair.com" = Production
  3. When creating invitations with the Makefile, you get both URLs:
  Invitation URL (Development):
    http://localhost:5173/invite?token=invite-123&redirect=http://localhost:8080

  Invitation URL (Production):
    https://auth.oilyourhair.com/invite?token=invite-123&redirect=https://oilyourhair.com
    - Use the first URL when testing locally
    - Use the second URL for the live public website

  Current Setup

  Right now you have both environments running:
  - Development: Docker containers accessible via localhost
  - Production: Same Docker containers accessible via Cloudflare Tunnel at oilyourhair.com

  They're using the same backend services, just accessed through different URLs. The frontend JavaScript automatically adapts based on which URL you're using!

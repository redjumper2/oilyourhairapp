The Google OAuth provider isn't configured. You have two options:

  Option 1: Set Up Google OAuth (Recommended)

  If you want Google login to work, you need to configure it:

  1. Get Google OAuth Credentials
    - Go to https://console.cloud.google.com/
    - Create a new project or select existing one
    - Enable "Google+ API"
    - Go to "Credentials" → "Create Credentials" → "OAuth 2.0 Client ID"
    - Application type: "Web application"
    - Authorized redirect URIs: http://localhost:9090/auth/google/callback
    - Copy the Client ID and Client Secret
  2. Add to .env file
  GOOGLE_CLIENT_ID=your-client-id-here.apps.googleusercontent.com
  GOOGLE_CLIENT_SECRET=your-client-secret-here
  GOOGLE_CALLBACK_URL=http://localhost:9090/auth/google/callback
  3. Restart auth-api
  docker compose restart auth-api

  Option 2: Disable Google Login (Quick Fix)

  If you don't need Google login, we can hide the button:

  Would you like me to:
  1. Help you set up Google OAuth (I'll guide you through it)
  2. Hide the Google login button (quick fix)
  3. Both (set it up properly but hide it for now)

  Which would you prefer?

───────────────────────────────────────────────
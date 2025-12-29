# Google OAuth Setup Guide

This guide shows how to configure Google OAuth for the auth module.

## Step 1: Create Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Note your project name

## Step 2: Enable Google+ API

1. In the Cloud Console, go to **APIs & Services** > **Library**
2. Search for "Google+ API"
3. Click **Enable**

## Step 3: Configure OAuth Consent Screen

1. Go to **APIs & Services** > **OAuth consent screen**
2. Select **External** (unless you have Google Workspace)
3. Click **Create**

**App information:**
- App name: `Your App Name` (e.g., "Oil Your Hair")
- User support email: Your email
- Developer contact email: Your email

**Scopes:**
- Add: `../auth/userinfo.email`
- Add: `../auth/userinfo.profile`

**Test users** (for development):
- Add your Gmail address for testing

Click **Save and Continue**

## Step 4: Create OAuth Credentials

1. Go to **APIs & Services** > **Credentials**
2. Click **Create Credentials** > **OAuth client ID**
3. Application type: **Web application**
4. Name: `Auth Module` (or any name)

**Authorized JavaScript origins:**
```
http://localhost:8080
http://localhost:3000
https://oilyourhair.com
https://www.oilyourhair.com
```

**Authorized redirect URIs:**
```
http://localhost:8080/api/v1/auth/google/callback
https://oilyourhair.com/api/v1/auth/google/callback
https://www.oilyourhair.com/api/v1/auth/google/callback
```

Click **Create**

## Step 5: Copy Credentials

You'll see a dialog with:
- **Client ID**: `1234567890-abcdefghijklmnop.apps.googleusercontent.com`
- **Client secret**: `GOCSPX-abc123xyz...`

**Copy both values** - you'll need them next.

## Step 6: Configure Auth Module

**With Docker (recommended):**

Edit `.env`:
```bash
GOOGLE_CLIENT_ID=1234567890-abcdefghijklmnop.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-abc123xyz...
GOOGLE_CALLBACK_URL=http://localhost:8080/api/v1/auth/google/callback
```

Restart services:
```bash
make docker-restart
```

**Without Docker:**

Edit `config.yaml`:
```yaml
google:
  client_id: 1234567890-abcdefghijklmnop.apps.googleusercontent.com
  client_secret: GOCSPX-abc123xyz...
  callback_url: http://localhost:8080/api/v1/auth/google/callback
```

Restart server:
```bash
make dev
```

## Step 7: Test OAuth

### Browser Test

1. Create a domain:
```bash
make domain-create DOMAIN=localhost NAME="Local Dev" EMAIL=admin@localhost
```

2. Open browser to:
```
http://localhost:8080/api/v1/auth/google
```

3. You should be redirected to Google sign-in
4. After signing in, you'll be redirected to frontend with JWT token

### API Test

```bash
# This will redirect you to Google
curl -L http://localhost:8080/api/v1/auth/google \
  -H "Host: localhost"
```

## Production Setup

For production domains:

1. Update **Authorized JavaScript origins** in Google Console:
```
https://oilyourhair.com
https://www.oilyourhair.com
```

2. Update **Authorized redirect URIs**:
```
https://oilyourhair.com/api/v1/auth/google/callback
https://www.oilyourhair.com/api/v1/auth/google/callback
```

3. Update `.env` or environment variables:
```bash
GOOGLE_CALLBACK_URL=https://oilyourhair.com/api/v1/auth/google/callback
FRONTEND_URL=https://oilyourhair.com
```

4. In Google Console, move app from **Testing** to **Production**:
   - Go to **OAuth consent screen**
   - Click **Publish App**
   - Submit for verification if needed

## Multiple Domains

If you have multiple domains (e.g., oilyourhair.com, fitnessclub.com), add all redirect URIs:

```
https://oilyourhair.com/api/v1/auth/google/callback
https://fitnessclub.com/api/v1/auth/google/callback
```

The auth module will automatically scope users to the correct domain based on the `Host` header.

## Troubleshooting

### Error: "redirect_uri_mismatch"

**Cause:** The redirect URI doesn't match what's configured in Google Console.

**Fix:**
1. Check your callback URL in `.env` or `config.yaml`
2. Verify it **exactly** matches one in Google Console (including protocol, domain, path)
3. Common mistake: `http://` vs `https://`

### Error: "access_denied"

**Cause:** User denied permissions or app not verified.

**Fix:**
1. Make sure you added test users in Google Console
2. Check OAuth consent screen is configured
3. Try with a different Google account

### OAuth works but user not created

**Check logs:**
```bash
make docker-logs-api
```

**Common issues:**
- Domain not found (create domain first)
- MongoDB connection issue
- Check `users` collection in MongoDB

### OAuth redirects to wrong URL

**Check environment variables:**
```bash
# Docker
docker exec auth-module-api env | grep FRONTEND_URL

# Local
echo $AUTH_APP_FRONTEND_URL
```

Should match your frontend URL.

## Security Notes

1. **Keep client secret secure** - Never commit to git
2. **Use HTTPS in production** - HTTP OAuth is insecure
3. **Verify domains** - Only add domains you control
4. **Rotate secrets regularly** - Create new credentials periodically
5. **Monitor usage** - Check Google Cloud Console for suspicious activity

## Testing with Multiple Domains

```bash
# Test with oilyourhair.com
curl -H "Host: oilyourhair.com" http://localhost:8080/api/v1/auth/google

# Test with fitnessclub.com
curl -H "Host: fitnessclub.com" http://localhost:8080/api/v1/auth/google
```

Each domain gets isolated users - same Google account can sign in to both domains separately.

## Next Steps

- Configure email settings for magic links
- Set up invitation system
- Build frontend integration
- See `API.md` for complete API documentation

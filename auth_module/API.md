## Auth Module API Documentation

Base URL: `http://localhost:8080/api/v1`

All API requests to protected endpoints must include domain context via `Host` header and JWT via `Authorization: Bearer <token>`.

---

## Authentication Endpoints

### Request Magic Link

Send a magic link email to user.

**Endpoint:** `POST /auth/magic-link/request`

**Headers:**
- `Host: oilyourhair.com` (domain context)

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Response:** `200 OK`
```json
{
  "message": "Magic link sent to your email"
}
```

---

### Verify Magic Link

Verify magic link token and get JWT.

**Endpoint:** `GET /auth/magic-link/verify?token=<token>`

**Response:** `200 OK`
```json
{
  "token": "eyJhbGc...",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@example.com",
    "domain": "oilyourhair.com",
    "role": "customer",
    "permissions": ["products.read", "cart.read", "cart.write"]
  }
}
```

---

### Verify Invitation

Check invitation validity and get details.

**Endpoint:** `GET /auth/invitation/verify?token=<token>`

**Response:** `200 OK`
```json
{
  "invitation_id": "507f1f77bcf86cd799439011",
  "role": "editor",
  "domain": "oilyourhair.com",
  "expires_at": "2025-12-30T12:00:00Z",
  "time_remaining": "23h45m",
  "email": "editor@example.com",
  "promo_code": "WELCOME20",
  "discount_percent": 20,
  "branding": {
    "company_name": "Oil Your Hair",
    "primary_color": "#2E7D32",
    "logo_url": "https://..."
  }
}
```

---

### Accept Invitation

Accept invitation and create user account.

**Endpoint:** `POST /auth/invitation/accept`

**Request Body:**
```json
{
  "token": "invitation-token-here",
  "email": "user@example.com",
  "auth_provider": "magic_link",
  "provider_id": ""
}
```

**Response:** `200 OK`
```json
{
  "token": "eyJhbGc...",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@example.com",
    "domain": "oilyourhair.com",
    "role": "editor",
    "permissions": ["products.read", "products.write", ...]
  }
}
```

---

### Google OAuth Login

Initiate Google OAuth flow.

**Endpoint:** `GET /auth/google`

**Headers:**
- `Host: oilyourhair.com` (domain context)

**Flow:**
1. User clicks "Sign in with Google"
2. Frontend redirects to: `GET /api/v1/auth/google`
3. User is redirected to Google for authentication
4. After approval, Google redirects back to callback URL

---

### Google OAuth Callback

Handle Google OAuth callback (browser redirect).

**Endpoint:** `GET /auth/google/callback`

**Note:** This is called by Google, not your frontend directly.

**Response:** Redirects to frontend with JWT
```
HTTP 302 Redirect
Location: http://localhost:3000/auth/callback?token=eyJhbGc...
```

Frontend should extract token from URL and store it.

---

### Google OAuth Callback (JSON)

Handle Google OAuth callback and return JSON (for API clients).

**Endpoint:** `GET /auth/google/callback/json?domain=<domain>`

**Response:** `200 OK`
```json
{
  "token": "eyJhbGc...",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "email": "user@gmail.com",
    "domain": "oilyourhair.com",
    "role": "customer",
    "permissions": ["products.read", "cart.read", "cart.write"]
  }
}
```

---

### Get Current User

Get authenticated user information.

**Endpoint:** `GET /auth/me`

**Headers:**
- `Authorization: Bearer <jwt-token>`
- `Host: oilyourhair.com`

**Response:** `200 OK`
```json
{
  "id": "507f1f77bcf86cd799439011",
  "email": "user@example.com",
  "domain": "oilyourhair.com",
  "role": "admin",
  "permissions": ["domain.settings.read", "domain.settings.write", ...]
}
```

---

## Admin Endpoints

All admin endpoints require:
- `Authorization: Bearer <jwt-token>` header
- `Host: <domain>` header
- User role: `admin`

### Get Domain Settings

**Endpoint:** `GET /admin/domain/settings`

**Response:** `200 OK`
```json
{
  "domain": "oilyourhair.com",
  "name": "Oil Your Hair",
  "status": "active",
  "settings": {
    "allowed_auth_providers": ["google", "magic_link"],
    "default_role": "customer",
    "require_email_verification": true
  },
  "branding": {
    "company_name": "Oil Your Hair",
    "primary_color": "#2E7D32",
    "logo_url": "https://...",
    "support_email": "support@oilyourhair.com"
  }
}
```

---

### Update Domain Settings

**Endpoint:** `PUT /admin/domain/settings`

**Request Body:**
```json
{
  "settings": {
    "allowed_auth_providers": ["google", "magic_link"],
    "default_role": "customer"
  },
  "branding": {
    "company_name": "Oil Your Hair",
    "primary_color": "#2E7D32",
    "logo_url": "https://cdn.example.com/logo.png",
    "support_email": "support@oilyourhair.com"
  }
}
```

**Response:** `200 OK`
```json
{
  "message": "Domain settings updated successfully"
}
```

---

### List Users

**Endpoint:** `GET /admin/users?role=<role>`

Query parameters:
- `role` (optional): Filter by role (`admin`, `editor`, `viewer`, `customer`)

**Response:** `200 OK`
```json
{
  "users": [
    {
      "id": "507f1f77bcf86cd799439011",
      "email": "admin@example.com",
      "role": "admin",
      "auth_provider": "google",
      "created_at": "2025-12-28T10:00:00Z",
      "last_login": "2025-12-28T14:30:00Z"
    }
  ],
  "count": 1
}
```

---

### Invite User

Create invitation and optionally send email.

**Endpoint:** `POST /admin/users/invite`

**Request Body:**
```json
{
  "email": "newuser@example.com",
  "role": "editor",
  "type": "email_with_qr",
  "single_use": true,
  "expires_in_hours": 72,
  "promo_code": "WELCOME20",
  "source": "admin_panel",
  "discount_percent": 20
}
```

**Fields:**
- `email` (optional): User email (required for email invitations)
- `role` (required): `admin`, `editor`, `viewer`, or `customer`
- `type` (required): `email`, `qr_code`, or `email_with_qr`
- `single_use` (optional): `true` for user-specific, `false` for promotional
- `max_uses` (optional): Limit for multi-use invitations
- `expires_in_hours` (optional): Custom expiry (defaults based on type)
- `promo_code`, `source`, `ref`, `discount_percent` (optional): Tracking metadata

**Response:** `200 OK`
```json
{
  "invitation_id": "507f1f77bcf86cd799439011",
  "url": "http://localhost:3000/invite?token=abc123...",
  "token": "abc123...",
  "expires_at": "2025-12-31T10:00:00Z",
  "qr_code": "data:image/png;base64,iVBORw0KG..."
}
```

---

### Update User

Update user role and permissions.

**Endpoint:** `PUT /admin/users/:id`

**Request Body:**
```json
{
  "role": "admin",
  "permissions": ["domain.settings.read", "domain.settings.write", ...]
}
```

**Response:** `200 OK`
```json
{
  "message": "User updated successfully"
}
```

---

### Delete User

Soft delete a user (cannot delete yourself or last admin).

**Endpoint:** `DELETE /admin/users/:id`

**Response:** `200 OK`
```json
{
  "message": "User deleted successfully"
}
```

---

## Error Responses

All endpoints may return error responses:

**400 Bad Request**
```json
{
  "error": "Invalid request body"
}
```

**401 Unauthorized**
```json
{
  "error": "Invalid or expired token"
}
```

**403 Forbidden**
```json
{
  "error": "Insufficient permissions"
}
```

**404 Not Found**
```json
{
  "error": "Resource not found"
}
```

**500 Internal Server Error**
```json
{
  "error": "Internal server error"
}
```

---

## Authentication Flow Examples

### Magic Link Flow

1. User enters email on `oilyourhair.com`
2. Frontend: `POST /auth/magic-link/request` with `Host: oilyourhair.com`
3. User receives email with magic link
4. User clicks link → Frontend: `GET /auth/magic-link/verify?token=xxx`
5. Frontend receives JWT, stores it, redirects to dashboard
6. All subsequent requests include `Authorization: Bearer <jwt>`

### Invitation Flow

1. Admin: `POST /admin/users/invite` to create invitation
2. Invitee receives email with link + QR code
3. Invitee clicks link → Frontend: `GET /auth/invitation/verify?token=xxx`
4. Frontend shows invitation details (role, company, promo)
5. User accepts → Frontend: `POST /auth/invitation/accept`
6. Frontend receives JWT, user is logged in

### Google OAuth Flow

1. User clicks "Sign in with Google" on `oilyourhair.com`
2. Frontend: `window.location.href = "http://localhost:8080/api/v1/auth/google"` with `Host: oilyourhair.com` header
3. Backend stores domain in session
4. User redirected to Google for authentication
5. User approves permissions
6. Google redirects to: `GET /auth/google/callback`
7. Backend verifies OAuth, creates/logs in user
8. Backend redirects to: `http://localhost:3000/auth/callback?token=<jwt>`
9. Frontend extracts token from URL, stores it
10. All subsequent requests include `Authorization: Bearer <jwt>`

**Alternative (API clients):**
Use `/auth/google/callback/json` to get JSON response instead of redirect

### Domain Isolation

JWT tokens are domain-scoped. If user has JWT for `oilyourhair.com`, they cannot use it on `fitnessclub.com`.

**Example JWT payload:**
```json
{
  "user_id": "507f1f77bcf86cd799439011",
  "email": "user@example.com",
  "domain": "oilyourhair.com",
  "role": "admin",
  "permissions": ["domain.settings.write", ...],
  "exp": 1735574400,
  "iat": 1735488000
}
```

Middleware validates: `JWT.domain === Host header`

## Auth Module API Documentation

Base URL: `http://localhost:8080/api/v1`

All API requests to protected endpoints must include domain context via `Host` header and JWT via `Authorization: Bearer <token>`.

---

## API Endpoints Summary

### Authentication Endpoints (Public)

- `POST /auth/magic-link/request` - Request a magic link email for passwordless login
- `GET /auth/magic-link/verify` - Verify magic link token and receive JWT
- `GET /auth/invitation/verify` - Check invitation validity and get details
- `POST /auth/invitation/accept` - Accept invitation and create user account
- `GET /auth/google` - Initiate Google OAuth login flow
- `GET /auth/google/callback` - Handle Google OAuth callback (browser redirect)
- `GET /auth/google/callback/json` - Handle Google OAuth callback (JSON response)
- `GET /auth/me` - Get current authenticated user information (requires JWT)

### Admin Endpoints (Protected - Admin Only)

**Domain Settings**
- `GET /admin/domain/settings` - Get domain settings and branding
- `PUT /admin/domain/settings` - Update domain settings and branding

**User Management**
- `GET /admin/users` - List all users in domain (optional role filter)
- `POST /admin/users/invite` - Create user invitation with QR code
- `PUT /admin/users/:id` - Update user role and permissions
- `DELETE /admin/users/:id` - Soft delete a user

**Permissions**
- `GET /admin/permissions` - Get all available permissions grouped by category
- `GET /admin/permissions/roles` - Get permissions for each role

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

### Get All Permissions

Get all available permissions grouped by category.

**Endpoint:** `GET /admin/permissions`

**Headers:**
- `Authorization: Bearer <jwt-token>`
- `Host: <domain>`

**Response:** `200 OK`
```json
{
  "groups": [
    {
      "name": "Domain Management",
      "description": "Manage domain settings and branding",
      "permissions": [
        "domain.settings.read",
        "domain.settings.write"
      ]
    },
    {
      "name": "User Management",
      "description": "Manage users, roles, and invitations",
      "permissions": [
        "users.read",
        "users.write",
        "users.delete",
        "users.invite"
      ]
    },
    {
      "name": "Product Management",
      "description": "Manage product catalog",
      "permissions": [
        "products.read",
        "products.write"
      ]
    },
    {
      "name": "Order Management",
      "description": "View and manage orders",
      "permissions": [
        "orders.read",
        "orders.write"
      ]
    },
    {
      "name": "Inventory Management",
      "description": "Manage stock and inventory",
      "permissions": [
        "inventory.read",
        "inventory.write"
      ]
    },
    {
      "name": "Shopping Cart",
      "description": "Customer shopping cart operations",
      "permissions": [
        "cart.read",
        "cart.write"
      ]
    }
  ],
  "total": 16
}
```

**Use case:** Build permission selector UI in admin panel.

---

### Get Permissions by Role

Get permissions for each role.

**Endpoint:** `GET /admin/permissions/roles`

**Headers:**
- `Authorization: Bearer <jwt-token>`
- `Host: <domain>`

**Response:** `200 OK`
```json
{
  "admin": {
    "count": 12,
    "permissions": [
      "domain.settings.read",
      "domain.settings.write",
      "users.read",
      "users.write",
      "users.delete",
      "users.invite",
      "products.read",
      "products.write",
      "orders.read",
      "orders.write",
      "inventory.read",
      "inventory.write"
    ]
  },
  "editor": {
    "count": 5,
    "permissions": [
      "products.read",
      "products.write",
      "orders.read",
      "inventory.read",
      "inventory.write"
    ]
  },
  "viewer": {
    "count": 3,
    "permissions": [
      "products.read",
      "orders.read",
      "inventory.read"
    ]
  },
  "customer": {
    "count": 4,
    "permissions": [
      "products.read",
      "cart.read",
      "cart.write",
      "orders.read"
    ]
  }
}
```

**Use case:** Display role comparison table in admin UI.

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

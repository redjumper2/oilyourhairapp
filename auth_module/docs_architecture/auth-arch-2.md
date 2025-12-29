=======================================

1. separate auth module usable by any API.  2. auth module needs to work for an org or lets call it domain (like oilyourhair.com), but same auth 
module instance needs to be able to work for other domains.  3. but when users login, we know the domain they are coming in from so will be assigned to
 that domain. 4. Since they might be part of a shop 1 domain or a fitness club domain, they are allowed to be part of more than 1 domain, but their 
auth setup remains independent per domain. So their API permission remains specific to the domain they are logging into.  Does this sound ok with your 
above assumptions for auth? 

● Ah, this is actually cleaner than what I was describing. Let me make sure I understand correctly:

  Domain-Scoped Auth

  You're building a white-label auth service where:

  - Domain = Tenant: oilyourhair.com is one domain/tenant, fitnessclub.com is another
  - Same auth backend serves all domains
  - Login context is determined by domain: When user hits https://oilyourhair.com/login, the auth module knows they're logging into the oilyourhair.com domain
  - No org picker needed: User's session is automatically scoped to the domain they logged in from
  - Cross-domain isolation: Same email (john@gmail.com) can exist in multiple domains with completely independent roles/permissions

  So if John logs into oilyourhair.com with Google → he gets a JWT with:
  {
    user_id: "abc123",
    domain: "oilyourhair.com",
    role: "admin",
    permissions: ["products.write", "orders.read"]
  }

  If John later logs into fitnessclub.com with Google → separate session, different JWT:
  {
    user_id: "def456",  // could be different user_id, or same user_id
    domain: "fitnessclub.com", 
    role: "viewer",
    permissions: ["classes.read"]
  }

  Key Questions This Clarifies

  1. How does the auth module know which domain?
  - HTTP request header (Host header: oilyourhair.com)
  - Or subdomain: oilyourhair.yourplatform.com
  - Or path-based: yourplatform.com/oilyourhair/login
  - Or domain passed as parameter in auth request

  2. User records - shared or separate?

  Option A: Separate user records per domain
  - users collection: {user_id, email, domain, ...}
  - john@gmail.com on oilyourhair.com = different document than john@gmail.com on fitnessclub.com
  - Simpler isolation, clear boundaries

  Option B: Shared user record, domain-scoped roles
  - users collection: {user_id, email, ...} (no domain)
  - user_domains collection: {user_id, domain, role, permissions, ...}
  - Same person, multiple domain memberships tracked
  - Allows potential future feature: "view all my domains"

  3. Auth flow per domain:

  User goes to oilyourhair.com/login:
  1. Click "Sign in with Google"
  2. OAuth callback includes domain context
  3. Auth module checks: does this Google email exist for domain "oilyourhair.com"?
    - Yes → log them in, issue JWT with domain
    - No → create new user record for this domain, assign default role (customer?), issue JWT

  Magic link for oilyourhair.com:
  1. User enters email on oilyourhair.com
  2. Auth module generates token with domain embedded: {token, email, domain: "oilyourhair.com", expires}
  3. User clicks link, auth module verifies token + domain match
  4. Issue JWT scoped to oilyourhair.com

  This Simplifies Permission Checks

  In your Echo middleware:
  - Extract JWT
  - Get domain from JWT (or verify against request Host header)
  - Check permissions scoped to that domain
  - All database queries automatically filtered by domain

  No org picker, no "switch organization" UI complexity. User's browser session is tied to the domain they're on.

  Does This Mean...?

  - Each domain gets its own product catalog? Yes - oilyourhair.com products are separate from fitnessclub.com products
  - Orders are domain-scoped? Yes - orders in oilyourhair.com domain only
  - Users can't "share" carts across domains? Correct - cart is per domain
  - The auth module is domain-agnostic? Yes - same code, just switches context based on incoming domain

  Is this the model? And are you doing custom domains (oilyourhair.com, fitnessclub.com) or subdomains (oilyourhair.yourplatform.com)?

===============================================================================================
● Perfect, total domain isolation. This is actually a really clean architecture - share nothing makes everything simpler.

  The Auth Module API Surface

  Since it's domain-agnostic and reusable, I'm thinking it exposes endpoints like:

  Auth endpoints (consumed by frontend):
  - POST /auth/magic-link/request - send magic link email
  - GET /auth/magic-link/verify?token=xyz - verify and create session
  - GET /auth/google - initiate Google OAuth
  - GET /auth/google/callback - handle OAuth callback
  - POST /auth/logout - invalidate session
  - GET /auth/me - get current user info

  Each request includes Host header, auth module extracts domain, scopes everything to that domain.

  MongoDB Collections (Option A - Separate Records)

  Core collections:
  users: {
    _id,
    email,
    domain,  // "oilyourhair.com"
    auth_provider,  // "google" | "magic_link"
    provider_id,  // Google user ID if applicable
    role,  // "admin" | "editor" | "viewer" | "customer"
    permissions: [],
    created_at,
    last_login
  }

  Index on: {email, domain} - unique compound index

  magic_link_tokens: {
    _id,
    email,
    domain,
    token,  // cryptographically random
    expires_at
  }

  Index on: {token} and TTL index on expires_at

  domains: {
    _id,
    domain,  // "oilyourhair.com"
    name,  // "Oil Your Hair"
    status,  // "active" | "suspended"
    settings: {
      allowed_auth_providers: ["google", "magic_link"],
      default_role: "customer",
      ...
    }
  }

  Unique index on: domain

  Domain Whitelisting & Security

  Critical security consideration: You need a domains collection to whitelist valid domains. Otherwise someone could spoof the Host header and create users for arbitrary domains.

  Flow:
  1. Request comes in with Host: oilyourhair.com
  2. Auth module queries domains collection: does this domain exist and is it active?
  3. If not found → reject request (403 or redirect to error page)
  4. If found → proceed with auth flow scoped to that domain

  Domain registration: You'll need an admin API (or manual process) to register new domains:
  - Add domain to domains collection
  - Configure DNS (point oilyourhair.com to your server)
  - Potentially verify domain ownership (like Let's Encrypt does)

  How Other APIs Consume This Auth Module

  Option 1: Shared authentication middleware

  Product API, Order API, etc. all use the same auth middleware:
  1. Extract JWT from request
  2. Verify JWT signature
  3. Extract domain from JWT claims
  4. Verify domain matches Host header (prevent JWT reuse across domains)
  5. Attach user info to request context
  6. Downstream handlers filter all queries by domain

  Option 2: Auth module as a service

  Auth is a separate microservice, other APIs call it to validate tokens. More complex but more decoupled.

  I'd lean toward Option 1 - shared middleware package that all your APIs import. Auth module issues JWTs, other APIs just validate them.

  First User Problem

  When a new domain is registered (oilyourhair.com), how does the first user become admin?

  Approach 1: Domain registration creates an initial admin user
  - When you register "oilyourhair.com" in domains collection, also create a user with role="admin"
  - Send them magic link to set up their account

  Approach 2: First user to sign up becomes admin
  - Check if any users exist for this domain
  - If count == 0, assign role="admin" to this user
  - All subsequent users get role="customer"

  Approach 3: Invite-only
  - Domain exists but no users can self-register
  - Admin invites users via email with pre-assigned roles
  - Requires you to manually create first admin for each domain

  Which feels right for your use case?

  Cross-Domain Attack Prevention

  With Host header-based routing, you need to prevent:

  JWT reuse across domains:
  - User gets JWT for oilyourhair.com
  - Tries to use same JWT on fitnessclub.com
  - Your middleware should check: JWT.domain == Host header, reject if mismatch

  Host header spoofing:
  - Not a real issue if domain whitelist is enforced
  - Even if someone sends fake Host header, they can't create users for unregistered domains

  Magic link token theft:
  - Tokens should be single-use (delete after verification)
  - Short expiry (5-15 minutes)
  - Domain must match when verifying

  Sound good? What's your thinking on the "first user" problem for new domains?all

  =====================================================================================================================

  ● Approach 1 is definitely the cleanest. Here's why I like it:

  Why Approach 1 Works Well

  Control & Security:
  - You explicitly control who gets admin access - no race condition where random person becomes admin
  - Domain registration is an intentional, tracked event
  - Clear audit trail: domain created → admin user created → admin invited

  Onboarding Flow:
  1. You (or domain owner via some registration portal) register "oilyourhair.com"
  2. Provide admin email: "admin@oilyourhair.com"
  3. System creates:
    - Domain record in domains collection
    - Admin user record in users collection with role="admin"
    - Magic link token
  4. Send email: "Your Oil Your Hair admin account is ready, click here to set up"
  5. Admin clicks, verifies magic link, gets logged in
  6. Admin can now invite other users or configure auth settings

  Prevents Squatting:
  - Without this, if you use Approach 2 (first signup = admin), someone could race to sign up before the real domain owner
  - With Approach 1, domain + admin email are registered together atomically

  Domain Registration Workflow

  You'll need some way to register domains. Two options:

  Manual/CLI approach (simple, start here):
  - You have a script or Makefile target: make register-domain DOMAIN=oilyourhair.com EMAIL=admin@oilyourhair.com
  - Inserts into MongoDB, sends magic link
  - Good for early stages when you're onboarding domains yourself

  Self-service portal (later):
  - Public site: yourplatform.com/register
  - Domain owner enters their domain, admin email
  - You verify domain ownership (DNS TXT record or email verification)
  - Auto-provisions domain + admin user
  - More scalable but needs verification to prevent abuse

  I'd start with the manual approach - gets you moving fast. You can build the self-service portal later once the auth module is solid.

  Admin's First Steps

  Once admin logs in for the first time, they probably want to:
  1. Customize domain settings - which auth providers to enable (Google, magic link, maybe add Facebook later)
  2. Invite team members - send magic links to editors/viewers
  3. Set up products - start building their catalog

  This means you'll eventually need some admin UI endpoints:
  - GET /admin/domain/settings - view domain config
  - PUT /admin/domain/settings - update auth providers, default role, etc.
  - POST /admin/users/invite - send invitation to email with pre-assigned role
  - GET /admin/users - list all users in this domain
  - DELETE /admin/users/:id - remove user

  But that's down the road. For now, just getting domain registration + admin creation working is the foundation.

  One More Thought: Domain Verification

  Since you're doing custom domains (not subdomains), you'll eventually need DNS verification to prevent someone registering a domain they don't own.

  Simple approach: Email verification
  - User registers "oilyourhair.com" with "admin@oilyourhair.com"
  - Send verification email to that address
  - They must click link to activate domain
  - Assumes they control the email on that domain

LATER - IMPORTANT maybe:
  Stronger approach: DNS TXT record
  - User adds TXT record: _yourplatform.oilyourhair.com = "verification-code-xyz"
  - Your system checks DNS before activating domain
  - Proves they control DNS for that domain

  Email verification is probably fine to start. You can add DNS verification later if needed.

  Does this match what you were envisioning for domain onboarding?

  ========================================================================================
MY RESPONSE:
> lets keep domain verification for later. For now, lets have us create the domain with an email address who will be made the admin for that domain 
receiving a magic link. We are the default administrators for the system.  So who we send the magic link to is the admin for that domain.  Also, lets 
go ahead and add the admin APIs you have listed above as that will be needed to manage the domain settings and users and invite users. 

  ========================================================================================
● Perfect, that makes the initial rollout simple - you're the gatekeeper for domain creation. No verification needed since you control the process.

  Admin API Architecture

  These admin endpoints need to be domain-scoped and role-protected. Here's how I'd think about them:

  Permission Model for Admin APIs

  Who can access admin endpoints?
  - Only users with role: "admin" for that domain
  - Middleware checks: user.role === "admin" && user.domain === Host header
  - Editors/viewers/customers get 403

  Granular permissions (optional, but useful):
  You could have admin sub-permissions:
  - domain.settings.read / domain.settings.write
  - users.invite / users.read / users.delete

  This lets you have "super admin" vs "limited admin" roles later if needed.

  Admin API Endpoints Breakdown

  Domain Settings Management:
  GET /admin/domain/settings
  - Returns current domain configuration
  - Shows: allowed auth providers, default role, domain status, etc.
  - Permission: admin role

  PUT /admin/domain/settings
  - Update domain configuration
  - Body: { allowed_auth_providers: ["google", "magic_link"], default_role: "customer" }
  - Updates `domains` collection for current domain
  - Permission: admin role

  User Management:
  GET /admin/users
  - List all users in current domain
  - Query params: ?role=editor, ?page=1, ?limit=50
  - Returns: users with email, role, last_login, created_at
  - Permission: admin role

  POST /admin/users/invite
  - Invite new user to domain with pre-assigned role
  - Body: { email: "editor@example.com", role: "editor", permissions: [...] }
  - Creates user record + magic link token
  - Sends invitation email
  - Permission: admin role (with users.invite permission)

  PUT /admin/users/:id
  - Update existing user's role/permissions
  - Body: { role: "viewer", permissions: [...] }
  - Only updates users in same domain (security check)
  - Permission: admin role

  DELETE /admin/users/:id
  - Remove user from domain
  - Soft delete or hard delete?
  - Can't delete yourself (prevent lockout)
  - Can't delete if only admin (prevent losing admin access)
  - Permission: admin role

  Safety Checks for User Deletion

  Prevent admin lockout:
  1. Can't delete yourself while logged in
  2. Can't delete the last admin user in a domain
  3. Query: if (user.role === 'admin') { count = users.count({domain, role: 'admin'}); if (count <= 1) reject }

  Invitation Flow

  When admin invites someone:
  1. Admin calls POST /admin/users/invite with email + role
  2. System creates user record: {email, domain, role, status: "invited", created_at}
  3. Generate magic link token with expiry
  4. Send email: "You've been invited to Oil Your Hair as an Editor. Click here to access your account."
  5. Invitee clicks link, verifies token
  6. Update user status: "invited" → "active"
  7. User gets logged in, JWT issued

  Alternative: Instead of creating user immediately, create an invitations collection:
  invitations: {
    email,
    domain,
    role,
    invited_by,  // admin user_id
    token,
    expires_at,
    status: "pending" | "accepted"
  }

  When they accept, create the actual user record. This keeps users collection cleaner (only active users).

  Which approach feels better - create user immediately (status: "invited") or separate invitations collection?

  Domain Settings Structure

  What kinds of settings would you want admins to configure?

  Auth settings:
  - allowed_auth_providers: ["google", "magic_link"] - which login methods are enabled
  - default_role: "customer" - role for self-signup users (if you enable that later)
 LATER - require_email_domain: "oilyourhair.com" - optional: only allow emails from specific domain

 LATER- Feature flags (for later):
  - features: { products: true, orders: true, analytics: false }
  - Lets you enable/disable modules per domain

  Branding (maybe):
  - branding: { logo_url, primary_color, company_name }
  - Could be used by frontend to customize UI per domain

  Start simple (just auth settings), expand later.

  Bootstrap Process

  When you manually create a domain as system admin:

  Option A: CLI script/Makefile
  make create-domain DOMAIN=oilyourhair.com EMAIL=admin@oilyourhair.com NAME="Oil Your Hair"
  Script does:
  1. Insert into domains collection
  2. Insert admin user into users collection
  3. Generate magic link token
  4. Print magic link URL (or send email)

LATER:  Option B: System admin API
  You could have super-admin endpoints (separate from domain admin):
  POST /system/domains
  - Requires system admin JWT (separate from domain admin)
  - Creates domain + admin user
  - Returns magic link

  This is only accessible to you (the platform operators).

  I'd start with Option A (Makefile/script) since it's just you creating domains initially. Simpler than building a whole system admin layer.

  Make sense? And which invitation approach - immediate user creation vs separate invitations collection?

==================================================================================
> Option A sounds best to start with.  Soft delete sounds best for now.  The invitations flow that you mentioned above sounds perfect.  Create user 
only if invitation is accepted.  Need to expire magic link in reasonable amount of time as well.  Lets default the settings to having google and 
magic_link for now by default.  lets not enforce emails from specific domains for now. (lets do it later). Branding also sounds good for the invites 
etc during auth flows. 

● Perfect, this is coming together nicely. Let me think through a few more details:

  Magic Link Expiry

  Reasonable timeouts:
  - Short-lived (5-15 minutes): More secure, but annoying if user doesn't check email immediately
  - Medium (1 hour): Good balance - user can finish what they're doing and still click link
  - Long (24 hours): Convenient but security risk if email compromised

  I'd suggest 15 minutes for regular magic links, 24 hours for invitations.

  Reasoning:
  - Regular login magic link = user is actively trying to log in right now, they'll click fast
  - Invitation = admin sends invite, recipient might not see email for hours, needs longer window
  - Both are single-use tokens anyway (deleted after verification)

  Sound reasonable?

  Soft Delete Implementation

  For soft delete, add a field to users:
  users: {
    ...
    deleted_at: null | ISODate,
    deleted_by: null | user_id  // which admin deleted them
  }

  Queries need to filter out deleted users:
  - GET /admin/users → query: {domain, deleted_at: null}
  - Auth flows → query: {email, domain, deleted_at: null}
  - Can add ?include_deleted=true param for admins to view deleted users

  Restoration (bonus feature for later):
  - POST /admin/users/:id/restore → set deleted_at: null
  - Useful if admin accidentally deletes someone

  Invitations Collection Structure

  invitations: {
    _id,
    email,
    domain,
    role,
    permissions: [],  // optional, defaults based on role
    invited_by,  // admin user_id who sent invite
    token,  // cryptographically random, single-use
    created_at,
    expires_at,  // created_at + 24 hours
    status: "pending" | "accepted" | "expired",
    accepted_at: null | ISODate
  }

  Indexes:
  - {token} - unique, for lookup during verification
  - {expires_at} - TTL index to auto-delete expired invitations
  - {email, domain, status} - check if pending invitation exists

  Prevent duplicate invitations:
  - Before creating invitation, check if user already exists: {email, domain, deleted_at: null}
  - Also check if pending invitation exists: {email, domain, status: "pending"}
  - If either exists, reject or return existing invitation link

  Branding in Auth Flows

  When invitation email is sent, or magic link email, you can template it with domain branding:

  Fetch branding from domain settings:
  domains: {
    domain,
    name,  // "Oil Your Hair"
    branding: {
      logo_url: "https://cdn.example.com/oilyourhair-logo.png",
      primary_color: "#2E7D32",  // green for hair oil brand
      company_name: "Oil Your Hair",
      support_email: "support@oilyourhair.com"
    },
    settings: {
      allowed_auth_providers: ["google", "magic_link"],  // default
      default_role: "customer"
    }
  }

  Email template uses branding:
  Subject: You've been invited to {company_name}

  Hi,

  {admin_name} has invited you to join {company_name} as an {role}.

  [Click here to accept invitation]
  (Button styled with primary_color)

  Logo: {logo_url}

  Auth pages use branding:
  - Login page shows domain's logo
  - Styled with primary color
  - "Sign in to {company_name}"

  This makes it feel white-labeled per domain.

  Default Domain Settings on Creation

  When you run make create-domain, it should initialize:
  {
    domain: "oilyourhair.com",
    name: "Oil Your Hair",  // passed as param
    status: "active",
    settings: {
      allowed_auth_providers: ["google", "magic_link"],  // default
      default_role: "customer",
      require_email_verification: true  // for magic links
    },
    branding: {
      company_name: "Oil Your Hair",  // same as name, can customize later
      primary_color: "#000000",  // default black
      logo_url: null,  // admin uploads later
      support_email: null
    },
    created_at: ISODate(),
    created_by: "system"  // since you're creating it manually
  }

  Admin can update branding via PUT /admin/domain/settings later.

  Soft Delete for Invitations Too?

  Should expired/declined invitations be soft deleted or hard deleted?

  Option 1: Hard delete expired
  - TTL index auto-deletes after expiry
  - Keeps collection clean
  - No audit trail

  Option 2: Soft delete/mark expired
  - Update status to "expired" instead of deleting
  - Admin can see invitation history
  - Can resend invitation if needed

  I'd lean toward Option 1 (hard delete) for invitations since they're temporary by nature. But keep audit log if you want history:

  invitation_logs: {
    email,
    domain,
    role,
    invited_by,
    sent_at,
    accepted_at,
    status: "accepted" | "expired" | "cancelled"
  }

  Insert record when invitation is sent, update when accepted/expired. Cheap audit trail without cluttering main collection.

  Worth it or overkill for now?

=============================================================================
● Excellent! We've got a solid architecture now. Let me summarize what we've designed:

  Architecture Summary

  Auth Module - Domain-Scoped, Multi-Tenant

  Core Collections:
  - domains - whitelist of registered domains with settings/branding
  - users - separate user records per domain, soft delete with deleted_at
  - invitations - pending invites, hard deleted after expiry/acceptance
  - invitation_logs - audit trail for all invitations sent
  - magic_link_tokens - temporary tokens for auth (15 min login, 24 hr invites)

  Auth Methods:
  - Google OAuth (via Goth)
  - Magic links (email-based)
  - Domain determined by Host header
  - JWT tokens with domain scoping

  Admin APIs (role: admin only):
  - Domain settings management (GET/PUT /admin/domain/settings)
  - User management (GET/POST/PUT/DELETE /admin/users)
  - Invitation system (POST /admin/users/invite)

  Security:
  - Domain whitelist validation
  - JWT domain verification against Host header
  - Single-use magic link tokens with expiry
  - Soft delete users, prevent admin lockout
  - Complete domain isolation (share nothing)

  Bootstrap:
  - Makefile target to create domain + admin user
  - Admin gets 24-hour magic link invitation
  - Default settings: Google + magic link enabled

  Branding:
  - Per-domain logo, colors, company name
  - Used in emails and auth pages for white-labeling

  Ready to Build?

  We've talked through the architecture pretty thoroughly. Are you ready for me to start generating code, or is there anything else you want to discuss/adjust first?

  If we're good to go, I'm thinking we start with:
  1. MongoDB collections/models setup
  2. Domain creation Makefile target
  3. Basic auth endpoints (magic link flow)
  4. Google OAuth integration
  5. Admin APIs

  What do you think - ready to build, or more planning needed?

  ======

  > had 1 more question - how can I present users with a card that contains a QR Code that takes them to a link that I can use to invite them for the 
domain, along with a special info in the code that tells me something about the user. 
  ⎿  Interrupted · What should Claude do instead?

> card or email with QRCode 
  ⎿  Interrupted · What should Claude do instead?

> for example a promotion 

● Ah, this is a cool feature - QR code invitations with tracking metadata. Let me think through how this would work:

  QR Code Invitation Flow

  Use case example:
  - You're running a promotion for Oil Your Hair at a trade show
  - You print cards with QR codes
  - Each QR code contains: invitation link + promo code + source tracking
  - Customer scans → lands on signup page
  - System knows: they came from "TradeShow2025" promotion, gets 20% off first order

  How to Encode Metadata in Invitation Links

  Invitation link structure:
  https://oilyourhair.com/invite?token=abc123&promo=TRADESHOW20&source=booth&ref=sarah

  QR code contains this full URL with encoded parameters:
  - token - unique invitation token (like magic link token)
  - promo - promotion code (TRADESHOW20 = 20% off)
  - source - attribution (where did they scan: booth, email, instagram, etc.)
  - ref - who invited them (sales rep name, affiliate, etc.)

  Enhanced Invitations Collection

  Expand the invitations collection to support this:
  invitations: {
    _id,
    email: null,  // null for QR code invitations (anonymous until claimed)
    domain,
    role: "customer",  // default role for promo invitations
    token,
    
    // New fields for QR/promo invitations
    invitation_type: "email" | "qr_code" | "link",
    metadata: {
      promo_code: "TRADESHOW20",
      source: "booth",
      ref: "sarah",
      discount_percent: 20,
      custom_data: {...}  // any other tracking info
    },

    expires_at,
    created_at,
    claimed_at: null,  // when user accepts
    claimed_by: null,  // user_id who claimed it
    status: "pending" | "claimed" | "expired"
  }

  Two Types of Invitations

  Email invitation (what we discussed earlier):
  - Admin invites specific email address
  - Email required, role assigned
  - Magic link sent via email
  - Single use

  QR code invitation (new):
  - Admin generates anonymous invitation
  - No email required initially
  - Contains promo/tracking metadata
  - Can be used multiple times OR single-use (configurable)
  - User scans, signs up with their email, invitation gets "claimed"

  Admin API for QR Code Generation

  New admin endpoint:
  POST /admin/invitations/generate-qr
  Body: {
    promo_code: "TRADESHOW20",
    source: "booth",
    ref: "sarah",
    discount_percent: 20,
    expires_at: "2025-02-01",  // promotion end date
    max_uses: 100,  // optional: limit how many people can use this QR
    single_use: false  // true = one person only, false = reusable
  }

  Response: {
    invitation_id,
    url: "https://oilyourhair.com/invite?token=abc123&promo=TRADESHOW20&source=booth",
    qr_code_data_url: "data:image/png;base64,..."  // QR code image
  }

  Backend generates:
  1. Invitation record with metadata
  2. QR code image (using library like qrcode package)
  3. Returns both URL and QR code image

  QR Code Generation

  Libraries:
  - Go: github.com/skip2/go-qrcode - simple QR code generation
  - Returns PNG image as bytes or base64 data URL

  Usage:
  Generate QR code containing the invitation URL
  Return as base64 data URL so frontend can display/download
  Admin can:
    - Download PNG for printing on physical cards
    - Embed in email template
    - Display on screen for scanning

  Physical Card vs Email Use Cases

  Physical card (trade show, in-store):
  - Print QR code on business cards, flyers, posters
  - Customer scans → signup page with promo pre-filled
  - Track: "100 signups from trade show booth"

  Email with QR code:
  - Send promotional email blast
  - Email contains both: clickable link + QR code image
  - Desktop users click link, mobile users scan QR
  - Track: "50 signups from email campaign"

  Social media:
  - Post QR code image on Instagram
  - Followers scan → signup with discount
  - Track: "25 signups from Instagram post"

  Tracking & Analytics

  When user claims QR invitation:
  1. User scans QR → lands on /invite?token=abc123&promo=TRADESHOW20
  2. Frontend shows signup form with promo code pre-filled
  3. User enters email, signs up (Google or magic link)
  4. Backend creates user, marks invitation as claimed
  5. Links user to invitation metadata

  Analytics queries:
  - How many users came from "booth" source?
  - Which promo codes converted best?
  - Which sales rep (ref) brought most customers?

  Store in invitation_logs:
  invitation_logs: {
    invitation_id,
    domain,
    promo_code,
    source,
    ref,
    claimed_by: user_id,
    claimed_at,
    user_email
  }

  Multi-Use vs Single-Use

  Single-use (personal invitation):
  - Generate unique QR for one person
  - After claimed, mark as "used", can't be reused
  - Good for: VIP invites, one-on-one sales

  Multi-use (promotional):
  - Same QR code can be used by multiple people
  - Track total claims
  - Optional max_uses limit (e.g., "first 100 customers")
  - Good for: marketing campaigns, events

  Admin UI for This

  Admin would see:
  GET /admin/invitations
  - List all invitations (email + QR code types)
  - Filter by status, source, promo code
  - View claims count for multi-use QR codes

  GET /admin/invitations/:id
  - View invitation details
  - See who claimed it (for multi-use, list of all users)
  - Download QR code image again

  Does this match what you had in mind? And should QR invitations default to multi-use or single-use?

  ==========================================================================================

  ● Ah perfect, so you need both single-use (user-specific) AND multi-use (promotional) QR codes. This actually unifies nicely:

  Three Invitation Scenarios

  1. Email Invitation with QR Code (User-Specific, Single-Use)

  Scenario: User asks "send me an invite"

  Flow:
  - Admin calls POST /admin/users/invite with {email: "john@example.com", role: "editor"}
  - Backend creates invitation with invitation_type: "email", single_use: true
  - Generates magic link token
  - Sends email containing:
    - Clickable link: for desktop users
    - QR code image: same link, for mobile/convenience
  - User can either click link OR scan QR code (both work, same token)
  - Once claimed, token is invalidated

  Email template:
  Hi John,

  You've been invited to Oil Your Hair as an Editor.

  [Click here to accept] ← clickable button

  Or scan this QR code:
  [QR code image]

  Expires in 24 hours.

  2. Single-Use QR Code Only (User-Specific, No Email)

  Scenario: User says "I don't use email much, can you just give me a QR code?"

  Flow:
  - Admin calls POST /admin/invitations/generate-qr with {email: "john@example.com", role: "editor", single_use: true}
  - Backend creates invitation, generates QR code
  - Admin downloads QR code PNG or shows it on screen
  - User scans in person → lands on signup page
  - Email is pre-filled from invitation, user completes signup
  - Token is single-use, invalidated after claim

  Use case: In-person onboarding, trade show booth where you're talking to specific person

  3. Multi-Use QR Code (Promotional, Anonymous)

  Scenario: Marketing campaign, trade show, social media

  Flow:
  - Admin calls POST /admin/invitations/generate-qr with {promo_code: "SUMMER20", source: "instagram", multi_use: true, max_uses: 500}
  - Backend creates invitation with email: null (anonymous)
  - Admin downloads QR code, posts on Instagram
  - 500 different people scan → all land on signup page with promo pre-filled
  - Each person enters their own email, signs up
  - Invitation tracks all claims in audit log

  Use case: Marketing, promotions, bulk user acquisition

  Unified Invitation Structure

  All three scenarios use the same invitations collection:

  invitations: {
    _id,
    domain,
    token,  // unique token in all cases
    
    // User-specific or anonymous
    email: "john@example.com" | null,  // null = anonymous/multi-use
    role,
    permissions,
    
    // Type and usage
    invitation_type: "email" | "qr_code" | "email_with_qr",
    single_use: true | false,
    max_uses: null | 500,  // for multi-use, optional limit
    uses_count: 0,  // track how many times claimed

    // Promotional metadata (optional)
    metadata: {
      promo_code: null | "SUMMER20",
      source: null | "instagram",
      ref: null | "sarah",
      discount_percent: null | 20
    },

    // Tracking
    created_by: admin_user_id,
    created_at,
    expires_at,
    status: "pending" | "claimed" | "expired" | "exhausted"  // exhausted = max_uses reached
  }

  Admin APIs - Unified

  Single endpoint for all invitation types:

  POST /admin/invitations/create
  Body: {
    // Required
    type: "email" | "qr_code" | "email_with_qr",
    role: "editor",
    
    // Optional - for user-specific invitations
    email: "john@example.com" | null,
    
    // Optional - for multi-use
    single_use: true | false,  // default: true if email provided, false if null
    max_uses: null | 500,
    
    // Optional - promotional metadata
    promo_code: null | "SUMMER20",
    source: null | "instagram",
    ref: null | "sarah",
    discount_percent: null | 20,
    
    // Optional - custom expiry
    expires_in_hours: 24 | 168  // default 24 for user-specific, longer for promos
  }

  Response: {
    invitation_id,
    url: "https://oilyourhair.com/invite?token=abc123",
    qr_code_data_url: "data:image/png;base64,..."  // always include QR code
  }

  Backend logic:
  - If type: "email" or type: "email_with_qr" → send email
  - If type: "qr_code" or type: "email_with_qr" → generate and return QR code
  - If email provided → defaults to single_use: true
  - If email: null → defaults to single_use: false (promotional)

  Claim Flow (User Accepts Invitation)

  User scans QR or clicks link → lands on:
  /invite?token=abc123

  Frontend:
  1. Calls GET /auth/invitation/verify?token=abc123
  2. Backend validates token, returns invitation details
  3. If email pre-assigned → show: "You're invited as john@example.com, choose login method"
  4. If anonymous → show: "Enter your email to accept this invitation"
  5. User completes signup (Google or magic link)

  Backend:
  1. Verify token exists and not expired
  2. If single_use: true:
    - Check status === "pending"
    - If already claimed → reject
  3. If single_use: false:
    - Check uses_count < max_uses (if max_uses set)
    - Increment uses_count
  4. Create user record
  5. Mark invitation claimed
  6. Log to invitation_logs

  Preventing Abuse on Multi-Use QR Codes

  Rate limiting:
  - Same email can't claim same multi-use invitation twice
  - Check: has this email already used this promo code?
  - Store in invitation_logs with unique index on {invitation_id, user_email}

  Max uses:
  - Admin sets max_uses: 500
  - After 500 claims, status → "exhausted", no more signups

  Expiry:
  - All invitations have expires_at
  - Background job or TTL index cleans up

  Does this cover both your use cases - specific user invitations AND promotional QR codes?

  =============================================================================================

  ● Perfect! Configurable expiry makes total sense - different campaigns need different timeframes.

  Configurable Expiry Options

  When creating invitation, admin can specify:

  Option 1: Relative expiry (hours/days from creation):
  POST /admin/invitations/create
  Body: {
    ...
    expires_in_hours: 168,  // 7 days
  }
  Backend calculates: expires_at = now() + 168 hours

  Option 2: Absolute expiry (specific date/time):
  POST /admin/invitations/create
  Body: {
    ...
    expires_at: "2025-02-15T23:59:59Z",  // campaign ends Feb 15
  }
  Backend uses this exact timestamp

  Option 3: Both supported (use whichever is provided):
  - If expires_in_hours provided → calculate from now
  - If expires_at provided → use absolute timestamp
  - If neither → use sensible defaults

  Sensible Defaults

  If admin doesn't specify expiry, default based on invitation type:

  Defaults:
  - User-specific email invitation: 24 hours
  - User-specific QR code: 72 hours (3 days, more flexible)
  - Multi-use promotional QR: 30 days (or no expiry?)
  - Trade show/event QR: admin must specify (tied to event date)

  Or keep it simple: default 24 hours for everything, admin customizes as needed.

  Enhanced Create Invitation API

  POST /admin/invitations/create
  Body: {
    type: "email" | "qr_code" | "email_with_qr",
    role: "editor",
    email: "john@example.com" | null,
    single_use: true | false,
    max_uses: null | 500,
    
    // Expiry - provide one or neither (uses default)
    expires_in_hours: 168,  // OR
    expires_at: "2025-02-15T23:59:59Z",

    // Promotional metadata
    promo_code: null | "SUMMER20",
    source: null | "instagram",
    ref: null | "sarah",
    discount_percent: null | 20
  }

  Backend validation:
  - Can't provide both expires_in_hours AND expires_at (pick one)
  - expires_at must be in the future
  - expires_in_hours must be > 0
  - If neither provided → default to 24 hours

  Admin UI Benefits

  This gives admins flexibility:

  VIP customer invitation:
  {
    email: "vip@example.com",
    role: "customer",
    expires_in_hours: 168,  // 7 days, no rush
    promo_code: "VIP50",
    discount_percent: 50
  }

  Flash sale QR code:
  {
    email: null,
    multi_use: true,
    max_uses: 1000,
    expires_at: "2025-01-31T23:59:59Z",  // sale ends Jan 31
    promo_code: "FLASH24",
    discount_percent: 24
  }

  Trade show booth:
  {
    email: null,
    multi_use: true,
    expires_at: "2025-03-15T18:00:00Z",  // event ends 6pm on March 15
    source: "tradeshow-booth",
    promo_code: "SHOW2025"
  }

  Display Expiry to Users

  When user lands on invitation page:
  /invite?token=abc123

  Frontend calls GET /auth/invitation/verify?token=abc123

  Response includes:
  {
    invitation_id,
    email: "john@example.com" | null,
    role: "editor",
    promo_code: "VIP50",
    expires_at: "2025-02-15T23:59:59Z",
    time_remaining: "6 days 14 hours"  // calculated
  }

  Frontend shows: "This invitation expires in 6 days. Accept now to get 50% off!"

  Creates urgency, improves conversion.

  Sound good? Want defaults to be 24 hours across the board, or different defaults per invitation type?

  =================================================================================================

  
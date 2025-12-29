
  The Permission Flow

  1. Domain Creation (What You Ran)

  ./auth-module domain create --domain=testdomain.com --admin-email=admin@testdomain.com

  When you created the domain, the code created an admin invitation with role: "admin". Let me show you the code:

  From cmd/domain.go:56-62:
  invitation, qrCodeURL, err := invitationService.CreateInvitation(ctx, &services.CreateInvitationRequest{
      Domain:          domainName,
      Email:           adminEmail,
      Role:            "admin",  // ← Hard-coded to admin role
      Type:            "email_with_qr",
      SingleUse:       true,
      ExpiresInHours:  cfg.Invitation.Defaults.EmailExpiryHours,
  })

  2. Invitation Acceptance (What You Just Did)

  When you accepted the invitation via the API, it read the role from the invitation and assigned permissions based on that role.

  From internal/services/auth.go:147-154:
  // Update user role/permissions from invitation
  user.Role = invitation.Role  // ← Gets "admin" from invitation
  user.Permissions = invitation.Permissions  // ← Gets permissions for admin role
  _, err = s.db.Users.UpdateOne(ctx,
      bson.M{"_id": user.ID},
      bson.M{"$set": bson.M{
          "role":        user.Role,
          "permissions": user.Permissions,
      }},
  )

  3. Permission Mapping (The Key Part!)

  The permissions for each role are defined in internal/models/user.go:28-59:

  func GetPermissionsForRole(role string) []string {
      switch role {
      case "admin":
          return []string{
              "domain.settings.read", "domain.settings.write",
              "users.read", "users.write", "users.delete", "users.invite",
              "products.read", "products.write",
              "orders.read", "orders.write",
              "inventory.read", "inventory.write",
          }
      case "editor":
          return []string{
              "products.read", "products.write",
              "orders.read",
              "inventory.read", "inventory.write",
          }
      case "viewer":
          return []string{
              "products.read",
              "orders.read",
              "inventory.read",
          }
      case "customer":
          return []string{
              "products.read",
              "cart.read", "cart.write",
              "orders.read",
          }
      default:
          return []string{}
      }
  }

  Summary: The Permission Chain

  Domain Creation
      ↓
  Creates Admin Invitation (role: "admin")
      ↓
  You Accept Invitation
      ↓
  GetPermissionsForRole("admin") is called
      ↓
  Returns admin permission array
      ↓
  Stored in database with user
      ↓
  JWT token includes these permissions

  Testing Different Roles

  Want to see how other roles work? Create invitations with different roles:

  Create an Editor Invitation

  curl -X POST http://localhost:8080/api/v1/admin/users/invite \
    -H "Authorization: Bearer $ADMIN_JWT" \
    -H "Host: testdomain.com" \
    -H "Content-Type: application/json" \
    -d '{
      "email": "editor@testdomain.com",
      "role": "editor",
      "type": "qr_code",
      "single_use": true
    }' | jq

  Editor will get only:
  [
    "products.read", "products.write",
    "orders.read",
    "inventory.read", "inventory.write"
  ]

  Create a Customer Invitation

  curl -X POST http://localhost:8080/api/v1/admin/users/invite \
    -H "Authorization: Bearer $ADMIN_JWT" \
    -H "Host: testdomain.com" \
    -H "Content-Type: application/json" \
    -d '{
      "role": "customer",
      "type": "qr_code",
      "single_use": false,
      "max_uses": 100
    }' | jq

  Customer will get only:
  [
    "products.read",
    "cart.read", "cart.write",
    "orders.read"
  ]

  Custom Permissions

  You can also override the default permissions when creating an invitation:

  curl -X POST http://localhost:8080/api/v1/admin/users/invite \
    -H "Authorization: Bearer $ADMIN_JWT" \
    -H "Host: testdomain.com" \
    -H "Content-Type: application/json" \
    -d '{
      "email": "custom@testdomain.com",
      "role": "editor",
      "permissions": ["products.read", "products.write", "analytics.read"],
      "type": "email_with_qr"
    }' | jq

  This will give them custom permissions instead of the default editor permissions.

  ---
  So to answer your question: The domain didn't "know" - the role determined the permissions via the GetPermissionsForRole() function, which is a hard-coded mapping in the code. When the first admin invitation was created during domain setup, it was assigned the "admin" role, which maps to the full set of admin permissions.

  Does that make sense? Want to test creating users with different roles?
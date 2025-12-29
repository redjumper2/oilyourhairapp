
  3. API Endpoints for Permissions

  Test with your admin JWT:

  # Get all permissions
  curl http://localhost:8080/api/v1/admin/permissions \
    -H "Authorization: Bearer $ADMIN_JWT" \
    -H "Host: testdomain.com" | jq

  Response:
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
      ...
    ],
    "total": 16
  }

  Get permissions by role:

  curl http://localhost:8080/api/v1/admin/permissions/roles \
    -H "Authorization: Bearer $ADMIN_JWT" \
    -H "Host: testdomain.com" | jq

  Response:
  {
    "admin": {
      "permissions": ["domain.settings.read", "domain.settings.write", ...],
      "count": 12
    },
    "editor": {
      "permissions": ["products.read", "products.write", ...],
      "count": 5
    },
    ...
  }

  How to Add New Permissions

  Now when you want to add a new permission:

  1. Add to internal/models/permissions.go:
  const (
      // ... existing permissions ...

      // New permission
      PermReportsRead Permission = "reports.read"
  )

  2. Add to a permission group:
  {
      Name:        "Reports",
      Description: "View reports and analytics",
      Permissions: []Permission{
          PermReportsRead,
      },
  },

  3. Add to roles that need it:
  "admin": {
      PermDomainSettingsRead,
      // ... existing ...
      PermReportsRead,  // ← Add here
  },

  Done! The permission is now:
  - ✅ Validated
  - ✅ Shown in CLI
  - ✅ Shown in API
  - ✅ Available to assign to users

  Summary
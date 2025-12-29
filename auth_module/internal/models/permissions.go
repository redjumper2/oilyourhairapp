package models

// Permission represents a single permission
type Permission string

// All available permissions in the system
// This is the SINGLE SOURCE OF TRUTH for permissions
const (
	// Domain permissions
	PermDomainSettingsRead  Permission = "domain.settings.read"
	PermDomainSettingsWrite Permission = "domain.settings.write"

	// User permissions
	PermUsersRead   Permission = "users.read"
	PermUsersWrite  Permission = "users.write"
	PermUsersDelete Permission = "users.delete"
	PermUsersInvite Permission = "users.invite"

	// Product permissions
	PermProductsRead  Permission = "products.read"
	PermProductsWrite Permission = "products.write"

	// Order permissions
	PermOrdersRead  Permission = "orders.read"
	PermOrdersWrite Permission = "orders.write"

	// Inventory permissions
	PermInventoryRead  Permission = "inventory.read"
	PermInventoryWrite Permission = "inventory.write"

	// Cart permissions
	PermCartRead  Permission = "cart.read"
	PermCartWrite Permission = "cart.write"

	// Analytics permissions (example for future)
	PermAnalyticsRead Permission = "analytics.read"
)

// PermissionGroup represents a category of permissions
type PermissionGroup struct {
	Name        string
	Description string
	Permissions []Permission
}

// AllPermissionGroups returns all permission groups organized by category
func AllPermissionGroups() []PermissionGroup {
	return []PermissionGroup{
		{
			Name:        "Domain Management",
			Description: "Manage domain settings and branding",
			Permissions: []Permission{
				PermDomainSettingsRead,
				PermDomainSettingsWrite,
			},
		},
		{
			Name:        "User Management",
			Description: "Manage users, roles, and invitations",
			Permissions: []Permission{
				PermUsersRead,
				PermUsersWrite,
				PermUsersDelete,
				PermUsersInvite,
			},
		},
		{
			Name:        "Product Management",
			Description: "Manage product catalog",
			Permissions: []Permission{
				PermProductsRead,
				PermProductsWrite,
			},
		},
		{
			Name:        "Order Management",
			Description: "View and manage orders",
			Permissions: []Permission{
				PermOrdersRead,
				PermOrdersWrite,
			},
		},
		{
			Name:        "Inventory Management",
			Description: "Manage stock and inventory",
			Permissions: []Permission{
				PermInventoryRead,
				PermInventoryWrite,
			},
		},
		{
			Name:        "Shopping Cart",
			Description: "Customer shopping cart operations",
			Permissions: []Permission{
				PermCartRead,
				PermCartWrite,
			},
		},
	}
}

// AllPermissions returns a flat list of all available permissions
func AllPermissions() []Permission {
	var all []Permission
	for _, group := range AllPermissionGroups() {
		all = append(all, group.Permissions...)
	}
	return all
}

// AllPermissionsAsStrings returns all permissions as strings
func AllPermissionsAsStrings() []string {
	perms := AllPermissions()
	result := make([]string, len(perms))
	for i, p := range perms {
		result[i] = string(p)
	}
	return result
}

// IsValidPermission checks if a permission string is valid
func IsValidPermission(perm string) bool {
	for _, p := range AllPermissions() {
		if string(p) == perm {
			return true
		}
	}
	return false
}

// ValidatePermissions checks if all permissions in a list are valid
func ValidatePermissions(permissions []string) ([]string, []string) {
	var valid []string
	var invalid []string

	for _, perm := range permissions {
		if IsValidPermission(perm) {
			valid = append(valid, perm)
		} else {
			invalid = append(invalid, perm)
		}
	}

	return valid, invalid
}

// PermissionsByRole returns permissions for each role using constants
func PermissionsByRole() map[string][]Permission {
	return map[string][]Permission{
		"admin": {
			PermDomainSettingsRead,
			PermDomainSettingsWrite,
			PermUsersRead,
			PermUsersWrite,
			PermUsersDelete,
			PermUsersInvite,
			PermProductsRead,
			PermProductsWrite,
			PermOrdersRead,
			PermOrdersWrite,
			PermInventoryRead,
			PermInventoryWrite,
		},
		"editor": {
			PermProductsRead,
			PermProductsWrite,
			PermOrdersRead,
			PermInventoryRead,
			PermInventoryWrite,
		},
		"viewer": {
			PermProductsRead,
			PermOrdersRead,
			PermInventoryRead,
		},
		"customer": {
			PermProductsRead,
			PermCartRead,
			PermCartWrite,
			PermOrdersRead,
		},
	}
}

// GetPermissionsForRoleTyped returns typed permissions for a role
func GetPermissionsForRoleTyped(role string) []Permission {
	perms, exists := PermissionsByRole()[role]
	if !exists {
		return []Permission{}
	}
	return perms
}

// PermissionsToStrings converts Permission slice to string slice
func PermissionsToStrings(permissions []Permission) []string {
	result := make([]string, len(permissions))
	for i, p := range permissions {
		result[i] = string(p)
	}
	return result
}

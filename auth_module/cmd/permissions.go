package cmd

import (
	"fmt"

	"github.com/sparque/auth_module/internal/models"
	"github.com/spf13/cobra"
)

// permissionsCmd represents the permissions command
var permissionsCmd = &cobra.Command{
	Use:   "permissions",
	Short: "Manage and view permissions",
	Long:  `View available permissions and role mappings`,
}

// permissionsListCmd lists all available permissions
var permissionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available permissions",
	Long:  `Display all available permissions grouped by category`,
	Run: func(cmd *cobra.Command, args []string) {
		listPermissions()
	},
}

// permissionsRolesCmd shows permissions by role
var permissionsRolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "Show permissions by role",
	Long:  `Display what permissions each role has`,
	Run: func(cmd *cobra.Command, args []string) {
		showRolePermissions()
	},
}

func init() {
	rootCmd.AddCommand(permissionsCmd)
	permissionsCmd.AddCommand(permissionsListCmd)
	permissionsCmd.AddCommand(permissionsRolesCmd)
}

func listPermissions() {
	fmt.Println("📋 Available Permissions")
	fmt.Println()

	groups := models.AllPermissionGroups()
	for _, group := range groups {
		fmt.Printf("## %s\n", group.Name)
		fmt.Printf("   %s\n", group.Description)
		fmt.Println()
		for _, perm := range group.Permissions {
			fmt.Printf("   - %s\n", perm)
		}
		fmt.Println()
	}

	fmt.Printf("Total permissions: %d\n", len(models.AllPermissions()))
}

func showRolePermissions() {
	fmt.Println("👥 Permissions by Role")
	fmt.Println()

	rolePerms := models.PermissionsByRole()
	roles := []string{"admin", "editor", "viewer", "customer"}

	for _, role := range roles {
		perms := rolePerms[role]
		fmt.Printf("## %s (%d permissions)\n", role, len(perms))
		for _, perm := range perms {
			fmt.Printf("   ✓ %s\n", perm)
		}
		fmt.Println()
	}
}

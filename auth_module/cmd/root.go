package cmd

import (
	"fmt"
	"os"

	"github.com/sparque/auth_module/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *config.Config
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "auth-module",
	Short: "Multi-tenant authentication module",
	Long: `A flexible, domain-scoped authentication module supporting:
- Google OAuth and magic link authentication
- Multi-tenant architecture with complete domain isolation
- Invitation system with QR code generation
- Role-based access control`,
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
}

// initConfig reads in config file and ENV variables
func initConfig() {
	var err error
	cfg, err = config.Load(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Config validation error: %v\n", err)
		os.Exit(1)
	}
}

// GetConfig returns the loaded configuration
func GetConfig() *config.Config {
	return cfg
}

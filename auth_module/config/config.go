package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig      `mapstructure:"server"`
	MongoDB    MongoDBConfig     `mapstructure:"mongodb"`
	JWT        JWTConfig         `mapstructure:"jwt"`
	Google     GoogleOAuthConfig `mapstructure:"google"`
	Email      EmailConfig       `mapstructure:"email"`
	MagicLink  MagicLinkConfig   `mapstructure:"magic_link"`
	Invitation InvitationConfig  `mapstructure:"invitation"`
	App        AppConfig         `mapstructure:"app"`
}

type ServerConfig struct {
	Port string
	Env  string
	Host string
}

type MongoDBConfig struct {
	URI      string
	Database string
}

type JWTConfig struct {
	Secret      string
	ExpiryHours int `mapstructure:"expiry_hours"`
}

type GoogleOAuthConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	CallbackURL  string `mapstructure:"callback_url"`
}

type EmailConfig struct {
	SMTP        SMTPConfig
	FromAddress string `mapstructure:"from_address"`
}

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

type MagicLinkConfig struct {
	ExpiryMinutes int `mapstructure:"expiry_minutes"`
}

type InvitationConfig struct {
	Defaults InvitationDefaults
}

type InvitationDefaults struct {
	EmailExpiryHours        int `mapstructure:"email_expiry_hours"`
	QRCodeExpiryHours       int `mapstructure:"qr_code_expiry_hours"`
	PromotionalExpiryHours  int `mapstructure:"promotional_expiry_hours"`
}

type AppConfig struct {
	URL         string
	FrontendURL string `mapstructure:"frontend_url"`
}

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set config file path
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// Default config locations
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("/etc/auth_module")
	}

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		log.Println("No config file found, using environment variables and defaults")
	}

	// Environment variables override config file
	v.SetEnvPrefix("AUTH")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Explicitly bind critical environment variables
	v.BindEnv("magic_link.expiry_minutes", "AUTH_MAGIC_LINK_EXPIRY_MINUTES")

	// Set defaults
	setDefaults(v)

	// Unmarshal into Config struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Debug logging
	log.Printf("🔧 Config loaded:")
	log.Printf("   magic_link.expiry_minutes from viper: %v", v.Get("magic_link.expiry_minutes"))
	log.Printf("   MagicLink.ExpiryMinutes in struct: %d", cfg.MagicLink.ExpiryMinutes)

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.env", "development")
	v.SetDefault("server.host", "0.0.0.0")

	// MongoDB defaults
	v.SetDefault("mongodb.uri", "mongodb://localhost:27017")
	v.SetDefault("mongodb.database", "auth_module")

	// JWT defaults
	v.SetDefault("jwt.expiry_hours", 24)

	// Email defaults
	v.SetDefault("email.smtp.host", "smtp.gmail.com")
	v.SetDefault("email.smtp.port", 587)
	v.SetDefault("email.from_address", "noreply@example.com")

	// Magic link defaults
	v.SetDefault("magic_link.expiry_minutes", 15)

	// Invitation defaults
	v.SetDefault("invitation.defaults.email_expiry_hours", 24)
	v.SetDefault("invitation.defaults.qr_code_expiry_hours", 72)
	v.SetDefault("invitation.defaults.promotional_expiry_hours", 720)

	// App defaults
	v.SetDefault("app.url", "http://localhost:8080")
	v.SetDefault("app.frontend_url", "http://localhost:3000")
}

// Validate checks if required configuration is present
func (c *Config) Validate() error {
	if c.JWT.Secret == "" {
		return fmt.Errorf("jwt.secret is required")
	}
	if c.MongoDB.URI == "" {
		return fmt.Errorf("mongodb.uri is required")
	}
	return nil
}

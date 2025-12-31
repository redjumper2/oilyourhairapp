package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config holds all configuration
type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	MongoDB MongoDBConfig `mapstructure:"mongodb"`
	JWT     JWTConfig     `mapstructure:"jwt"`
	CORS    CORSConfig    `mapstructure:"cors"`
	Auth    AuthConfig    `mapstructure:"auth"`
}

// ServerConfig holds server settings
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Env  string `mapstructure:"env"`
}

// MongoDBConfig holds MongoDB connection settings
type MongoDBConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
}

// JWTConfig holds JWT settings (shared secret with auth_module)
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

// CORSConfig holds CORS settings
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

// AuthConfig holds auth module connection settings
type AuthConfig struct {
	BaseURL        string `mapstructure:"base_url"`
	DomainsDB      string `mapstructure:"domains_db"`       // auth_module database name
	DomainsCollection string `mapstructure:"domains_collection"` // domains collection name
}

var cfg *Config

// Load reads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", "9091")
	v.SetDefault("server.env", "development")
	v.SetDefault("mongodb.uri", "mongodb://localhost:27017")
	v.SetDefault("mongodb.database", "products_module")
	v.SetDefault("cors.allowed_origins", []string{"http://localhost:3000"})
	v.SetDefault("auth.base_url", "http://localhost:8080")
	v.SetDefault("auth.domains_db", "auth_module")
	v.SetDefault("auth.domains_collection", "domains")

	// Read config file if provided
	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			log.Printf("⚠️  Config file not found (%s), using defaults and environment variables", configPath)
		} else {
			log.Printf("✅ Config file loaded: %s", configPath)
		}
	}

	// Environment variables override config file
	v.SetEnvPrefix("PRODUCTS")
	v.AutomaticEnv()

	// Explicitly bind critical environment variables
	v.BindEnv("server.port", "PRODUCTS_SERVER_PORT")
	v.BindEnv("server.env", "PRODUCTS_SERVER_ENV")
	v.BindEnv("mongodb.uri", "PRODUCTS_MONGODB_URI")
	v.BindEnv("mongodb.database", "PRODUCTS_MONGODB_DATABASE")
	v.BindEnv("jwt.secret", "PRODUCTS_JWT_SECRET")
	v.BindEnv("cors.allowed_origins", "PRODUCTS_CORS_ALLOWED_ORIGINS")
	v.BindEnv("auth.base_url", "PRODUCTS_AUTH_BASE_URL")
	v.BindEnv("auth.domains_db", "PRODUCTS_AUTH_DOMAINS_DB")

	// Unmarshal config
	cfg = &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate required fields
	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required (must match auth_module)")
	}

	log.Printf("✅ Configuration loaded:")
	log.Printf("   Server: %s:%s (%s)", cfg.Server.Host, cfg.Server.Port, cfg.Server.Env)
	log.Printf("   MongoDB: %s/%s", cfg.MongoDB.URI, cfg.MongoDB.Database)
	log.Printf("   Auth DB: %s/%s", cfg.Auth.DomainsDB, cfg.Auth.DomainsCollection)

	return cfg, nil
}

// Get returns the global config instance
func Get() *Config {
	return cfg
}

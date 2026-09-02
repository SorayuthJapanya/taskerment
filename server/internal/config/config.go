package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env  string
	Port string
	TimeZone string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret  string
	AccessTTL  time.Duration
	RefreshTTL time.Duration

	CORSOrigins []string
}

// Lead reads .env
func Load() (*Config, error) {
	_ = godotenv.Load() // missing .env is fine in production

	cfg := &Config{
		Env:         getEnv("APP_ENV", "development"),
		Port:        getEnv("APP_PORT", "8080"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "taskmanager"),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", "taskmanager"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
		TimeZone: 	 getEnv("APP_TIMEZONE", "UTC"),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		CORSOrigins: []string{getEnv("CORS_ORIGIN", "http://localhost:5173")},
	}

	var err error

	// Validate Access TTL
	if cfg.AccessTTL, err = getDuration("JWT_ACCESS_TTL", 15*time.Minute); err != nil {
		return  nil, err
	}

	// Validate Refresh TTL
	if cfg.AccessTTL, err = getDuration("JWT_REFRESH_TTL", 168*time.Minute); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return  nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}
	if c.DBPassword == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return  v
	}
	return  fallback
}

func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
        c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode, c.TimeZone,
    )
}

func getDuration(key string, fallback time.Duration) (time.Duration, error) {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback, nil
	}

	d, err := time.ParseDuration(raw)
	if err != nil {
		return  0, fmt.Errorf("Invalid duration for %s: %w", key, err)
	}

	return d, nil
}

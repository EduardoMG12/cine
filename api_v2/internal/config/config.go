package config

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Email    EmailConfig    `mapstructure:"email"`
	TMDb     TMDbConfig     `mapstructure:"tmdb"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type DatabaseConfig struct {
	URL             string `mapstructure:"url"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	Expiration int    `mapstructure:"expiration"` // hours
}

type EmailConfig struct {
	SMTPHost     string `mapstructure:"smtp_host"`
	SMTPPort     int    `mapstructure:"smtp_port"`
	SMTPUsername string `mapstructure:"smtp_username"`
	SMTPPassword string `mapstructure:"smtp_password"`
	FromEmail    string `mapstructure:"from_email"`
	FromName     string `mapstructure:"from_name"`
}

type TMDbConfig struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")

	// Enable reading from environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CINE")

	// Set defaults
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		slog.Warn("Config file not found, using environment variables and defaults", "error", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")

	// Database defaults
	viper.SetDefault("database.url", "postgres://cineverse:password@postgres:5432/cineverse?sslmode=disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", 300) // 5 minutes

	// Redis defaults
	viper.SetDefault("redis.addr", "redis:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// JWT defaults
	viper.SetDefault("jwt.secret", "your-secret-key-change-in-production")
	viper.SetDefault("jwt.expiration", 24) // 24 hours

	// Email defaults (placeholder values)
	viper.SetDefault("email.smtp_host", "smtp.gmail.com")
	viper.SetDefault("email.smtp_port", 587)
	viper.SetDefault("email.smtp_username", "")
	viper.SetDefault("email.smtp_password", "")
	viper.SetDefault("email.from_email", "noreply@cineverse.com")
	viper.SetDefault("email.from_name", "CineVerse")

	// TMDb defaults
	viper.SetDefault("tmdb.api_key", "")
	viper.SetDefault("tmdb.base_url", "https://api.themoviedb.org/3")
}

// NewValidator creates and returns a new validator instance
func NewValidator() *validator.Validate {
	return validator.New()
}

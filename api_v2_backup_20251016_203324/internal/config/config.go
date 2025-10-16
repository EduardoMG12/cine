package config

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Redis       RedisConfig       `mapstructure:"redis"`
	JWT         JWTConfig         `mapstructure:"jwt"`
	Email       EmailConfig       `mapstructure:"email"`
	TMDb        TMDbConfig        `mapstructure:"tmdb"`
	Application ApplicationConfig `mapstructure:"application"`
	I18n        I18nConfig        `mapstructure:"i18n"`
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

type ApplicationConfig struct {
	DefaultTheme         string `mapstructure:"default_theme"`
	SessionDurationHours int    `mapstructure:"session_duration_hours"`
	TokenLength          int    `mapstructure:"token_length"`
	BcryptCost           int    `mapstructure:"bcrypt_cost"`
	Environment          string `mapstructure:"environment"`
	Debug                bool   `mapstructure:"debug"`
	LogLevel             string `mapstructure:"log_level"`
	EnableSwagger        bool   `mapstructure:"enable_swagger"`
	BaseURL              string `mapstructure:"base_url"`
	SwaggerURL           string `mapstructure:"swagger_url"`
}

type I18nConfig struct {
	DefaultLanguage    string   `mapstructure:"default_language"`
	SupportedLanguages []string `mapstructure:"supported_languages"`
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

	// Application defaults
	viper.SetDefault("application.default_theme", "light")
	viper.SetDefault("application.session_duration_hours", 24)
	viper.SetDefault("application.token_length", 32)
	viper.SetDefault("application.bcrypt_cost", 12)
	viper.SetDefault("application.environment", "development")
	viper.SetDefault("application.debug", true)
	viper.SetDefault("application.log_level", "info")
	viper.SetDefault("application.enable_swagger", true)
	viper.SetDefault("application.base_url", "http://localhost:8080")
	viper.SetDefault("application.swagger_url", "http://localhost:8080/swagger/doc.json")

	// I18n defaults
	viper.SetDefault("i18n.default_language", "en")
	viper.SetDefault("i18n.supported_languages", []string{"en", "pt", "es"})
}

// NewValidator creates and returns a new validator instance
func NewValidator() *validator.Validate {
	return validator.New()
}

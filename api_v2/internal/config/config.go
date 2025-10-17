package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	JWT      JWTConfig      `json:"jwt"`
	TMDb     TMDbConfig     `json:"tmdb"`
	Redis    RedisConfig    `json:"redis"`
}

type ServerConfig struct {
	Port            string        `json:"port"`
	Host            string        `json:"host"`
	Environment     string        `json:"environment"`
	ReadTimeout     time.Duration `json:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout"`
}

type DatabaseConfig struct {
	Host            string `json:"host"`
	Port            int    `json:"port"`
	Name            string `json:"name"`
	User            string `json:"user"`
	Password        string `json:"password"`
	SSLMode         string `json:"ssl_mode"`
	MaxOpenConns    int    `json:"max_open_conns"`
	MaxIdleConns    int    `json:"max_idle_conns"`
	ConnMaxLifetime int    `json:"conn_max_lifetime"` // in minutes
}

type JWTConfig struct {
	Secret         string        `json:"-"`
	ExpirationTime time.Duration `json:"expiration_time"`
	RefreshTime    time.Duration `json:"refresh_time"`
	Issuer         string        `json:"issuer"`
}

type TMDbConfig struct {
	APIKey       string        `json:"-"`
	BaseURL      string        `json:"base_url"`
	ImageBaseURL string        `json:"image_base_url"`
	CacheTTL     time.Duration `json:"cache_ttl"`
	RateLimit    int           `json:"rate_limit"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"-"`
	DB       int    `json:"db"`
}

func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "8080"),
			Host:            getEnv("HOST", "0.0.0.0"),
			Environment:     getEnv("ENVIRONMENT", "development"),
			ReadTimeout:     getEnvDuration("READ_TIMEOUT", "15s"),
			WriteTimeout:    getEnvDuration("WRITE_TIMEOUT", "15s"),
			ShutdownTimeout: getEnvDuration("SHUTDOWN_TIMEOUT", "30s"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvInt("DB_PORT", 5432),
			Name:            getEnv("DB_NAME", "cineverse"),
			User:            getEnv("DB_USER", "cineverse"),
			Password:        getEnv("DB_PASSWORD", "cineverse"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: getEnvInt("DB_CONN_MAX_LIFETIME", 5),
		},
		JWT: JWTConfig{
			Secret:         getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this"),
			ExpirationTime: getEnvDuration("JWT_EXPIRATION", "24h"),
			RefreshTime:    getEnvDuration("JWT_REFRESH", "7d"),
			Issuer:         getEnv("JWT_ISSUER", "cineverse-api"),
		},
		TMDb: TMDbConfig{
			APIKey:       getEnv("TMDB_API_KEY", ""),
			BaseURL:      getEnv("TMDB_BASE_URL", "https://api.themoviedb.org/3"),
			ImageBaseURL: getEnv("TMDB_IMAGE_BASE_URL", "https://image.tmdb.org/t/p/"),
			CacheTTL:     getEnvDuration("TMDB_CACHE_TTL", "24h"),
			RateLimit:    getEnvInt("TMDB_RATE_LIMIT", 40),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
	}

	return config, config.Validate()
}

func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if c.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}

	if c.JWT.Secret == "" || c.JWT.Secret == "your-super-secret-jwt-key-change-this" {
		log.Println("WARNING: Using default JWT secret. Set JWT_SECRET environment variable for production!")
	}

	if c.TMDb.APIKey == "" {
		log.Println("WARNING: TMDb API key not configured. Movie data will not be available!")
	}

	return nil
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

func (c *ServerConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (c *ServerConfig) IsProduction() bool {
	return c.Environment == "production"
}

func (c *ServerConfig) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *RedisConfig) GetRedisAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvDuration(key, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}

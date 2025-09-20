package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Domain   string
	Port     string
	BaseURL  string
	AuthToken string

	DatabaseURL string
	RedisURL    string

	RefreshCron           string
	HTTPTimeout           time.Duration
	MaxConcurrentFetches  int

	EnableDirectDownload bool
	StorageBackend       string
	DownloadDir          string

	S3Endpoint   string
	S3Bucket     string
	S3AccessKey  string
	S3SecretKey  string
	S3Region     string

	LogLevel  string
	LogFormat string

	CorsOrigins                string
	RateLimitRequestsPerMinute int
}

func Load() *Config {
	return &Config{
		Domain:    getEnv("DOMAIN", "localhost"),
		Port:      getEnv("PORT", "8080"),
		BaseURL:   getEnv("BASE_URL", "http://localhost:8080"),
		AuthToken: getEnv("AUTH_TOKEN", "change-me"),

		DatabaseURL: getEnv("DB_URL", "postgres://alldl:alldl@localhost:5432/alldownloads?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379/0"),

		RefreshCron:           getEnv("REFRESH_CRON", "@every 6h"),
		HTTPTimeout:           getDurationEnv("HTTP_TIMEOUT", 15*time.Second),
		MaxConcurrentFetches:  getIntEnv("MAX_CONCURRENT_FETCHES", 6),

		EnableDirectDownload: getBoolEnv("ENABLE_DIRECT_DOWNLOAD", false),
		StorageBackend:       getEnv("STORAGE_BACKEND", "local"),
		DownloadDir:          getEnv("DOWNLOAD_DIR", "/data"),

		S3Endpoint:  getEnv("S3_ENDPOINT", "http://localhost:9000"),
		S3Bucket:    getEnv("S3_BUCKET", "downloads"),
		S3AccessKey: getEnv("S3_ACCESS_KEY", "minioadmin"),
		S3SecretKey: getEnv("S3_SECRET_KEY", "minioadmin"),
		S3Region:    getEnv("S3_REGION", "us-east-1"),

		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "json"),

		CorsOrigins:                getEnv("CORS_ORIGINS", "http://localhost:3000"),
		RateLimitRequestsPerMinute: getIntEnv("RATE_LIMIT_REQUESTS_PER_MINUTE", 60),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	AppName string
	AppEnv  string
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	JWTSecret  string

	QdrantKey      string
	QdrantHost     string
	QdrantPort     int
	QdrantGRPCPort int

	MinioUser   string
	MinioPass   string
	MinioPort   int
	MinioHost   string
	MinioBucket string

	GeminiKey string

	UnidocKey string

	Logger *logrus.Logger
}

func LoadConfig() *Config {
	// Load .env file if present (safe to ignore in prod)
	_ = godotenv.Load()

	// Configure logger
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	appEnv := getEnv("APP_ENV", "development")

	if appEnv == "production" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
			PadLevelText:    true,
		})
		logger.SetLevel(logrus.DebugLevel)
	}

	cfg := &Config{
		AppName:        getEnv("APP_NAME", "myapp"),
		AppEnv:         appEnv,
		AppPort:        getEnv("APP_PORT", "8080"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "myapp_db"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		JWTSecret:      getEnv("JWT_SECRET", "super-secret-key"),
		QdrantKey:      getEnv("QDRANT_API_KEY", ""),
		QdrantHost:     getEnv("QDRANT_HOST", "localhost"),
		QdrantPort:     getEnv("QDRANT_PORT", 6333),
		QdrantGRPCPort: getEnv("QDRANT_GRPC_PORT", 6334),
		MinioUser:      getEnv("MINIO_ROOT_USER", ""),
		MinioPass:      getEnv("MINIO_ROOT_PASSWORD", ""),
		MinioPort:      getEnv("MINIO_API_PORT", 9000),
		MinioHost:      getEnv("MINIO_API_HOST", "localhost"),
		GeminiKey:      getEnv("GEMINI_API_KEY", ""),
		MinioBucket:    "cvbucket",
		UnidocKey:      getEnv("UNIDOC_KEY", ""),
		Logger:         logger,
	}

	cfg.Logger.Debugf("Config loaded: env=%s port=%s db=%s", cfg.AppEnv, cfg.AppPort, cfg.DBName)
	return cfg
}

// DatabaseDSN returns a Postgres DSN string (for GORM or pgx)
func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort, c.DBSSLMode,
	)
}

func getEnv[T any](key string, fallback T) T {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	switch any(fallback).(type) {
	case int:
		i, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}
		return any(i).(T)
	case string:
		// remove accidental quotes or whitespace
		value = strings.Trim(value, `"' `)
		return any(value).(T)
	default:
		return fallback
	}
}

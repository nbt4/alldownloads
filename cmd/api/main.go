package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/your-username/alldownloads/internal/api"
	"github.com/your-username/alldownloads/internal/config"
	"github.com/your-username/alldownloads/internal/jobs"
	"github.com/your-username/alldownloads/internal/middleware"
	"github.com/your-username/alldownloads/internal/store"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.Load()

	logger, err := createLogger(cfg.LogLevel, cfg.LogFormat)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	postgresStore, err := store.NewPostgresStore(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer postgresStore.Close()

	jobQueue, err := jobs.NewQueue(cfg.RedisURL, logger)
	if err != nil {
		logger.Fatal("Failed to create job queue", zap.Error(err))
	}
	defer jobQueue.Close()

	handler := api.NewHandler(postgresStore, jobQueue, logger, cfg.AuthToken)

	if cfg.LogFormat == "json" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	rateLimiter := middleware.NewIPRateLimiter(cfg.RateLimitRequestsPerMinute)

	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.CORS(cfg.CorsOrigins))
	router.Use(rateLimiter.RateLimit())
	router.Use(api.PrometheusMetrics())

	v1 := router.Group("/api")
	{
		v1.GET("/health", handler.HealthCheck)
		v1.GET("/products", handler.GetProducts)
		v1.GET("/products/:id", handler.GetProduct)
		v1.POST("/refresh", handler.RefreshProducts)
	}

	router.GET("/metrics", api.MetricsHandler())

	srv := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		logger.Info("Starting server", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

func createLogger(level string, format string) (*zap.Logger, error) {
	var config zap.Config

	if format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	return config.Build()
}
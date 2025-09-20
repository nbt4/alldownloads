package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/your-username/alldownloads/internal/config"
	"github.com/your-username/alldownloads/internal/jobs"
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

	worker := jobs.NewWorker(postgresStore, jobQueue, logger, cfg.MaxConcurrentFetches)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := worker.Start(ctx); err != nil {
			logger.Error("worker failed", zap.Error(err))
		}
	}()

	scheduler := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(logger.Sugar())))

	_, err = scheduler.AddFunc(cfg.RefreshCron, func() {
		logger.Info("scheduled refresh triggered")

		products, err := postgresStore.GetProducts(context.Background())
		if err != nil {
			logger.Error("failed to get products for scheduled refresh", zap.Error(err))
			return
		}

		for _, product := range products {
			job := &store.FetchJob{
				ProductID: product.ID,
				Status:    store.JobStatusPending,
			}

			if err := postgresStore.CreateFetchJob(context.Background(), job); err != nil {
				logger.Error("failed to create scheduled fetch job", zap.Error(err), zap.String("product_id", product.ID))
				continue
			}

			if err := jobQueue.Enqueue(context.Background(), job.ID); err != nil {
				logger.Error("failed to enqueue scheduled job", zap.Error(err), zap.String("job_id", job.ID))
				continue
			}
		}

		logger.Info("scheduled refresh completed", zap.Int("jobs_queued", len(products)))
	})

	if err != nil {
		logger.Fatal("Failed to add cron job", zap.Error(err))
	}

	scheduler.Start()
	logger.Info("scheduler started", zap.String("cron", cfg.RefreshCron))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down worker...")

	scheduler.Stop()
	cancel()
	wg.Wait()

	logger.Info("worker exited")
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
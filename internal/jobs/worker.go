package jobs

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/your-username/alldownloads/internal/api"
	"github.com/your-username/alldownloads/internal/sources"
	"github.com/your-username/alldownloads/internal/store"
	"go.uber.org/zap"
)

type Worker struct {
	store      *store.PostgresStore
	queue      *Queue
	fetchers   map[string]sources.Fetcher
	logger     *zap.Logger
	maxWorkers int
}

func NewWorker(store *store.PostgresStore, queue *Queue, logger *zap.Logger, maxWorkers int) *Worker {
	fetchers := map[string]sources.Fetcher{
		"ubuntu":             sources.NewUbuntuFetcher(),
		"debian":             sources.NewDebianFetcher(),
		"arch":               sources.NewArchFetcher(),
		"kali":               sources.NewKaliFetcher(),
		"windows":            sources.NewWindowsFetcher(),
		"termius":            sources.NewTermiusFetcher(),
		"telegram":           sources.NewTelegramFetcher(),
		"whatsapp":           sources.NewWhatsAppFetcher(),
		"tailscale":          sources.NewTailscaleFetcher(),
		"nextcloud":          sources.NewNextcloudFetcher(),
		"chrome":             sources.NewChromeFetcher(),
		"firefox":            sources.NewFirefoxFetcher(),
		"brave":              sources.NewBraveFetcher(),
		"vscode":             sources.NewVSCodeFetcher(),
		"notepadplusplus":    sources.NewNotepadPlusPlusFetcher(),
		"powershell":         sources.NewPowerShellFetcher(),
	}

	return &Worker{
		store:      store,
		queue:      queue,
		fetchers:   fetchers,
		logger:     logger,
		maxWorkers: maxWorkers,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	w.logger.Info("starting worker", zap.Int("max_workers", w.maxWorkers))

	var wg sync.WaitGroup

	for i := 0; i < w.maxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			w.workerLoop(ctx, workerID)
		}(i)
	}

	wg.Wait()
	w.logger.Info("all workers stopped")
	return nil
}

func (w *Worker) workerLoop(ctx context.Context, workerID int) {
	logger := w.logger.With(zap.Int("worker_id", workerID))
	logger.Info("worker started")

	for {
		select {
		case <-ctx.Done():
			logger.Info("worker stopping due to context cancellation")
			return
		default:
			message, err := w.queue.Dequeue(ctx, 5*time.Second)
			if err != nil {
				logger.Error("failed to dequeue job", zap.Error(err))
				continue
			}

			if message == nil {
				continue
			}

			logger.Info("processing job", zap.String("job_id", message.ID))

			if err := w.processJob(ctx, message); err != nil {
				logger.Error("job processing failed", zap.Error(err), zap.String("job_id", message.ID))

				if retryErr := w.queue.RetryJob(ctx, message); retryErr != nil {
					logger.Error("failed to retry job", zap.Error(retryErr), zap.String("job_id", message.ID))
				}

				api.IncrementFetchJobMetric("failed")
			} else {
				logger.Info("job completed successfully", zap.String("job_id", message.ID))

				if err := w.queue.MarkCompleted(ctx, message.ID); err != nil {
					logger.Error("failed to mark job as completed", zap.Error(err), zap.String("job_id", message.ID))
				}

				api.IncrementFetchJobMetric("completed")
			}
		}
	}
}

func (w *Worker) processJob(ctx context.Context, message *JobMessage) error {
	job := &store.FetchJob{
		ID:        message.ID,
		Status:    store.JobStatusRunning,
		StartedAt: &[]time.Time{time.Now()}[0],
	}

	if err := w.store.UpdateFetchJob(ctx, job); err != nil {
		return fmt.Errorf("failed to update job status to running: %w", err)
	}

	jobFromDB, err := w.store.GetFetchJob(ctx, message.ID)
	if err != nil {
		return fmt.Errorf("failed to get job from database: %w", err)
	}

	product, err := w.store.GetProduct(ctx, jobFromDB.ProductID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	if product == nil {
		return fmt.Errorf("product not found: %s", jobFromDB.ProductID)
	}

	fetcher, exists := w.fetchers[product.ID]
	if !exists {
		return fmt.Errorf("no fetcher available for product: %s", product.ID)
	}

	versions, err := fetcher.Fetch(ctx)
	if err != nil {
		job.Status = store.JobStatusFailed
		job.Error = err.Error()
		completedAt := time.Now()
		job.CompletedAt = &completedAt

		if updateErr := w.store.UpdateFetchJob(ctx, job); updateErr != nil {
			w.logger.Error("failed to update failed job", zap.Error(updateErr), zap.String("job_id", message.ID))
		}

		return fmt.Errorf("fetcher failed: %w", err)
	}

	for _, version := range versions {
		version.ProductID = product.ID
		if err := w.store.CreateOrUpdateProductVersion(ctx, version); err != nil {
			w.logger.Error("failed to save product version", zap.Error(err), zap.String("version", version.Version))
		}
	}

	if err := w.store.MarkLatestVersions(ctx, product.ID); err != nil {
		w.logger.Error("failed to mark latest versions", zap.Error(err), zap.String("product_id", product.ID))
	}

	job.Status = store.JobStatusCompleted
	completedAt := time.Now()
	job.CompletedAt = &completedAt

	if err := w.store.UpdateFetchJob(ctx, job); err != nil {
		return fmt.Errorf("failed to update job status to completed: %w", err)
	}

	api.SetProductVersionsMetric(product.ID, float64(len(versions)))

	w.logger.Info("product updated", zap.String("product_id", product.ID), zap.Int("versions", len(versions)))

	return nil
}
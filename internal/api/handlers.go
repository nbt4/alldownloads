package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/your-username/alldownloads/internal/jobs"
	"github.com/your-username/alldownloads/internal/store"
	"go.uber.org/zap"
)

type Handler struct {
	store     *store.PostgresStore
	jobQueue  *jobs.Queue
	logger    *zap.Logger
	authToken string
}

func NewHandler(store *store.PostgresStore, jobQueue *jobs.Queue, logger *zap.Logger, authToken string) *Handler {
	return &Handler{
		store:     store,
		jobQueue:  jobQueue,
		logger:    logger,
		authToken: authToken,
	}
}

func (h *Handler) GetProducts(c *gin.Context) {
	ctx := c.Request.Context()

	products, err := h.store.GetProducts(ctx)
	if err != nil {
		h.logger.Error("failed to get products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func (h *Handler) GetProduct(c *gin.Context) {
	ctx := c.Request.Context()
	productID := c.Param("id")

	productWithVersions, err := h.store.GetProductWithVersions(ctx, productID)
	if err != nil {
		h.logger.Error("failed to get product", zap.Error(err), zap.String("product_id", productID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		return
	}

	if productWithVersions == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, productWithVersions)
}

func (h *Handler) RefreshProducts(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "Bearer "+h.authToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization"})
		return
	}

	ctx := c.Request.Context()

	products, err := h.store.GetProducts(ctx)
	if err != nil {
		h.logger.Error("failed to get products for refresh", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh products"})
		return
	}

	var queuedJobs []string
	for _, product := range products {
		job := &store.FetchJob{
			ProductID: product.ID,
			Status:    store.JobStatusPending,
		}

		if err := h.store.CreateFetchJob(ctx, job); err != nil {
			h.logger.Error("failed to create fetch job", zap.Error(err), zap.String("product_id", product.ID))
			continue
		}

		if err := h.jobQueue.Enqueue(ctx, job.ID); err != nil {
			h.logger.Error("failed to enqueue job", zap.Error(err), zap.String("job_id", job.ID))
			continue
		}

		queuedJobs = append(queuedJobs, job.ID)
	}

	h.logger.Info("refresh initiated", zap.Int("jobs_queued", len(queuedJobs)))

	c.JSON(http.StatusOK, gin.H{
		"message":     "Refresh initiated",
		"jobs_queued": len(queuedJobs),
		"job_ids":     queuedJobs,
	})
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"version":   "0.1.0",
	})
}
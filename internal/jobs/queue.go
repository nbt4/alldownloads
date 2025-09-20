package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	QueueName      = "fetch_jobs"
	ProcessingSet  = "fetch_jobs:processing"
	RetryLimit     = 3
	RetryDelay     = 5 * time.Minute
)

type Queue struct {
	client *redis.Client
	logger *zap.Logger
}

type JobMessage struct {
	ID        string    `json:"id"`
	Retries   int       `json:"retries"`
	CreatedAt time.Time `json:"created_at"`
}

func NewQueue(redisURL string, logger *zap.Logger) (*Queue, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Queue{
		client: client,
		logger: logger,
	}, nil
}

func (q *Queue) Close() error {
	return q.client.Close()
}

func (q *Queue) Enqueue(ctx context.Context, jobID string) error {
	message := JobMessage{
		ID:        jobID,
		Retries:   0,
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal job message: %w", err)
	}

	if err := q.client.LPush(ctx, QueueName, data).Err(); err != nil {
		return fmt.Errorf("failed to enqueue job: %w", err)
	}

	q.logger.Info("job enqueued", zap.String("job_id", jobID))
	return nil
}

func (q *Queue) Dequeue(ctx context.Context, timeout time.Duration) (*JobMessage, error) {
	result, err := q.client.BRPop(ctx, timeout, QueueName).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to dequeue job: %w", err)
	}

	if len(result) != 2 {
		return nil, fmt.Errorf("unexpected result format from Redis")
	}

	var message JobMessage
	if err := json.Unmarshal([]byte(result[1]), &message); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job message: %w", err)
	}

	if err := q.client.SAdd(ctx, ProcessingSet, message.ID).Err(); err != nil {
		q.logger.Error("failed to add job to processing set", zap.Error(err), zap.String("job_id", message.ID))
	}

	return &message, nil
}

func (q *Queue) MarkCompleted(ctx context.Context, jobID string) error {
	if err := q.client.SRem(ctx, ProcessingSet, jobID).Err(); err != nil {
		return fmt.Errorf("failed to remove job from processing set: %w", err)
	}

	q.logger.Info("job completed", zap.String("job_id", jobID))
	return nil
}

func (q *Queue) RetryJob(ctx context.Context, message *JobMessage) error {
	if message.Retries >= RetryLimit {
		q.logger.Error("job exceeded retry limit", zap.String("job_id", message.ID), zap.Int("retries", message.Retries))
		return q.MarkCompleted(ctx, message.ID)
	}

	message.Retries++

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal retry job message: %w", err)
	}

	if err := q.client.LPush(ctx, QueueName, data).Err(); err != nil {
		return fmt.Errorf("failed to retry job: %w", err)
	}

	if err := q.client.SRem(ctx, ProcessingSet, message.ID).Err(); err != nil {
		q.logger.Error("failed to remove job from processing set during retry", zap.Error(err), zap.String("job_id", message.ID))
	}

	q.logger.Info("job retried", zap.String("job_id", message.ID), zap.Int("retries", message.Retries))
	return nil
}

func (q *Queue) GetQueueLength(ctx context.Context) (int64, error) {
	return q.client.LLen(ctx, QueueName).Result()
}

func (q *Queue) GetProcessingCount(ctx context.Context) (int64, error) {
	return q.client.SCard(ctx, ProcessingSet).Result()
}
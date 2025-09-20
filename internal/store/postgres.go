package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/google/uuid"
)

type PostgresStore struct {
	db *pgxpool.Pool
}

func NewPostgresStore(databaseURL string) (*PostgresStore, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Close() {
	s.db.Close()
}

func (s *PostgresStore) GetProducts(ctx context.Context) ([]Product, error) {
	query := `
		SELECT id, name, vendor, category, description, icon_url, website_url, created_at, updated_at
		FROM products
		ORDER BY vendor, name
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Vendor, &p.Category, &p.Description, &p.IconURL, &p.WebsiteURL, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}

	return products, nil
}

func (s *PostgresStore) GetProduct(ctx context.Context, id string) (*Product, error) {
	query := `
		SELECT id, name, vendor, category, description, icon_url, website_url, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var p Product
	err := s.db.QueryRow(ctx, query, id).Scan(&p.ID, &p.Name, &p.Vendor, &p.Category, &p.Description, &p.IconURL, &p.WebsiteURL, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &p, nil
}

func (s *PostgresStore) GetProductWithVersions(ctx context.Context, id string) (*ProductWithVersions, error) {
	product, err := s.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, nil
	}

	versions, err := s.GetProductVersions(ctx, id)
	if err != nil {
		return nil, err
	}

	return &ProductWithVersions{
		Product:  *product,
		Versions: versions,
	}, nil
}

func (s *PostgresStore) GetProductVersions(ctx context.Context, productID string) ([]ProductVersion, error) {
	query := `
		SELECT id, product_id, version, platform, architecture, download_url, checksum, checksum_type,
		       file_size, filename, is_latest, etag, last_fetched, created_at, updated_at
		FROM product_versions
		WHERE product_id = $1
		ORDER BY is_latest DESC, created_at DESC
	`

	rows, err := s.db.Query(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query product versions: %w", err)
	}
	defer rows.Close()

	var versions []ProductVersion
	for rows.Next() {
		var v ProductVersion
		err := rows.Scan(&v.ID, &v.ProductID, &v.Version, &v.Platform, &v.Architecture,
			&v.DownloadURL, &v.Checksum, &v.ChecksumType, &v.FileSize, &v.Filename,
			&v.IsLatest, &v.ETag, &v.LastFetched, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product version: %w", err)
		}
		versions = append(versions, v)
	}

	return versions, nil
}

func (s *PostgresStore) CreateProduct(ctx context.Context, product *Product) error {
	if product.ID == "" {
		product.ID = uuid.New().String()
	}
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	query := `
		INSERT INTO products (id, name, vendor, category, description, icon_url, website_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := s.db.Exec(ctx, query, product.ID, product.Name, product.Vendor, product.Category,
		product.Description, product.IconURL, product.WebsiteURL, product.CreatedAt, product.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

func (s *PostgresStore) UpdateProduct(ctx context.Context, product *Product) error {
	product.UpdatedAt = time.Now()

	query := `
		UPDATE products
		SET name = $2, vendor = $3, category = $4, description = $5, icon_url = $6, website_url = $7, updated_at = $8
		WHERE id = $1
	`

	_, err := s.db.Exec(ctx, query, product.ID, product.Name, product.Vendor, product.Category,
		product.Description, product.IconURL, product.WebsiteURL, product.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

func (s *PostgresStore) CreateOrUpdateProductVersion(ctx context.Context, version *ProductVersion) error {
	if version.ID == "" {
		version.ID = uuid.New().String()
	}
	version.LastFetched = time.Now()
	version.UpdatedAt = time.Now()
	if version.CreatedAt.IsZero() {
		version.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO product_versions (id, product_id, version, platform, architecture, download_url,
		                            checksum, checksum_type, file_size, filename, is_latest, etag,
		                            last_fetched, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (product_id, version, platform, architecture)
		DO UPDATE SET
			download_url = EXCLUDED.download_url,
			checksum = EXCLUDED.checksum,
			checksum_type = EXCLUDED.checksum_type,
			file_size = EXCLUDED.file_size,
			filename = EXCLUDED.filename,
			is_latest = EXCLUDED.is_latest,
			etag = EXCLUDED.etag,
			last_fetched = EXCLUDED.last_fetched,
			updated_at = EXCLUDED.updated_at
	`

	_, err := s.db.Exec(ctx, query, version.ID, version.ProductID, version.Version, version.Platform,
		version.Architecture, version.DownloadURL, version.Checksum, version.ChecksumType,
		version.FileSize, version.Filename, version.IsLatest, version.ETag,
		version.LastFetched, version.CreatedAt, version.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create or update product version: %w", err)
	}

	return nil
}

func (s *PostgresStore) MarkLatestVersions(ctx context.Context, productID string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "UPDATE product_versions SET is_latest = false WHERE product_id = $1", productID)
	if err != nil {
		return fmt.Errorf("failed to unmark latest versions: %w", err)
	}

	query := `
		UPDATE product_versions
		SET is_latest = true
		WHERE id IN (
			SELECT DISTINCT ON (platform, architecture) id
			FROM product_versions
			WHERE product_id = $1
			ORDER BY platform, architecture, created_at DESC
		)
	`

	_, err = tx.Exec(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("failed to mark latest versions: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *PostgresStore) CreateFetchJob(ctx context.Context, job *FetchJob) error {
	if job.ID == "" {
		job.ID = uuid.New().String()
	}
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()

	query := `
		INSERT INTO fetch_jobs (id, product_id, status, started_at, completed_at, error, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := s.db.Exec(ctx, query, job.ID, job.ProductID, job.Status, job.StartedAt,
		job.CompletedAt, job.Error, job.CreatedAt, job.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create fetch job: %w", err)
	}

	return nil
}

func (s *PostgresStore) UpdateFetchJob(ctx context.Context, job *FetchJob) error {
	job.UpdatedAt = time.Now()

	query := `
		UPDATE fetch_jobs
		SET status = $2, started_at = $3, completed_at = $4, error = $5, updated_at = $6
		WHERE id = $1
	`

	_, err := s.db.Exec(ctx, query, job.ID, job.Status, job.StartedAt,
		job.CompletedAt, job.Error, job.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update fetch job: %w", err)
	}

	return nil
}

func (s *PostgresStore) GetFetchJob(ctx context.Context, id string) (*FetchJob, error) {
	query := `
		SELECT id, product_id, status, started_at, completed_at, error, created_at, updated_at
		FROM fetch_jobs
		WHERE id = $1
	`

	var job FetchJob
	err := s.db.QueryRow(ctx, query, id).Scan(&job.ID, &job.ProductID, &job.Status,
		&job.StartedAt, &job.CompletedAt, &job.Error, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get fetch job: %w", err)
	}

	return &job, nil
}
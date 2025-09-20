package store

import (
	"time"
)

type Product struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Vendor      string    `json:"vendor" db:"vendor"`
	Category    string    `json:"category" db:"category"`
	Description string    `json:"description" db:"description"`
	IconURL     string    `json:"icon_url" db:"icon_url"`
	WebsiteURL  string    `json:"website_url" db:"website_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ProductVersion struct {
	ID          string    `json:"id" db:"id"`
	ProductID   string    `json:"product_id" db:"product_id"`
	Version     string    `json:"version" db:"version"`
	Platform    string    `json:"platform" db:"platform"`
	Architecture string   `json:"architecture" db:"architecture"`
	DownloadURL string    `json:"download_url" db:"download_url"`
	Checksum    string    `json:"checksum" db:"checksum"`
	ChecksumType string   `json:"checksum_type" db:"checksum_type"`
	FileSize    int64     `json:"file_size" db:"file_size"`
	Filename    string    `json:"filename" db:"filename"`
	IsLatest    bool      `json:"is_latest" db:"is_latest"`
	ETag        string    `json:"etag" db:"etag"`
	LastFetched time.Time `json:"last_fetched" db:"last_fetched"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ProductWithVersions struct {
	Product  Product          `json:"product"`
	Versions []ProductVersion `json:"versions"`
}

type FetchJob struct {
	ID          string    `json:"id" db:"id"`
	ProductID   string    `json:"product_id" db:"product_id"`
	Status      string    `json:"status" db:"status"`
	StartedAt   *time.Time `json:"started_at" db:"started_at"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
	Error       string    `json:"error" db:"error"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

const (
	JobStatusPending   = "pending"
	JobStatusRunning   = "running"
	JobStatusCompleted = "completed"
	JobStatusFailed    = "failed"
)

const (
	CategoryOS   = "os"
	CategoryApp  = "app"
	CategoryTool = "tool"
)

const (
	PlatformWindows = "windows"
	PlatformLinux   = "linux"
	PlatformMacOS   = "macos"
	PlatformWeb     = "web"
)

const (
	ArchAMD64 = "amd64"
	ArchARM64 = "arm64"
	Arch386   = "386"
	ArchARM   = "arm"
)
package sources

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/your-username/alldownloads/internal/store"
)

type ArchFetcher struct {
	client *HTTPClient
}

func NewArchFetcher() *ArchFetcher {
	return &ArchFetcher{
		client: NewHTTPClient(),
	}
}

func (f *ArchFetcher) Fetch(ctx context.Context) ([]*store.ProductVersion, error) {
	var versions []*store.ProductVersion

	baseURL := "https://archlinux.org/download/"
	resp, err := f.client.Get(ctx, baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Arch Linux download page: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	html := string(body)

	mirrorRegex := regexp.MustCompile(`href="([^"]*archlinux-[^"]*\.iso)"`)
	mirrorMatches := mirrorRegex.FindAllStringSubmatch(html, -1)

	if len(mirrorMatches) == 0 {
		return nil, fmt.Errorf("no Arch Linux ISO links found")
	}

	downloadURL := mirrorMatches[0][1]
	if !strings.HasPrefix(downloadURL, "http") {
		return nil, fmt.Errorf("invalid download URL: %s", downloadURL)
	}

	filename := extractFilename(downloadURL)

	versionRegex := regexp.MustCompile(`archlinux-([0-9]{4}\.[0-9]{2}\.[0-9]{2})`)
	versionMatch := versionRegex.FindStringSubmatch(filename)

	version := "latest"
	if len(versionMatch) >= 2 {
		version = versionMatch[1]
	}

	fileSize := getFileSizeFromURL(ctx, f.client, downloadURL)

	checksum, err := f.fetchChecksum(ctx, downloadURL)
	if err != nil {
		checksum = ""
	}

	pv := &store.ProductVersion{
		Version:      version,
		Platform:     store.PlatformLinux,
		Architecture: store.ArchAMD64,
		DownloadURL:  downloadURL,
		Checksum:     checksum,
		ChecksumType: "sha256",
		FileSize:     fileSize,
		Filename:     filename,
		IsLatest:     true,
	}

	versions = append(versions, pv)

	return versions, nil
}

func (f *ArchFetcher) fetchChecksum(ctx context.Context, isoURL string) (string, error) {
	checksumURL := strings.Replace(isoURL, ".iso", ".iso.sha256", 1)

	resp, err := f.client.Get(ctx, checksumURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("checksum file not found: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	checksum := strings.TrimSpace(string(body))
	parts := strings.Fields(checksum)
	if len(parts) > 0 {
		return parts[0], nil
	}

	return checksum, nil
}
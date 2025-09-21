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

	// Use direct mirror URLs since main page has magnet links
	mirrors := []string{
		"https://mirror.rackspace.com/archlinux/iso/latest/archlinux-x86_64.iso",
		"https://mirrors.kernel.org/archlinux/iso/latest/archlinux-x86_64.iso",
	}

	var downloadURL string
	var fileSize int64

	// Try mirrors until we find one that works
	for _, mirror := range mirrors {
		resp, err := f.client.Head(ctx, mirror)
		if err == nil && resp.StatusCode == 200 {
			downloadURL = mirror
			if resp.ContentLength > 0 {
				fileSize = resp.ContentLength
			}
			resp.Body.Close()
			break
		}
		if resp != nil {
			resp.Body.Close()
		}
	}

	if downloadURL == "" {
		return nil, fmt.Errorf("no working Arch Linux mirrors found")
	}

	filename := "archlinux-x86_64.iso"
	version := "latest"

	// Try to get actual version from filename if mirror redirects
	if resp, err := f.client.Head(ctx, downloadURL); err == nil {
		if finalURL := resp.Header.Get("Location"); finalURL != "" {
			if versionRegex := regexp.MustCompile(`archlinux-([0-9]{4}\.[0-9]{2}\.[0-9]{2})`); versionRegex.MatchString(finalURL) {
				matches := versionRegex.FindStringSubmatch(finalURL)
				if len(matches) >= 2 {
					version = matches[1]
					filename = extractFilename(finalURL)
				}
			}
		}
		resp.Body.Close()
	}

	if fileSize == 0 {
		fileSize = getFileSizeFromURL(ctx, f.client, downloadURL)
	}

	pv := &store.ProductVersion{
		Version:      version,
		Platform:     store.PlatformLinux,
		Architecture: store.ArchAMD64,
		DownloadURL:  downloadURL,
		Checksum:     "",
		ChecksumType: "",
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
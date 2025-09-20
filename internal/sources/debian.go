package sources

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/your-username/alldownloads/internal/store"
)

type DebianFetcher struct {
	client *HTTPClient
}

func NewDebianFetcher() *DebianFetcher {
	return &DebianFetcher{
		client: NewHTTPClient(),
	}
}

func (f *DebianFetcher) Fetch(ctx context.Context) ([]*store.ProductVersion, error) {
	var versions []*store.ProductVersion

	baseURL := "https://cdimage.debian.org/debian-cd/current/amd64/iso-cd/"
	resp, err := f.client.Get(ctx, baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Debian ISOs: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	html := string(body)

	isoRegex := regexp.MustCompile(`href="(debian-[^"]*\.iso)"`)
	isoMatches := isoRegex.FindAllStringSubmatch(html, -1)

	versionRegex := regexp.MustCompile(`debian-([0-9]+\.[0-9]+(?:\.[0-9]+)?)`)

	for _, match := range isoMatches {
		if len(match) < 2 {
			continue
		}

		filename := match[1]
		downloadURL := baseURL + filename

		versionMatch := versionRegex.FindStringSubmatch(filename)
		if len(versionMatch) < 2 {
			continue
		}

		version := versionMatch[1]
		arch := store.ArchAMD64

		if strings.Contains(filename, "i386") {
			arch = store.Arch386
		}

		fileSize := getFileSizeFromURL(ctx, f.client, downloadURL)

		checksum, err := f.fetchChecksum(ctx, baseURL, filename)
		if err != nil {
			checksum = ""
		}

		pv := &store.ProductVersion{
			Version:      version,
			Platform:     store.PlatformLinux,
			Architecture: arch,
			DownloadURL:  downloadURL,
			Checksum:     checksum,
			ChecksumType: "sha256",
			FileSize:     fileSize,
			Filename:     filename,
			IsLatest:     false,
		}

		versions = append(versions, pv)
	}

	return versions, nil
}

func (f *DebianFetcher) fetchChecksum(ctx context.Context, baseURL, filename string) (string, error) {
	checksumURL := baseURL + "SHA256SUMS"
	resp, err := f.client.Get(ctx, checksumURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 && strings.Contains(parts[1], filename) {
			return parts[0], nil
		}
	}

	return "", fmt.Errorf("checksum not found for %s", filename)
}
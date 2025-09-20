package sources

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/your-username/alldownloads/internal/store"
)

type UbuntuFetcher struct {
	client *HTTPClient
}

func NewUbuntuFetcher() *UbuntuFetcher {
	return &UbuntuFetcher{
		client: NewHTTPClient(),
	}
}

func (f *UbuntuFetcher) Fetch(ctx context.Context) ([]*store.ProductVersion, error) {
	var versions []*store.ProductVersion

	releasesURL := "http://releases.ubuntu.com/"
	resp, err := f.client.Get(ctx, releasesURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Ubuntu releases: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	html := string(body)

	versionRegex := regexp.MustCompile(`href="([0-9]+\.[0-9]+(?:\.[0-9]+)?)/?"`)
	versionMatches := versionRegex.FindAllStringSubmatch(html, -1)

	for _, match := range versionMatches {
		if len(match) < 2 {
			continue
		}

		version := match[1]
		if shouldSkipUbuntuVersion(version) {
			continue
		}

		versionVersions, err := f.fetchVersionDetails(ctx, version)
		if err != nil {
			continue
		}

		versions = append(versions, versionVersions...)
	}

	return versions, nil
}

func (f *UbuntuFetcher) fetchVersionDetails(ctx context.Context, version string) ([]*store.ProductVersion, error) {
	var versions []*store.ProductVersion

	versionURL := fmt.Sprintf("http://releases.ubuntu.com/%s/", version)
	resp, err := f.client.Get(ctx, versionURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	html := string(body)

	isoRegex := regexp.MustCompile(`href="(ubuntu-[^"]*\.iso)"`)
	isoMatches := isoRegex.FindAllStringSubmatch(html, -1)

	for _, match := range isoMatches {
		if len(match) < 2 {
			continue
		}

		filename := match[1]
		downloadURL := versionURL + filename

		platform := store.PlatformLinux
		arch := store.ArchAMD64

		if strings.Contains(filename, "desktop-amd64") {
			arch = store.ArchAMD64
		} else if strings.Contains(filename, "desktop-i386") {
			arch = store.Arch386
		} else if strings.Contains(filename, "server-amd64") {
			arch = store.ArchAMD64
		} else if strings.Contains(filename, "live-server-amd64") {
			arch = store.ArchAMD64
		}

		fileSize := getFileSizeFromURL(ctx, f.client, downloadURL)

		checksum, err := f.fetchChecksum(ctx, versionURL, filename)
		if err != nil {
			checksum = ""
		}

		pv := &store.ProductVersion{
			Version:      version,
			Platform:     platform,
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

func (f *UbuntuFetcher) fetchChecksum(ctx context.Context, baseURL, filename string) (string, error) {
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

func shouldSkipUbuntuVersion(version string) bool {
	skipVersions := []string{"14.04", "16.04", "18.04", "19.04", "19.10", "21.04", "21.10"}
	for _, skip := range skipVersions {
		if version == skip {
			return true
		}
	}
	return false
}
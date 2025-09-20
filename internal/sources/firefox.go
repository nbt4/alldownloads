package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/your-username/alldownloads/internal/store"
)

type FirefoxFetcher struct {
	client *HTTPClient
}

type FirefoxRelease struct {
	Products struct {
		Firefox struct {
			Releases map[string]FirefoxReleaseDetails `json:"releases"`
		} `json:"firefox"`
	} `json:"products"`
}

type FirefoxReleaseDetails struct {
	Version string                 `json:"version"`
	Files   []FirefoxReleaseFile   `json:"files"`
}

type FirefoxReleaseFile struct {
	OS       string `json:"os"`
	Language string `json:"language"`
	URL      string `json:"url"`
}

func NewFirefoxFetcher() *FirefoxFetcher {
	return &FirefoxFetcher{
		client: NewHTTPClient(),
	}
}

func (f *FirefoxFetcher) Fetch(ctx context.Context) ([]*store.ProductVersion, error) {
	var versions []*store.ProductVersion

	apiURL := "https://product-details.mozilla.org/1.0/firefox_versions.json"
	resp, err := f.client.GetJSON(ctx, apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Firefox versions: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var firefoxVersions map[string]string
	if err := json.Unmarshal(body, &firefoxVersions); err != nil {
		return nil, fmt.Errorf("failed to parse Firefox versions: %w", err)
	}

	latestVersion, exists := firefoxVersions["LATEST_FIREFOX_VERSION"]
	if !exists {
		return nil, fmt.Errorf("latest Firefox version not found")
	}

	downloadURLs := map[string]map[string]string{
		store.PlatformWindows: {
			store.ArchAMD64: fmt.Sprintf("https://download.mozilla.org/?product=firefox-latest&os=win64&lang=en-US"),
			store.Arch386:   fmt.Sprintf("https://download.mozilla.org/?product=firefox-latest&os=win&lang=en-US"),
		},
		store.PlatformMacOS: {
			store.ArchAMD64: fmt.Sprintf("https://download.mozilla.org/?product=firefox-latest&os=osx&lang=en-US"),
		},
		store.PlatformLinux: {
			store.ArchAMD64: fmt.Sprintf("https://download.mozilla.org/?product=firefox-latest&os=linux64&lang=en-US"),
		},
	}

	for platform, archMap := range downloadURLs {
		for arch, downloadURL := range archMap {
			fileSize := getFileSizeFromURL(ctx, f.client, downloadURL)

			filename := "firefox-installer"
			switch platform {
			case store.PlatformWindows:
				filename = "Firefox Setup.exe"
			case store.PlatformMacOS:
				filename = "Firefox.dmg"
			case store.PlatformLinux:
				filename = "firefox.tar.bz2"
			}

			pv := &store.ProductVersion{
				Version:      latestVersion,
				Platform:     platform,
				Architecture: arch,
				DownloadURL:  downloadURL,
				Checksum:     "",
				ChecksumType: "",
				FileSize:     fileSize,
				Filename:     filename,
				IsLatest:     true,
			}

			versions = append(versions, pv)
		}
	}

	return versions, nil
}
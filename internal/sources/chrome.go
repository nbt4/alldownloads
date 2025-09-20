package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/your-username/alldownloads/internal/store"
)

type ChromeFetcher struct {
	client *HTTPClient
}

type ChromeVersionAPI struct {
	Versions []ChromeVersion `json:"versions"`
}

type ChromeVersion struct {
	Version string `json:"version"`
}

func NewChromeFetcher() *ChromeFetcher {
	return &ChromeFetcher{
		client: NewHTTPClient(),
	}
}

func (f *ChromeFetcher) Fetch(ctx context.Context) ([]*store.ProductVersion, error) {
	var versions []*store.ProductVersion

	apiURL := "https://versionhistory.googleapis.com/v1/chrome/platforms/win/channels/stable/versions"
	resp, err := f.client.GetJSON(ctx, apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Chrome versions: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var chromeAPI ChromeVersionAPI
	if err := json.Unmarshal(body, &chromeAPI); err != nil {
		return nil, fmt.Errorf("failed to parse Chrome API response: %w", err)
	}

	if len(chromeAPI.Versions) == 0 {
		return nil, fmt.Errorf("no Chrome versions found")
	}

	latestVersion := chromeAPI.Versions[0].Version

	downloadURLs := map[string]map[string]string{
		store.PlatformWindows: {
			store.ArchAMD64: fmt.Sprintf("https://dl.google.com/chrome/install/googlechromestandaloneenterprise64.msi"),
			store.Arch386:   fmt.Sprintf("https://dl.google.com/chrome/install/googlechromestandaloneenterprise.msi"),
		},
		store.PlatformMacOS: {
			store.ArchAMD64: fmt.Sprintf("https://dl.google.com/chrome/mac/stable/GGRO/googlechrome.dmg"),
		},
		store.PlatformLinux: {
			store.ArchAMD64: fmt.Sprintf("https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb"),
		},
	}

	for platform, archMap := range downloadURLs {
		for arch, downloadURL := range archMap {
			fileSize := getFileSizeFromURL(ctx, f.client, downloadURL)
			filename := extractFilename(downloadURL)

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
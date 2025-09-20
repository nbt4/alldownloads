package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/your-username/alldownloads/internal/store"
)

type VSCodeFetcher struct {
	client *HTTPClient
}

type VSCodeUpdate struct {
	URL     string `json:"url"`
	Version string `json:"name"`
}

func NewVSCodeFetcher() *VSCodeFetcher {
	return &VSCodeFetcher{
		client: NewHTTPClient(),
	}
}

func (f *VSCodeFetcher) Fetch(ctx context.Context) ([]*store.ProductVersion, error) {
	var versions []*store.ProductVersion

	apiURL := "https://update.code.visualstudio.com/api/update/win32-x64-user/stable/latest"
	resp, err := f.client.GetJSON(ctx, apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch VS Code update info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var update VSCodeUpdate
	if err := json.Unmarshal(body, &update); err != nil {
		return nil, fmt.Errorf("failed to parse VS Code update: %w", err)
	}

	downloadURLs := map[string]map[string]string{
		store.PlatformWindows: {
			store.ArchAMD64: fmt.Sprintf("https://update.code.visualstudio.com/latest/win32-x64-user/stable"),
			store.Arch386:   fmt.Sprintf("https://update.code.visualstudio.com/latest/win32-user/stable"),
		},
		store.PlatformMacOS: {
			store.ArchAMD64: fmt.Sprintf("https://update.code.visualstudio.com/latest/darwin/stable"),
		},
		store.PlatformLinux: {
			store.ArchAMD64: fmt.Sprintf("https://update.code.visualstudio.com/latest/linux-x64/stable"),
		},
	}

	for platform, archMap := range downloadURLs {
		for arch, downloadURL := range archMap {
			fileSize := getFileSizeFromURL(ctx, f.client, downloadURL)

			filename := "vscode-installer"
			switch platform {
			case store.PlatformWindows:
				filename = "VSCodeUserSetup.exe"
			case store.PlatformMacOS:
				filename = "VSCode-darwin.zip"
			case store.PlatformLinux:
				filename = "code.tar.gz"
			}

			pv := &store.ProductVersion{
				Version:      update.Version,
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
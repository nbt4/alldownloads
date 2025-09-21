package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/your-username/alldownloads/internal/store"
)

type KaliFetcher struct {
	client *HTTPClient
}

func NewKaliFetcher() *KaliFetcher {
	return &KaliFetcher{client: NewHTTPClient()}
}

func (f *KaliFetcher) Fetch(ctx context.Context) ([]*store.ProductVersion, error) {
	var versions []*store.ProductVersion

	// Kali Linux ISO download URLs
	downloadURLs := map[string]map[string]string{
		store.PlatformLinux: {
			store.ArchAMD64: "https://cdimage.kali.org/kali-2025.2/kali-linux-2025.2-installer-amd64.iso",
		},
	}

	version := "2025.2"

	for platform, archMap := range downloadURLs {
		for arch, downloadURL := range archMap {
			fileSize := getFileSizeFromURL(ctx, f.client, downloadURL)

			filename := extractFilename(downloadURL)

			pv := &store.ProductVersion{
				Version:      version,
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

type WindowsFetcher struct {
	client *HTTPClient
}

func NewWindowsFetcher() *WindowsFetcher {
	return &WindowsFetcher{client: NewHTTPClient()}
}

func (f *WindowsFetcher) Fetch(ctx context.Context) ([]*store.ProductVersion, error) {
	var versions []*store.ProductVersion

	// Windows 11 ISO download URLs (Microsoft Media Creation Tool)
	downloadURLs := map[string]map[string]string{
		store.PlatformWindows: {
			store.ArchAMD64: "https://go.microsoft.com/fwlink/?LinkId=691209",
		},
	}

	version := "Windows 11"

	for platform, archMap := range downloadURLs {
		for arch, downloadURL := range archMap {
			fileSize := getFileSizeFromURL(ctx, f.client, downloadURL)

			filename := "MediaCreationToolW11.exe"

			pv := &store.ProductVersion{
				Version:      version,
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

type TermiusFetcher struct {
	client *HTTPClient
}

func NewTermiusFetcher() *TermiusFetcher {
	return &TermiusFetcher{client: NewHTTPClient()}
}

func (f *TermiusFetcher) Fetch(ctx context.Context) ([]*store.ProductVersion, error) {
	var versions []*store.ProductVersion

	// Termius direct download URLs (actual download links)
	downloadURLs := map[string]map[string]string{
		store.PlatformWindows: {
			store.ArchAMD64: "https://autoupdate.termius.com/windows/Termius.exe",
		},
		store.PlatformMacOS: {
			store.ArchAMD64: "https://autoupdate.termius.com/mac/Termius.dmg",
		},
		store.PlatformLinux: {
			store.ArchAMD64: "https://autoupdate.termius.com/linux/Termius.AppImage",
		},
	}

	// Since Termius doesn't have a public API, we'll use a generic version
	version := "latest"

	for platform, archMap := range downloadURLs {
		for arch, downloadURL := range archMap {
			fileSize := getFileSizeFromURL(ctx, f.client, downloadURL)

			filename := "termius-installer"
			switch platform {
			case store.PlatformWindows:
				filename = "Termius.exe"
			case store.PlatformMacOS:
				filename = "Termius.dmg"
			case store.PlatformLinux:
				filename = "Termius.AppImage"
			}

			pv := &store.ProductVersion{
				Version:      version,
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

type TelegramFetcher struct {
	client *HTTPClient
}

func NewTelegramFetcher() *TelegramFetcher {
	return &TelegramFetcher{client: NewHTTPClient()}
}

func (f *TelegramFetcher) Fetch(ctx context.Context) ([]*store.ProductVersion, error) {
	var versions []*store.ProductVersion

	// Get latest release from GitHub API
	apiURL := "https://api.github.com/repos/telegramdesktop/tdesktop/releases/latest"
	resp, err := f.client.GetJSON(ctx, apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Telegram release info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
			Size               int64  `json:"size"`
		} `json:"assets"`
	}
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("failed to parse Telegram release: %w", err)
	}

	// Find matching assets with simple, clear patterns
	for _, asset := range release.Assets {
		var platform, arch string
		matched := false

		// Windows x64
		if asset.Name == "tsetup-x64."+strings.TrimPrefix(release.TagName, "v")+".exe" {
			platform = store.PlatformWindows
			arch = store.ArchAMD64
			matched = true
		}
		// Windows x86 (regular tsetup.exe)
		if asset.Name == "tsetup."+strings.TrimPrefix(release.TagName, "v")+".exe" {
			platform = store.PlatformWindows
			arch = store.Arch386
			matched = true
		}
		// macOS
		if asset.Name == "tsetup."+strings.TrimPrefix(release.TagName, "v")+".dmg" {
			platform = store.PlatformMacOS
			arch = store.ArchAMD64
			matched = true
		}
		// Linux
		if asset.Name == "tsetup."+strings.TrimPrefix(release.TagName, "v")+".tar.xz" {
			platform = store.PlatformLinux
			arch = store.ArchAMD64
			matched = true
		}

		if matched {
			pv := &store.ProductVersion{
				Version:      strings.TrimPrefix(release.TagName, "v"),
				Platform:     platform,
				Architecture: arch,
				DownloadURL:  asset.BrowserDownloadURL,
				Checksum:     "",
				ChecksumType: "",
				FileSize:     asset.Size,
				Filename:     asset.Name,
				IsLatest:     true,
			}
			versions = append(versions, pv)
		}
	}

	return versions, nil
}

type WhatsAppFetcher struct {
	client *HTTPClient
}

func NewWhatsAppFetcher() *WhatsAppFetcher {
	return &WhatsAppFetcher{client: NewHTTPClient()}
}

func (f *WhatsAppFetcher) Fetch(ctx context.Context) ([]*store.ProductVersion, error) {
	var versions []*store.ProductVersion

	// WhatsApp Desktop direct download URLs
	downloadURLs := map[string]map[string]string{
		store.PlatformWindows: {
			store.ArchAMD64: "https://apps.microsoft.com/detail/whatsapp/9NKSQGP7F2NH",
		},
		store.PlatformMacOS: {
			store.ArchAMD64: "https://web.whatsapp.com/desktop/mac_native/release",
		},
	}

	version := "latest"

	for platform, archMap := range downloadURLs {
		for arch, downloadURL := range archMap {
			fileSize := getFileSizeFromURL(ctx, f.client, downloadURL)

			filename := "whatsapp-installer"
			switch platform {
			case store.PlatformWindows:
				filename = "WhatsApp.msix"
			case store.PlatformMacOS:
				filename = "WhatsApp.dmg"
			}

			pv := &store.ProductVersion{
				Version:      version,
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

type TailscaleFetcher struct {
	client *HTTPClient
}

func NewTailscaleFetcher() *TailscaleFetcher {
	return &TailscaleFetcher{client: NewHTTPClient()}
}

func (f *TailscaleFetcher) Fetch(ctx context.Context) ([]*store.ProductVersion, error) {
	var versions []*store.ProductVersion

	// Tailscale direct download URLs
	downloadURLs := map[string]map[string]string{
		store.PlatformWindows: {
			store.ArchAMD64: "https://pkgs.tailscale.com/stable/tailscale-setup-latest.exe",
		},
		store.PlatformMacOS: {
			store.ArchAMD64: "https://pkgs.tailscale.com/stable/Tailscale-latest-macos.pkg",
		},
		store.PlatformLinux: {
			store.ArchAMD64: "https://pkgs.tailscale.com/stable/tailscale_latest_amd64.deb",
		},
	}

	version := "latest"

	for platform, archMap := range downloadURLs {
		for arch, downloadURL := range archMap {
			fileSize := getFileSizeFromURL(ctx, f.client, downloadURL)

			filename := "tailscale-installer"
			switch platform {
			case store.PlatformWindows:
				filename = "tailscale-setup-latest.exe"
			case store.PlatformMacOS:
				filename = "Tailscale-latest-macos.pkg"
			case store.PlatformLinux:
				filename = "tailscale_latest_amd64.deb"
			}

			pv := &store.ProductVersion{
				Version:      version,
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

type NextcloudFetcher struct {
	client *HTTPClient
}

func NewNextcloudFetcher() *NextcloudFetcher {
	return &NextcloudFetcher{client: NewHTTPClient()}
}

func (f *NextcloudFetcher) Fetch(ctx context.Context) ([]*store.ProductVersion, error) {
	var versions []*store.ProductVersion

	// Get latest release from GitHub API
	apiURL := "https://api.github.com/repos/nextcloud-releases/desktop/releases/latest"
	resp, err := f.client.GetJSON(ctx, apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Nextcloud release info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
			Size               int64  `json:"size"`
		} `json:"assets"`
	}
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("failed to parse Nextcloud release: %w", err)
	}

	// Map of platform/arch to asset name patterns
	assetPatterns := map[string]map[string]string{
		store.PlatformWindows: {
			store.ArchAMD64: "-x64.msi",
		},
		store.PlatformMacOS: {
			store.ArchAMD64: ".pkg",
		},
		store.PlatformLinux: {
			store.ArchAMD64: "-x86_64.AppImage",
		},
	}

	// Find matching assets
	for platform, archMap := range assetPatterns {
		for arch, pattern := range archMap {
			for _, asset := range release.Assets {
				matched := strings.Contains(asset.Name, pattern) && !strings.Contains(asset.Name, ".asc") && !strings.Contains(asset.Name, ".tbz")

				// Additional filtering for platform-specific files
				if matched {
					if platform == store.PlatformMacOS && strings.Contains(asset.Name, "vfs") {
						// Skip VFS version for now, use regular version
						continue
					}
				}

				if matched {
					pv := &store.ProductVersion{
						Version:      strings.TrimPrefix(release.TagName, "v"),
						Platform:     platform,
						Architecture: arch,
						DownloadURL:  asset.BrowserDownloadURL,
						Checksum:     "",
						ChecksumType: "",
						FileSize:     asset.Size,
						Filename:     asset.Name,
						IsLatest:     true,
					}
					versions = append(versions, pv)
					break
				}
			}
		}
	}

	return versions, nil
}

type BraveFetcher struct {
	client *HTTPClient
}

func NewBraveFetcher() *BraveFetcher {
	return &BraveFetcher{client: NewHTTPClient()}
}

func (f *BraveFetcher) Fetch(ctx context.Context) ([]*store.ProductVersion, error) {
	var versions []*store.ProductVersion

	// Get latest release from GitHub API
	apiURL := "https://api.github.com/repos/brave/brave-browser/releases/latest"
	resp, err := f.client.GetJSON(ctx, apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Brave release info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
			Size               int64  `json:"size"`
		} `json:"assets"`
	}
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("failed to parse Brave release: %w", err)
	}

	// Map of platform/arch to asset name patterns
	assetPatterns := map[string]map[string]string{
		store.PlatformWindows: {
			store.ArchAMD64: "BraveBrowserSetup.exe",
			store.Arch386:   "brave-v.*-win32-ia32.zip",
		},
		store.PlatformMacOS: {
			store.ArchAMD64: "Brave-Browser-universal.dmg",
		},
		store.PlatformLinux: {
			store.ArchAMD64: "brave-browser_.*_amd64.deb",
		},
	}

	// Find matching assets
	for platform, archMap := range assetPatterns {
		for arch, pattern := range archMap {
			for _, asset := range release.Assets {
				matched := false
				if pattern == "BraveBrowserSetup.exe" {
					matched = asset.Name == "BraveBrowserSetup.exe"
				} else {
					matched = strings.Contains(asset.Name, strings.Split(pattern, ".*")[0])
				}

				if matched && !strings.Contains(asset.Name, "sha256") && !strings.Contains(asset.Name, "asc") {
					pv := &store.ProductVersion{
						Version:      strings.TrimPrefix(release.TagName, "v"),
						Platform:     platform,
						Architecture: arch,
						DownloadURL:  asset.BrowserDownloadURL,
						Checksum:     "",
						ChecksumType: "",
						FileSize:     asset.Size,
						Filename:     asset.Name,
						IsLatest:     true,
					}
					versions = append(versions, pv)
					break
				}
			}
		}
	}

	return versions, nil
}

type NotepadPlusPlusFetcher struct {
	client *HTTPClient
}

func NewNotepadPlusPlusFetcher() *NotepadPlusPlusFetcher {
	return &NotepadPlusPlusFetcher{client: NewHTTPClient()}
}

func (f *NotepadPlusPlusFetcher) Fetch(ctx context.Context) ([]*store.ProductVersion, error) {
	var versions []*store.ProductVersion

	// Get latest release from GitHub API
	apiURL := "https://api.github.com/repos/notepad-plus-plus/notepad-plus-plus/releases/latest"
	resp, err := f.client.GetJSON(ctx, apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Notepad++ release info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
			Size               int64  `json:"size"`
		} `json:"assets"`
	}
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("failed to parse Notepad++ release: %w", err)
	}

	// Find installer assets (avoid signatures and checksums)
	for _, asset := range release.Assets {
		var platform, arch string
		matched := false

		// Windows x64 installer
		if strings.Contains(asset.Name, "Installer.x64.exe") && !strings.Contains(asset.Name, ".sig") {
			platform = store.PlatformWindows
			arch = store.ArchAMD64
			matched = true
		}
		// Windows x86 installer (just "Installer.exe" without x64)
		if strings.Contains(asset.Name, "Installer.exe") && !strings.Contains(asset.Name, "x64") && !strings.Contains(asset.Name, "arm64") && !strings.Contains(asset.Name, ".sig") {
			platform = store.PlatformWindows
			arch = store.Arch386
			matched = true
		}

		if matched {
			pv := &store.ProductVersion{
				Version:      strings.TrimPrefix(release.TagName, "v"),
				Platform:     platform,
				Architecture: arch,
				DownloadURL:  asset.BrowserDownloadURL,
				Checksum:     "",
				ChecksumType: "",
				FileSize:     asset.Size,
				Filename:     asset.Name,
				IsLatest:     true,
			}
			versions = append(versions, pv)
		}
	}

	return versions, nil
}

type PowerShellFetcher struct {
	client *HTTPClient
}

func NewPowerShellFetcher() *PowerShellFetcher {
	return &PowerShellFetcher{client: NewHTTPClient()}
}

func (f *PowerShellFetcher) Fetch(ctx context.Context) ([]*store.ProductVersion, error) {
	var versions []*store.ProductVersion

	// Get latest release from GitHub API
	apiURL := "https://api.github.com/repos/PowerShell/PowerShell/releases/latest"
	resp, err := f.client.GetJSON(ctx, apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PowerShell release info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
			Size               int64  `json:"size"`
		} `json:"assets"`
	}
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("failed to parse PowerShell release: %w", err)
	}

	// Find installer packages
	for _, asset := range release.Assets {
		var platform, arch string
		matched := false

		// Windows MSI installers
		if strings.Contains(asset.Name, "PowerShell-") && strings.HasSuffix(asset.Name, ".msi") {
			platform = store.PlatformWindows
			if strings.Contains(asset.Name, "win-x64") {
				arch = store.ArchAMD64
				matched = true
			} else if strings.Contains(asset.Name, "win-x86") {
				arch = store.Arch386
				matched = true
			}
		}
		// macOS PKG installers
		if strings.Contains(asset.Name, "powershell-") && strings.Contains(asset.Name, "osx") && strings.HasSuffix(asset.Name, ".pkg") {
			platform = store.PlatformMacOS
			arch = store.ArchAMD64
			matched = true
		}
		// Linux DEB packages
		if strings.Contains(asset.Name, "powershell_") && strings.Contains(asset.Name, "deb_amd64") && strings.HasSuffix(asset.Name, ".deb") {
			platform = store.PlatformLinux
			arch = store.ArchAMD64
			matched = true
		}

		if matched {
			pv := &store.ProductVersion{
				Version:      strings.TrimPrefix(release.TagName, "v"),
				Platform:     platform,
				Architecture: arch,
				DownloadURL:  asset.BrowserDownloadURL,
				Checksum:     "",
				ChecksumType: "",
				FileSize:     asset.Size,
				Filename:     asset.Name,
				IsLatest:     true,
			}
			versions = append(versions, pv)
		}
	}

	return versions, nil
}

type OfficeFetcher struct {
	client *HTTPClient
}

func NewOfficeFetcher() *OfficeFetcher {
	return &OfficeFetcher{client: NewHTTPClient()}
}

func (f *OfficeFetcher) Fetch(ctx context.Context) ([]*store.ProductVersion, error) {
	var versions []*store.ProductVersion

	// MS Office download URLs (Microsoft 365 Apps)
	downloadURLs := map[string]map[string]string{
		store.PlatformWindows: {
			store.ArchAMD64: "https://c2rsetup.officeapps.live.com/c2r/download.aspx?ProductreleaseID=O365ProPlusRetail&platform=x64&language=en-us",
			store.Arch386:   "https://c2rsetup.officeapps.live.com/c2r/download.aspx?ProductreleaseID=O365ProPlusRetail&platform=x86&language=en-us",
		},
		store.PlatformMacOS: {
			store.ArchAMD64: "https://go.microsoft.com/fwlink/?linkid=525133",
		},
	}

	// Office uses a rolling release model
	version := "Microsoft 365"

	for platform, archMap := range downloadURLs {
		for arch, downloadURL := range archMap {
			fileSize := getFileSizeFromURL(ctx, f.client, downloadURL)

			filename := "office-installer"
			switch platform {
			case store.PlatformWindows:
				if arch == store.ArchAMD64 {
					filename = "OfficeSetup-x64.exe"
				} else {
					filename = "OfficeSetup-x86.exe"
				}
			case store.PlatformMacOS:
				filename = "Microsoft_Office.pkg"
			}

			pv := &store.ProductVersion{
				Version:      version,
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
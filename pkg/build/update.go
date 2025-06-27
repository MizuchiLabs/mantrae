package build

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	RepoURL = "https://api.github.com/repos/mizuchilabs/mantrae"
)

type releaseAsset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
	ID          int    `json:"id"`
	Size        int    `json:"size"`
}

type release struct {
	Name      string          `json:"name"`
	Tag       string          `json:"tag_name"`
	Published string          `json:"published_at"`
	URL       string          `json:"html_url"`
	Body      string          `json:"body"`
	Assets    []*releaseAsset `json:"assets"`
	ID        int             `json:"id"`
}

func Update(update bool) {
	if runningInDocker() {
		slog.Info("Running in docker, skipping update")
		return
	}

	latest, err := fetchLatestRelease()
	if err != nil {
		slog.Error("Update failed", "Error", err)
		return
	}

	if !update {
		if compareVersions(
			strings.TrimPrefix(Version, "v"),
			strings.TrimPrefix(latest.Tag, "v"),
		) <= 0 {
			slog.Info("You are running the latest version!")
			return
		}
		slog.Info("New version available!", "latest", latest.Tag, "current", Version)
		return
	}

	asset := latest.findBinary(filepath.Base(os.Args[0]))
	if asset == nil {
		slog.Info("Unsupported platform", "platform", runtime.GOOS+"/"+runtime.GOARCH)
		return
	}
	exec, err := os.Executable()
	if err != nil {
		slog.Error("Update failed", "Error", err)
		return
	}
	if err := os.Remove(exec); err != nil {
		slog.Error("Failed to remove current executable", "Error", err)
		return
	}

	slog.Info("Downloading...", "release", latest.Tag, "binary", asset.Name)
	if err := downloadFile(asset.DownloadURL, exec); err != nil {
		slog.Error("Failed to download", "Error", err)
		return
	}

	slog.Info("Update success!")
}

func runningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}

func fetchLatestRelease() (*release, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s/releases/latest", RepoURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("(%d) failed to send latest release request", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := &release{}
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}

	return result, nil
}

func downloadFile(url string, dest string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("(%d) failed to send download file request", res.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() {
		if err = out.Close(); err != nil {
			slog.Error("failed to close output file", "error", err)
		}
	}()

	if _, err := io.Copy(out, res.Body); err != nil {
		return err
	}

	if err := out.Chmod(0755); err != nil {
		return err
	}

	return nil
}

func (r *release) findBinary(name string) *releaseAsset {
	var assetName string

	switch runtime.GOOS {
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			assetName = name + "_linux_amd64"
		case "arm64":
			assetName = name + "_linux_arm64"
		case "arm":
			assetName = name + "_linux_armv7"
		}
	case "darwin":
		switch runtime.GOARCH {
		case "amd64":
			assetName = name + "_darwin_amd64"
		case "arm64":
			assetName = name + "_darwin_arm64"
		}
	case "windows":
		switch runtime.GOARCH {
		case "amd64":
			assetName = name + "_windows_amd64"
		case "arm64":
			assetName = name + "_windows_arm64"
		}
	}

	for _, asset := range r.Assets {
		if assetName == asset.Name {
			return asset
		}
	}

	return nil
}

func compareVersions(a, b string) int {
	aSplit := strings.Split(a, ".")
	aTotal := len(aSplit)

	bSplit := strings.Split(b, ".")
	bTotal := len(bSplit)

	limit := max(aTotal, bTotal)

	for i := range limit {
		var x, y int

		if i < aTotal {
			x, _ = strconv.Atoi(aSplit[i])
		}

		if i < bTotal {
			y, _ = strconv.Atoi(bSplit[i])
		}

		if x < y {
			return 1 // b is newer
		}

		if x > y {
			return -1 // a is newer
		}
	}

	return 0 // equal
}

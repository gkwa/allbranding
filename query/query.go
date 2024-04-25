package query

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/hashicorp/go-version"
)

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	BrowserDownloadURL string `json:"browser_download_url"`
}

func Run(releasesURL, assetRegex string, noCache, parseHarder bool, ignoreRegex []string) {
	var releasesData []byte
	var err error

	cacheDir := filepath.Join(os.TempDir(), "allbranding")
	cacheFileName := fmt.Sprintf("releases_%x.json", sha256.Sum256([]byte(releasesURL)))
	cacheFile := filepath.Join(cacheDir, cacheFileName)

	slog.Debug("cache path", "path", cacheFile)

	if !noCache {
		if _, err := os.Stat(cacheFile); err == nil {
			fileInfo, err := os.Stat(cacheFile)
			if err != nil {
				slog.Error("failed to get file info", "error", err)
				return
			}

			if time.Since(fileInfo.ModTime()) < 1*time.Hour {
				releasesData, err = os.ReadFile(cacheFile)
				if err != nil {
					slog.Error("failed to read cache file", "error", err)
					return
				}
			}
		}
	}

	if releasesData == nil {
		resp, err := http.Get(releasesURL)
		if err != nil {
			slog.Error("failed to fetch releases", "error", err)
			return
		}
		defer resp.Body.Close()

		releasesData, err = io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("failed to read response body", "error", err)
			return
		}

		if !noCache {
			err = os.MkdirAll(cacheDir, os.ModePerm)
			if err != nil {
				slog.Error("failed to create cache directory", "error", err)
				return
			}

			err = os.WriteFile(cacheFile, releasesData, 0o644)
			if err != nil {
				slog.Error("failed to write cache file", "error", err)
				return
			}
		}
	}

	var releases []Release
	err = json.NewDecoder(bytes.NewReader(releasesData)).Decode(&releases)
	if err != nil {
		slog.Error("failed to decode releases JSON", "error", err)
		return
	}

	ignorePatterns := make([]*regexp.Regexp, len(ignoreRegex))
	for i, pattern := range ignoreRegex {
		ignorePatterns[i] = regexp.MustCompile(pattern)
	}

	filteredReleases := make([]Release, 0, len(releases))
	for _, release := range releases {
		ignore := false
		for _, pattern := range ignorePatterns {
			if pattern.MatchString(release.TagName) {
				ignore = true
				break
			}
		}
		if !ignore {
			filteredReleases = append(filteredReleases, release)
		}
	}

	sort.Slice(filteredReleases, func(i, j int) bool {
		vi, err := version.NewVersion(parseVersion(filteredReleases[i].TagName, parseHarder))
		if err != nil {
			slog.Warn("invalid version", "version", filteredReleases[i].TagName, "error", err)
			return false
		}
		vj, err := version.NewVersion(parseVersion(filteredReleases[j].TagName, parseHarder))
		if err != nil {
			slog.Warn("invalid version", "version", filteredReleases[j].TagName, "error", err)
			return false
		}
		slog.Debug("comparing versions", "version1", vi.String(), "version2", vj.String())
		return vi.GreaterThan(vj)
	})

	re, err := regexp.Compile(assetRegex)
	if err != nil {
		slog.Error("invalid asset regex", "regex", assetRegex, "error", err)
		return
	}

	var version, assetURL string
	for _, release := range filteredReleases {
		for _, asset := range release.Assets {
			if re.MatchString(filepath.Base(asset.BrowserDownloadURL)) {
				version = release.TagName
				assetURL = asset.BrowserDownloadURL
				break
			}
		}
		if version != "" {
			break
		}
	}

	result := map[string]string{
		"version":              version,
		"browser_download_url": assetURL,
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		slog.Error("failed to marshal result JSON", "error", err)
		return
	}

	fmt.Println(string(resultJSON))
}

func parseVersion(version string, parseHarder bool) string {
	if parseHarder {
		re := regexp.MustCompile(`[^0-9.]+`)
		return re.ReplaceAllString(version, "")
	}
	return version
}

package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type releaseInfo struct {
	Version       string
	DownloadURL   string
	AssetName     string
	VscodeZipURL  string
}

type ghRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func fetchLatestRelease() (*releaseInfo, error) {
	resp, err := http.Get("https://api.github.com/repos/hannogeo/emerald/releases/latest")
	if err != nil {
		return nil, fmt.Errorf("failed to reach GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var release ghRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("failed to parse release info: %w", err)
	}

	tag := strings.TrimPrefix(release.TagName, "v")

	info := &releaseInfo{Version: tag}

	for _, asset := range release.Assets {
		switch asset.Name {
		case "emerald.exe":
			info.DownloadURL = asset.BrowserDownloadURL
			info.AssetName = asset.Name
		case "emerald-vscode.zip":
			info.VscodeZipURL = asset.BrowserDownloadURL
		}
	}

	if info.DownloadURL == "" {
		return nil, fmt.Errorf("no emerald.exe asset found in latest release")
	}

	return info, nil
}

func compareVersions(a, b string) int {
	sa := strings.Split(a, ".")
	sb := strings.Split(b, ".")
	maxLen := len(sa)
	if len(sb) > maxLen {
		maxLen = len(sb)
	}
	for i := 0; i < maxLen; i++ {
		var na, nb int
		if i < len(sa) {
			na, _ = strconv.Atoi(sa[i])
		}
		if i < len(sb) {
			nb, _ = strconv.Atoi(sb[i])
		}
		if na < nb {
			return -1
		}
		if na > nb {
			return 1
		}
	}
	return 0
}

func downloadFile(url, dest string) error {
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	_, err = io.Copy(out, resp.Body)
	return err
}

func replaceBinary(currentPath, newPath string) error {
	dir := filepath.Dir(currentPath)
	scriptPath := filepath.Join(dir, "_emerald_update.bat")

	content := fmt.Sprintf(`@echo off
:wait
ping -n 2 127.0.0.1 > nul 2>&1
del "%s" > nul 2>&1
if exist "%s" goto wait
ren "%s" "emerald.exe" > nul 2>&1
if exist "%s" (
  del "%%~f0"
  start "" "emerald.exe" version
)
`, currentPath, currentPath, newPath, currentPath)

	os.WriteFile(scriptPath, []byte(content), 0755)

	cmd := exec.Command("cmd", "/c", "start", "/b", "", scriptPath)
	return cmd.Start()
}

func installVscodeExtension(zipURL string) {
	installDir := filepath.Join(os.Getenv("USERPROFILE"), ".emerald")
	extDir := filepath.Join(installDir, "vscode-emerald")
	zipPath := filepath.Join(installDir, "emerald-vscode.zip")

	if err := downloadFile(zipURL, zipPath); err != nil {
		fmt.Println("Warning: could not download VS Code extension.")
		return
	}

	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		fmt.Println("Warning: could not read VS Code extension package.")
		os.Remove(zipPath)
		return
	}
	defer zr.Close()

	os.RemoveAll(extDir)
	os.MkdirAll(extDir, 0755)

	for _, f := range zr.File {
		fpath := filepath.Join(extDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(fpath), 0755)
		r, err := f.Open()
		if err != nil {
			continue
		}
		out, err := os.Create(fpath)
		if err != nil {
			r.Close()
			continue
		}
		io.Copy(out, r)
		out.Close()
		r.Close()
	}

	os.Remove(zipPath)

	cmd := exec.Command("code", "--install-extension", extDir)
	if err := cmd.Run(); err != nil {
		fmt.Println("VS Code not found. Extension files are at:", extDir)
		fmt.Println("Install VS Code, then run: code --install-extension \"" + extDir + "\"")
		return
	}

	fmt.Println("VS Code extension installed.")
}

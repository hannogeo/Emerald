package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const repoURL = "https://api.github.com/repos/hannogeo/emerald/releases/latest"

type ghRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func main() {
	fmt.Println("Emerald Installer")
	fmt.Println("=================")

	installDir := filepath.Join(os.Getenv("USERPROFILE"), ".emerald")
	fmt.Printf("Installing to: %s\n", installDir)

	err := os.MkdirAll(installDir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Fetching latest release...")
	resp, err := http.Get(repoURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reaching GitHub: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "GitHub API returned status %d\n", resp.StatusCode)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response: %v\n", err)
		os.Exit(1)
	}

	var release ghRelease
	if err := json.Unmarshal(body, &release); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing release: %v\n", err)
		os.Exit(1)
	}

	var downloadURL, assetName string
	for _, asset := range release.Assets {
		if asset.Name == "emerald.exe" {
			downloadURL = asset.BrowserDownloadURL
			assetName = asset.Name
			break
		}
	}

	if downloadURL == "" {
		fmt.Fprintln(os.Stderr, "No .exe asset found in latest release")
		os.Exit(1)
	}

	fmt.Printf("Downloading %s...\n", assetName)

	exePath := filepath.Join(installDir, "emerald.exe")
	out, err := os.Create(exePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(1)
	}

	dlResp, err := http.Get(downloadURL)
	if err != nil {
		out.Close()
		os.Remove(exePath)
		fmt.Fprintf(os.Stderr, "Error downloading: %v\n", err)
		os.Exit(1)
	}
	defer dlResp.Body.Close()

	_, err = io.Copy(out, dlResp.Body)
	out.Close()
	if err != nil {
		os.Remove(exePath)
		fmt.Fprintf(os.Stderr, "Error saving file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Download complete!")

	path := os.Getenv("PATH")
	if !strings.Contains(path, installDir) {
		fmt.Println("Adding to PATH...")
		cmd := exec.Command("powershell", "-NoProfile", "-Command",
			fmt.Sprintf(`[Environment]::SetEnvironmentVariable("Path", [Environment]::GetEnvironmentVariable("Path", "User") + ";%s", "User")`, installDir))
		err = cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not add to PATH automatically.\n")
			fmt.Fprintf(os.Stderr, "Please add %s to your PATH manually.\n", installDir)
		} else {
			fmt.Println("Added to PATH. You may need to restart your terminal.")
		}
	} else {
		fmt.Println("Already in PATH.")
	}

	fmt.Printf("\nEmerald v%s installed successfully!\n", strings.TrimPrefix(release.TagName, "v"))
	fmt.Println("Run 'emerald version' to verify.")
}

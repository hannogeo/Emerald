package main

import (
	"fmt"
	"os"

	"emerald/interpreter"
	"emerald/lexer"
	"emerald/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Emerald v%s\n", Version)
		fmt.Fprintln(os.Stderr, "Usage: emerald <file.emld>")
		fmt.Fprintln(os.Stderr, "       emerald update")
		fmt.Fprintln(os.Stderr, "       emerald version")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "update":
		runUpdate()
	case "version":
		versionCmd()
	default:
		runFile(os.Args[1])
	}
}

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	l := lexer.NewLexer(string(content))
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, errMsg := range p.Errors() {
			fmt.Fprintln(os.Stderr, "Parse error:", errMsg)
		}
		os.Exit(1)
	}

	interp := interpreter.NewInterpreter()
	err = interp.Eval(program)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Runtime error:", err)
		os.Exit(1)
	}
}

func versionCmd() {
	fmt.Printf("Emerald v%s\n", Version)

	latest, err := fetchLatestRelease()
	if err != nil {
		return
	}

	if compareVersions(Version, latest.Version) < 0 {
		fmt.Printf("A newer version is available: v%s\n", latest.Version)
		fmt.Println("Run 'emerald update' to upgrade.")
	}
}

func runUpdate() {
	fmt.Printf("Emerald v%s - checking for updates...\n", Version)

	latest, err := fetchLatestRelease()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking for updates: %v\n", err)
		os.Exit(1)
	}

	cmp := compareVersions(Version, latest.Version)
	if cmp >= 0 {
		fmt.Println("You already have the latest version.")
		return
	}

	fmt.Printf("New version available: v%s\n", latest.Version)
	fmt.Printf("Downloading %s...\n", latest.AssetName)

	exePath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting executable path: %v\n", err)
		os.Exit(1)
	}

	err = downloadFile(latest.DownloadURL, exePath+".new")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error downloading update: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Download complete. Installing...")

	if latest.VscodeZipURL != "" {
		installVscodeExtension(latest.VscodeZipURL)
	}

	err = replaceBinary(exePath, exePath+".new")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error installing update: %v\n", err)
		fmt.Println("The downloaded file is at:", exePath+".new")
		os.Exit(1)
	}

	fmt.Println("Update complete! You are now running Emerald v" + latest.Version)
}

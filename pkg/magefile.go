//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/sh"
)

// Build builds the binary for the current platform
func Build() error {
	fmt.Println("Building grafana-plugin-api for current platform...")
	return sh.Run("go", "build", "-o", "grafana-plugin-api", "./cmd")
}

// BuildAll builds binaries for all supported platforms
func BuildAll() error {
	binDir := filepath.Join("..", "grafana-hover-tracker-panel", "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return err
	}

	platforms := []struct {
		goos   string
		goarch string
		output string
	}{
		{"linux", "amd64", "grafana-plugin-api-linux-amd64"},
		{"linux", "arm64", "grafana-plugin-api-linux-arm64"},
		{"darwin", "amd64", "grafana-plugin-api-darwin-amd64"},
		{"darwin", "arm64", "grafana-plugin-api-darwin-arm64"},
		{"windows", "amd64", "grafana-plugin-api-windows-amd64.exe"},
	}

	for _, platform := range platforms {
		fmt.Printf("Building for %s/%s...\n", platform.goos, platform.goarch)

		env := map[string]string{
			"GOOS":   platform.goos,
			"GOARCH": platform.goarch,
			"CGO_ENABLED": "0",
		}

		outputPath := filepath.Join(binDir, platform.output)

		if err := sh.RunWith(env, "go", "build", "-o", outputPath, "./cmd"); err != nil {
			fmt.Printf("Warning: Failed to build for %s/%s: %v\n", platform.goos, platform.goarch, err)
			continue
		}

		fmt.Printf("✓ Built %s\n", platform.output)
	}

	return nil
}

// Install copies the binary to the panel bin folder
func Install() error {
	binDir := filepath.Join("..", "grafana-hover-tracker-panel", "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return err
	}

	if err := Build(); err != nil {
		return err
	}

	src := "grafana-plugin-api"
	if runtime.GOOS == "windows" {
		src += ".exe"
	}
	dst := filepath.Join(binDir, src)

	fmt.Printf("Copying binary to %s...\n", dst)
	return sh.Copy(dst, src)
}

// Clean removes build artifacts
func Clean() error {
	fmt.Println("Cleaning build artifacts...")
	if err := sh.Run("go", "clean"); err != nil {
		return err
	}

	// Remove binary
	if err := os.Remove("grafana-plugin-api"); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// Tidy runs go mod tidy
func Tidy() error {
	fmt.Println("Running go mod tidy...")
	return sh.Run("go", "mod", "tidy")
}

// Test runs all tests
func Test() error {
	fmt.Println("Running tests...")
	return sh.Run("go", "test", "-v", "./...")
}

// TestCoverage runs tests with coverage
func TestCoverage() error {
	fmt.Println("Running tests with coverage...")
	if err := sh.Run("go", "test", "-v", "-coverprofile=coverage.out", "./..."); err != nil {
		return err
	}
	return sh.Run("go", "tool", "cover", "-html=coverage.out", "-o", "coverage.html")
}

// Lint runs golangci-lint
func Lint() error {
	fmt.Println("Running linter...")
	return sh.Run("golangci-lint", "run")
}

// Dev runs the server in development mode
func Dev() error {
	fmt.Println("Starting development server...")
	return sh.Run("go", "run", "./cmd")
}

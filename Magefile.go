// Magefile.go
//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/grafana/grafana-plugin-sdk-go/build"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
var Default = BuildAll

// BuildAll builds the backend plugin for all platforms and generates manifest
func BuildAll() error {
	mg.Deps(Clean)

	if err := Build(); err != nil {
		return err
	}

	return build.GenerateManifestFile()
}

// Build builds the backend plugin for all platforms
func Build() error {
	platforms := []struct {
		os   string
		arch string
	}{
		{"linux", "amd64"},
		{"linux", "arm64"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
		{"windows", "amd64"},
	}

	for _, p := range platforms {
		if err := buildBinary(p.os, p.arch); err != nil {
			return err
		}
	}

	return nil
}

func buildBinary(goos, goarch string) error {
	fmt.Printf("Building for %s/%s...\n", goos, goarch)

	env := map[string]string{
		"GOOS":   goos,
		"GOARCH": goarch,
		"CGO_ENABLED": "0",
	}

	output := filepath.Join("dist", fmt.Sprintf("gpx_hover-hover-panel_%s_%s", goos, goarch))
	if goos == "windows" {
		output += ".exe"
	}

	// Get build info
	buildInfo, err := build.GetBuildInfo()
	if err != nil {
		return fmt.Errorf("failed to get build info: %w", err)
	}

	// Build with ldflags to inject build info
	ldflags := buildInfo.MakeLdFlags()

	return sh.RunWith(env, "go", "build", "-ldflags", ldflags, "-o", output, "./pkg/cmd")
}

// Clean removes build artifacts
func Clean() error {
	fmt.Println("Cleaning...")
	os.RemoveAll("dist")
	return os.MkdirAll("dist", 0755)
}

// BuildDev builds backend for current platform only
func BuildDev() error {
	fmt.Println("Building for current platform...")
	return sh.Run("go", "build", "-o", "dist/gpx_hover-hover-panel", "./pkg/cmd")
}

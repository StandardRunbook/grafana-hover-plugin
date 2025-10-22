// Magefile.go
//go:build mage
// +build mage

package main

import (
	"fmt"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Build builds the backend plugin for all platforms
func Build() error {
	mg.Deps(Clean)

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

	return sh.RunWith(env, "go", "build", "-o", output, "./pkg/cmd")
}

// Clean removes build artifacts
func Clean() error {
	fmt.Println("Cleaning...")
	return sh.Rm("dist/gpx_hover-hover-panel_*")
}

// BuildDev builds backend for current platform only
func BuildDev() error {
	fmt.Println("Building for current platform...")
	return sh.Run("go", "build", "-o", "dist/gpx_hover-hover-panel", "./pkg/cmd")
}

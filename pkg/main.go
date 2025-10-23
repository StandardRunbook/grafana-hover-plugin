package main

import (
	"os"

	"github.com/StandardRunbook/grafana-hover-plugin/internal/plugin"

	"github.com/grafana/grafana-plugin-sdk-go/backend/app"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/build"
)

func main() {
	// Log build info
	buildInfo, err := build.GetBuildInfo()
	if err != nil {
		log.DefaultLogger.Warn("Failed to get build info", "error", err)
	} else {
		log.DefaultLogger.Info("Starting plugin", "version", buildInfo.Version, "pluginID", buildInfo.PluginID)
	}

	// Start listening to requests sent from Grafana
	// Manage automatically manages life cycle of app instances
	if err := app.Manage("hover-hover-panel", plugin.NewApp, app.ManageOpts{}); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}

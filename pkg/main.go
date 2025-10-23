package main

import (
	"context"
	"os"

	"github.com/StandardRunbook/grafana-hover-plugin/internal/plugin"

	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func main() {
	log.DefaultLogger.Info("Starting hover-hover-panel backend")

	// Create a single global instance for panel plugin
	app, err := plugin.NewPanelApp(context.Background())
	if err != nil {
		log.DefaultLogger.Error("Failed to create plugin", "error", err)
		os.Exit(1)
	}

	// Start listening to requests sent from Grafana
	// For panel plugins with backend, use datasource.Serve
	if err := datasource.Serve(datasource.ServeOpts{
		CallResourceHandler: app,
	}); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}

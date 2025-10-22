package main

import (
	"os"

	"grafana-plugin-api/internal/plugin"

	"github.com/grafana/grafana-plugin-sdk-go/backend/app"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func main() {
	// Start listening to requests sent from Grafana
	// Manage automatically manages life cycle of app instances
	if err := app.Manage("hover-hover-panel", plugin.NewApp, app.ManageOpts{}); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}

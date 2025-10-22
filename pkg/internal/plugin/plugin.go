package plugin

import (
	"context"
	"net/http"

	"grafana-plugin-api/internal/api"
	"grafana-plugin-api/internal/config"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
)

// Make sure App implements required interfaces
var (
	_ backend.CallResourceHandler = (*App)(nil)
)

// App is the backend plugin implementation
type App struct {
	backend.CallResourceHandler
	handler *api.Handler
}

// NewApp creates a new instance of the app plugin
func NewApp(ctx context.Context, settings backend.AppInstanceSettings) (instancemgmt.Instance, error) {
	log.DefaultLogger.Info("Creating new app instance")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.DefaultLogger.Error("Failed to load configuration", "error", err)
		// Don't fail plugin startup, just log the error
		cfg = &config.Config{
			Server: config.ServerConfig{
				Host: "0.0.0.0",
				Port: 8080,
			},
		}
	}

	// Create API handler
	handler := api.NewHandler(cfg)

	// Try to verify tables but don't fail if it doesn't work
	if err := handler.VerifyTables(); err != nil {
		log.DefaultLogger.Warn("Failed to verify ClickHouse tables", "error", err)
		log.DefaultLogger.Info("Plugin will return mock data when ClickHouse is unavailable")
	}

	app := &App{
		handler: handler,
	}

	// Setup resource handler
	mux := http.NewServeMux()
	mux.HandleFunc("/query_logs", app.handleQueryLogs)
	app.CallResourceHandler = httpadapter.New(mux)

	return app, nil
}

// Dispose is called when the app instance is being disposed
func (a *App) Dispose() {
	log.DefaultLogger.Info("Disposing app instance")
}

// handleQueryLogs handles the query_logs resource call
func (a *App) handleQueryLogs(w http.ResponseWriter, r *http.Request) {
	log.DefaultLogger.Debug("Handling query_logs request")
	a.handler.QueryLogs(w, r)
}

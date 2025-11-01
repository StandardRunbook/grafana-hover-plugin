package main

import (
	"log"
	"net/http"

	"github.com/StandardRunbook/grafana-hover-plugin/internal/api"
	"github.com/StandardRunbook/grafana-hover-plugin/internal/config"
)

func main() {
	log.Println("🚀 Starting Hover Plugin Standalone Server...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Printf("⚠️  Failed to load config: %v. Using defaults.", err)
		cfg = &config.Config{
			Server: config.ServerConfig{
				Host: "127.0.0.1",
				Port: 8080,
			},
			ClickHouse: config.ClickHouseConfig{
				URL:      "localhost:9000",
				Database: "default",
				User:     "default",
				Password: "",
			},
		}
	}

	log.Printf("📝 Config: Server=%s, ClickHouse=%s", cfg.Server.GetAddress(), cfg.ClickHouse.URL)

	// Create handler (it will create analyzer and connect to ClickHouse internally)
	handler := api.NewHandler(cfg)

	// Setup routes
	http.HandleFunc("/analyze", handler.QueryLogs)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	addr := cfg.Server.GetAddress()
	log.Printf("🎯 Server listening on http://%s", addr)
	log.Println("📊 Endpoints:")
	log.Println("   POST /analyze - Analyze logs with KL divergence")
	log.Println("   GET  /health  - Health check")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}

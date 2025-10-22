package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"grafana-plugin-api/internal/config"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2"
)

type Client struct {
	db *sql.DB
}

type TemplateCount struct {
	TemplateID string
	Count      uint64
}

type TemplateRepresentative struct {
	TemplateID        string
	RepresentativeLogs []string
}

func NewClient(cfg *config.ClickHouseConfig) (*Client, error) {
	// Build connection string
	opts := &clickhouse.Options{
		Addr: []string{cfg.URL},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.User,
			Password: cfg.Password,
		},
	}

	db := clickhouse.OpenDB(opts)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	return &Client{db: db}, nil
}

func (c *Client) Close() error {
	return c.db.Close()
}

// VerifyTables checks if required tables exist, creating them if missing
func (c *Client) VerifyTables() error {
	requiredTables := []string{"log_template_ids", "log_template_representatives"}
	var missingTables []string

	for _, tableName := range requiredTables {
		query := fmt.Sprintf("SELECT 1 FROM %s LIMIT 0", tableName)
		_, err := c.db.Query(query)
		if err != nil {
			if containsError(err, "UNKNOWN_TABLE") || containsError(err, "doesn't exist") {
				log.Printf("✗ Table '%s' does not exist", tableName)
				missingTables = append(missingTables, tableName)
			} else {
				return fmt.Errorf("error checking table '%s': %w", tableName, err)
			}
		} else {
			log.Printf("✓ Table '%s' exists", tableName)
		}
	}

	if len(missingTables) > 0 {
		log.Printf("Missing tables: %v. Creating them now...", missingTables)
		if err := c.createTables(); err != nil {
			return err
		}

		// Verify tables were created
		for _, tableName := range missingTables {
			query := fmt.Sprintf("SELECT 1 FROM %s LIMIT 0", tableName)
			if _, err := c.db.Query(query); err != nil {
				return fmt.Errorf("failed to verify table '%s' after creation: %w", tableName, err)
			}
		}
		log.Println("✓ Successfully verified all created tables")
	}

	log.Println("✓ All required ClickHouse tables exist")
	return nil
}

func (c *Client) createTables() error {
	schemaPath := "schema/clickhouse_schema.sql"
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file at %s: %w", schemaPath, err)
	}

	schema := string(schemaBytes)

	// Split by semicolon and execute each statement
	statements := splitSQL(schema)
	log.Printf("Found %d SQL statements in schema", len(statements))

	for i, statement := range statements {
		if statement == "" {
			continue
		}

		log.Printf("Executing SQL statement %d of %d...", i+1, len(statements))

		if _, err := c.db.Exec(statement); err != nil {
			return fmt.Errorf("failed to execute statement %d: %w\nSQL: %s", i+1, err, statement)
		}
	}

	log.Println("✓ Successfully executed all schema statements")
	return nil
}

// GetTemplateCounts retrieves template ID counts for a given time window
func (c *Client) GetTemplateCounts(ctx context.Context, org, dashboard, panelTitle, metricName string, startTime, endTime time.Time) (map[string]uint64, error) {
	query := `
		SELECT
			template_id,
			count(*) as count
		FROM log_template_ids
		WHERE org = ?
			AND dashboard = ?
			AND panel_title = ?
			AND metric_name = ?
			AND timestamp >= ?
			AND timestamp < ?
		GROUP BY template_id
	`

	rows, err := c.db.QueryContext(ctx, query, org, dashboard, panelTitle, metricName, startTime, endTime)
	if err != nil {
		if containsError(err, "UNKNOWN_TABLE") {
			return nil, fmt.Errorf("table 'log_template_ids' does not exist. Please restart the service to auto-create tables")
		}
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]uint64)
	for rows.Next() {
		var tc TemplateCount
		if err := rows.Scan(&tc.TemplateID, &tc.Count); err != nil {
			return nil, err
		}
		counts[tc.TemplateID] = tc.Count
	}

	return counts, rows.Err()
}

// GetRepresentativeLogs retrieves representative logs for specific template IDs
func (c *Client) GetRepresentativeLogs(ctx context.Context, org, dashboard, panelTitle, metricName string, templateIDs []string) (map[string][]string, error) {
	if len(templateIDs) == 0 {
		return make(map[string][]string), nil
	}

	query := `
		SELECT
			template_id,
			representative_logs
		FROM log_template_representatives
		WHERE org = ?
			AND dashboard = ?
			AND panel_title = ?
			AND metric_name = ?
			AND template_id IN (?)
	`

	// ClickHouse requires array format for IN clause
	rows, err := c.db.QueryContext(ctx, query, org, dashboard, panelTitle, metricName, templateIDs)
	if err != nil {
		if containsError(err, "UNKNOWN_TABLE") {
			return nil, fmt.Errorf("table 'log_template_representatives' does not exist. Please restart the service to auto-create tables")
		}
		return nil, err
	}
	defer rows.Close()

	representatives := make(map[string][]string)
	for rows.Next() {
		var tr TemplateRepresentative
		if err := rows.Scan(&tr.TemplateID, &tr.RepresentativeLogs); err != nil {
			return nil, err
		}
		representatives[tr.TemplateID] = tr.RepresentativeLogs
	}

	return representatives, rows.Err()
}

// Helper functions

func containsError(err error, substr string) bool {
	if err == nil {
		return false
	}
	return contains(err.Error(), substr)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func splitSQL(schema string) []string {
	var statements []string
	var current string

	for _, line := range splitLines(schema) {
		trimmed := trimSpace(line)

		// Skip empty lines and comments
		if trimmed == "" || hasPrefix(trimmed, "--") {
			continue
		}

		current += line + "\n"

		// Check if statement ends with semicolon
		if hasSuffix(trimSpace(current), ";") {
			statements = append(statements, trimSpace(current[:len(current)-1]))
			current = ""
		}
	}

	if trimSpace(current) != "" {
		statements = append(statements, trimSpace(current))
	}

	return statements
}

func splitLines(s string) []string {
	var lines []string
	current := ""

	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(s[i])
		}
	}

	if current != "" {
		lines = append(lines, current)
	}

	return lines
}

func trimSpace(s string) string {
	start := 0
	end := len(s)

	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}

	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}

	return s[start:end]
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

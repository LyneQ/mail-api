package test

import (
	"github.com/lyneq/mailapi/config"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file for testing
	tempFile := "config_test.ini"
	content := `[AllowedDomains]
domains = example.com, test.com

[Database]
driver = sqlite
path = ./db/test.db`

	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary config file: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {

		}
	}(tempFile)

	// Override the config file path for testing
	originalOpen := config.OsOpen
	defer func() { config.OsOpen = originalOpen }()

	config.OsOpen = func(name string) (*os.File, error) {
		return os.Open(tempFile)
	}

	// Load the configuration
	err = config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Check if the allowed domains were loaded correctly
	domains := config.GetAllowedDomains()
	if len(domains) != 2 {
		t.Errorf("Expected 2 domains, got %d", len(domains))
	}

	if domains[0] != "example.com" {
		t.Errorf("Expected first domain to be example.com, got %s", domains[0])
	}

	if domains[1] != "test.com" {
		t.Errorf("Expected second domain to be test.com, got %s", domains[1])
	}

	// Check if the database configuration was loaded correctly
	driver := config.GetDatabaseDriver()
	if driver != "sqlite" {
		t.Errorf("Expected database driver to be sqlite, got %s", driver)
	}

	path := config.GetDatabasePath()
	if path != "./db/test.db" {
		t.Errorf("Expected database path to be ./db/test.db, got %s", path)
	}
}

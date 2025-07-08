package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Config holds all configuration values
type Config struct {
	AllowedDomains []string
	Database       DatabaseConfig
	Api            ApiConfig
}

// DatabaseConfig holds database configuration values
type DatabaseConfig struct {
	Driver string
	Path   string
}

type ApiConfig struct {
	port string
}

var (
	// AppConfig is the global configuration instance
	AppConfig Config

	// OsOpen is a variable that holds the function to open files
	// It can be overridden for testing
	OsOpen = os.Open
)

// LoadConfig loads configuration from the .ini file
func LoadConfig() error {
	file, err := OsOpen("config/config.ini")
	if err != nil {
		return fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()

	var currentSection string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}

		// Check if this is a section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.TrimPrefix(strings.TrimSuffix(line, "]"), "[")
			continue
		}

		// Parse key-value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if currentSection == "AllowedDomains" {
			if key == "domains" {
				// Split the comma-separated list of domains
				domains := strings.Split(value, ",")
				for i, domain := range domains {
					domains[i] = strings.TrimSpace(domain)
				}
				AppConfig.AllowedDomains = domains
			}
		} else if currentSection == "Database" {
			switch key {
			case "driver":
				AppConfig.Database.Driver = value
			case "path":
				AppConfig.Database.Path = value
			}
		} else if currentSection == "Api" {
			switch key {
			case "port":
				AppConfig.Api.port = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	fmt.Printf("[App] Configuration loaded from file %v\n", file.Name())
	return nil
}

// GetAllowedDomains returns the list of allowed domains
func GetAllowedDomains() []string {
	return AppConfig.AllowedDomains
}

// GetDatabaseDriver returns the database driver
func GetDatabaseDriver() string {
	return AppConfig.Database.Driver
}

// GetDatabasePath returns the database path
func GetDatabasePath() string {
	return AppConfig.Database.Path
}

func GetAPIPort() string {
	return AppConfig.Api.port
}

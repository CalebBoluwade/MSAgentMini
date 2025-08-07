package constants

import (
	"github.com/google/uuid"
	"log"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

// GetEnvWithDefault retrieves an environment variable with a fallback default value
func GetEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		// log.Printf("Found env %s=%s", key, value)
		return value
	}

	upperKey := strings.ToUpper(key)
	if upperKey != key {
		if value, exists := os.LookupEnv(upperKey); exists {
			log.Printf("Found env %s=%s (via uppercase conversion)", upperKey, value)
			return value
		}
	}

	// Try a lowercase version if different
	lowerKey := strings.ToLower(key)
	if lowerKey != key {
		if value, exists := os.LookupEnv(lowerKey); exists {
			log.Printf("Found env %s=%s (via lowercase conversion)", lowerKey, value)
			return value
		}
	}

	log.Printf("Env key '%s' not found, using default: '%s'", key, defaultValue)

	return defaultValue
}

const (
	AgentVersion = "1.0.0"
)

var AgentID = GetEnvWithDefault("AGENT_ID", uuid.New().String())

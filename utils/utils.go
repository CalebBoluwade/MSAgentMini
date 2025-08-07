package utils

import (
	"log"
	"net"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

// GetHostIP gets the preferred outbound ip of this machine.
func GetHostIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		// Fallback to hostname if UDP dial fails
		addrs, err := net.LookupHost("localhost")
		if err == nil && len(addrs) > 0 {
			return addrs[0]
		}
		return "N/A"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// GetEnvWithDefault retrieves an environment variable with a fallback default value
func GetEnvWithDefault(key, defaultValue string) string {
	// Try the exact case first
	if value := os.Getenv(key); value != "" {
		// log.Printf("Found env %s=%s", key, value)
		return value
	}

	// Try an uppercase version if different
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

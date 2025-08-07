package main

import (
	"MSAgent/api"
	"MSAgent/collectors"
	"MSAgent/database"
	"MSAgent/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// ===== MONITORING AGENT =====

// startMonitoring starts a goroutine to collect metrics periodically.
func startMonitoring(ctx context.Context, db *gorm.DB, interval time.Duration) {
	log.Printf("Starting monitoring agent with a collection interval of %s", interval)

	// Collect metrics immediately on start
	go collectMetrics(db)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go collectMetrics(db)
		case <-ctx.Done():
			log.Println("Stopping monitoring agent.")
			return
		}
	}
}

func collectMetrics(db *gorm.DB) {
	log.Println("Collecting Host System Metrics...")
	collectors.CollectCPUInfo(db)
	collectors.CollectMemoryInfo(db)
	collectors.CollectDiskSpace(db, false)
	collectors.CollectNetworkInfo(db)
	log.Println("Metrics collection complete.. Agent has completed routine procedure")
}

func main() {
	log.Println("Initializing Agent Service")

	// Create a context that we can cancel to gracefully shut down the monitoring agent
	ctx, cancel := context.WithCancel(context.Background())

	var db = database.InitializeDB()

	// Start the monitoring agent in a separate goroutine
	go startMonitoring(ctx, db, 1*time.Minute)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// Add middleware
	app.Use(cors.New())
	app.Use(logger.New())

	apiHandler := &api.AgentAPIHandler{
		DB: db,
	}

	// Setup API routes
	apiHandler.SetupRoutes(app)

	PortStr := utils.GetEnvWithDefault("AGENT_API_PORT", "30025")
	agentPort, portErr := strconv.Atoi(PortStr)

	if portErr != nil {
		log.Fatalf("Invalid Agent API Port: %v", portErr)
	}

	// Start the server
	go func() {
		log.Println(fmt.Sprintf("Starting Agent API server on :%d", agentPort))

		if err := app.Listen(fmt.Sprintf(":%d", agentPort)); err != nil {
			log.Fatalf("Error starting Agent API server: %v", err)
		}
	}()

	// Wait for a shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown signal received, shutting down gracefully...")

	// Cancel the context to stop the monitoring agent
	cancel()

	// Shutdown the Fiber app
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Error shutting down Agent server: %v", err)
	}

	log.Println("-------------------- AGENT STOPPED --------------------")
}

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"github.com/Aditya-Pimpalkar/clarity/internal/api"
	"github.com/Aditya-Pimpalkar/clarity/internal/middleware"
	"github.com/Aditya-Pimpalkar/clarity/internal/repository"
	"github.com/Aditya-Pimpalkar/clarity/internal/kafka"
	"github.com/Aditya-Pimpalkar/clarity/internal/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using environment variables")
	}

	// Build configuration
	config := loadConfig()

	// Print startup banner
	printBanner(config)

	// Connect to ClickHouse
	log.Println("🔌 Connecting to ClickHouse...")
	repo, err := repository.NewClickHouseRepository("localhost:9000")
	if err != nil {
		log.Fatal("❌ Failed to connect to ClickHouse:", err)
	}
	defer repo.Close()
	log.Println("✅ Connected to ClickHouse")

	// Test database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := repo.Ping(ctx); err != nil {
		log.Fatal("❌ ClickHouse ping failed:", err)
	}
	log.Println("✅ Database health check passed")

	// Initialize Kafka producer (optional - continues if unavailable)
	log.Println("🔌 Connecting to Kafka...")
	kafkaConfig := &kafka.ProducerConfig{
		Brokers: []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
		Topic:   getEnv("KAFKA_TOPIC", "llm-traces"),
	}
	
	kafkaProducer, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		log.Printf("⚠️  Kafka unavailable (continuing without events): %v", err)
		kafkaProducer = nil // Continue without Kafka
	} else {
		defer kafkaProducer.Close()
		log.Println("✅ Kafka producer connected")
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:               config.AppName,
		ServerHeader:          "LLM-Observability",
		DisableStartupMessage: true,
		ErrorHandler:          customErrorHandler,
		ReadTimeout:           time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout:          time.Duration(config.WriteTimeout) * time.Second,
		IdleTimeout:           time.Duration(config.IdleTimeout) * time.Second,
	})

	// Setup middleware
	log.Println("🔧 Setting up middleware...")
	setupMiddleware(app, config)
	log.Println("✅ Middleware configured")

	// Setup routes
	log.Println("🛣️  Setting up routes...")
	setupRoutes(app, repo, kafkaProducer, config)
	log.Println("✅ Routes configured")

	// Start server in goroutine
	go func() {
		log.Printf("🚀 Server starting on port %s", config.Port)
		log.Printf("📊 Health: http://localhost:%s/health", config.Port)
		log.Printf("📚 API: http://localhost:%s/api/v1", config.Port)
		log.Printf("🔐 Auth: JWT and API Key supported")
		log.Println("✨ Press Ctrl+C to shutdown")

		if err := app.Listen("0.0.0.0:" + config.Port); err != nil {
			log.Fatal("❌ Server failed to start:", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Fatal("❌ Server forced to shutdown:", err)
	}

	log.Println("✅ Server gracefully stopped")
}

// Config holds application configuration
type Config struct {
	AppName       string
	Port          string
	Environment   string
	ClickHouseDSN string
	JWTSecret     string
	CORSOrigins   string
	ReadTimeout   int
	WriteTimeout  int
	IdleTimeout   int
}

// loadConfig loads configuration from environment
func loadConfig() Config {
	return Config{
		AppName:       getEnv("APP_NAME", "LLM Observability Platform"),
		Port:          getEnv("PORT", "8080"),
		Environment:   getEnv("ENV", "development"),
		ClickHouseDSN: buildClickHouseDSN(),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		CORSOrigins:   getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173"),
		ReadTimeout:   getEnvInt("READ_TIMEOUT", 10),
		WriteTimeout:  getEnvInt("WRITE_TIMEOUT", 10),
		IdleTimeout:   getEnvInt("IDLE_TIMEOUT", 120),
	}
}

// setupMiddleware configures all middleware
func setupMiddleware(app *fiber.App, config Config) {
	// Recovery (must be first)
	app.Use(recover.New(recover.Config{
		EnableStackTrace: config.Environment == "development",
	}))

	// Request logger
	app.Use(middleware.RequestLogger())

	// CORS
	app.Use(middleware.CORSConfig())

	// Compression
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// Global rate limiter
	app.Use(middleware.GlobalRateLimiter())
}

// setupRoutes configures all routes with appropriate middleware
func setupRoutes(app *fiber.App, repo repository.Repository, kafkaProducer *kafka.Producer, config Config) {
	// Create services
	traceService := services.NewTraceService(repo, kafkaProducer)
	analyticsService := services.NewAnalyticsService(repo)
	userService := services.NewUserService(repo)

	// Create handlers
	healthHandler := api.NewHealthHandler(repo)
	traceHandler := api.NewTraceHandler(traceService)
	analyticsHandler := api.NewAnalyticsHandler(analyticsService)
	authHandler := api.NewAuthHandler(userService)

	// Public routes (no authentication)
	setupPublicRoutes(app, healthHandler, authHandler)

	// API key routes (SDK ingestion + analytics)
	setupAPIKeyRoutes(app, traceHandler, analyticsHandler)

	// JWT routes (dashboard)
	setupAuthenticatedRoutes(app, traceHandler, analyticsHandler, userService)
}

// setupPublicRoutes configures public endpoints
func setupPublicRoutes(app *fiber.App, healthHandler *api.HealthHandler, authHandler *api.AuthHandler) {
	// Health checks
	app.Get("/health", healthHandler.GetHealth)
	app.Get("/ready", healthHandler.GetReadiness)
	app.Get("/live", healthHandler.GetLiveness)

	// API info
	app.Get("/api/v1", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name":    "LLM Observability Platform API",
			"version": "1.0.0",
			"status":  "operational",
			"docs":    "https://docs.yourdomain.com",
			"authentication": fiber.Map{
				"jwt":     "Authorization: Bearer <token>",
				"api_key": "X-API-Key: <key>",
			},
			"endpoints": fiber.Map{
				"health":    "/health",
				"traces":    "/api/v1/traces",
				"analytics": "/api/v1/analytics",
				"auth":      "/api/v1/auth",
			},
		})
	})

	// Authentication endpoints
	auth := app.Group("/api/v1/auth")
	auth.Post("/login", authHandler.Login)
}

// setupAPIKeyRoutes configures API key protected routes
func setupAPIKeyRoutes(app *fiber.App, traceHandler *api.TraceHandler, analyticsHandler *api.AnalyticsHandler) {
	apiKey := app.Group("/api/v1",
		middleware.APIKeyAuth(),
		middleware.APIKeyRateLimiter(),
	)

	// Trace ingestion (SDK usage)
	apiKey.Post("/traces", traceHandler.CreateTrace)
	apiKey.Post("/traces/batch", traceHandler.CreateTraceBatch)

	// Analytics (also accessible via API key for programmatic access)
	analytics := apiKey.Group("/analytics")
	analytics.Get("/dashboard", analyticsHandler.GetDashboard)
	analytics.Get("/costs", analyticsHandler.GetCostAnalysis)
	analytics.Get("/performance", analyticsHandler.GetPerformanceMetrics)
	analytics.Get("/models", analyticsHandler.GetModelComparison)

	// ADD THESE LINES - Trace reading (for frontend)
	apiKey.Get("/traces", traceHandler.ListTraces)
	apiKey.Get("/traces/:id", traceHandler.GetTrace)
}

// setupAuthenticatedRoutes configures JWT protected routes
func setupAuthenticatedRoutes(app *fiber.App, traceHandler *api.TraceHandler,
	analyticsHandler *api.AnalyticsHandler, userService *services.UserService) {

	auth := app.Group("/api/v1",
		middleware.AuthMiddleware(),
		middleware.StrictRateLimiter(),
	)

	// Traces (read operations)
	auth.Get("/traces", traceHandler.ListTraces)
	auth.Get("/traces/:id", traceHandler.GetTrace)

	// Analytics
	analytics := auth.Group("/analytics")
	analytics.Get("/dashboard", analyticsHandler.GetDashboard)
	analytics.Get("/costs", analyticsHandler.GetCostAnalysis)
	analytics.Get("/performance", analyticsHandler.GetPerformanceMetrics)
	analytics.Get("/models", analyticsHandler.GetModelComparison)

	// User endpoints
	authHandler := api.NewAuthHandler(userService)
	auth.Get("/auth/me", authHandler.GetCurrentUser)
	auth.Post("/auth/api-keys", authHandler.GenerateAPIKey)

	// Admin routes
	admin := auth.Group("/admin", middleware.RequireRole("admin"))
	admin.Get("/stats", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Admin statistics"})
	})
}

// buildClickHouseDSN builds the ClickHouse connection string
func buildClickHouseDSN() string {
	host := getEnv("CLICKHOUSE_HOST", "localhost")
	port := getEnv("CLICKHOUSE_PORT", "9000")
	database := getEnv("CLICKHOUSE_DATABASE", "llm_observability")
	user := getEnv("CLICKHOUSE_USER", "default")
	password := getEnv("CLICKHOUSE_PASSWORD", "")

	dsn := "clickhouse://"
	if user != "" {
		dsn += user
		if password != "" {
			dsn += ":" + password
		}
		dsn += "@"
	}
	dsn += host + ":" + port + "/" + database

	return dsn
}

// customErrorHandler handles errors consistently
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	requestID := middleware.GetRequestID(c)
	log.Printf("❌ [ERROR] Request ID: %s | %v", requestID, err)

	return c.Status(code).JSON(fiber.Map{
		"error":      message,
		"request_id": requestID,
		"timestamp":  time.Now().Format(time.RFC3339),
	})
}

// printBanner prints startup banner
func printBanner(config Config) {
	banner := `
╔═══════════════════════════════════════════════════╗
║                                                   ║
║     LLM Observability Platform                    ║
║     Production-Ready API Server                   ║
║                                                   ║
╚═══════════════════════════════════════════════════╝
`
	log.Println(banner)
	log.Printf("Environment: %s", config.Environment)
	log.Printf("Version: 1.0.0")
	log.Println("─────────────────────────────────────────────────────")
}

// getEnv gets environment variable with default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets integer environment variable with default
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
// CI test
// CI test

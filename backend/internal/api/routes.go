package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Aditya-Pimpalkar/clarity/internal/repository"
	"github.com/Aditya-Pimpalkar/clarity/internal/services"
)

// SetupRoutes configures all API routes
func SetupRoutes(app *fiber.App, repo repository.Repository) {
	// Create services
	traceService := services.NewTraceService(repo, nil)
	analyticsService := services.NewAnalyticsService(repo)

	// Create handlers
	traceHandler := NewTraceHandler(traceService)
	analyticsHandler := NewAnalyticsHandler(analyticsService)
	healthHandler := NewHealthHandler(repo)

	// Health check routes
	app.Get("/health", healthHandler.GetHealth)
	app.Get("/ready", healthHandler.GetReadiness)
	app.Get("/live", healthHandler.GetLiveness)

	// API v1 routes
	v1 := app.Group("/api/v1")

	// API info
	v1.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name":    "LLM Observability Platform API",
			"version": "1.0.0",
			"endpoints": fiber.Map{
				"traces":    "/api/v1/traces",
				"analytics": "/api/v1/analytics",
				"health":    "/health",
			},
		})
	})

	// Trace routes
	traces := v1.Group("/traces")
	traces.Post("/", traceHandler.CreateTrace)
	traces.Post("/batch", traceHandler.CreateTraceBatch)
	traces.Get("/", traceHandler.ListTraces)
	traces.Get("/:id", traceHandler.GetTrace)

	// Analytics routes
	analytics := v1.Group("/analytics")
	analytics.Get("/dashboard", analyticsHandler.GetDashboard)
	analytics.Get("/costs", analyticsHandler.GetCostAnalysis)
	analytics.Get("/performance", analyticsHandler.GetPerformanceMetrics)
	analytics.Get("/models", analyticsHandler.GetModelComparison)
}

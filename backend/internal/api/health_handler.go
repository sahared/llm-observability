package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/Aditya-Pimpalkar/clarity/internal/models"
	"github.com/Aditya-Pimpalkar/clarity/internal/repository"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	repo repository.Repository
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(repo repository.Repository) *HealthHandler {
	return &HealthHandler{
		repo: repo,
	}
}

// GetHealth handles GET /health
func (h *HealthHandler) GetHealth(c *fiber.Ctx) error {
	services := make(map[string]string)

	// Check ClickHouse connection
	if err := h.repo.Ping(c.Context()); err != nil {
		services["clickhouse"] = "unhealthy: " + err.Error()
	} else {
		services["clickhouse"] = "healthy"
	}

	// Determine overall status
	status := "ok"
	statusCode := fiber.StatusOK

	for _, serviceStatus := range services {
		if serviceStatus != "healthy" {
			status = "degraded"
			statusCode = fiber.StatusServiceUnavailable
			break
		}
	}

	response := models.HealthResponse{
		Status:    status,
		Version:   "1.0.0",
		Services:  services,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return c.Status(statusCode).JSON(response)
}

// GetReadiness handles GET /ready
func (h *HealthHandler) GetReadiness(c *fiber.Ctx) error {
	if err := h.repo.Ping(c.Context()); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"ready":   false,
			"message": "Database connection failed",
		})
	}

	return c.JSON(fiber.Map{
		"ready":   true,
		"message": "Service is ready",
	})
}

// GetLiveness handles GET /live
func (h *HealthHandler) GetLiveness(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"alive": true,
	})
}

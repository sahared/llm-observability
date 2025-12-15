package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Aditya-Pimpalkar/clarity/internal/middleware"
	"github.com/Aditya-Pimpalkar/clarity/internal/services"
)

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetDashboard handles GET /api/v1/analytics/dashboard
func (h *AnalyticsHandler) GetDashboard(c *fiber.Ctx) error {
	timeRange := c.Query("time_range", "24h")

	// Get org ID from context (set by API key or JWT middleware)
	orgID := middleware.GetOrgID(c)

	// If no org ID from auth, try query parameter or use default
	if orgID == "" {
		orgID = c.Query("organization_id", "org-test-123")
	}

	stats, err := h.analyticsService.GetDashboard(c.Context(), timeRange, orgID)
	if err != nil {
		return InternalErrorResponse(c, "Failed to get dashboard stats")
	}

	return SuccessResponse(c, stats)
}

// GetCostAnalysis handles GET /api/v1/analytics/costs
func (h *AnalyticsHandler) GetCostAnalysis(c *fiber.Ctx) error {
	return SuccessResponse(c, fiber.Map{
		"message": "Cost analysis endpoint - coming soon!",
	})
}

// GetPerformanceMetrics handles GET /api/v1/analytics/performance
func (h *AnalyticsHandler) GetPerformanceMetrics(c *fiber.Ctx) error {
	return SuccessResponse(c, fiber.Map{
		"message": "Performance metrics endpoint - coming soon!",
	})
}

// GetModelComparison handles GET /api/v1/analytics/models
func (h *AnalyticsHandler) GetModelComparison(c *fiber.Ctx) error {
	return SuccessResponse(c, fiber.Map{
		"message": "Model comparison endpoint - coming soon!",
	})
}

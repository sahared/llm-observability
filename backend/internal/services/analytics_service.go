package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Aditya-Pimpalkar/clarity/internal/models"
	"github.com/Aditya-Pimpalkar/clarity/internal/repository"
)

// AnalyticsService handles business logic for analytics and metrics
type AnalyticsService struct {
	repo repository.Repository
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(repo repository.Repository) *AnalyticsService {
	return &AnalyticsService{
		repo: repo,
	}
}

// ============================================================================
// ORIGINAL METHODS (Keep these - they work with your existing repository)
// ============================================================================

// GetDashboardSummary returns a comprehensive dashboard summary
func (s *AnalyticsService) GetDashboardSummary(ctx context.Context, orgID, projectID string, timeRange string) (*DashboardSummary, error) {
	// Validate inputs
	if orgID == "" || projectID == "" {
		return nil, fmt.Errorf("organization_id and project_id are required")
	}

	// Parse time range
	startTime, endTime, err := s.parseTimeRange(timeRange)
	if err != nil {
		return nil, fmt.Errorf("invalid time range: %w", err)
	}

	// Get metric summary
	metricSummary, err := s.repo.GetMetricSummary(ctx, orgID, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get metric summary: %w", err)
	}

	// Get cost breakdown
	costBreakdown, err := s.repo.GetCostBreakdown(ctx, orgID, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get cost breakdown: %w", err)
	}

	// Get model usage
	modelUsage, err := s.repo.GetModelUsage(ctx, orgID, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get model usage: %w", err)
	}

	// Calculate additional insights
	insights := s.generateInsights(metricSummary, costBreakdown, modelUsage)

	return &DashboardSummary{
		TimeRange:     timeRange,
		StartTime:     startTime,
		EndTime:       endTime,
		MetricSummary: metricSummary,
		CostBreakdown: costBreakdown,
		ModelUsage:    modelUsage,
		Insights:      insights,
	}, nil
}

// parseTimeRange converts time range string to start and end times
func (s *AnalyticsService) parseTimeRange(timeRange string) (time.Time, time.Time, error) {
	now := time.Now()
	var startTime time.Time

	switch timeRange {
	case "1h", "last_hour":
		startTime = now.Add(-1 * time.Hour)
	case "24h", "last_24h", "today":
		startTime = now.Add(-24 * time.Hour)
	case "7d", "last_7d", "week":
		startTime = now.Add(-7 * 24 * time.Hour)
	case "30d", "last_30d", "month":
		startTime = now.Add(-30 * 24 * time.Hour)
	case "90d", "last_90d":
		startTime = now.Add(-90 * 24 * time.Hour)
	default:
		return time.Time{}, time.Time{}, fmt.Errorf("unsupported time range: %s", timeRange)
	}

	return startTime, now, nil
}

// generateInsights generates insights from the data
func (s *AnalyticsService) generateInsights(
	summary *models.MetricSummary,
	costBreakdown []*models.CostBreakdown,
	modelUsage []*models.ModelUsage,
) []Insight {
	insights := []Insight{}

	// Insight 1: Error rate
	if summary.ErrorRate > 5.0 {
		insights = append(insights, Insight{
			Type:        "warning",
			Category:    "reliability",
			Title:       "High Error Rate",
			Description: fmt.Sprintf("Error rate is %.2f%%, which is above the recommended 5%% threshold", summary.ErrorRate),
			Severity:    "high",
		})
	} else if summary.ErrorRate > 1.0 {
		insights = append(insights, Insight{
			Type:        "info",
			Category:    "reliability",
			Title:       "Elevated Error Rate",
			Description: fmt.Sprintf("Error rate is %.2f%%, consider investigating recent changes", summary.ErrorRate),
			Severity:    "medium",
		})
	}

	// Insight 2: Cost concentration
	if len(costBreakdown) > 0 && costBreakdown[0].Percentage > 80.0 {
		insights = append(insights, Insight{
			Type:        "info",
			Category:    "cost",
			Title:       "Cost Concentration",
			Description: fmt.Sprintf("%.1f%% of costs come from %s. Consider model optimization or caching", costBreakdown[0].Percentage, costBreakdown[0].Model),
			Severity:    "low",
		})
	}

	// Insight 3: Latency
	if summary.P95LatencyMs > 2000 {
		insights = append(insights, Insight{
			Type:        "warning",
			Category:    "performance",
			Title:       "High P95 Latency",
			Description: fmt.Sprintf("P95 latency is %.0fms. Users may experience slow responses", summary.P95LatencyMs),
			Severity:    "medium",
		})
	}

	// Insight 4: Cost per request
	if summary.AvgCostPerRequest > 0.01 {
		insights = append(insights, Insight{
			Type:        "info",
			Category:    "cost",
			Title:       "High Average Cost",
			Description: fmt.Sprintf("Average cost per request is $%.4f. Consider optimizing prompts or using cheaper models", summary.AvgCostPerRequest),
			Severity:    "low",
		})
	}

	// Insight 5: Success rate achievement
	if summary.SuccessRate > 99.0 {
		insights = append(insights, Insight{
			Type:        "success",
			Category:    "reliability",
			Title:       "Excellent Reliability",
			Description: fmt.Sprintf("Success rate is %.2f%% - great job!", summary.SuccessRate),
			Severity:    "info",
		})
	}

	// Insight 6: Model usage patterns
	if len(modelUsage) > 0 {
		// Find the most used model
		mostUsed := modelUsage[0]
		for _, usage := range modelUsage {
			if usage.CallCount > mostUsed.CallCount {
				mostUsed = usage
			}
		}

		if mostUsed.CallCount > 100 {
			insights = append(insights, Insight{
				Type:        "info",
				Category:    "usage",
				Title:       "High Model Usage",
				Description: fmt.Sprintf("%s is your most used model with %d calls. Ensure you're getting the best value", mostUsed.Model, mostUsed.CallCount),
				Severity:    "info",
			})
		}
	}

	// Insight 7: Token efficiency
	if len(modelUsage) > 0 {
		for _, usage := range modelUsage {
			if usage.CallCount > 0 {
				avgTokensPerCall := float64(usage.TotalTokens) / float64(usage.CallCount)
				if avgTokensPerCall > 2000 {
					insights = append(insights, Insight{
						Type:        "info",
						Category:    "efficiency",
						Title:       "High Token Usage",
						Description: fmt.Sprintf("%s averages %.0f tokens per call. Consider prompt optimization", usage.Model, avgTokensPerCall),
						Severity:    "low",
					})
					break // Only show one token efficiency insight
				}
			}
		}
	}

	// Insight 8: Model diversity
	if len(modelUsage) == 1 && modelUsage[0].CallCount > 50 {
		insights = append(insights, Insight{
			Type:        "info",
			Category:    "optimization",
			Title:       "Single Model Usage",
			Description: fmt.Sprintf("You're only using %s. Consider testing other models for cost/performance optimization", modelUsage[0].Model),
			Severity:    "low",
		})
	} else if len(modelUsage) >= 3 {
		insights = append(insights, Insight{
			Type:        "success",
			Category:    "optimization",
			Title:       "Good Model Diversity",
			Description: fmt.Sprintf("Using %d different models - great job optimizing for different use cases!", len(modelUsage)),
			Severity:    "info",
		})
	}

	return insights
}

// GetCostAnalysis returns detailed cost analysis
func (s *AnalyticsService) GetCostAnalysis(ctx context.Context, orgID, projectID string, startTime, endTime time.Time) (*CostAnalysis, error) {
	// Get cost breakdown by model
	breakdown, err := s.repo.GetCostBreakdown(ctx, orgID, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get cost breakdown: %w", err)
	}

	// Calculate total cost
	var totalCost float64
	for _, item := range breakdown {
		totalCost += item.TotalCost
	}

	// Calculate daily average
	days := endTime.Sub(startTime).Hours() / 24
	if days == 0 {
		days = 1
	}
	dailyAverage := totalCost / days

	// Project monthly cost
	monthlyProjection := dailyAverage * 30

	// Find most expensive model
	var mostExpensiveModel string
	var highestCost float64
	if len(breakdown) > 0 {
		mostExpensiveModel = breakdown[0].Model
		highestCost = breakdown[0].TotalCost
	}

	return &CostAnalysis{
		TotalCost:          totalCost,
		DailyAverage:       dailyAverage,
		MonthlyProjection:  monthlyProjection,
		CostBreakdown:      breakdown,
		MostExpensiveModel: mostExpensiveModel,
		HighestCost:        highestCost,
	}, nil
}

// GetPerformanceMetrics returns performance analysis
func (s *AnalyticsService) GetPerformanceMetrics(ctx context.Context, orgID, projectID string, startTime, endTime time.Time) (*PerformanceMetrics, error) {
	summary, err := s.repo.GetMetricSummary(ctx, orgID, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get metric summary: %w", err)
	}

	// Calculate additional metrics
	var status string
	var recommendation string

	// Determine overall performance status
	if summary.P95LatencyMs < 500 && summary.ErrorRate < 1.0 {
		status = "excellent"
		recommendation = "Performance is optimal. Continue monitoring."
	} else if summary.P95LatencyMs < 1000 && summary.ErrorRate < 5.0 {
		status = "good"
		recommendation = "Performance is acceptable but could be improved."
	} else if summary.P95LatencyMs < 2000 && summary.ErrorRate < 10.0 {
		status = "fair"
		recommendation = "Performance issues detected. Consider optimization."
	} else {
		status = "poor"
		recommendation = "Critical performance issues. Immediate attention required."
	}

	return &PerformanceMetrics{
		Summary:        summary,
		Status:         status,
		Recommendation: recommendation,
	}, nil
}

// GetModelComparison compares performance and cost across models
func (s *AnalyticsService) GetModelComparison(ctx context.Context, orgID, projectID string, startTime, endTime time.Time) ([]*ModelComparison, error) {
	modelUsage, err := s.repo.GetModelUsage(ctx, orgID, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get model usage: %w", err)
	}

	comparisons := make([]*ModelComparison, len(modelUsage))
	for i, usage := range modelUsage {
		// Calculate metrics per request
		avgCostPerRequest := 0.0
		avgTokensPerRequest := 0.0
		if usage.CallCount > 0 {
			avgCostPerRequest = usage.TotalCost / float64(usage.CallCount)
			avgTokensPerRequest = float64(usage.TotalTokens) / float64(usage.CallCount)
		}

		// Calculate efficiency score (lower is better)
		// Score = (normalized_latency + normalized_cost) / 2
		efficiencyScore := (usage.AvgLatency / 1000.0) + (avgCostPerRequest * 100.0)

		comparisons[i] = &ModelComparison{
			Model:               usage.Model,
			TotalCalls:          usage.CallCount,
			TotalCost:           usage.TotalCost,
			AvgCostPerRequest:   avgCostPerRequest,
			AvgLatency:          usage.AvgLatency,
			AvgTokensPerRequest: avgTokensPerRequest,
			EfficiencyScore:     efficiencyScore,
		}
	}

	return comparisons, nil
}

// ============================================================================
// NEW DASHBOARD METHOD (For frontend dashboard API)
// ============================================================================

// GetDashboard returns dashboard stats matching frontend expectations
func (s *AnalyticsService) GetDashboard(ctx context.Context, timeRange string, orgID string) (*DashboardStats, error) {
	// Parse time range
	startTime, endTime, err := s.parseTimeRange(timeRange)
	if err != nil {
		return nil, fmt.Errorf("invalid time range: %w", err)
	}

	duration := endTime.Sub(startTime)

	// Get current period stats
	totalTraces, err := s.getTotalTraces(ctx, orgID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	totalCost, err := s.getTotalCost(ctx, orgID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	totalTokens, err := s.getTotalTokens(ctx, orgID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	avgLatency, err := s.getAvgLatency(ctx, orgID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	errorRate, successRate, err := s.getErrorRate(ctx, orgID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// Get previous period for trends
	prevStart := startTime.Add(-duration)
	prevEnd := startTime

	prevTraces, _ := s.getTotalTraces(ctx, orgID, prevStart, prevEnd)
	prevCost, _ := s.getTotalCost(ctx, orgID, prevStart, prevEnd)
	prevTokens, _ := s.getTotalTokens(ctx, orgID, prevStart, prevEnd)
	prevLatency, _ := s.getAvgLatency(ctx, orgID, prevStart, prevEnd)

	// Calculate trends
	trends := TrendData{
		Traces:  calculatePercentChange(float64(prevTraces), float64(totalTraces)),
		Cost:    calculatePercentChange(prevCost, totalCost),
		Tokens:  calculatePercentChange(float64(prevTokens), float64(totalTokens)),
		Latency: calculatePercentChange(prevLatency, avgLatency),
	}

	// Get top models
	topModels, err := s.getTopModels(ctx, orgID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// Get cost by day
	costByDay, err := s.getCostByDay(ctx, orgID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// Get traces by status
	tracesByStatus, err := s.getTracesByStatus(ctx, orgID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	return &DashboardStats{
		TotalTraces:    totalTraces,
		TotalCost:      totalCost,
		TotalTokens:    totalTokens,
		AvgLatency:     avgLatency,
		ErrorRate:      errorRate,
		SuccessRate:    successRate,
		Trends:         trends,
		TopModels:      topModels,
		CostByDay:      costByDay,
		TracesByStatus: tracesByStatus,
	}, nil
}

// ============================================================================
// HELPER METHODS FOR DASHBOARD
// ============================================================================

func (s *AnalyticsService) getTotalTraces(ctx context.Context, orgID string, start, end time.Time) (int64, error) {
	query := `
		SELECT count() as total
		FROM llm_observability.traces
		WHERE organization_id = ?
		AND timestamp >= ?
		AND timestamp < ?
	`

	rows, err := s.repo.Query(ctx, query, orgID, start, end)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var total int64
	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return 0, err
		}
	}
	return total, nil
}

func (s *AnalyticsService) getTotalCost(ctx context.Context, orgID string, start, end time.Time) (float64, error) {
	query := `
		SELECT sum(total_cost_usd) as total_cost
		FROM llm_observability.traces
		WHERE organization_id = ?
		AND timestamp >= ?
		AND timestamp < ?
	`

	rows, err := s.repo.Query(ctx, query, orgID, start, end)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var totalCost float64
	if rows.Next() {
		if err := rows.Scan(&totalCost); err != nil {
			return 0, err
		}
	}
	return totalCost, nil
}

func (s *AnalyticsService) getTotalTokens(ctx context.Context, orgID string, start, end time.Time) (int64, error) {
	query := `
		SELECT sum(total_tokens) as total_tokens
		FROM llm_observability.traces
		WHERE organization_id = ?
		AND timestamp >= ?
		AND timestamp < ?
	`

	rows, err := s.repo.Query(ctx, query, orgID, start, end)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var totalTokens int64
	if rows.Next() {
		if err := rows.Scan(&totalTokens); err != nil {
			return 0, err
		}
	}
	return totalTokens, nil
}

// func (s *AnalyticsService) getAvgLatency(ctx context.Context, orgID string, start, end time.Time) (float64, error) {
// 	query := `
// 		SELECT avg(duration_ms) as avg_latency
// 		FROM llm_observability.traces
// 		WHERE organization_id = ?
// 		AND timestamp >= ?
// 		AND timestamp < ?
// 	`

// 	rows, err := s.repo.Query(ctx, query, orgID, start, end)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer rows.Close()

// 	var avgLatency float64
// 	if rows.Next() {
// 		if err := rows.Scan(&avgLatency); err != nil {
// 			return 0, err
// 		}
// 	}
// 	return avgLatency, nil
// }

func (s *AnalyticsService) getAvgLatency(ctx context.Context, orgID string, start, end time.Time) (float64, error) {
	query := `
		SELECT 
			count() as trace_count,
			avg(duration_ms) as avg_latency
		FROM llm_observability.traces
		WHERE organization_id = ?
		AND timestamp >= ?
		AND timestamp < ?
	`

	rows, err := s.repo.Query(ctx, query, orgID, start, end)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var traceCount int64
	var avgLatency float64
	if rows.Next() {
		if err := rows.Scan(&traceCount, &avgLatency); err != nil {
			return 0, err
		}
	}

	// If no traces, return 0
	if traceCount == 0 {
		return 0, nil
	}

	return avgLatency, nil
}

// func (s *AnalyticsService) getErrorRate(ctx context.Context, orgID string, start, end time.Time) (float64, float64, error) {
// 	query := `
// 		SELECT
// 			countIf(status = 'error') * 100.0 / count() as error_rate,
// 			countIf(status = 'success') * 100.0 / count() as success_rate
// 		FROM llm_observability.traces
// 		WHERE organization_id = ?
// 		AND timestamp >= ?
// 		AND timestamp < ?
// 	`

// 	rows, err := s.repo.Query(ctx, query, orgID, start, end)
// 	if err != nil {
// 		return 0, 0, err
// 	}
// 	defer rows.Close()

// 	var errorRate, successRate float64
// 	if rows.Next() {
// 		if err := rows.Scan(&errorRate, &successRate); err != nil {
// 			return 0, 0, err
// 		}
// 	}
// 	return errorRate, successRate, nil
// }

func (s *AnalyticsService) getErrorRate(ctx context.Context, orgID string, start, end time.Time) (float64, float64, error) {
	query := `
		SELECT 
			count() as total_count
		FROM llm_observability.traces
		WHERE organization_id = ?
		AND timestamp >= ?
		AND timestamp < ?
	`

	rows, err := s.repo.Query(ctx, query, orgID, start, end)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	var totalCount int64
	if rows.Next() {
		if err := rows.Scan(&totalCount); err != nil {
			return 0, 0, err
		}
	}

	// If no traces, return 0% for both
	if totalCount == 0 {
		return 0.0, 0.0, nil
	}

	// Now get error and success counts
	query2 := `
		SELECT 
			countIf(status = 'error') as errors,
			countIf(status = 'success') as successes
		FROM llm_observability.traces
		WHERE organization_id = ?
		AND timestamp >= ?
		AND timestamp < ?
	`

	rows2, err := s.repo.Query(ctx, query2, orgID, start, end)
	if err != nil {
		return 0, 0, err
	}
	defer rows2.Close()

	var errorCount, successCount int64
	if rows2.Next() {
		if err := rows2.Scan(&errorCount, &successCount); err != nil {
			return 0, 0, err
		}
	}

	errorRate := (float64(errorCount) / float64(totalCount)) * 100.0
	successRate := (float64(successCount) / float64(totalCount)) * 100.0

	return errorRate, successRate, nil
}

func calculatePercentChange(old, new float64) float64 {
	if old == 0 {
		if new > 0 {
			return 100.0
		}
		return 0.0
	}
	change := ((new - old) / old) * 100.0
	// Prevent NaN and Inf
	if change != change || change > 1000000 || change < -1000000 {
		return 0.0
	}
	return change
}

func (s *AnalyticsService) getTopModels(ctx context.Context, orgID string, start, end time.Time) ([]ModelStats, error) {
	query := `
		SELECT 
			model,
			count() as count,
			sum(total_cost_usd) as total_cost
		FROM llm_observability.traces
		WHERE organization_id = ?
		AND timestamp >= ?
		AND timestamp < ?
		GROUP BY model
		ORDER BY count DESC
		LIMIT 10
	`

	rows, err := s.repo.Query(ctx, query, orgID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []ModelStats
	for rows.Next() {
		var m ModelStats
		if err := rows.Scan(&m.Model, &m.Count, &m.Cost); err != nil {
			return nil, err
		}
		models = append(models, m)
	}

	return models, nil
}

func (s *AnalyticsService) getCostByDay(ctx context.Context, orgID string, start, end time.Time) ([]DailyCost, error) {
	query := `
		SELECT 
			toDate(timestamp) as date,
			sum(total_cost_usd) as cost
		FROM llm_observability.traces
		WHERE organization_id = ?
		AND timestamp >= ?
		AND timestamp < ?
		GROUP BY date
		ORDER BY date ASC
	`

	rows, err := s.repo.Query(ctx, query, orgID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var costs []DailyCost
	for rows.Next() {
		var c DailyCost
		var date time.Time
		if err := rows.Scan(&date, &c.Cost); err != nil {
			return nil, err
		}
		c.Date = date.Format("2006-01-02")
		costs = append(costs, c)
	}

	return costs, nil
}

func (s *AnalyticsService) getTracesByStatus(ctx context.Context, orgID string, start, end time.Time) ([]StatusCount, error) {
	query := `
		SELECT 
			status,
			count() as count
		FROM llm_observability.traces
		WHERE organization_id = ?
		AND timestamp >= ?
		AND timestamp < ?
		GROUP BY status
		ORDER BY count DESC
	`

	rows, err := s.repo.Query(ctx, query, orgID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []StatusCount
	for rows.Next() {
		var s StatusCount
		if err := rows.Scan(&s.Status, &s.Count); err != nil {
			return nil, err
		}
		statuses = append(statuses, s)
	}

	return statuses, nil
}

// ============================================================================
// TYPE DEFINITIONS
// ============================================================================

// DashboardSummary contains all data for the main dashboard
type DashboardSummary struct {
	TimeRange     string                  `json:"time_range"`
	StartTime     time.Time               `json:"start_time"`
	EndTime       time.Time               `json:"end_time"`
	MetricSummary *models.MetricSummary   `json:"metric_summary"`
	CostBreakdown []*models.CostBreakdown `json:"cost_breakdown"`
	ModelUsage    []*models.ModelUsage    `json:"model_usage"`
	Insights      []Insight               `json:"insights"`
}

// Insight represents an actionable insight
type Insight struct {
	Type        string `json:"type"`     // success, info, warning, error
	Category    string `json:"category"` // cost, performance, reliability
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"` // info, low, medium, high, critical
}

// CostAnalysis contains detailed cost information
type CostAnalysis struct {
	TotalCost          float64                 `json:"total_cost"`
	DailyAverage       float64                 `json:"daily_average"`
	MonthlyProjection  float64                 `json:"monthly_projection"`
	CostBreakdown      []*models.CostBreakdown `json:"cost_breakdown"`
	MostExpensiveModel string                  `json:"most_expensive_model"`
	HighestCost        float64                 `json:"highest_cost"`
}

// PerformanceMetrics contains performance analysis
type PerformanceMetrics struct {
	Summary        *models.MetricSummary `json:"summary"`
	Status         string                `json:"status"` // excellent, good, fair, poor
	Recommendation string                `json:"recommendation"`
}

// ModelComparison compares different models
type ModelComparison struct {
	Model               string  `json:"model"`
	TotalCalls          int64   `json:"total_calls"`
	TotalCost           float64 `json:"total_cost"`
	AvgCostPerRequest   float64 `json:"avg_cost_per_request"`
	AvgLatency          float64 `json:"avg_latency"`
	AvgTokensPerRequest float64 `json:"avg_tokens_per_request"`
	EfficiencyScore     float64 `json:"efficiency_score"` // Lower is better
}

// DashboardStats types for new dashboard API response
type DashboardStats struct {
	TotalTraces    int64         `json:"total_traces"`
	TotalCost      float64       `json:"total_cost"`
	TotalTokens    int64         `json:"total_tokens"`
	AvgLatency     float64       `json:"avg_latency"`
	ErrorRate      float64       `json:"error_rate"`
	SuccessRate    float64       `json:"success_rate"`
	Trends         TrendData     `json:"trends"`
	TopModels      []ModelStats  `json:"top_models"`
	CostByDay      []DailyCost   `json:"cost_by_day"`
	TracesByStatus []StatusCount `json:"traces_by_status"`
}

type TrendData struct {
	Traces  float64 `json:"traces"`
	Cost    float64 `json:"cost"`
	Tokens  float64 `json:"tokens"`
	Latency float64 `json:"latency"`
}

type ModelStats struct {
	Model string  `json:"model"`
	Count int64   `json:"count"`
	Cost  float64 `json:"cost"`
}

type DailyCost struct {
	Date string  `json:"date"`
	Cost float64 `json:"cost"`
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

package models

import "time"

// TraceQuery represents query parameters for listing traces
type TraceQuery struct {
    OrganizationID string    `json:"organization_id"`
    ProjectID      string    `json:"project_id,omitempty"`
    UserID         string    `json:"user_id,omitempty"`
    Model          string    `json:"model,omitempty"`
    Provider       string    `json:"provider,omitempty"`
    Status         string    `json:"status,omitempty"`
    StartTime      time.Time `json:"start_time"`
    EndTime        time.Time `json:"end_time"`
    Limit          int       `json:"limit"`
    Offset         int       `json:"offset"`
}

// Metric represents a single metric data point
type Metric struct {
    MetricName       string                 `json:"metric_name" ch:"metric_name"`
    MetricValue      float64                `json:"metric_value" ch:"metric_value"`
    Timestamp        time.Time              `json:"timestamp" ch:"timestamp"`
    OrganizationID   string                 `json:"organization_id" ch:"organization_id"`
    ProjectID        string                 `json:"project_id,omitempty" ch:"project_id"`
    Model            string                 `json:"model" ch:"model"`
    Provider         string                 `json:"provider" ch:"provider"`
    Tags             map[string]interface{} `json:"tags,omitempty"`
    RequestCount     int64                  `json:"request_count" ch:"request_count"`
    TotalTokens      int64                  `json:"total_tokens" ch:"total_tokens"`
    PromptTokens     int64                  `json:"prompt_tokens" ch:"prompt_tokens"`
    CompletionTokens int64                  `json:"completion_tokens" ch:"completion_tokens"`
    TotalCost        float64                `json:"total_cost" ch:"total_cost"`
    AvgLatencyMs     float64                `json:"avg_latency_ms" ch:"avg_latency_ms"`
    P50LatencyMs     float64                `json:"p50_latency_ms" ch:"p50_latency_ms"`
    P95LatencyMs     float64                `json:"p95_latency_ms" ch:"p95_latency_ms"`
    P99LatencyMs     float64                `json:"p99_latency_ms" ch:"p99_latency_ms"`
    ErrorCount       int64                  `json:"error_count" ch:"error_count"`
    ErrorRate        float64                `json:"error_rate" ch:"error_rate"`
    Metadata         map[string]interface{} `json:"metadata,omitempty" ch:"metadata"`
}

// CostBreakdown represents cost analysis breakdown
type CostBreakdown struct {
    Model          string               `json:"model,omitempty"`
    TotalCost      float64              `json:"total_cost"`
    TotalCalls     int64                `json:"total_calls"`
    AvgCost        float64              `json:"avg_cost"`
    Percentage     float64              `json:"percentage"`
    ByModel        []ModelCost          `json:"by_model"`
    ByProvider     []ProviderCost       `json:"by_provider"`
    ByProject      []ProjectCost        `json:"by_project"`
    DailyCosts     []DailyCost          `json:"daily_costs"`
    TopExpensive   []ExpensiveTrace     `json:"top_expensive"`
}

// ModelCost represents cost breakdown by model
type ModelCost struct {
    Model         string  `json:"model"`
    Provider      string  `json:"provider"`
    TotalCost     float64 `json:"total_cost"`
    RequestCount  int64   `json:"request_count"`
    AvgCostPerReq float64 `json:"avg_cost_per_request"`
    TotalTokens   int64   `json:"total_tokens"`
}

// ProjectCost represents cost breakdown by project
type ProjectCost struct {
    ProjectID    string  `json:"project_id"`
    ProjectName  string  `json:"project_name"`
    TotalCost    float64 `json:"total_cost"`
    RequestCount int64   `json:"request_count"`
}

// DailyCost represents daily cost data
type DailyCost struct {
    Date         string  `json:"date"`
    TotalCost    float64 `json:"total_cost"`
    RequestCount int64   `json:"request_count"`
}

// ExpensiveTrace represents an expensive trace
type ExpensiveTrace struct {
    TraceID   string  `json:"trace_id"`
    Model     string  `json:"model"`
    Provider  string  `json:"provider"`
    Cost      float64 `json:"cost"`
    Tokens    int64   `json:"tokens"`
    Timestamp string  `json:"timestamp"`
}

// MetricQuery represents query parameters for metrics
type MetricQuery struct {
    MetricName     string    `json:"metric_name,omitempty"`
    OrganizationID string    `json:"organization_id"`
    ProjectID      string    `json:"project_id,omitempty"`
    Model          string    `json:"model,omitempty"`
    Provider       string    `json:"provider,omitempty"`
    StartTime      time.Time `json:"start_time"`
    EndTime        time.Time `json:"end_time"`
    Granularity    string    `json:"granularity"` // hour, day, week
    Limit          int       `json:"limit"`
}

// MetricSummary represents aggregated metric summary
type MetricSummary struct {
    Period            string               `json:"period"`
    TotalRequests     int64                `json:"total_requests"`
    TotalTokens       int64                `json:"total_tokens"`
    TotalCost         float64              `json:"total_cost"`
    TotalCostUSD      float64              `json:"total_cost_usd"`
    AvgCostPerRequest float64              `json:"avg_cost_per_request"`
    AvgLatencyMs      float64              `json:"avg_latency_ms"`
    P50LatencyMs      float64              `json:"p50_latency_ms"`
    P95LatencyMs      float64              `json:"p95_latency_ms"`
    P99LatencyMs      float64              `json:"p99_latency_ms"`
    ErrorRate         float64              `json:"error_rate"`
    SuccessRate       float64              `json:"success_rate"`
    TopModels         []ModelUsage         `json:"top_models"`
    ByProvider        []ProviderMetrics    `json:"by_provider"`
    TimeSeries        []TimeSeriesPoint    `json:"time_series"`
}

// ProviderMetrics represents metrics grouped by provider
type ProviderMetrics struct {
    Provider      string  `json:"provider"`
    RequestCount  int64   `json:"request_count"`
    TotalCost     float64 `json:"total_cost"`
    AvgLatencyMs  float64 `json:"avg_latency_ms"`
    ErrorRate     float64 `json:"error_rate"`
}

// TimeSeriesPoint represents a single point in time series data
type TimeSeriesPoint struct {
    Timestamp    time.Time `json:"timestamp"`
    RequestCount int64     `json:"request_count"`
    TotalCost    float64   `json:"total_cost"`
    AvgLatency   float64   `json:"avg_latency"`
    ErrorRate    float64   `json:"error_rate"`
}

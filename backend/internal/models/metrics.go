package models

import "time"

type MetricsSummary struct {
    TotalTraces      int64            `json:"total_traces"`
    TotalSpans       int64            `json:"total_spans"`
    TotalTokens      int64            `json:"total_tokens"`
    TotalCost        float64          `json:"total_cost"`
    AvgLatencyMs     float64          `json:"avg_latency_ms"`
    P95LatencyMs     float64          `json:"p95_latency_ms"`
    P99LatencyMs     float64          `json:"p99_latency_ms"`
    ErrorRate        float64          `json:"error_rate"`
    SuccessRate      float64          `json:"success_rate"`
    TopModels        []ModelUsage     `json:"top_models"`
    CostByProvider   []ProviderCost   `json:"cost_by_provider"`
}

type ModelUsage struct {
    Model        string  `json:"model"`
    Provider     string  `json:"provider"`
    Count        int64   `json:"count"`
    CallCount    int64   `json:"call_count"`
    TotalCost    float64 `json:"total_cost"`
    AvgTokens    float64 `json:"avg_tokens"`
    TotalTokens  int64   `json:"total_tokens"`
    AvgLatency   float64 `json:"avg_latency"`
}

type ProviderCost struct {
    Provider  string  `json:"provider"`
    TotalCost float64 `json:"total_cost"`
    Count     int64   `json:"count"`
}

type MetricsQuery struct {
    OrganizationID string    `json:"organization_id"`
    ProjectID      string    `json:"project_id,omitempty"`
    StartTime      time.Time `json:"start_time"`
    EndTime        time.Time `json:"end_time"`
    Model          string    `json:"model,omitempty"`
    Provider       string    `json:"provider,omitempty"`
}

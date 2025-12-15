package models

import "time"

type Trace struct {
    TraceID        string                 `json:"trace_id" ch:"trace_id"`
    OrganizationID string                 `json:"organization_id" ch:"organization_id"`
    ProjectID      string                 `json:"project_id" ch:"project_id"`
    TraceType      string                 `json:"trace_type" ch:"trace_type"`
    Model          string                 `json:"model" ch:"model"`
    Provider       string                 `json:"provider" ch:"provider"`
    UserID         string                 `json:"user_id,omitempty" ch:"user_id"`
    Status         string                 `json:"status" ch:"status"`
    TotalTokens    int                    `json:"total_tokens" ch:"total_tokens"`
    TotalCost      float64                `json:"total_cost" ch:"total_cost"`
    TotalCostUSD   float64                `json:"total_cost_usd" ch:"total_cost_usd"`
    DurationMs     int64                  `json:"duration_ms" ch:"duration_ms"`
    Metadata       map[string]string      `json:"metadata,omitempty"`
    Tags           map[string]string      `json:"tags,omitempty"`
    Timestamp      time.Time              `json:"timestamp" ch:"timestamp"`
    CreatedAt      time.Time              `json:"created_at" ch:"created_at"`
    Spans          []Span                 `json:"spans,omitempty"`
}

type Span struct {
    SpanID           string            `json:"span_id" ch:"span_id"`
    TraceID          string            `json:"trace_id" ch:"trace_id"`
    ParentSpanID     string            `json:"parent_span_id,omitempty" ch:"parent_span_id"`
    Name             string            `json:"name" ch:"name"`
    Model            string            `json:"model" ch:"model"`
    Provider         string            `json:"provider" ch:"provider"`
    Input            string            `json:"input" ch:"input"`
    Output           string            `json:"output" ch:"output"`
    PromptTokens     int               `json:"prompt_tokens" ch:"prompt_tokens"`
    CompletionTokens int               `json:"completion_tokens" ch:"completion_tokens"`
    TotalTokens      int               `json:"total_tokens" ch:"total_tokens"`
    Cost             float64           `json:"cost" ch:"cost"`
    CostUSD          float64           `json:"cost_usd" ch:"cost_usd"`
    DurationMs       int64             `json:"duration_ms" ch:"duration_ms"`
    StartTime        time.Time         `json:"start_time" ch:"start_time"`
    EndTime          time.Time         `json:"end_time" ch:"end_time"`
    Status           string            `json:"status" ch:"status"`
    ErrorMessage     string            `json:"error_message,omitempty" ch:"error_message"`
    Metadata         map[string]string `json:"metadata,omitempty"`
    Tags             map[string]string `json:"tags,omitempty"`
    CreatedAt        time.Time         `json:"created_at" ch:"created_at"`
}

type CreateTraceRequest struct {
    OrganizationID string            `json:"organization_id" validate:"required"`
    ProjectID      string            `json:"project_id" validate:"required"`
    TraceType      string            `json:"trace_type" validate:"required"`
    Model          string            `json:"model" validate:"required"`
    Provider       string            `json:"provider" validate:"required"`
    UserID         string            `json:"user_id,omitempty"`
    Metadata       map[string]string `json:"metadata,omitempty"`
    Spans          []SpanData        `json:"spans" validate:"required,min=1"`
}

type SpanData struct {
    Name             string            `json:"name" validate:"required"`
    ParentSpanID     string            `json:"parent_span_id,omitempty"`
    Model            string            `json:"model" validate:"required"`
    Provider         string            `json:"provider" validate:"required"`
    Input            string            `json:"input" validate:"required"`
    Output           string            `json:"output" validate:"required"`
    PromptTokens     int               `json:"prompt_tokens" validate:"min=0"`
    CompletionTokens int               `json:"completion_tokens" validate:"min=0"`
    DurationMs       int64             `json:"duration_ms" validate:"min=0"`
    Status           string            `json:"status" validate:"required"`
    ErrorMessage     string            `json:"error_message,omitempty"`
    Metadata         map[string]string `json:"metadata,omitempty"`
}

// TraceRequest is for creating traces via API
type TraceRequest struct {
    OrganizationID   string            `json:"organization_id" validate:"required"`
    ProjectID        string            `json:"project_id,omitempty"`
    Model            string            `json:"model" validate:"required"`
    Provider         string            `json:"provider" validate:"required"`
    TraceType        string            `json:"trace_type,omitempty"`
    UserID           string            `json:"user_id,omitempty"`
    Input            string            `json:"input,omitempty"`
    Output           string            `json:"output,omitempty"`
    PromptTokens     int               `json:"prompt_tokens,omitempty"`
    CompletionTokens int               `json:"completion_tokens,omitempty"`
    Latency          int64             `json:"latency,omitempty"`
    Status           string            `json:"status,omitempty"`
    ErrorMessage     string            `json:"error_message,omitempty"`
    Tags             map[string]string `json:"tags,omitempty"`
    Metadata         map[string]string `json:"metadata,omitempty"`
    Spans            []SpanRequest     `json:"spans,omitempty"`
}

// SpanRequest represents a span in TraceRequest
type SpanRequest struct {
    Name             string            `json:"name" validate:"required"`
    ParentSpanID     string            `json:"parent_span_id,omitempty"`
    Model            string            `json:"model" validate:"required"`
    Provider         string            `json:"provider" validate:"required"`
    Input            string            `json:"input" validate:"required"`
    Output           string            `json:"output" validate:"required"`
    PromptTokens     int               `json:"prompt_tokens" validate:"min=0"`
    CompletionTokens int               `json:"completion_tokens" validate:"min=0"`
    DurationMs       int64             `json:"duration_ms" validate:"min=0"`
    Status           string            `json:"status" validate:"required"`
    ErrorMessage     string            `json:"error_message,omitempty"`
    Tags             map[string]string `json:"tags,omitempty"`
    Metadata         map[string]string `json:"metadata,omitempty"`
}

// TraceResponse is returned after creating a trace
type TraceResponse struct {
    TraceID        string    `json:"trace_id"`
    OrganizationID string    `json:"organization_id"`
    ProjectID      string    `json:"project_id,omitempty"`
    Model          string    `json:"model"`
    Provider       string    `json:"provider"`
    Status         string    `json:"status"`
    TotalTokens    int       `json:"total_tokens"`
    TotalCost      float64   `json:"total_cost"`
    DurationMs     int64     `json:"duration_ms"`
    Timestamp      string    `json:"timestamp"`
    CreatedAt      string    `json:"created_at"`
    Message        string    `json:"message,omitempty"`
}

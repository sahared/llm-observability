package models

// ErrorResponse represents an error response
type ErrorResponse struct {
    Error     string                 `json:"error"`
    Message   string                 `json:"message,omitempty"`
    Details   map[string]interface{} `json:"details,omitempty"`
    RequestID string                 `json:"request_id,omitempty"`
    Timestamp string                 `json:"timestamp"`
}

// HealthResponse represents health check response
type HealthResponse struct {
    Status    string            `json:"status"`
    Version   string            `json:"version,omitempty"`
    Timestamp string            `json:"timestamp"`
    Services  map[string]string `json:"services,omitempty"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    Total      int64       `json:"total"`
    TotalCount int64       `json:"total_count"`
    Page       int         `json:"page"`
    PageSize   int         `json:"page_size"`
    TotalPages int         `json:"total_pages"`
}

// BatchTraceRequest represents a batch trace creation request
type BatchTraceRequest struct {
    Traces []TraceRequest `json:"traces" validate:"required,min=1,max=100"`
}

// BatchTraceResponse represents response for batch trace creation
type BatchTraceResponse struct {
    Success  int              `json:"success"`
    Failed   int              `json:"failed"`
    Accepted int              `json:"accepted"`
    Rejected int              `json:"rejected"`
    Traces   []TraceResponse  `json:"traces,omitempty"`
    Errors   []string         `json:"errors,omitempty"`
    Message  string           `json:"message"`
}

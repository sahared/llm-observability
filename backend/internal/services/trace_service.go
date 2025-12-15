package services

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/Aditya-Pimpalkar/clarity/internal/models"
    "github.com/Aditya-Pimpalkar/clarity/internal/repository"
    "github.com/Aditya-Pimpalkar/clarity/internal/kafka"
)

type TraceService struct {
    repo     repository.Repository
    producer *kafka.Producer
}

func NewTraceService(repo repository.Repository, producer *kafka.Producer) *TraceService {
    return &TraceService{
        repo:     repo,
        producer: producer,
    }
}

func (s *TraceService) CreateTrace(ctx context.Context, req *models.TraceRequest) (*models.TraceResponse, error) {
    // Validate request
    if err := s.validateTraceRequest(req); err != nil {
        return nil, err
    }

    // Generate trace ID
    traceID := uuid.New().String()
    now := time.Now()

    // Process spans
    var spans []models.Span
    var totalTokens int
    var totalDuration int64
    var totalCost float64

    for i, spanReq := range req.Spans {
        spanID := uuid.New().String()
        
        // Calculate span timing
        startTime := now.Add(time.Duration(i) * time.Millisecond * 100)
        endTime := startTime.Add(time.Duration(spanReq.DurationMs) * time.Millisecond)

        // Calculate cost for this span
        cost := s.calculateCost(spanReq.Model, spanReq.Provider, int(spanReq.PromptTokens), int(spanReq.CompletionTokens))

        span := models.Span{
            SpanID:           spanID,
            TraceID:          traceID,
            ParentSpanID:     spanReq.ParentSpanID,
            Name:             spanReq.Name,
            StartTime:        startTime,
            EndTime:          endTime,
            DurationMs:       int64(spanReq.DurationMs),
            Model:            spanReq.Model,
            Provider:         spanReq.Provider,
            Input:            spanReq.Input,
            Output:           spanReq.Output,
            PromptTokens:     int(spanReq.PromptTokens),
            CompletionTokens: int(spanReq.CompletionTokens),
            TotalTokens:      int(spanReq.PromptTokens + spanReq.CompletionTokens),
            CostUSD:          cost,
            Status:           spanReq.Status,
            Metadata:         spanReq.Tags,
        }

        spans = append(spans, span)
        totalTokens += span.TotalTokens
        totalDuration += span.DurationMs
        totalCost += cost

        // Publish span created event to Kafka
        if s.producer != nil {
            _ = s.producer.PublishSpanCreated(ctx, spanID, traceID, span.DurationMs, span.TotalTokens)
        }
    }

    // Determine overall status
    status := s.determineTraceStatus(spans)

    // Create trace
    trace := &models.Trace{
        TraceID:        traceID,
        OrganizationID: req.OrganizationID,
        ProjectID:      req.ProjectID,
        Timestamp:      now,
        TraceType:      req.TraceType,
        DurationMs:     totalDuration,
        Status:         status,
        TotalCostUSD:   totalCost,
        TotalTokens:    totalTokens,
        Model:          req.Model,
        Provider:       req.Provider,
        UserID:         req.UserID,
        Metadata:       req.Metadata,
        Spans:          spans,
    }

    // Save to database
    if err := s.repo.SaveTrace(ctx, trace); err != nil {
        return nil, fmt.Errorf("failed to save trace: %w", err)
    }

    // Publish trace created event to Kafka
    if s.producer != nil {
        _ = s.producer.PublishTraceCreated(
            ctx,
            traceID,
            req.OrganizationID,
            req.ProjectID,
            req.Model,
            req.Provider,
            totalTokens,
            totalCost,
        )
    }

    return &models.TraceResponse{
        TraceID:        traceID,
        OrganizationID: req.OrganizationID,
        Model:          req.Model,
        Provider:       req.Provider,
        Status:         "accepted",
        TotalTokens:    totalTokens,
        TotalCost:      totalCost,
        DurationMs:     totalDuration,
        Timestamp:      now.Format(time.RFC3339),
        CreatedAt:      now.Format(time.RFC3339),
        Message:        "Trace created successfully",
    }, nil
}

func (s *TraceService) validateTraceRequest(req *models.TraceRequest) error {
    if req.OrganizationID == "" {
        return fmt.Errorf("organization_id is required")
    }
    if len(req.Spans) == 0 {
        return fmt.Errorf("at least one span is required")
    }
    if req.TraceType != "single_call" && req.TraceType != "multi_step" && req.TraceType != "streaming" {
        return fmt.Errorf("invalid trace_type: must be single_call, multi_step, or streaming")
    }
    return nil
}

func (s *TraceService) calculateCost(model, provider string, promptTokens, completionTokens int) float64 {
    // Pricing per 1M tokens
    pricing := map[string]map[string][2]float64{
        "openai": {
            "gpt-4":          {30.0, 60.0},   // prompt, completion per 1M tokens
            "gpt-4-turbo":    {10.0, 30.0},
            "gpt-3.5-turbo":  {0.5, 1.5},
        },
        "anthropic": {
            "claude-3-opus":   {15.0, 75.0},
            "claude-3-sonnet": {3.0, 15.0},
            "claude-3-haiku":  {0.25, 1.25},
        },
    }

    rates, ok := pricing[provider][model]
    if !ok {
        return 0.0
    }

    promptCost := (float64(promptTokens) / 1000000.0) * rates[0]
    completionCost := (float64(completionTokens) / 1000000.0) * rates[1]
    
    return promptCost + completionCost
}

func (s *TraceService) determineTraceStatus(spans []models.Span) string {
    if len(spans) == 0 {
        return "unknown"
    }

    for _, span := range spans {
        if span.Status == "error" || span.Status == "failed" {
            return "error"
        }
        if span.Status == "timeout" {
            return "timeout"
        }
    }

    return "success"
}

// GetTrace retrieves a single trace by ID
func (s *TraceService) GetTrace(ctx context.Context, traceID string) (*models.Trace, error) {
    return s.repo.GetTraceByID(ctx, traceID)
}

// GetTraces retrieves multiple traces with filtering
func (s *TraceService) GetTraces(ctx context.Context, query *models.TraceQuery) ([]*models.Trace, int64, error) {
    traces, err := s.repo.GetTraces(ctx, query)
    if err != nil {
        return nil, 0, err
    }

    count, err := s.repo.GetTraceCount(ctx, query)
    if err != nil {
        return traces, 0, err
    }

    return traces, count, nil
}

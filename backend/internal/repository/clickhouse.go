package repository

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "time"

    "github.com/ClickHouse/clickhouse-go/v2"
    "github.com/Aditya-Pimpalkar/clarity/internal/models"
)

type ClickHouseRepository struct {
    conn clickhouse.Conn
}

func NewClickHouseRepository(addr string) (*ClickHouseRepository, error) {
    conn, err := clickhouse.Open(&clickhouse.Options{
        Addr: []string{addr},
        Auth: clickhouse.Auth{
            Database: "llm_observability",
            Username: "default",
            Password: "",
        },
        Settings: clickhouse.Settings{
            "max_execution_time": 60,
        },
    })
    if err != nil {
        return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
    }

    if err := conn.Ping(context.Background()); err != nil {
        return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
    }

    return &ClickHouseRepository{conn: conn}, nil
}

// SaveTrace stores a trace in ClickHouse
func (r *ClickHouseRepository) SaveTrace(ctx context.Context, trace *models.Trace) error {
    metadataJSON := "{}"
    if trace.Metadata != nil {
        if data, err := json.Marshal(trace.Metadata); err == nil {
            metadataJSON = string(data)
        }
    }

    query := `
        INSERT INTO traces (
            trace_id, organization_id, project_id, timestamp, 
            trace_type, duration_ms, status, total_cost_usd, 
            total_tokens, model, provider, user_id, metadata
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

    err := r.conn.Exec(ctx, query,
        trace.TraceID,
        trace.OrganizationID,
        trace.ProjectID,
        trace.Timestamp,
        trace.TraceType,
        trace.DurationMs,
        trace.Status,
        trace.TotalCostUSD,
        trace.TotalTokens,
        trace.Model,
        trace.Provider,
        trace.UserID,
        metadataJSON,
    )

    if err != nil {
        return fmt.Errorf("failed to insert trace: %w", err)
    }

    // Save spans
    return r.SaveSpans(ctx, trace.Spans)
}

// SaveSpans stores spans in ClickHouse
func (r *ClickHouseRepository) SaveSpans(ctx context.Context, spans []models.Span) error {
    if len(spans) == 0 {
        return nil
    }

    batch, err := r.conn.PrepareBatch(ctx, `
        INSERT INTO spans (
            span_id, trace_id, parent_span_id, name, start_time, end_time,
            duration_ms, model, provider, input, output, prompt_tokens,
            completion_tokens, total_tokens, cost_usd, metadata
        )
    `)
    if err != nil {
        return fmt.Errorf("failed to prepare batch: %w", err)
    }

    for _, span := range spans {
        metadataJSON := "{}"
        if span.Metadata != nil {
            if data, err := json.Marshal(span.Metadata); err == nil {
                metadataJSON = string(data)
            }
        }

        err := batch.Append(
            span.SpanID,
            span.TraceID,
            span.ParentSpanID,
            span.Name,
            span.StartTime,
            span.EndTime,
            span.DurationMs,
            span.Model,
            span.Provider,
            span.Input,
            span.Output,
            span.PromptTokens,
            span.CompletionTokens,
            span.TotalTokens,
            span.CostUSD,
            metadataJSON,
        )
        if err != nil {
            return fmt.Errorf("failed to append span: %w", err)
        }
    }

    return batch.Send()
}


// GetSpansByTraceID retrieves all spans for a trace
func (r *ClickHouseRepository) GetSpansByTraceID(ctx context.Context, traceID string) ([]models.Span, error) {
    query := `
        SELECT 
            span_id, trace_id, parent_span_id, name, start_time, end_time,
            duration_ms, model, provider, input, output, prompt_tokens,
            completion_tokens, total_tokens, cost_usd, metadata
        FROM spans
        WHERE trace_id = ?
        ORDER BY start_time ASC
    `

    rows, err := r.conn.Query(ctx, query, traceID)
    if err != nil {
        return nil, fmt.Errorf("failed to query spans: %w", err)
    }
    defer rows.Close()

    var spans []models.Span
    for rows.Next() {
        var span models.Span
        var metadataJSON string

        err := rows.Scan(
            &span.SpanID,
            &span.TraceID,
            &span.ParentSpanID,
            &span.Name,
            &span.StartTime,
            &span.EndTime,
            &span.DurationMs,
            &span.Model,
            &span.Provider,
            &span.Input,
            &span.Output,
            &span.PromptTokens,
            &span.CompletionTokens,
            &span.TotalTokens,
            &span.CostUSD,
            &metadataJSON,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan span: %w", err)
        }

        if metadataJSON != "" && metadataJSON != "{}" {
            json.Unmarshal([]byte(metadataJSON), &span.Metadata)
        }

        spans = append(spans, span)
    }

    return spans, nil
}

// GetTraces retrieves traces with filtering


// SaveMetric stores a metric
func (r *ClickHouseRepository) SaveMetric(ctx context.Context, metric *models.Metric) error {
    tagsJSON := "{}"
    if metric.Tags != nil {
        if data, err := json.Marshal(metric.Tags); err == nil {
            tagsJSON = string(data)
        }
    }

    query := `
        INSERT INTO metrics (
            timestamp, organization_id, project_id, 
            metric_name, metric_value, tags
        ) VALUES (?, ?, ?, ?, ?, ?)
    `

    return r.conn.Exec(ctx, query,
        metric.Timestamp,
        metric.OrganizationID,
        metric.ProjectID,
        metric.MetricName,
        metric.MetricValue,
        tagsJSON,
    )
}

// GetMetricSummary retrieves aggregated metrics

// Close closes the ClickHouse connection
func (r *ClickHouseRepository) Close() error {
    if r.conn != nil {
        return r.conn.Close()
    }
    return nil
}

// Ping checks if ClickHouse is reachable
func (r *ClickHouseRepository) Ping(ctx context.Context) error {
    return r.conn.Ping(ctx)
}

// CreateOrganization - stub for now (Phase 2)
func (r *ClickHouseRepository) CreateOrganization(ctx context.Context, org *models.Organization) error {
    return fmt.Errorf("not implemented yet - Phase 2 feature")
}

// GetOrganization - stub for now (Phase 2)
func (r *ClickHouseRepository) GetOrganization(ctx context.Context, id string) (*models.Organization, error) {
    return nil, fmt.Errorf("not implemented yet - Phase 2 feature")
}

// CreateProject - stub for now (Phase 2)
func (r *ClickHouseRepository) CreateProject(ctx context.Context, project *models.Project) error {
    return fmt.Errorf("not implemented yet - Phase 2 feature")
}

// GetProject - stub for now (Phase 2)
func (r *ClickHouseRepository) GetProject(ctx context.Context, id string) (*models.Project, error) {
    return nil, fmt.Errorf("not implemented yet - Phase 2 feature")
}

// CreateUser - stub for now (Phase 2)
func (r *ClickHouseRepository) CreateUser(ctx context.Context, user *models.User) error {
    return fmt.Errorf("not implemented yet - Phase 2 feature")
}

// GetUserByEmail - stub for now (Phase 2)
func (r *ClickHouseRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
    return nil, fmt.Errorf("not implemented yet - Phase 2 feature")
}

// CreateAPIKey - stub for now (Phase 2)
func (r *ClickHouseRepository) CreateAPIKey(ctx context.Context, apiKey *models.APIKey) error {
    return fmt.Errorf("not implemented yet - Phase 2 feature")
}

// GetAPIKey - stub for now (Phase 2)


// GetCostBreakdown - stub for now (Phase 4)
func (r *ClickHouseRepository) GetCostBreakdown(ctx context.Context, orgID, projectID string, startTime, endTime time.Time) ([]*models.CostBreakdown, error) {
    return nil, fmt.Errorf("not implemented yet - Phase 4 feature")
}

// GetMetricSummary retrieves aggregated metrics
func (r *ClickHouseRepository) GetMetricSummary(ctx context.Context, orgID, projectID string, startTime, endTime time.Time) (*models.MetricSummary, error) {
    period := fmt.Sprintf("%s to %s", startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))

    var totalRequests int64
    err := r.conn.QueryRow(ctx, `
        SELECT count() FROM traces 
        WHERE organization_id = ? 
        AND timestamp >= ? AND timestamp <= ?
    `, orgID, startTime, endTime).Scan(&totalRequests)
    
    if err != nil {
        totalRequests = 0
    }

    var totalTokens int64
    var totalCost float64
    err = r.conn.QueryRow(ctx, `
        SELECT 
            sum(total_tokens) as total_tokens,
            sum(total_cost_usd) as total_cost
        FROM traces 
        WHERE organization_id = ? 
        AND timestamp >= ? AND timestamp <= ?
    `, orgID, startTime, endTime).Scan(&totalTokens, &totalCost)
    
    if err != nil {
        totalTokens = 0
        totalCost = 0
    }

    var avgLatency, p50, p95, p99 float64
    err = r.conn.QueryRow(ctx, `
        SELECT 
            avg(duration_ms),
            quantile(0.50)(duration_ms),
            quantile(0.95)(duration_ms),
            quantile(0.99)(duration_ms)
        FROM traces 
        WHERE organization_id = ? 
        AND timestamp >= ? AND timestamp <= ?
    `, orgID, startTime, endTime).Scan(&avgLatency, &p50, &p95, &p99)
    
    if err != nil {
        avgLatency, p50, p95, p99 = 0, 0, 0, 0
    }

    var successCount int64
    r.conn.QueryRow(ctx, `
        SELECT count() FROM traces 
        WHERE organization_id = ? 
        AND timestamp >= ? AND timestamp <= ?
        AND status = 'success'
    `, orgID, startTime, endTime).Scan(&successCount)

    successRate := float64(0)
    errorRate := float64(0)
    if totalRequests > 0 {
        successRate = (float64(successCount) / float64(totalRequests)) * 100
        errorRate = 100 - successRate
    }

    avgCostPerReq := float64(0)
    if totalRequests > 0 {
        avgCostPerReq = totalCost / float64(totalRequests)
    }

    return &models.MetricSummary{
        Period:            period,
        TotalRequests:     totalRequests,
        TotalTokens:       totalTokens,
        TotalCost:         totalCost,
        TotalCostUSD:      totalCost,
        AvgCostPerRequest: avgCostPerReq,
        AvgLatencyMs:      avgLatency,
        P50LatencyMs:      p50,
        P95LatencyMs:      p95,
        P99LatencyMs:      p99,
        ErrorRate:         errorRate,
        SuccessRate:       successRate,
        TopModels:         []models.ModelUsage{},
        ByProvider:        []models.ProviderMetrics{},
        TimeSeries:        []models.TimeSeriesPoint{},
    }, nil
}

// GetMetrics - stub for now (Phase 4)
func (r *ClickHouseRepository) GetMetrics(ctx context.Context, query *models.MetricQuery) ([]*models.Metric, error) {
    return nil, fmt.Errorf("not implemented yet - Phase 4 feature")
}

// GetModelUsage - stub for now (Phase 4)
func (r *ClickHouseRepository) GetModelUsage(ctx context.Context, orgID, projectID string, startTime, endTime time.Time) ([]*models.ModelUsage, error) {
    return nil, fmt.Errorf("not implemented yet - Phase 4 feature")
}

// GetProviderMetrics - stub for now (Phase 4)
func (r *ClickHouseRepository) GetProviderMetrics(ctx context.Context, orgID, projectID string, startTime, endTime time.Time) ([]*models.ProviderMetrics, error) {
    return nil, fmt.Errorf("not implemented yet - Phase 4 feature")
}

// GetTimeSeries - stub for now (Phase 4)
func (r *ClickHouseRepository) GetTimeSeries(ctx context.Context, orgID, projectID string, startTime, endTime time.Time, interval string) ([]*models.TimeSeriesPoint, error) {
    return nil, fmt.Errorf("not implemented yet - Phase 4 feature")
}

// SaveSpan - single span save (we have SaveSpans for batch)
func (r *ClickHouseRepository) SaveSpan(ctx context.Context, span *models.Span) error {
    return r.SaveSpans(ctx, []models.Span{*span})
}

// GetUserByID - stub for now (Phase 2)
func (r *ClickHouseRepository) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
    return nil, fmt.Errorf("not implemented yet - Phase 2 feature")
}

// GetProjectsByOrg - stub for now (Phase 2)
func (r *ClickHouseRepository) GetProjectsByOrg(ctx context.Context, orgID string) ([]*models.Project, error) {
    return nil, fmt.Errorf("not implemented yet - Phase 2 feature")
}

// Query - raw query execution
func (r *ClickHouseRepository) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    return nil, fmt.Errorf("raw queries not supported - use specific methods")
}
// GetAPIKey - simple demo implementation (Phase 2 will add real DB)
func (r *ClickHouseRepository) GetAPIKey(ctx context.Context, key string) (*models.APIKey, error) {
    // For development - accept "demo-key"
    if key == "demo-key" {
        return &models.APIKey{
            ID:             "demo-key-id",
            Key:            "demo-key",
            OrganizationID: "org-demo",
            ProjectID:      "proj-demo",
            Name:           "Demo API Key",
            IsActive:       true,
            CreatedAt:      time.Now(),
        }, nil
    }
    return nil, fmt.Errorf("invalid API key")
}
// GetTraces retrieves traces with filtering
func (r *ClickHouseRepository) GetTraces(ctx context.Context, query *models.TraceQuery) ([]*models.Trace, error) {
    sql := `
        SELECT 
            trace_id, organization_id, project_id, timestamp,
            trace_type, duration_ms, status, total_cost_usd,
            total_tokens, model, provider, user_id
        FROM traces
        WHERE organization_id = ?
    `

    args := []interface{}{query.OrganizationID}

    if query.ProjectID != "" {
        sql += " AND project_id = ?"
        args = append(args, query.ProjectID)
    }

    if !query.StartTime.IsZero() {
        sql += " AND timestamp >= ?"
        args = append(args, query.StartTime)
    }

    if !query.EndTime.IsZero() {
        sql += " AND timestamp <= ?"
        args = append(args, query.EndTime)
    }

    if query.Status != "" {
        sql += " AND status = ?"
        args = append(args, query.Status)
    }

    sql += " ORDER BY timestamp DESC LIMIT ? OFFSET ?"
    args = append(args, query.Limit, query.Offset)

    rows, err := r.conn.Query(ctx, sql, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to query traces: %w", err)
    }
    defer rows.Close()

    var traces []*models.Trace
    for rows.Next() {
        var trace models.Trace
        var durationMs uint32
        var totalTokens uint32
        
        err := rows.Scan(
            &trace.TraceID,
            &trace.OrganizationID,
            &trace.ProjectID,
            &trace.Timestamp,
            &trace.TraceType,
            &durationMs,
            &trace.Status,
            &trace.TotalCostUSD,
            &totalTokens,
            &trace.Model,
            &trace.Provider,
            &trace.UserID,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan trace: %w", err)
        }
        
        // Convert types
        trace.DurationMs = int64(durationMs)
        trace.TotalTokens = int(totalTokens)
        
        traces = append(traces, &trace)
    }

    return traces, nil
}

// GetTraceByID retrieves a trace by ID
func (r *ClickHouseRepository) GetTraceByID(ctx context.Context, traceID string) (*models.Trace, error) {
    var trace models.Trace
    var metadataJSON string
    var durationMs uint32
    var totalTokens uint32

    query := `
        SELECT 
            trace_id, organization_id, project_id, timestamp,
            trace_type, duration_ms, status, total_cost_usd,
            total_tokens, model, provider, user_id, metadata
        FROM traces
        WHERE trace_id = ?
    `

    err := r.conn.QueryRow(ctx, query, traceID).Scan(
        &trace.TraceID,
        &trace.OrganizationID,
        &trace.ProjectID,
        &trace.Timestamp,
        &trace.TraceType,
        &durationMs,
        &trace.Status,
        &trace.TotalCostUSD,
        &totalTokens,
        &trace.Model,
        &trace.Provider,
        &trace.UserID,
        &metadataJSON,
    )

    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("trace not found")
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get trace: %w", err)
    }

    // Convert types
    trace.DurationMs = int64(durationMs)
    trace.TotalTokens = int(totalTokens)

    if metadataJSON != "" && metadataJSON != "{}" {
        json.Unmarshal([]byte(metadataJSON), &trace.Metadata)
    }

    spans, err := r.GetSpansByTraceID(ctx, traceID)
    if err == nil {
        trace.Spans = spans
    }

    return &trace, nil
}

// GetTraceCount returns total count for pagination
func (r *ClickHouseRepository) GetTraceCount(ctx context.Context, query *models.TraceQuery) (int64, error) {
    sql := "SELECT count() FROM traces WHERE organization_id = ?"
    args := []interface{}{query.OrganizationID}

    if query.ProjectID != "" {
        sql += " AND project_id = ?"
        args = append(args, query.ProjectID)
    }

    if !query.StartTime.IsZero() {
        sql += " AND timestamp >= ?"
        args = append(args, query.StartTime)
    }

    if !query.EndTime.IsZero() {
        sql += " AND timestamp <= ?"
        args = append(args, query.EndTime)
    }

    var count uint64
    err := r.conn.QueryRow(ctx, sql, args...).Scan(&count)
    return int64(count), err
}

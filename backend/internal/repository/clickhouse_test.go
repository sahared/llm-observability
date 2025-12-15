package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/Aditya-Pimpalkar/clarity/internal/models"
)

// getTestDSN returns the ClickHouse connection string for testing
func getTestDSN() string {
	dsn := os.Getenv("CLICKHOUSE_DSN")
	if dsn == "" {
		dsn = "clickhouse://localhost:9000/llm_observability"
	}
	return dsn
}

// TestNewClickHouseRepository tests connection creation
func TestNewClickHouseRepository(t *testing.T) {
	repo, err := NewClickHouseRepository(getTestDSN())
	if err != nil {
		t.Skipf("Skipping test: ClickHouse not available: %v", err)
	}
	defer repo.Close()

	// Test ping
	ctx := context.Background()
	if err := repo.Ping(ctx); err != nil {
		t.Fatalf("Failed to ping ClickHouse: %v", err)
	}
}

// TestSaveAndGetTrace tests saving and retrieving a trace
func TestSaveAndGetTrace(t *testing.T) {
	repo, err := NewClickHouseRepository(getTestDSN())
	if err != nil {
		t.Skipf("Skipping test: ClickHouse not available: %v", err)
	}
	defer repo.Close()

	ctx := context.Background()

	// Create a test trace
	trace := &models.Trace{
		TraceID:        uuid.New().String(),
		OrganizationID: "test-org-" + uuid.New().String(),
		ProjectID:      "test-proj-" + uuid.New().String(),
		Timestamp:      time.Now(),
		TraceType:      "single_call",
		DurationMs:     150,
		Status:         "success",
		TotalCostUSD:   0.0024,
		TotalTokens:    500,
		Model:          "gpt-4",
		Provider:       "openai",
		UserID:         "test-user",
		Metadata:       map[string]string{"temperature": "0.7"},
		Spans: []models.Span{
			{
				SpanID:           uuid.New().String(),
				TraceID:          "", // Will be set to trace.TraceID
				Name:             "llm_call",
				StartTime:        time.Now(),
				EndTime:          time.Now().Add(150 * time.Millisecond),
				DurationMs:       150,
				Model:            "gpt-4",
				Provider:         "openai",
				Input:            "Hello, world!",
				Output:           "Hi there!",
				PromptTokens:     10,
				CompletionTokens: 5,
				TotalTokens:      15,
				CostUSD:          0.0024,
				Status:           "success",
				Metadata:         map[string]string{},
			},
		},
	}

	// Set span trace ID
	trace.Spans[0].TraceID = trace.TraceID

	// Save the trace
	if err := repo.SaveTrace(ctx, trace); err != nil {
		t.Fatalf("Failed to save trace: %v", err)
	}

	// Retrieve the trace
	retrieved, err := repo.GetTraceByID(ctx, trace.TraceID)
	if err != nil {
		t.Fatalf("Failed to get trace: %v", err)
	}

	// Verify data
	if retrieved.TraceID != trace.TraceID {
		t.Errorf("TraceID mismatch: got %s, want %s", retrieved.TraceID, trace.TraceID)
	}

	if retrieved.Model != trace.Model {
		t.Errorf("Model mismatch: got %s, want %s", retrieved.Model, trace.Model)
	}

	if len(retrieved.Spans) != len(trace.Spans) {
		t.Errorf("Span count mismatch: got %d, want %d", len(retrieved.Spans), len(trace.Spans))
	}
}

// TestGetTraces tests filtering traces
func TestGetTraces(t *testing.T) {
	repo, err := NewClickHouseRepository(getTestDSN())
	if err != nil {
		t.Skipf("Skipping test: ClickHouse not available: %v", err)
	}
	defer repo.Close()

	ctx := context.Background()
	orgID := "test-org-" + uuid.New().String()
	projectID := "test-proj-" + uuid.New().String()

	// Create and save multiple test traces
	for i := 0; i < 3; i++ {
		trace := &models.Trace{
			TraceID:        uuid.New().String(),
			OrganizationID: orgID,
			ProjectID:      projectID,
			Timestamp:      time.Now(),
			TraceType:      "single_call",
			DurationMs:     int64(100 + i*50),
			Status:         "success",
			TotalCostUSD:   0.001 * float64(i+1),
			TotalTokens:    int(100 * (i + 1)),
			Model:          "gpt-4",
			Provider:       "openai",
			UserID:         "test-user",
			Metadata:       map[string]string{},
			Spans:          []models.Span{},
		}

		if err := repo.SaveTrace(ctx, trace); err != nil {
			t.Fatalf("Failed to save trace %d: %v", i, err)
		}

		// Small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)
	}

	// Query traces
	query := &models.TraceQuery{
		OrganizationID: orgID,
		ProjectID:      projectID,
		Limit:          10,
		Offset:         0,
	}

	traces, err := repo.GetTraces(ctx, query)
	if err != nil {
		t.Fatalf("Failed to get traces: %v", err)
	}

	if len(traces) != 3 {
		t.Errorf("Expected 3 traces, got %d", len(traces))
	}
}

// TestGetMetricSummary tests metric aggregation
func TestGetMetricSummary(t *testing.T) {
	repo, err := NewClickHouseRepository(getTestDSN())
	if err != nil {
		t.Skipf("Skipping test: ClickHouse not available: %v", err)
	}
	defer repo.Close()

	ctx := context.Background()
	orgID := "test-org-" + uuid.New().String()
	projectID := "test-proj-" + uuid.New().String()
	startTime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now()

	// Create test traces
	for i := 0; i < 5; i++ {
		trace := &models.Trace{
			TraceID:        uuid.New().String(),
			OrganizationID: orgID,
			ProjectID:      projectID,
			Timestamp:      time.Now().Add(time.Duration(-i*10) * time.Minute),
			TraceType:      "single_call",
			DurationMs:     int64(100 + i*20),
			Status:         "success",
			TotalCostUSD:   0.001 * float64(i+1),
			TotalTokens:    int(100 * (i + 1)),
			Model:          "gpt-4",
			Provider:       "openai",
			UserID:         "test-user",
			Metadata:       map[string]string{},
			Spans:          []models.Span{},
		}

		if err := repo.SaveTrace(ctx, trace); err != nil {
			t.Fatalf("Failed to save trace %d: %v", i, err)
		}
	}

	// Get metric summary
	summary, err := repo.GetMetricSummary(ctx, orgID, projectID, startTime, endTime)
	if err != nil {
		t.Fatalf("Failed to get metric summary: %v", err)
	}

	if summary.TotalRequests != 5 {
		t.Errorf("Expected 5 requests, got %d", summary.TotalRequests)
	}

	if summary.TotalCostUSD == 0 {
		t.Error("Expected non-zero total cost")
	}

	t.Logf("Metric Summary: %+v", summary)
}

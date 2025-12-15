package services

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Aditya-Pimpalkar/clarity/internal/models"
	"github.com/Aditya-Pimpalkar/clarity/internal/repository"
)

// Mock repository for testing
type mockRepository struct {
	saveTraceFunc  func(ctx context.Context, trace *models.Trace) error
	saveMetricFunc func(ctx context.Context, metric *models.Metric) error
}

func (m *mockRepository) SaveTrace(ctx context.Context, trace *models.Trace) error {
	if m.saveTraceFunc != nil {
		return m.saveTraceFunc(ctx, trace)
	}
	return nil
}

func (m *mockRepository) SaveMetric(ctx context.Context, metric *models.Metric) error {
	if m.saveMetricFunc != nil {
		return m.saveMetricFunc(ctx, metric)
	}
	return nil
}

// Implement other required methods with correct signatures
func (m *mockRepository) GetTraceByID(ctx context.Context, traceID string) (*models.Trace, error) {
	return nil, nil
}

func (m *mockRepository) GetTraces(ctx context.Context, query *models.TraceQuery) ([]*models.Trace, error) {
	return nil, nil
}

func (m *mockRepository) GetTraceCount(ctx context.Context, query *models.TraceQuery) (int64, error) {
	return 0, nil
}

func (m *mockRepository) SaveSpan(ctx context.Context, span *models.Span) error {
	return nil
}

func (m *mockRepository) GetSpansByTraceID(ctx context.Context, traceID string) ([]models.Span, error) {
	return nil, nil
}

func (m *mockRepository) GetMetrics(ctx context.Context, query *models.MetricQuery) ([]*models.Metric, error) {
	return nil, nil
}

func (m *mockRepository) GetMetricSummary(ctx context.Context, orgID, projectID string, startTime, endTime time.Time) (*models.MetricSummary, error) {
	return &models.MetricSummary{}, nil
}

func (m *mockRepository) GetCostBreakdown(ctx context.Context, orgID, projectID string, startTime, endTime time.Time) ([]*models.CostBreakdown, error) {
	return nil, nil
}

func (m *mockRepository) GetModelUsage(ctx context.Context, orgID, projectID string, startTime, endTime time.Time) ([]*models.ModelUsage, error) {
	return nil, nil
}

func (m *mockRepository) CreateUser(ctx context.Context, user *models.User) error {
	return nil
}

func (m *mockRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, repository.ErrNotFound
}

func (m *mockRepository) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	return nil, nil
}

func (m *mockRepository) CreateOrganization(ctx context.Context, org *models.Organization) error {
	return nil
}

func (m *mockRepository) GetOrganization(ctx context.Context, orgID string) (*models.Organization, error) {
	return nil, nil
}

func (m *mockRepository) CreateProject(ctx context.Context, project *models.Project) error {
	return nil
}

func (m *mockRepository) GetProject(ctx context.Context, projectID string) (*models.Project, error) {
	return nil, nil
}

func (m *mockRepository) GetProjectsByOrg(ctx context.Context, orgID string) ([]*models.Project, error) {
	return nil, nil
}

func (m *mockRepository) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (m *mockRepository) Ping(ctx context.Context) error {
	return nil
}

func (m *mockRepository) Close() error {
	return nil
}

// TestCreateTrace tests the CreateTrace method
func TestCreateTrace(t *testing.T) {
	mock := &mockRepository{}
	service := NewTraceService(mock)

	ctx := context.Background()

	// Test valid request
	req := &models.TraceRequest{
		OrganizationID: "org-123",
		ProjectID:      "proj-456",
		TraceType:      "single_call",
		Model:          "gpt-4",
		Provider:       "openai",
		Spans: []models.SpanRequest{
			{
				Name:             "test_span",
				Model:            "gpt-4",
				Provider:         "openai",
				PromptTokens:     100,
				CompletionTokens: 50,
				DurationMs:       150,
				Status:           "success",
			},
		},
	}

	resp, err := service.CreateTrace(ctx, req)
	if err != nil {
		t.Fatalf("CreateTrace failed: %v", err)
	}

	if resp.TraceID == "" {
		t.Error("Expected trace ID to be generated")
	}

	if resp.Status != "accepted" {
		t.Errorf("Expected status 'accepted', got '%s'", resp.Status)
	}
}

// TestValidateTraceRequest tests request validation
func TestValidateTraceRequest(t *testing.T) {
	service := NewTraceService(&mockRepository{})

	tests := []struct {
		name    string
		req     *models.TraceRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: &models.TraceRequest{
				OrganizationID: "org-123",
				ProjectID:      "proj-456",
				TraceType:      "single_call",
				Model:          "gpt-4",
				Provider:       "openai",
				Spans: []models.SpanRequest{
					{
						Name:     "test",
						Model:    "gpt-4",
						Provider: "openai",
						Status:   "success",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing organization_id",
			req: &models.TraceRequest{
				ProjectID: "proj-456",
				Model:     "gpt-4",
				Provider:  "openai",
			},
			wantErr: true,
		},
		{
			name: "missing spans",
			req: &models.TraceRequest{
				OrganizationID: "org-123",
				ProjectID:      "proj-456",
				Model:          "gpt-4",
				Provider:       "openai",
				Spans:          []models.SpanRequest{},
			},
			wantErr: true,
		},
		{
			name: "invalid trace type",
			req: &models.TraceRequest{
				OrganizationID: "org-123",
				ProjectID:      "proj-456",
				TraceType:      "invalid_type",
				Model:          "gpt-4",
				Provider:       "openai",
				Spans: []models.SpanRequest{
					{Name: "test", Model: "gpt-4", Provider: "openai", Status: "success"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateTraceRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateTraceRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestCalculateCost tests cost calculation
func TestCalculateCost(t *testing.T) {
	service := NewTraceService(&mockRepository{})

	tests := []struct {
		name             string
		model            string
		provider         string
		promptTokens     uint32
		completionTokens uint32
		wantCost         float64
	}{
		{
			name:             "gpt-4",
			model:            "gpt-4",
			provider:         "openai",
			promptTokens:     1000,
			completionTokens: 500,
			wantCost:         0.06, // (1000/1000 * 0.03) + (500/1000 * 0.06)
		},
		{
			name:             "claude-3-sonnet",
			model:            "claude-3-sonnet",
			provider:         "anthropic",
			promptTokens:     1000,
			completionTokens: 500,
			wantCost:         0.0105, // (1000/1000 * 0.003) + (500/1000 * 0.015)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := service.calculateCost(tt.model, tt.provider, int(tt.promptTokens), int(tt.completionTokens))

			// Allow small floating point differences
			if cost < tt.wantCost*0.99 || cost > tt.wantCost*1.01 {
				t.Errorf("calculateCost() = %v, want %v", cost, tt.wantCost)
			}
		})
	}
}

// TestDetermineTraceStatus tests status determination
func TestDetermineTraceStatus(t *testing.T) {
	service := NewTraceService(&mockRepository{})

	tests := []struct {
		name       string
		spans      []models.Span
		wantStatus string
	}{
		{
			name: "all success",
			spans: []models.Span{
				{Status: "success"},
				{Status: "success"},
			},
			wantStatus: "success",
		},
		{
			name: "one error",
			spans: []models.Span{
				{Status: "success"},
				{Status: "error"},
			},
			wantStatus: "error",
		},
		{
			name: "one timeout",
			spans: []models.Span{
				{Status: "success"},
				{Status: "timeout"},
			},
			wantStatus: "timeout",
		},
		{
			name:       "empty spans",
			spans:      []models.Span{},
			wantStatus: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := service.determineTraceStatus(tt.spans)
			if status != tt.wantStatus {
				t.Errorf("determineTraceStatus() = %v, want %v", status, tt.wantStatus)
			}
		})
	}
}

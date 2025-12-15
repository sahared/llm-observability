package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Aditya-Pimpalkar/clarity/internal/models"
)

// Repository defines the interface for data persistence
// Any database implementation must satisfy this interface
type Repository interface {
	// Trace operations
	SaveTrace(ctx context.Context, trace *models.Trace) error
	GetTraceByID(ctx context.Context, traceID string) (*models.Trace, error)
	GetTraces(ctx context.Context, query *models.TraceQuery) ([]*models.Trace, error)
	GetTraceCount(ctx context.Context, query *models.TraceQuery) (int64, error)

	// Span operations
	SaveSpan(ctx context.Context, span *models.Span) error
	GetSpansByTraceID(ctx context.Context, traceID string) ([]models.Span, error)

	// Metrics operations
	SaveMetric(ctx context.Context, metric *models.Metric) error
	GetMetrics(ctx context.Context, query *models.MetricQuery) ([]*models.Metric, error)
	GetMetricSummary(ctx context.Context, orgID, projectID string, startTime, endTime time.Time) (*models.MetricSummary, error)

	// Analytics operations
	GetCostBreakdown(ctx context.Context, orgID, projectID string, startTime, endTime time.Time) ([]*models.CostBreakdown, error)
	GetModelUsage(ctx context.Context, orgID, projectID string, startTime, endTime time.Time) ([]*models.ModelUsage, error)

	// User operations
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, userID string) (*models.User, error)

	// Organization operations
	CreateOrganization(ctx context.Context, org *models.Organization) error
	GetOrganization(ctx context.Context, orgID string) (*models.Organization, error)

	// Project operations
	CreateProject(ctx context.Context, project *models.Project) error
	GetProject(ctx context.Context, projectID string) (*models.Project, error)
	GetProjectsByOrg(ctx context.Context, orgID string) ([]*models.Project, error)

	// Raw query operations (for analytics)
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

	// Health check
	Ping(ctx context.Context) error
	Close() error
}

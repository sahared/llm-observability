package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/Aditya-Pimpalkar/clarity/internal/models"
	"github.com/Aditya-Pimpalkar/clarity/internal/repository"
)

// UserService handles business logic for users and authentication
type UserService struct {
	repo repository.Repository
}

// NewUserService creates a new user service
func NewUserService(repo repository.Repository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, email, name, orgID, role string) (*models.User, error) {
	// Validate inputs
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if orgID == "" {
		return nil, fmt.Errorf("organization_id is required")
	}

	// Set default role
	if role == "" {
		role = "member"
	}

	// Validate role
	validRoles := map[string]bool{
		"admin":  true,
		"member": true,
		"viewer": true,
	}
	if !validRoles[role] {
		return nil, fmt.Errorf("invalid role: %s", role)
	}

	// Check if user already exists
	existing, err := s.repo.GetUserByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Create user
	user := &models.User{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		Email:          email,
		Name:           name,
		Role:           role,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// CreateOrganization creates a new organization
func (s *UserService) CreateOrganization(ctx context.Context, name, plan string) (*models.Organization, error) {
	// Validate inputs
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	// Set default plan
	if plan == "" {
		plan = "free"
	}

	// Validate plan
	validPlans := map[string]bool{
		"free":       true,
		"pro":        true,
		"enterprise": true,
	}
	if !validPlans[plan] {
		return nil, fmt.Errorf("invalid plan: %s", plan)
	}

	// Create organization
	org := &models.Organization{
		ID:        uuid.New().String(),
		Name:      name,
		Plan:      plan,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateOrganization(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	return org, nil
}

// GetOrganization retrieves an organization by ID
func (s *UserService) GetOrganization(ctx context.Context, orgID string) (*models.Organization, error) {
	if orgID == "" {
		return nil, fmt.Errorf("organization_id is required")
	}

	org, err := s.repo.GetOrganization(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	return org, nil
}

// CreateProject creates a new project
func (s *UserService) CreateProject(ctx context.Context, orgID, name, description string) (*models.Project, error) {
	// Validate inputs
	if orgID == "" {
		return nil, fmt.Errorf("organization_id is required")
	}
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	// Verify organization exists
	_, err := s.repo.GetOrganization(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}

	// Create project
	project := &models.Project{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		Name:           name,
		Description:    description,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.CreateProject(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return project, nil
}

// GetProject retrieves a project by ID
func (s *UserService) GetProject(ctx context.Context, projectID string) (*models.Project, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project_id is required")
	}

	project, err := s.repo.GetProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return project, nil
}

// GetProjectsByOrganization retrieves all projects for an organization
func (s *UserService) GetProjectsByOrganization(ctx context.Context, orgID string) ([]*models.Project, error) {
	if orgID == "" {
		return nil, fmt.Errorf("organization_id is required")
	}

	projects, err := s.repo.GetProjectsByOrg(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	return projects, nil
}

// ValidateUserAccess checks if a user has access to an organization/project
func (s *UserService) ValidateUserAccess(ctx context.Context, userID, orgID, projectID string) error {
	// Get user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Check organization access
	if user.OrganizationID != orgID {
		return fmt.Errorf("user does not have access to organization %s", orgID)
	}

	// If project ID provided, verify it belongs to the organization
	if projectID != "" {
		project, err := s.repo.GetProject(ctx, projectID)
		if err != nil {
			return fmt.Errorf("project not found: %w", err)
		}

		if project.OrganizationID != orgID {
			return fmt.Errorf("project does not belong to organization")
		}
	}

	return nil
}

// CheckUserRole checks if a user has a specific role or higher
func (s *UserService) CheckUserRole(ctx context.Context, userID string, requiredRole string) (bool, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("user not found: %w", err)
	}

	// Role hierarchy: admin > member > viewer
	roleLevel := map[string]int{
		"viewer": 1,
		"member": 2,
		"admin":  3,
	}

	userLevel := roleLevel[user.Role]
	requiredLevel := roleLevel[requiredRole]

	return userLevel >= requiredLevel, nil
}

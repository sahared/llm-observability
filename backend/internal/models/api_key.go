package models

import "time"

type APIKey struct {
    ID             string    `json:"id"`
    Key            string    `json:"key"`
    Name           string    `json:"name"`
    OrganizationID string    `json:"organization_id"`
    ProjectID      string    `json:"project_id,omitempty"`
    UserID         string    `json:"user_id"`
    IsActive       bool      `json:"is_active"`
    ExpiresAt      time.Time `json:"expires_at,omitempty"`
    LastUsedAt     time.Time `json:"last_used_at,omitempty"`
    CreatedAt      time.Time `json:"created_at"`
}

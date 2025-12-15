package models

import "time"

type Project struct {
    ID             string    `json:"id"`
    Name           string    `json:"name"`
    Description    string    `json:"description,omitempty"`
    OrganizationID string    `json:"organization_id"`
    CreatedBy      string    `json:"created_by"`
    IsActive       bool      `json:"is_active"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

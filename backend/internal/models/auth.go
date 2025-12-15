package models

import "time"

type User struct {
    ID             string    `json:"id" ch:"id"`
    Name           string    `json:"name" ch:"name"`
    Email          string    `json:"email" ch:"email"`
    PasswordHash   string    `json:"-" ch:"password_hash"`
    OrganizationID string    `json:"organization_id" ch:"organization_id"`
    Role           string    `json:"role" ch:"role"`
    CreatedAt      time.Time `json:"created_at" ch:"created_at"`
    UpdatedAt      time.Time `json:"updated_at" ch:"updated_at"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
    Token        string `json:"token"`
    RefreshToken string `json:"refresh_token,omitempty"`
    User         User   `json:"user"`
    ExpiresAt    string `json:"expires_at,omitempty"`
}

type Organization struct {
    ID        string    `json:"id" ch:"id"`
    Name      string    `json:"name" ch:"name"`
    Plan      string    `json:"plan" ch:"plan"`
    CreatedAt time.Time `json:"created_at" ch:"created_at"`
    UpdatedAt time.Time `json:"updated_at" ch:"updated_at"`
}

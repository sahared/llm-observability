package api

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/Aditya-Pimpalkar/clarity/internal/middleware"
	"github.com/Aditya-Pimpalkar/clarity/internal/models"
	"github.com/Aditya-Pimpalkar/clarity/internal/services"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	userService *services.UserService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := BindJSON(c, &req); err != nil {
		return BadRequestResponse(c, "Invalid request body")
	}

	// Validate credentials
	if req.Email == "" || req.Password == "" {
		return BadRequestResponse(c, "Email and password are required")
	}

	// Get user from database
	user, err := h.userService.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
			"code":  "INVALID_CREDENTIALS",
		})
	}

	// In production, verify password hash here
	// For demo, we'll skip password verification
	// bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))

	// Generate JWT token
	token, err := middleware.GenerateToken(
		user.ID,
		user.OrganizationID,
		user.Role,
		user.Email,
	)
	if err != nil {
		return InternalErrorResponse(c, "Failed to generate token")
	}

	// Return response
	return SuccessResponse(c, models.LoginResponse{
		Token:        token,
		RefreshToken: "", // Implement refresh tokens in production
		User:         *user,
	})
}

// GetCurrentUser handles GET /api/v1/auth/me
func (h *AuthHandler) GetCurrentUser(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
			"code":  "NOT_AUTHENTICATED",
		})
	}

	user, err := h.userService.GetUser(c.Context(), userID)
	if err != nil {
		return NotFoundResponse(c, "User not found")
	}

	return SuccessResponse(c, user)
}

// GenerateAPIKey handles POST /api/v1/auth/api-keys
func (h *AuthHandler) GenerateAPIKey(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	orgID := middleware.GetOrgID(c)

	if userID == "" || orgID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
			"code":  "NOT_AUTHENTICATED",
		})
	}

	// Parse request
	type APIKeyRequest struct {
		Name      string `json:"name"`
		ProjectID string `json:"project_id"`
	}

	var req APIKeyRequest
	if err := BindJSON(c, &req); err != nil {
		return BadRequestResponse(c, "Invalid request body")
	}

	// Validate required fields
	if req.Name == "" {
		return BadRequestResponse(c, "Name is required")
	}

	if req.ProjectID == "" {
		return BadRequestResponse(c, "Project ID is required")
	}

	// Generate API key (in production, use crypto/rand and store hash in DB)
	apiKey := "llm_" + generateRandomString(32)

	// Return response
	return CreatedResponse(c, fiber.Map{
		"api_key":    apiKey,
		"name":       req.Name,
		"project_id": req.ProjectID,
		"created_at": time.Now().Format(time.RFC3339),
		"note":       "Store this key securely. It won't be shown again.",
	})
}

// generateRandomString generates a cryptographically secure random string
func generateRandomString(length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to a demo key if random generation fails
		return "demo_key_" + time.Now().Format("20060102150405")
	}
	return hex.EncodeToString(bytes)
}

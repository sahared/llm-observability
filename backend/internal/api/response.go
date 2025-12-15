package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Aditya-Pimpalkar/clarity/internal/models"
)

// Response helpers for consistent API responses

// SuccessResponse sends a successful response
func SuccessResponse(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// CreatedResponse sends a 201 Created response
func CreatedResponse(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *fiber.Ctx, statusCode int, message string, details interface{}) error {
	response := models.ErrorResponse{
		Error: message,
	}

	if details != nil {
		if detailsMap, ok := details.(map[string]interface{}); ok {
			response.Details = detailsMap
		}
	}

	return c.Status(statusCode).JSON(response)
}

// BadRequestResponse sends a 400 Bad Request response
func BadRequestResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, message, nil)
}

// UnauthorizedResponse sends a 401 Unauthorized response
func UnauthorizedResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusUnauthorized, message, nil)
}

// ForbiddenResponse sends a 403 Forbidden response
func ForbiddenResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusForbidden, message, nil)
}

// NotFoundResponse sends a 404 Not Found response
func NotFoundResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusNotFound, message, nil)
}

// InternalErrorResponse sends a 500 Internal Server Error response
func InternalErrorResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusInternalServerError, message, nil)
}

// PaginatedResponse sends a paginated response
func PaginatedResponse(c *fiber.Ctx, data interface{}, total int64, page, pageSize int) error {
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return c.JSON(models.PaginatedResponse{
		Data:       data,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// BindJSON binds and validates JSON request body
func BindJSON(c *fiber.Ctx, v interface{}) error {
	if err := c.BodyParser(v); err != nil {
		return err
	}
	// You can add validation here using go-playground/validator
	return nil
}

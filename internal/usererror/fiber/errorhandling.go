package fiber

import (
	"errors"
	"fmt"

	"github.com/abramad-labs/irbankmock/internal/usererror"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

type UserErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type NonUserErrorResponse struct {
	Success   bool   `json:"success"`
	Error     string `json:"error"`
	ErrorId   string `json:"errorId"`
	RequestId string `json:"requestId,omitempty"`
}

func FiberUserErrorHandling(c *fiber.Ctx, err error) error {
	var userError *usererror.UserError
	if errors.As(err, &userError) {
		status := fiber.StatusInternalServerError
		if userError.Status != nil {
			status = *userError.Status
		}
		return c.Status(status).JSON(&UserErrorResponse{
			Success: false,
			Error:   userError.Error(),
		})
	}
	errUuid := uuid.NewString()
	log.Errorf("server error [%s]: %s", errUuid, err.Error())
	return c.Status(fiber.StatusInternalServerError).JSON(&NonUserErrorResponse{
		Success:   false,
		ErrorId:   errUuid,
		Error:     fmt.Sprintf("Server error. Error id: %s", errUuid),
		RequestId: c.Locals("requestid").(string),
	})
}

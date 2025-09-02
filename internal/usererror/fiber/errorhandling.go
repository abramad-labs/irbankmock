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

	code := fiber.StatusInternalServerError
	message := ""
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}
	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	var responseError string
	if message == "" {
		responseError = fmt.Sprintf("Server error. Error id: %s", errUuid)
	} else {
		responseError = message
	}

	return c.Status(code).JSON(&NonUserErrorResponse{
		Success:   false,
		ErrorId:   errUuid,
		Error:     responseError,
		RequestId: c.Locals("requestid").(string),
	})
}

package api

import (
	"errors"

	"github.com/SamVitebsk/sharetrip-contract/gen"

	"github.com/gofiber/fiber/v2"
)

const (
	errorCodeInvalidRequest = "invalidRequest"
	errorCodeInternalError  = "internalError"
	errorCodeNotFound       = "notFound"
	errorCodeNotImplemented = "notImplemented"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	code := errorCodeInternalError
	message := "внутренняя ошибка сервера"

	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		status = fiberError.Code
		message = fiberError.Message

		switch status {
		case fiber.StatusBadRequest:
			code = errorCodeInvalidRequest
		case fiber.StatusNotFound:
			code = errorCodeNotFound
		}
	}

	return ctx.Status(status).JSON(gen.ErrorResponse{
		Code:    code,
		Message: message,
	})
}

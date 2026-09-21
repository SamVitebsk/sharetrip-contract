package api

import (
	"errors"
	"log"

	"github.com/SamVitebsk/sharetrip-contract/gen"
	"github.com/SamVitebsk/sharetrip-contract/internal/service"

	"github.com/gofiber/fiber/v2"
)

const (
	errorCodeInvalidRequest = "invalidRequest"
	errorCodeInternalError  = "internalError"
	errorCodeNotFound       = "notFound"
	errorCodeConflict       = "conflict"
	errorCodeNotImplemented = "notImplemented"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	code := errorCodeInternalError
	message := "внутренняя ошибка сервера"

	var fiberErr *fiber.Error

	switch {
	case errors.Is(err, service.ErrInvalidInput):
		status, code, message = fiber.StatusBadRequest, errorCodeInvalidRequest, service.ErrInvalidInput.Error()
	case errors.Is(err, service.ErrNotFound):
		status, code, message = fiber.StatusNotFound, errorCodeNotFound, service.ErrNotFound.Error()
	case errors.Is(err, service.ErrConflict):
		status, code, message = fiber.StatusConflict, errorCodeConflict, service.ErrConflict.Error()
	case errors.As(err, &fiberErr) && fiberErr.Code == fiber.StatusNotFound:
		status, code, message = fiberErr.Code, errorCodeNotFound, fiberErr.Message
	case errors.As(err, &fiberErr) && fiberErr.Code == fiber.StatusConflict:
		status, code, message = fiberErr.Code, errorCodeConflict, fiberErr.Message
	case errors.As(err, &fiberErr) && fiberErr.Code == fiber.StatusNotImplemented:
		status, code, message = fiberErr.Code, errorCodeNotImplemented, fiberErr.Message
	case errors.As(err, &fiberErr) && fiberErr.Code >= 400 && fiberErr.Code < 500:
		status, code, message = fiberErr.Code, errorCodeInvalidRequest, fiberErr.Message
	}

	if status == fiber.StatusInternalServerError {
		log.Printf("HTTP 500: method=%s route=%q error=%v", ctx.Method(), ctx.Route().Path, err)
	}

	return ctx.Status(status).JSON(gen.ErrorResponse{
		Code:    code,
		Message: message,
	})
}

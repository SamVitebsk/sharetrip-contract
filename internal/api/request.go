package api

import (
	"mime"

	"github.com/gofiber/fiber/v2"
)

// parseJSONBody проверяет заголовок Content-Type и десериализует JSON-тело запроса.
// При несоответствии Content-Type возвращает 415, при невалидном теле — 400.
// Вынесена в общий хелпер для повторного использования будущими POST/PUT-эндпоинтами с телом запроса.
func parseJSONBody(ctx *fiber.Ctx, out any) error {
	mediaType, _, err := mime.ParseMediaType(ctx.Get(fiber.HeaderContentType))
	if err != nil || mediaType != fiber.MIMEApplicationJSON {
		return fiber.NewError(fiber.StatusUnsupportedMediaType, "ожидается Content-Type application/json")
	}
	if err := ctx.BodyParser(out); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "некорректное тело запроса")
	}

	return nil
}

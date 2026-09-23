package api

import (
	"github.com/SamVitebsk/sharetrip-contract/gen"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) CheckServiceAvailability(ctx *fiber.Ctx) error {
	var request gen.CheckServiceAvailabilityRequest
	if err := parseJSONBody(ctx, &request); err != nil {
		return err
	}

	availabilityResponse, err := s.service.CheckServiceAvailability(ctx.UserContext(), toCheckServiceAvailabilityRequest(request))
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(toCheckServiceAvailabilityResponse(availabilityResponse))
}

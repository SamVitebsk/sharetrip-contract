package api

import (
	"github.com/SamVitebsk/sharetrip-contract/gen"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) CreateContract(ctx *fiber.Ctx) error {
	var request gen.CreateContractRequest
	if err := parseJSONBody(ctx, &request); err != nil {
		return err
	}

	result, err := s.service.CreateContract(ctx.UserContext(), toCreateContractCommand(request))
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(toCreateContractResponse(result))
}

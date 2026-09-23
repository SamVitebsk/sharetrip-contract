package api

import (
	"github.com/SamVitebsk/sharetrip-contract/gen"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) SignContract(ctx *fiber.Ctx, contractId gen.ContractId) error {
	response, err := s.service.SignContract(ctx.UserContext(), toSignContractRequest(contractId))
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(toSignContractResponse(response))
}

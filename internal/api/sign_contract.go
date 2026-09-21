package api

import (
	"github.com/SamVitebsk/sharetrip-contract/gen"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) SignContract(ctx *fiber.Ctx, contractId gen.ContractId) error {
	result, err := s.service.SignContract(ctx.UserContext(), toSignContractCommand(contractId))
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(toSignContractResponse(result))
}

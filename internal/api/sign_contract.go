package api

import (
	"time"

	"github.com/SamVitebsk/sharetrip-contract/gen"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (s *Server) SignContract(ctx *fiber.Ctx, contractId gen.ContractId) error {
	if contractId == uuid.Nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(gen.ErrorResponse{
			Code:    errorCodeInvalidRequest,
			Message: "некорректный идентификатор договора",
		})
	}

	now := time.Now().UTC()
	contract := gen.Contract{
		Id:        contractId,
		CompanyId: uuid.Nil,
		Status:    gen.Active,
		ValidFrom: now,
		Services:  []gen.ContractService{},
		CreatedAt: now,
		UpdatedAt: now,
	}

	return ctx.Status(fiber.StatusOK).JSON(gen.SignContractResponse{Contract: contract})
}

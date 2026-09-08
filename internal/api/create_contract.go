package api

import (
	"time"

	"github.com/SamVitebsk/sharetrip-contract/gen"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (s *Server) CreateContract(ctx *fiber.Ctx) error {
	var request gen.CreateContractRequest
	if err := ctx.BodyParser(&request); err != nil ||
		request.CompanyId == uuid.Nil ||
		request.ValidFrom.IsZero() ||
		!validContractServices(request.Services) {
		return ctx.Status(fiber.StatusBadRequest).JSON(gen.ErrorResponse{
			Code:    errorCodeInvalidRequest,
			Message: "некорректное тело запроса",
		})
	}

	now := time.Now().UTC()
	contract := gen.Contract{
		Id:         uuid.New(),
		CompanyId:  request.CompanyId,
		Status:     gen.Draft,
		ValidFrom:  request.ValidFrom,
		ValidUntil: request.ValidUntil,
		Services:   request.Services,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	return ctx.Status(fiber.StatusCreated).JSON(gen.CreateContractResponse{Contract: contract})
}

func validContractServices(services []gen.ContractService) bool {
	if len(services) == 0 {
		return false
	}

	for _, service := range services {
		if !service.ServiceCode.Valid() {
			return false
		}
	}

	return true
}

package api

import (
	"github.com/SamVitebsk/sharetrip-contract/gen"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) GetContract(ctx *fiber.Ctx, contractId gen.ContractId) error {
	return notImplemented(ctx)
}

func (s *Server) GetCompanyActiveContract(ctx *fiber.Ctx, companyId gen.CompanyId) error {
	return notImplemented(ctx)
}

func (s *Server) SuspendContract(ctx *fiber.Ctx, contractId gen.ContractId) error {
	return notImplemented(ctx)
}

func (s *Server) ResumeContract(ctx *fiber.Ctx, contractId gen.ContractId) error {
	return notImplemented(ctx)
}

func (s *Server) TerminateContract(ctx *fiber.Ctx, contractId gen.ContractId) error {
	return notImplemented(ctx)
}

func (s *Server) CheckServiceAvailability(ctx *fiber.Ctx, companyId gen.CompanyId, serviceCode gen.ServiceCode) error {
	return notImplemented(ctx)
}

func notImplemented(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusNotImplemented).JSON(gen.ErrorResponse{
		Code:    errorCodeNotImplemented,
		Message: "операция пока не реализована",
	})
}

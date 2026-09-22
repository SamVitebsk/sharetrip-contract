package api

import (
	"github.com/SamVitebsk/sharetrip-contract/gen"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) GetContract(ctx *fiber.Ctx, contractId gen.ContractId) error {
	return notImplemented()
}

func (s *Server) GetClientActiveContract(ctx *fiber.Ctx, clientId gen.ClientId) error {
	return notImplemented()
}

func (s *Server) SuspendContract(ctx *fiber.Ctx, contractId gen.ContractId) error {
	return notImplemented()
}

func (s *Server) ResumeContract(ctx *fiber.Ctx, contractId gen.ContractId) error {
	return notImplemented()
}

func (s *Server) TerminateContract(ctx *fiber.Ctx, contractId gen.ContractId) error {
	return notImplemented()
}

func notImplemented() error {
	return fiber.NewError(fiber.StatusNotImplemented, "операция пока не реализована")
}

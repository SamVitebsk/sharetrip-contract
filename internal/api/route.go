package api

import (
	"errors"

	"github.com/SamVitebsk/sharetrip-contract/gen"
	"github.com/SamVitebsk/sharetrip-contract/internal/service"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	service *service.Service
}

var _ gen.ServerInterface = (*Server)(nil)

func NewServer(service *service.Service) (*Server, error) {
	if service == nil {
		return nil, errors.New("сервис договоров обязателен")
	}
	return &Server{service: service}, nil
}

func RegisterRoutes(router fiber.Router, server *Server) {
	gen.RegisterHandlers(router, server)
}

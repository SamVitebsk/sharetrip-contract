package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SamVitebsk/sharetrip-contract/internal/api"
	"github.com/SamVitebsk/sharetrip-contract/internal/env"
	"github.com/SamVitebsk/sharetrip-contract/internal/repository"
	"github.com/SamVitebsk/sharetrip-contract/internal/service"
	"github.com/joho/godotenv"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run() error {
	_ = godotenv.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverPort, err := env.Int("HTTP_PORT", 9190)
	if err != nil {
		return err
	}
	if serverPort < 1 || serverPort > 65535 {
		return errors.New("HTTP_PORT должен быть в диапазоне от 1 до 65535")
	}

	port, err := env.Int("PGPORT", 6544)
	if err != nil {
		return err
	}
	if port < 1 || port > 65535 {
		return errors.New("PGPORT должен быть в диапазоне от 1 до 65535")
	}
	cfg := repository.Config{
		Host:     env.String("PGHOST", "localhost"),
		Port:     port,
		User:     env.String("PGUSER", "postgres"),
		Password: env.String("PGPASSWORD", "postgres"),
		DBName:   env.String("PGDATABASE", "sharetrip_contract"),
		SSLMode:  env.String("PGSSLMODE", "disable"),
	}
	pool, err := repository.NewPool(ctx, cfg.DSN())
	if err != nil {
		return errors.New("не удалось подключиться к PostgreSQL: проверьте PG*-настройки и доступность БД")
	}
	defer pool.Close()

	repo, err := repository.NewRepoPg(pool)
	if err != nil {
		return err
	}
	txRunner := func(ctx context.Context, operation func(context.Context, service.RepositoryTx) error) error {
		return repo.WithinContractTx(ctx, func(transactionCtx context.Context, contracts *repository.ContractTx) error {
			return operation(transactionCtx, contracts)
		})
	}

	svc, err := service.New(txRunner)
	if err != nil {
		return err
	}

	server, err := api.NewServer(svc)
	if err != nil {
		return err
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: api.ErrorHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	})
	app.Use(logger.New())
	app.Use(func(request *fiber.Ctx) error {
		request.SetUserContext(ctx)
		return request.Next()
	})

	specification, err := os.ReadFile("./api/contract.yaml")
	if err != nil {
		return fmt.Errorf("read OpenAPI specification: %w", err)
	}
	app.Use(swagger.New(swagger.Config{
		BasePath:    "/",
		FilePath:    "./api/contract.yaml",
		FileContent: specification,
		Path:        "docs",
		Title:       "ShareTrip Contract Service API documentation",
	}))
	app.Static("/api/paths", "./api/paths")

	api.RegisterRoutes(app, server)

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- app.Listen(fmt.Sprintf(":%d", serverPort))
	}()

	select {
	case err := <-serverErrors:
		return err
	case <-ctx.Done():
		return app.ShutdownWithTimeout(10 * time.Second)
	}
}

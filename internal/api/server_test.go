package api_test

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/SamVitebsk/sharetrip-contract/internal/api"
	"github.com/SamVitebsk/sharetrip-contract/internal/repository"
	"github.com/SamVitebsk/sharetrip-contract/internal/service"
)

var (
	testCtx       context.Context
	testDB        *sql.DB
	testPool      *pgxpool.Pool
	testApp       *fiber.App
	testContainer *postgres.PostgresContainer
)

func TestMain(m *testing.M) {
	testCtx = context.Background()

	var err error
	testContainer, err = postgres.Run(
		testCtx,
		"postgres:16",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Fatalf("Ошибка запуска контейнера postgres: %v", err)
	}
	defer func() {
		if err := testContainer.Terminate(testCtx); err != nil {
			log.Fatalf("Ошибка остановки контейнера: %v", err)
		}
	}()

	connStr, err := testContainer.ConnectionString(testCtx, "sslmode=disable")
	if err != nil {
		log.Fatalf("Ошибка получения строки подключения: %v", err)
	}

	testDB, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Ошибка sql.Open: %v", err)
	}
	defer func() { _ = testDB.Close() }()

	if err := testDB.Ping(); err != nil {
		log.Fatalf("База недоступна: %v", err)
	}

	goose.SetBaseFS(nil)
	if err := goose.Up(testDB, "../../migrations"); err != nil {
		log.Fatalf("Ошибка миграций goose: %v", err)
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Ошибка конфига пула: %v", err)
	}
	testPool, err = pgxpool.NewWithConfig(testCtx, config)
	if err != nil {
		log.Fatalf("Ошибка создания пула: %v", err)
	}
	defer testPool.Close()

	repoPg, err := repository.NewRepoPg(testPool)
	if err != nil {
		log.Fatalf("Ошибка создания репозитория: %v", err)
	}

	txRunner := func(ctx context.Context, operation func(context.Context, service.RepositoryTx) error) error {
		return repoPg.WithinContractTx(ctx, func(transactionCtx context.Context, contracts *repository.ContractTx) error {
			return operation(transactionCtx, contracts)
		})
	}

	srv, err := service.New(txRunner)
	if err != nil {
		log.Fatalf("Ошибка создания сервиса: %v", err)
	}

	apiServer, err := api.NewServer(srv)
	if err != nil {
		log.Fatalf("Ошибка создания API: %v", err)
	}

	testApp = fiber.New(fiber.Config{
		ErrorHandler: api.ErrorHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})
	api.RegisterRoutes(testApp, apiServer)

	m.Run()
}

package api_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// insertContract вставляет тестовый договор в БД в обход бизнес-логики.
func insertContract(t *testing.T, status string, services map[string]bool) (uuid.UUID, uuid.UUID) {
	t.Helper()

	contractID := uuid.New()
	clientID := uuid.New()
	now := time.Now().UTC()

	const contractQuery = `
		INSERT INTO contracts (id, client_id, status, valid_from, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := testPool.Exec(context.Background(), contractQuery, contractID, clientID, status, now.Add(-time.Hour), now, now)
	require.NoError(t, err, "ошибка вставки тестового договора")

	const serviceQuery = `
		INSERT INTO contract_services (contract_id, service_code, allowed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	for code, allowed := range services {
		_, err := testPool.Exec(context.Background(), serviceQuery, contractID, code, allowed, now, now)
		require.NoError(t, err, "ошибка вставки тестовой услуги договора")
	}

	return clientID, contractID
}

// insertExpiredContract вставляет просроченный договор (для теста contractExpired)
func insertExpiredContract(t *testing.T, services map[string]bool) (uuid.UUID, uuid.UUID) {
	t.Helper()

	contractID := uuid.New()
	clientID := uuid.New()
	now := time.Now().UTC()
	validFrom := now.Add(-time.Hour * 48)
	validUntil := now.Add(-time.Hour * 24)

	const contractQuery = `
		INSERT INTO contracts (id, client_id, status, valid_from, valid_until, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := testPool.Exec(context.Background(), contractQuery, contractID, clientID, "active", validFrom, validUntil, now, now)
	require.NoError(t, err, "ошибка вставки просроченного договора")

	const serviceQuery = `
		INSERT INTO contract_services (contract_id, service_code, allowed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	for code, allowed := range services {
		_, err := testPool.Exec(context.Background(), serviceQuery, contractID, code, allowed, now, now)
		require.NoError(t, err, "ошибка вставки тестовой услуги договора")
	}

	return clientID, contractID
}

// insertFutureContract вставляет договор, который начнет действовать только в будущем.
func insertFutureContract(t *testing.T, services map[string]bool) (uuid.UUID, uuid.UUID) {
	t.Helper()

	contractID := uuid.New()
	clientID := uuid.New()
	now := time.Now().UTC()
	validFrom := now.Add(time.Hour * 24)

	const contractQuery = `
		INSERT INTO contracts (id, client_id, status, valid_from, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := testPool.Exec(context.Background(), contractQuery, contractID, clientID, "active", validFrom, now, now)
	require.NoError(t, err, "ошибка вставки будущего договора")

	const serviceQuery = `
		INSERT INTO contract_services (contract_id, service_code, allowed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	for code, allowed := range services {
		_, err := testPool.Exec(context.Background(), serviceQuery, contractID, code, allowed, now, now)
		require.NoError(t, err, "ошибка вставки тестовой услуги договора")
	}

	return clientID, contractID
}

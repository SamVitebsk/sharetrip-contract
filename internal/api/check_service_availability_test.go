package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/SamVitebsk/sharetrip-contract/gen"
)

func TestServer_CheckServiceAvailability(t *testing.T) {
	t.Run("Успех: услуга включена в активном договоре", func(t *testing.T) {
		t.Parallel()
		clientID, contractID := insertContract(t, "active", map[string]bool{
			"trip_creation": true,
		})

		actual := sendCheckRequest(t, clientID, "tripCreation")

		expected := gen.CheckServiceAvailabilityResponse{
			Allowed:    true,
			Reason:     nil,
			ContractId: &contractID,
		}
		require.Equal(t, expected, actual)
	})

	t.Run("Отказ: услуга прямым флагом отключена в договоре", func(t *testing.T) {
		t.Parallel()
		clientID, contractID := insertContract(t, "active", map[string]bool{
			"notifications": false,
		})

		actual := sendCheckRequest(t, clientID, "notifications")

		reason := gen.ServiceNotAllowed
		expected := gen.CheckServiceAvailabilityResponse{
			Allowed:    false,
			Reason:     &reason,
			ContractId: &contractID,
		}
		require.Equal(t, expected, actual)
	})

	t.Run("Отказ: договор для данного клиента не существует", func(t *testing.T) {
		t.Parallel()
		fakeClientID := uuid.New()

		actual := sendCheckRequest(t, fakeClientID, "tripCreation")

		reason := gen.ContractNotFound
		expected := gen.CheckServiceAvailabilityResponse{
			Allowed:    false,
			Reason:     &reason,
			ContractId: nil,
		}
		require.Equal(t, expected, actual)
	})

	t.Run("Отказ: срок действия договора уже истёк", func(t *testing.T) {
		t.Parallel()
		clientID, contractID := insertExpiredContract(t, map[string]bool{
			"trip_creation": true,
		})

		actual := sendCheckRequest(t, clientID, "tripCreation")

		reason := gen.ContractExpired
		expected := gen.CheckServiceAvailabilityResponse{
			Allowed:    false,
			Reason:     &reason,
			ContractId: &contractID,
		}
		require.Equal(t, expected, actual)
	})

	t.Run("Отказ: статус договора не активен (например, draft)", func(t *testing.T) {
		t.Parallel()
		clientID, contractID := insertContract(t, "draft", map[string]bool{
			"trip_creation": true,
		})

		actual := sendCheckRequest(t, clientID, "tripCreation")

		reason := gen.ContractNotActive
		expected := gen.CheckServiceAvailabilityResponse{
			Allowed:    false,
			Reason:     &reason,
			ContractId: &contractID,
		}
		require.Equal(t, expected, actual)
	})

	t.Run("Отказ: дата начала договора еще не наступила", func(t *testing.T) {
		t.Parallel()
		clientID, contractID := insertFutureContract(t, map[string]bool{
			"trip_creation": true,
		})

		actual := sendCheckRequest(t, clientID, "tripCreation")

		reason := gen.ContractNotStarted
		expected := gen.CheckServiceAvailabilityResponse{
			Allowed:    false,
			Reason:     &reason,
			ContractId: &contractID,
		}
		require.Equal(t, expected, actual)
	})

	t.Run("Отказ: передан некорректный формат UUID для клиента", func(t *testing.T) {
		t.Parallel()
		reqBody := map[string]string{
			"client_id":    "not-a-valid-uuid",
			"service_code": "tripCreation",
		}
		bodyData, err := json.Marshal(reqBody)
		require.NoError(t, err, "ошибка маршалинга запроса с ошибкой")
		req := httptest.NewRequest(http.MethodPost, "/contracts/check-service", bytes.NewReader(bodyData))
		req.Header.Set("Content-Type", "application/json")

		resp, err := testApp.Test(req)
		require.NoError(t, err, "ошибка вызова fiber app")
		defer func() { _ = resp.Body.Close() }()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var errResp gen.ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		require.NoError(t, err, "ошибка анмаршалинга ответа с ошибкой")
		require.Equal(t, "invalidRequest", errResp.Code)
	})

	t.Run("Отказ: передан несуществующий код услуги", func(t *testing.T) {
		t.Parallel()
		reqBody := map[string]string{
			"client_id":    uuid.New().String(),
			"service_code": "magicService",
		}
		bodyData, err := json.Marshal(reqBody)
		require.NoError(t, err, "ошибка маршалинга запроса с ошибкой")
		req := httptest.NewRequest(http.MethodPost, "/contracts/check-service", bytes.NewReader(bodyData))
		req.Header.Set("Content-Type", "application/json")

		resp, err := testApp.Test(req)
		require.NoError(t, err, "ошибка вызова fiber app")
		defer func() { _ = resp.Body.Close() }()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func sendCheckRequest(t *testing.T, clientID uuid.UUID, serviceCode string) gen.CheckServiceAvailabilityResponse {
	t.Helper()

	reqBody := gen.CheckServiceAvailabilityRequest{
		ClientId:    clientID,
		ServiceCode: gen.ServiceCode(serviceCode),
	}
	bodyData, err := json.Marshal(reqBody)
	require.NoError(t, err, "ошибка маршалинга HTTP-запроса")

	req := httptest.NewRequest(http.MethodPost, "/contracts/check-service", bytes.NewReader(bodyData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := testApp.Test(req)
	require.NoError(t, err, "ошибка вызова fiber.Test")
	defer func() { _ = resp.Body.Close() }()
	require.Equal(t, http.StatusOK, resp.StatusCode, "ожидался код 200 OK")

	var actual gen.CheckServiceAvailabilityResponse
	err = json.NewDecoder(resp.Body).Decode(&actual)
	require.NoError(t, err, "ошибка декодирования ответа API")

	return actual
}

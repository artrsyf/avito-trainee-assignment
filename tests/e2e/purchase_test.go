package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPurchaseScenarios(t *testing.T) {
	cfg := setupTestEnvironment()

	/*
		Предполагаем, что у нас есть предопределенные предметы:
		"powerbank" - стоимость 200 монет
		"pink-hoody" - стоимость 500 монет
	*/

	t.Run("Successful item purchase", func(t *testing.T) {
		token := authUser(t, cfg, "test_buyer", "password")
		initialBalance := getBalance(t, cfg, token)

		req := httptest.NewRequest("GET", "/api/buy/powerbank", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		cfg.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, initialBalance-200, getBalance(t, cfg, token))
	})

	t.Run("Insufficient balance", func(t *testing.T) {
		token := authUser(t, cfg, "poor_buyer", "password")
		initialBalance := getBalance(t, cfg, token)

		req := httptest.NewRequest("GET", "/api/buy/pink-hoody", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		cfg.Router.ServeHTTP(rr, req)

		req = httptest.NewRequest("GET", "/api/buy/powerbank", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr = httptest.NewRecorder()
		cfg.Router.ServeHTTP(rr, req)

		updatedBalance := getBalance(t, cfg, token)
		assert.Equal(t, initialBalance-700, updatedBalance)

		req = httptest.NewRequest("GET", "/api/buy/pink-hoody", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr = httptest.NewRecorder()
		cfg.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, updatedBalance, getBalance(t, cfg, token))

		var response map[string]string
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, "not enough balance", response["errors"])
	})

	t.Run("Purchase non-existent item", func(t *testing.T) {
		token := authUser(t, cfg, "test_buyer2", "password")

		req := httptest.NewRequest("GET", "/api/buy/non_existent_item", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		cfg.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Multiple purchases", func(t *testing.T) {
		token := authUser(t, cfg, "repeat_buyer", "password")
		initialBalance := getBalance(t, cfg, token)

		req := httptest.NewRequest("GET", "/api/buy/powerbank", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		cfg.Router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, initialBalance-200, getBalance(t, cfg, token))

		req = httptest.NewRequest("GET", "/api/buy/powerbank", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr = httptest.NewRecorder()
		cfg.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, initialBalance-400, getBalance(t, cfg, token))
	})

	t.Run("Invalid purchase parameters", func(t *testing.T) {
		testCases := []struct {
			name      string
			itemPath  string
			errorCode int
			errorMsg  string
		}{
			{"Empty item name", "", http.StatusNotFound, ""},
			{"Special characters", "invalid@item!", http.StatusNotFound, "item not found"},
		}

		token := authUser(t, cfg, "test_buyer3", "password")

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest("GET", "/api/buy/"+tc.itemPath, nil)
				req.Header.Set("Authorization", "Bearer "+token)
				rr := httptest.NewRecorder()
				cfg.Router.ServeHTTP(rr, req)

				var response map[string]string
				json.Unmarshal(rr.Body.Bytes(), &response)

				assert.Equal(t, tc.errorCode, rr.Code)
				assert.Equal(t, tc.errorMsg, response["errors"])
			})
		}
	})

	t.Run("Purchase history verification", func(t *testing.T) {
		token := authUser(t, cfg, "history_user", "password")

		req := httptest.NewRequest("GET", "/api/buy/powerbank", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		cfg.Router.ServeHTTP(httptest.NewRecorder(), req)

		req = httptest.NewRequest("GET", "/api/buy/powerbank", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		cfg.Router.ServeHTTP(httptest.NewRecorder(), req)

		req = httptest.NewRequest("GET", "/api/buy/pink-hoody", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		cfg.Router.ServeHTTP(httptest.NewRecorder(), req)

		req = httptest.NewRequest("GET", "/api/info", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		cfg.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var infoResponse struct {
			Coins     uint `json:"coins"`
			Inventory []struct {
				PurchaseTypeName string `json:"type"`
				Quantity         uint   `json:"quantity"`
			} `json:"inventory"`
		}
		json.Unmarshal(rr.Body.Bytes(), &infoResponse)

		assert.Equal(t, infoResponse.Coins, uint(100))
		assert.Len(t, infoResponse.Inventory, 2)

		assert.Contains(t, []string{"pink-hoody", "powerbank"}, infoResponse.Inventory[0].PurchaseTypeName)
		assert.Contains(t, []string{"pink-hoody", "powerbank"}, infoResponse.Inventory[1].PurchaseTypeName)

		assert.Contains(t, []uint{1, 2}, infoResponse.Inventory[0].Quantity)
		assert.Contains(t, []uint{1, 2}, infoResponse.Inventory[1].Quantity)
	})
}

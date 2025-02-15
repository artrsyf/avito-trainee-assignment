package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionScenarios(t *testing.T) {
	cfg := setupTestEnvironment()

	t.Run("Successful coins transfer", func(t *testing.T) {
		senderToken := authUser(t, cfg, "sender_user", "password")
		receiverToken := authUser(t, cfg, "receiver_user", "password")

		senderBalanceBefore := getBalance(t, cfg, senderToken)
		receiverBalanceBefore := getBalance(t, cfg, receiverToken)

		payload := map[string]interface{}{
			"toUser": "receiver_user",
			"amount": 100,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+senderToken)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		cfg.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		assert.Equal(t, senderBalanceBefore-100, getBalance(t, cfg, senderToken))
		assert.Equal(t, receiverBalanceBefore+100, getBalance(t, cfg, receiverToken))
	})

	t.Run("Insufficient balance", func(t *testing.T) {
		senderToken := authUser(t, cfg, "poor_user", "password")
		_ = authUser(t, cfg, "another_user", "password")

		payload := map[string]interface{}{
			"toUser": "another_user",
			"amount": 1500,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+senderToken)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		cfg.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		var response map[string]string
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, "not enough balance", response["errors"])
	})

	t.Run("Invalid amount values", func(t *testing.T) {
		testCases := []struct {
			name     string
			amount   interface{}
			errorMsg string
		}{
			{"Zero amount", 0, "Amount is required"},
			{"Negative amount", -50, "bad request"},
			{"String amount", "invalid", "bad request"},
		}

		senderToken := authUser(t, cfg, "test_user_transaction", "password")

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				payload := map[string]interface{}{
					"toUser": "receiver_user",
					"amount": tc.amount,
				}
				body, _ := json.Marshal(payload)

				req := httptest.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(body))
				req.Header.Set("Authorization", "Bearer "+senderToken)
				req.Header.Set("Content-Type", "application/json")

				rr := httptest.NewRecorder()
				cfg.Router.ServeHTTP(rr, req)

				assert.Equal(t, http.StatusBadRequest, rr.Code)

				var response map[string]string
				json.Unmarshal(rr.Body.Bytes(), &response)
				assert.Contains(t, response["errors"], tc.errorMsg)
			})
		}
	})

	t.Run("Invalid receiver username", func(t *testing.T) {
		testCases := []struct {
			name     string
			username string
			errorMsg string
		}{
			{"Short username", "ab", "ReceiverUsername is too short"},
			{"Long username", "very_long_username_that_exceeds_maximum_allowed_length", "ReceiverUsername is too long"},
			{"Non-existent user", "non_existent_user", "can't find such user"},
			{"Transfer to yourself", "test_user_transaction", "money transfer to yourself is not allowed"},
		}

		senderToken := authUser(t, cfg, "test_user_transaction", "password")

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				payload := map[string]interface{}{
					"toUser": tc.username,
					"amount": 100,
				}
				body, _ := json.Marshal(payload)

				req := httptest.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(body))
				req.Header.Set("Authorization", "Bearer "+senderToken)
				req.Header.Set("Content-Type", "application/json")

				rr := httptest.NewRecorder()
				cfg.Router.ServeHTTP(rr, req)

				assert.Equal(t, http.StatusBadRequest, rr.Code)

				var response map[string]string
				json.Unmarshal(rr.Body.Bytes(), &response)
				assert.Contains(t, response["errors"], tc.errorMsg)
			})
		}
	})
}

func authUser(t *testing.T, cfg *TestConfig, username, password string) string {
	payload := map[string]string{
		"username": username,
		"password": password,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/auth", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	cfg.Router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	var authResponse struct {
		Token string `json:"token"`
	}
	json.Unmarshal(rr.Body.Bytes(), &authResponse)

	return authResponse.Token
}

func getBalance(t *testing.T, cfg *TestConfig, token string) uint {
	req := httptest.NewRequest("GET", "/api/info", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	cfg.Router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	var infoResponse struct {
		Coins uint `json:"coins"`
	}
	json.Unmarshal(rr.Body.Bytes(), &infoResponse)

	return infoResponse.Coins
}

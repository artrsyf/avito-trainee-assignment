package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/dto"
	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/usecase"
	"github.com/artrsyf/avito-trainee-assignment/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type TransactionHandler struct {
	transactionUC usecase.TransactionUsecaseI
	validator     *validator.Validate
	logger        *logrus.Logger
}

func NewTransactionHandler(transactionUsecase usecase.TransactionUsecaseI, validator *validator.Validate, logger *logrus.Logger) *TransactionHandler {
	return &TransactionHandler{
		transactionUC: transactionUsecase,
		logger:        logger,
		validator:     validator,
	}
}

func (h *TransactionHandler) SendCoins(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Incoming SendCoins request")

	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.WithError(err).Error("Failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"errors": "bad request"})
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			h.logger.WithError(err).Warn("Failed to close request body")
		}
	}()

	sendCoinsRequest := &dto.SendCoinsRequest{}
	err = json.Unmarshal(body, sendCoinsRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"errors": "bad request"})
		return
	}

	if err := sendCoinsRequest.ValidateSendCoinsRequest(h.validator); err != nil {
		h.logger.WithError(err).Warn("Failed validation for send coins request")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"errors": err.Error()})
		return
	}

	senderUsername, ok := ctx.Value(middleware.UsernameContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"errors": "internal error"})
	}

	transactionEntity := &entity.Transaction{
		SenderUsername:   senderUsername,
		ReceiverUsername: sendCoinsRequest.ReceiverUsername,
		Amount:           sendCoinsRequest.Amount,
	}

	err = h.transactionUC.Create(ctx, transactionEntity)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"stack": string(debug.Stack()),
		}).Debug("Transaction create error handling")

		switch err {
		case entity.ErrNotEnoughBalance:
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"errors": "not enough balance"})

		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"errors": "internal error"})
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

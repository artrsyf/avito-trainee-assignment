package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/dto"
	transaction "github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/usecase"
	userEntity "github.com/artrsyf/avito-trainee-assignment/internal/user/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/middleware"
	jsonresponse "github.com/artrsyf/avito-trainee-assignment/pkg/json_response"
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

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.WithError(err).Error("Failed to read request body")
		jsonresponse.JsonResponse(
			w,
			http.StatusBadRequest,
			map[string]string{"errors": "bad request"},
		)
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
		jsonresponse.JsonResponse(
			w,
			http.StatusBadRequest,
			map[string]string{"errors": "bad request"},
		)
		return
	}

	if err := sendCoinsRequest.ValidateSendCoinsRequest(h.validator); err != nil {
		h.logger.WithError(err).Warn("Failed validation for send coins request")
		jsonresponse.JsonResponse(
			w,
			http.StatusBadRequest,
			map[string]string{"errors": err.Error()},
		)
		return
	}

	senderUsername, ok := ctx.Value(middleware.UsernameContextKey).(string)
	if !ok {
		jsonresponse.JsonResponse(
			w,
			http.StatusInternalServerError,
			map[string]string{"errors": "internal error"},
		)
		return
	}

	if sendCoinsRequest.ReceiverUsername == senderUsername {
		jsonresponse.JsonResponse(
			w,
			http.StatusBadRequest,
			map[string]string{"errors": "money transfer to yourself is not allowed"},
		)
		return
	}

	transactionEntity := &transaction.Transaction{
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
		case transaction.ErrNotEnoughBalance:
			jsonresponse.JsonResponse(
				w,
				http.StatusBadRequest,
				map[string]string{"errors": "not enough balance"},
			)
		case userEntity.ErrIsNotExist:
			jsonresponse.JsonResponse(
				w,
				http.StatusBadRequest,
				map[string]string{"errors": "can't find such user"},
			)
		default:
			jsonresponse.JsonResponse(
				w,
				http.StatusInternalServerError,
				map[string]string{"errors": "internal error"},
			)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

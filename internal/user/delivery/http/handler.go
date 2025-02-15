package http

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/artrsyf/avito-trainee-assignment/internal/user/usecase"
	"github.com/artrsyf/avito-trainee-assignment/middleware"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	userUC usecase.UserUsecaseI
	logger *logrus.Logger
}

func NewTransactionHandler(userUsecase usecase.UserUsecaseI, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		userUC: userUsecase,
		logger: logger,
	}
}

func (h *UserHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Incoming GetInfo request")

	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	senderUserID, ok := ctx.Value(middleware.UserIDContextKey).(uint)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"errors": "internal error"})
	}

	getInfoResponse, err := h.userUC.GetInfoById(ctx, senderUserID)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"stack": string(debug.Stack()),
		}).Error("GetInfoById error handling")

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"errors": "internal error"})
	}

	response, err := json.Marshal(getInfoResponse)
	if err != nil {
		h.logger.WithError(err).Error("Failed to marshal auth response")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"errors": "internal error"})
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		h.logger.WithError(err).Error("Failed to write get info response")
	}
}

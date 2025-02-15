package http

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/artrsyf/avito-trainee-assignment/internal/user/usecase"
	"github.com/artrsyf/avito-trainee-assignment/middleware"
	JSONResponse "github.com/artrsyf/avito-trainee-assignment/pkg/json_response"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	userUC usecase.UserUsecaseI
	logger *logrus.Logger
}

func NewUserHandler(userUsecase usecase.UserUsecaseI, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		userUC: userUsecase,
		logger: logger,
	}
}

func (h *UserHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Incoming GetInfo request")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	senderUserID, ok := ctx.Value(middleware.UserIDContextKey).(uint)
	if !ok {
		JSONResponse.JSONResponse(
			w,
			http.StatusInternalServerError,
			map[string]string{"errors": "internal error"},
		)
		return
	}

	getInfoResponse, err := h.userUC.GetInfoByID(ctx, senderUserID)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"stack": string(debug.Stack()),
		}).Error("GetInfoById error handling")

		JSONResponse.JSONResponse(
			w,
			http.StatusInternalServerError,
			map[string]string{"errors": "internal error"},
		)
	}

	response, err := json.Marshal(getInfoResponse)
	if err != nil {
		h.logger.WithError(err).Error("Failed to marshal auth response")
		JSONResponse.JSONResponse(
			w,
			http.StatusInternalServerError,
			map[string]string{"errors": "internal error"},
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		h.logger.WithError(err).Error("Failed to write get info response")
	}
}

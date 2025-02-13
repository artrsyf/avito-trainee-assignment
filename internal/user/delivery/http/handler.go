package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/artrsyf/avito-trainee-assignment/internal/user/usecase"
	"github.com/artrsyf/avito-trainee-assignment/middleware"
)

type UserHandler struct {
	userUC usecase.UserUsecaseI
}

func NewTransactionHandler(userUsecase usecase.UserUsecaseI) *UserHandler {
	return &UserHandler{
		userUC: userUsecase,
	}
}

func (h *UserHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	senderUserID := ctx.Value(middleware.UserIDContextKey).(uint)

	getInfoResponse, err := h.userUC.GetInfoById(ctx, senderUserID)
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	response, err := json.Marshal(getInfoResponse)
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

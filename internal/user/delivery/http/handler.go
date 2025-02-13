package http

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	senderUserID := r.Context().Value(middleware.UserIDContextKey).(uint)

	getInfoResponse, err := h.userUC.GetInfoById(senderUserID)
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

package http

import (
	"fmt"
	"net/http"

	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/dto"
	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/usecase"
	"github.com/artrsyf/avito-trainee-assignment/middleware"
	"github.com/gorilla/mux"
)

type PurchaseHandler struct {
	purchaseUC usecase.PurchaseUsecaseI
}

func NewPurchaseHandler(purchaseUsecase usecase.PurchaseUsecaseI) *PurchaseHandler {
	return &PurchaseHandler{
		purchaseUC: purchaseUsecase,
	}
}

func (h *PurchaseHandler) BuyItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	purchaseTypeName := vars["item"]

	/*TODO check senderUsername work*/
	customerUserID := r.Context().Value(middleware.UserIDContextKey).(uint)

	purchaseItemRequest := &dto.PurchaseItemRequest{
		PurchaseTypeName: purchaseTypeName,
		UserID:           customerUserID,
	}

	err := h.purchaseUC.Create(purchaseItemRequest)
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

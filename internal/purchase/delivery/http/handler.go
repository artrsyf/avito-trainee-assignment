package http

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/dto"
	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/usecase"
	"github.com/artrsyf/avito-trainee-assignment/middleware"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type PurchaseHandler struct {
	purchaseUC usecase.PurchaseUsecaseI
	logger     *logrus.Logger
}

func NewPurchaseHandler(purchaseUsecase usecase.PurchaseUsecaseI, logger *logrus.Logger) *PurchaseHandler {
	return &PurchaseHandler{
		purchaseUC: purchaseUsecase,
		logger:     logger,
	}
}

func (h *PurchaseHandler) BuyItem(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Incoming BuyItem request")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	purchaseTypeName := vars["item"]

	customerUserID, ok := ctx.Value(middleware.UserIDContextKey).(uint)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"errors": "internal error"})
	}

	purchaseItemRequest := &dto.PurchaseItemRequest{
		PurchaseTypeName: purchaseTypeName,
		UserID:           customerUserID,
	}

	err := h.purchaseUC.Create(ctx, purchaseItemRequest)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"stack": string(debug.Stack()),
		}).Debug("Purchase create error handling")
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

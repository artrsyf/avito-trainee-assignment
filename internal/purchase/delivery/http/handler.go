package http

import (
	"context"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/dto"
	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/usecase"
	"github.com/artrsyf/avito-trainee-assignment/middleware"
	JSONResponse "github.com/artrsyf/avito-trainee-assignment/pkg/json_response"
)

type PurchaseHandler struct {
	purchaseUC usecase.PurchaseUsecaseI
	validate   *validator.Validate
	logger     *logrus.Logger
}

func NewPurchaseHandler(
	purchaseUsecase usecase.PurchaseUsecaseI,
	validate *validator.Validate,
	logger *logrus.Logger,
) *PurchaseHandler {
	return &PurchaseHandler{
		purchaseUC: purchaseUsecase,
		validate:   validate,
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
		JSONResponse.JSONResponse(
			w,
			http.StatusInternalServerError,
			map[string]string{"errors": "internal error"},
		)
		return
	}

	purchaseItemRequest := &dto.PurchaseItemRequest{
		PurchaseTypeName: purchaseTypeName,
		UserID:           customerUserID,
	}
	if err := purchaseItemRequest.ValidatePurchaseRequest(h.validate); err != nil {
		h.logger.WithError(err).Warn("Failed validation for purchase request")
		JSONResponse.JSONResponse(
			w,
			http.StatusBadRequest,
			map[string]string{"errors": err.Error()},
		)
		return
	}

	err := h.purchaseUC.Create(ctx, purchaseItemRequest)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"stack": string(debug.Stack()),
		}).Debug("Purchase create error handling")

		switch err {
		case entity.ErrNotEnoughBalance:
			JSONResponse.JSONResponse(
				w,
				http.StatusBadRequest,
				map[string]string{"errors": "not enough balance"},
			)
		case entity.ErrNotExistedProduct:
			JSONResponse.JSONResponse(
				w,
				http.StatusNotFound,
				map[string]string{"errors": "item not found"},
			)
		default:
			JSONResponse.JSONResponse(
				w,
				http.StatusInternalServerError,
				map[string]string{"errors": "internal error"},
			)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

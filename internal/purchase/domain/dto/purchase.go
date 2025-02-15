package dto

import (
	"errors"

	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/entity"
	"github.com/go-playground/validator/v10"
)

type PurchaseItemRequest struct {
	PurchaseTypeName string `validate:"required"`
	UserID           uint   `validate:"required,gt=0"`
}

func (req *PurchaseItemRequest) ValidatePurchaseRequest(validate *validator.Validate) error {
	err := validate.Struct(req)

	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		switch err.Tag() {
		case "required":
			return errors.New(field + " is required")
		case "gt":
			return errors.New("something went wrong")
		default:
			return errors.New("something went wrong")
		}
	}
	return nil
}

func PurchaseItemRequestToEntity(purchaseItemRequest *PurchaseItemRequest) *entity.Purchase {
	return &entity.Purchase{
		PurchaserId:      purchaseItemRequest.UserID,
		PurchaseTypeName: purchaseItemRequest.PurchaseTypeName,
	}
}

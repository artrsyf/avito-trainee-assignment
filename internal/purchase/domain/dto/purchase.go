package dto

import (
	"errors"

	"github.com/go-playground/validator/v10"

	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/entity"
)

type PurchaseItemRequest struct {
	PurchaseTypeName string `validate:"required,min=1"`
	UserID           uint   `validate:"required,gt=0"`
}

func (req *PurchaseItemRequest) ValidatePurchaseRequest(validate *validator.Validate) error {
	err := validate.Struct(req)
	if err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, err := range validationErrs {
				field := err.Field()

				switch err.Tag() {
				case "required":
					return errors.New(field + " is required")
				case "gt":
					return errors.New(field + " must be greater than 0")
				default:
					return errors.New(field + " is invalid")
				}
			}
		}
		return err
	}
	return nil
}

func PurchaseItemRequestToEntity(purchaseItemRequest *PurchaseItemRequest) *entity.Purchase {
	return &entity.Purchase{
		PurchaserID:      purchaseItemRequest.UserID,
		PurchaseTypeName: purchaseItemRequest.PurchaseTypeName,
	}
}

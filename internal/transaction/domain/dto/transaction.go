package dto

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type SendCoinsRequest struct {
	ReceiverUsername string `json:"toUser" validate:"required,min=3,max=50"`
	Amount           uint   `json:"amount" validate:"required,gt=0"`
}

func (req *SendCoinsRequest) ValidateSendCoinsRequest(validate *validator.Validate) error {
	err := validate.Struct(req)
	if err != nil {
		// Преобразуем ошибку в тип validator.ValidationErrors
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, err := range validationErrs {
				field := err.Field()
				switch err.Tag() {
				case "required":
					return errors.New(field + " is required")
				case "min":
					return errors.New(field + " is too short")
				case "max":
					return errors.New(field + " is too long")
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

package dto

import "github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/entity"

type PurchaseItemRequest struct {
	PurchaseTypeName string
	UserID           uint
}

func PurchaseItemRequestToEntity(purchaseItemRequest *PurchaseItemRequest) *entity.Purchase {
	return &entity.Purchase{
		PurchaserId:      purchaseItemRequest.UserID,
		PurchaseTypeName: purchaseItemRequest.PurchaseTypeName,
	}
}

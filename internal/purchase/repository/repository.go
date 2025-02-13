package repository

import (
	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/model"
)

type PurchaseRepositoryI interface {
	Create(purchase *entity.Purchase) (*model.Purchase, error)
	GetProductByType(purchaseTypeName string) (*model.PurchaseType, error)
	GetPurchasesByUserId(userID uint) (entity.Inventory, error)
}

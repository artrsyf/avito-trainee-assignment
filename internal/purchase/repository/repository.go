package repository

import (
	"context"

	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/model"
)

type PurchaseRepositoryI interface {
	Create(ctx context.Context, purchase *entity.Purchase) (*model.Purchase, error)
	GetProductByType(ctx context.Context, purchaseTypeName string) (*model.PurchaseType, error)
	GetPurchasesByUserId(ctx context.Context, userID uint) (entity.Inventory, error)
}

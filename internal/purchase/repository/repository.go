package repository

import (
	"context"

	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/model"
	"github.com/artrsyf/avito-trainee-assignment/pkg/uow"
)

//go:generate mockgen -source=repository.go -destination=mock_repository/purchase_mock.go -package=mock_repository MockPurchaseRepository
type PurchaseRepositoryI interface {
	Create(ctx context.Context, uow uow.Executor, purchase *entity.Purchase) (*model.Purchase, error)
	GetProductByType(ctx context.Context, purchaseTypeName string) (*model.PurchaseType, error)
	GetPurchasesByUserID(ctx context.Context, userID uint) (entity.Inventory, error)
}

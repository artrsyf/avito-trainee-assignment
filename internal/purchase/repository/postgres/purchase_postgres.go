package postgres

import (
	"database/sql"

	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/model"
)

type PurchasePostgresRepository struct {
	DB *sql.DB
}

func NewPurchasePostgresRepository(db *sql.DB) *PurchasePostgresRepository {
	return &PurchasePostgresRepository{
		DB: db,
	}
}

func (repo *PurchasePostgresRepository) Create(purchase *entity.Purchase) (*model.Purchase, error) {
	createdPurchase := model.Purchase{}
	var purchaseTypeID int

	err := repo.DB.
		QueryRow("SELECT id FROM purchase_types WHERE name = $1", purchase.PurchaseTypeName).Scan(&purchaseTypeID)
	if err != nil {
		return nil, err
	}

	err = repo.DB.QueryRow("INSERT INTO purchases (purchaser_id, purchase_type_id) VALUES ($1, $2) RETURNING id, purchaser_id, purchase_type_id", purchase.PurchaserId, purchaseTypeID).
		Scan(&createdPurchase.ID, &createdPurchase.PurchaserId, &createdPurchase.PurchaseTypeId)
	if err != nil {
		return nil, err
	}

	return &createdPurchase, nil
}

func (repo *PurchasePostgresRepository) GetProductByType(purchaseTypeName string) (*model.PurchaseType, error) {
	purchaseType := model.PurchaseType{}

	err := repo.DB.
		QueryRow("SELECT id, name, cost FROM purchase_types WHERE name = $1", purchaseTypeName).
		Scan(&purchaseType.ID, &purchaseType.Name, &purchaseType.Cost)
	if err != nil {
		return nil, err
	}

	return &purchaseType, nil
}

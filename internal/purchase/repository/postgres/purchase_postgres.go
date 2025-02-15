package postgres

import (
	"context"
	"database/sql"

	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/model"
	"github.com/sirupsen/logrus"
)

type PurchasePostgresRepository struct {
	DB     *sql.DB
	logger *logrus.Logger
}

func NewPurchasePostgresRepository(db *sql.DB, logger *logrus.Logger) *PurchasePostgresRepository {
	return &PurchasePostgresRepository{
		DB:     db,
		logger: logger,
	}
}

func (repo *PurchasePostgresRepository) Create(
	ctx context.Context,
	purchase *entity.Purchase,
) (*model.Purchase, error) {
	createdPurchase := model.Purchase{}
	var purchaseTypeID int

	err := repo.DB.QueryRowContext(
		ctx,
		"SELECT id FROM purchase_types WHERE name = $1",
		purchase.PurchaseTypeName,
	).Scan(&purchaseTypeID)
	if err == sql.ErrNoRows {
		repo.logger.WithError(err).Error("Failed to select not existed product")
		return nil, entity.ErrNotExistedProduct
	}
	if err != nil {
		repo.logger.WithError(err).Error("Failed to select purchase type id")
		return nil, err
	}

	err = repo.DB.QueryRowContext(
		ctx,
		`INSERT INTO purchases (purchaser_id, purchase_type_id) 
		VALUES ($1, $2) 
		RETURNING id, purchaser_id, purchase_type_id`,
		purchase.PurchaserID, purchaseTypeID,
	).Scan(
		&createdPurchase.ID,
		&createdPurchase.PurchaserID,
		&createdPurchase.PurchaseTypeID,
	)

	if err != nil {
		repo.logger.WithError(err).Error("Failed to insert purchase")
		return nil, err
	}

	return &createdPurchase, nil
}

func (repo *PurchasePostgresRepository) GetProductByType(
	ctx context.Context,
	purchaseTypeName string,
) (*model.PurchaseType, error) {
	purchaseType := model.PurchaseType{}

	err := repo.DB.QueryRowContext(
		ctx,
		"SELECT id, name, cost FROM purchase_types WHERE name = $1",
		purchaseTypeName,
	).Scan(&purchaseType.ID, &purchaseType.Name, &purchaseType.Cost)
	if err == sql.ErrNoRows {
		repo.logger.WithError(err).Error("Failed to select not existed product")
		return nil, entity.ErrNotExistedProduct
	}
	if err != nil {
		repo.logger.WithError(err).Error("Failed to select product by type")
		return nil, err
	}

	return &purchaseType, nil
}

func (repo *PurchasePostgresRepository) GetPurchasesByUserID(
	ctx context.Context,
	userID uint,
) (entity.Inventory, error) {
	rows, err := repo.DB.QueryContext(ctx, `
		SELECT pt.name, COUNT(p.id) as quantity
		FROM purchases p
		JOIN purchase_types pt ON p.purchase_type_id = pt.id
		WHERE p.purchaser_id = $1
		GROUP BY pt.name`, userID)
	if err != nil {
		repo.logger.WithError(err).Error("Failed to select user inventory")
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.WithError(err).Warn("Failed to close rows selecting user inventory")
		}
	}()

	inventory := entity.Inventory{}
	for rows.Next() {
		currentPurchaseGroup := entity.PurchaseGroup{}
		err := rows.Scan(&currentPurchaseGroup.PurchaseTypeName, &currentPurchaseGroup.Quantity)
		if err != nil {
			repo.logger.WithError(err).Error("Failed to select purchase group")
			return nil, err
		}

		inventory = append(inventory, currentPurchaseGroup)
	}

	return inventory, nil
}

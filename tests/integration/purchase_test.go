package integration

import (
	"context"
	"testing"

	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/dto"
	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/entity"
	purchaseRepo "github.com/artrsyf/avito-trainee-assignment/internal/purchase/repository/postgres"
	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/usecase"
	userRepo "github.com/artrsyf/avito-trainee-assignment/internal/user/repository/postgres"
	uow "github.com/artrsyf/avito-trainee-assignment/pkg/uow/postgres"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestPurchaseUsecase_Integration(t *testing.T) {
	userRepo := userRepo.NewUserPostgresRepository(DB, logrus.New())
	purchaseRepo := purchaseRepo.NewPurchasePostgresRepository(DB, logrus.New())

	uow := uow.NewFactory(DB)

	uc := usecase.NewPurchaseUsecase(purchaseRepo, userRepo, uow, logrus.New())
	ctx := context.Background()

	t.Run("successful purchase", func(t *testing.T) {
		SetupTestData(t, DB)
		userID := CreateTestUser(t, "user1", 1000)
		CreatePurchaseType(t, DB, "premium", 500)

		req := &dto.PurchaseItemRequest{
			UserID:           userID,
			PurchaseTypeName: "premium",
		}

		err := uc.Create(ctx, req)
		require.NoError(t, err)

		user, err := userRepo.GetByID(ctx, userID)
		require.NoError(t, err)
		require.Equal(t, uint(500), user.Coins)

		var count int
		err = DB.QueryRowContext(ctx,
			"SELECT COUNT(*) FROM purchases WHERE purchaser_id = $1", userID).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("insufficient balance", func(t *testing.T) {
		SetupTestData(t, DB)
		userID := CreateTestUser(t, "user2", 300)
		CreatePurchaseType(t, DB, "vip", 500)

		req := &dto.PurchaseItemRequest{
			UserID:           userID,
			PurchaseTypeName: "vip",
		}

		err := uc.Create(ctx, req)
		require.ErrorIs(t, err, entity.ErrNotEnoughBalance)

		user, err := userRepo.GetByID(ctx, userID)
		require.NoError(t, err)
		require.Equal(t, uint(300), user.Coins)
	})

	t.Run("product not found", func(t *testing.T) {
		SetupTestData(t, DB)
		userID := CreateTestUser(t, "user3", 1000)

		req := &dto.PurchaseItemRequest{
			UserID:           userID,
			PurchaseTypeName: "unknown",
		}

		err := uc.Create(ctx, req)
		require.Error(t, err)
	})

	t.Run("user not found", func(t *testing.T) {
		SetupTestData(t, DB)
		CreatePurchaseType(t, DB, "basic", 100)

		req := &dto.PurchaseItemRequest{
			UserID:           999,
			PurchaseTypeName: "basic",
		}

		err := uc.Create(ctx, req)
		require.Error(t, err)
	})

	t.Run("transaction rollback on error", func(t *testing.T) {
		SetupTestData(t, DB)
		userID := CreateTestUser(t, "user4", 1000)
		CreatePurchaseType(t, DB, "test", 800)

		_, err := DB.Exec("DROP TABLE purchases")
		require.NoError(t, err)

		req := &dto.PurchaseItemRequest{
			UserID:           userID,
			PurchaseTypeName: "test",
		}

		err = uc.Create(ctx, req)
		require.Error(t, err)

		user, err := userRepo.GetByID(ctx, userID)
		require.NoError(t, err)
		require.Equal(t, uint(1000), user.Coins)

		createPurchasesTable(t)
	})
}

func createPurchasesTable(t *testing.T) {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS purchases (
			id SERIAL PRIMARY KEY,
			purchaser_id INTEGER NOT NULL,
			purchase_type_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	require.NoError(t, err)
}

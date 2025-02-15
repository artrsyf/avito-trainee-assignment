package integration

import (
	"context"
	"testing"

	purchaseEntity "github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/entity"
	purchaseRepo "github.com/artrsyf/avito-trainee-assignment/internal/purchase/repository/postgres"
	transactionEntity "github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/entity"
	transactionRepo "github.com/artrsyf/avito-trainee-assignment/internal/transaction/repository/postgres"
	userEntity "github.com/artrsyf/avito-trainee-assignment/internal/user/domain/entity"
	userRepo "github.com/artrsyf/avito-trainee-assignment/internal/user/repository/postgres"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/usecase"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestUserUsecase_Integration(t *testing.T) {
	userRepo := userRepo.NewUserPostgresRepository(DB, logrus.New())
	purchaseRepo := purchaseRepo.NewPurchasePostgresRepository(DB, logrus.New())
	transactionRepo := transactionRepo.NewTransactionPostgresRepository(DB, logrus.New())

	uc := usecase.NewUserUsecase(purchaseRepo, transactionRepo, userRepo, logrus.New())
	ctx := context.Background()

	SetupTestData(t, DB)

	t.Run("get user info with empty history", func(t *testing.T) {
		userID := CreateTestUser(t, "user1", 1000)

		info, err := uc.GetInfoByID(ctx, userID)
		require.NoError(t, err)

		require.Equal(t, uint(1000), info.Coins)
		require.Empty(t, info.Inventory)
		require.Empty(t, info.CoinHistory.SentHistory)
		require.Empty(t, info.CoinHistory.ReceivedHistory)
	})

	t.Run("get user info with purchases", func(t *testing.T) {
		userID := CreateTestUser(t, "user2", 500)
		CreatePurchaseType(t, DB, "item1", uint(100))
		CreatePurchaseType(t, DB, "item2", uint(200))

		createPurchase(t, userID, "item1")
		createPurchase(t, userID, "item1")
		createPurchase(t, userID, "item2")

		info, err := uc.GetInfoByID(ctx, userID)
		require.NoError(t, err)

		require.Len(t, info.Inventory, 2)
		require.Equal(t, uint(2), findPurchaseQuantity(info.Inventory, "item1"))
		require.Equal(t, uint(1), findPurchaseQuantity(info.Inventory, "item2"))
	})

	t.Run("get user info with transactions", func(t *testing.T) {
		senderID := CreateTestUser(t, "sender", 1000)
		receiver1ID := CreateTestUser(t, "receiver1", 0)
		receiver2ID := CreateTestUser(t, "receiver2", 0)

		createTransaction(t, senderID, receiver1ID, 100)
		createTransaction(t, senderID, receiver1ID, 200)
		createTransaction(t, senderID, receiver2ID, 300)

		senderInfo, err := uc.GetInfoByID(ctx, senderID)
		require.NoError(t, err)

		require.Len(t, senderInfo.CoinHistory.SentHistory, 2)
		require.Equal(t, uint(300), findSentAmount(senderInfo.CoinHistory.SentHistory, "receiver1"))
		require.Equal(t, uint(300), findSentAmount(senderInfo.CoinHistory.SentHistory, "receiver2"))

		receiverInfo, err := uc.GetInfoByID(ctx, receiver1ID)
		require.NoError(t, err)

		require.Len(t, receiverInfo.CoinHistory.ReceivedHistory, 1)
		require.Equal(t, uint(300), findReceivedAmount(receiverInfo.CoinHistory.ReceivedHistory, "sender"))
	})

	t.Run("user not found", func(t *testing.T) {
		_, err := uc.GetInfoByID(ctx, 9999)
		require.Error(t, err)
		require.Equal(t, err, userEntity.ErrIsNotExist)
	})
}

func createPurchase(t *testing.T, userID uint, itemType string) {
	_, err := DB.Exec(`
		INSERT INTO purchases (purchaser_id, purchase_type_id)
		SELECT $1, id FROM purchase_types WHERE name = $2
	`, userID, itemType)
	require.NoError(t, err)
}

func createTransaction(t *testing.T, senderID, receiverID, amount uint) {
	_, err := DB.Exec(
		"INSERT INTO transactions (sender_user_id, receiver_user_id, amount) VALUES ($1, $2, $3)",
		senderID, receiverID, amount,
	)
	require.NoError(t, err)
}

func findPurchaseQuantity(inventory purchaseEntity.Inventory, name string) uint {
	for _, item := range inventory {
		if item.PurchaseTypeName == name {
			return item.Quantity
		}
	}
	return 0
}

func findSentAmount(history transactionEntity.SentHistory, receiver string) uint {
	for _, tx := range history {
		if tx.ReceiverUsername == receiver {
			return tx.Amount
		}
	}
	return 0
}

func findReceivedAmount(history transactionEntity.ReceivedHistory, sender string) uint {
	for _, tx := range history {
		if tx.SenderUsername == sender {
			return tx.Amount
		}
	}
	return 0
}

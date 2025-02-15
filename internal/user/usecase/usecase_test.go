package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"

	purchaseEntity "github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/entity"
	mockPurchase "github.com/artrsyf/avito-trainee-assignment/internal/purchase/repository/mock_repository"
	transactionEntity "github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/entity"
	mockTransaction "github.com/artrsyf/avito-trainee-assignment/internal/transaction/repository/mock_repository"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/model"
	mockUser "github.com/artrsyf/avito-trainee-assignment/internal/user/repository/mock_repository"
)

func TestUserUsecase_GetInfoById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockUser.NewMockUserRepositoryI(ctrl)
	mockPurchaseRepo := mockPurchase.NewMockPurchaseRepositoryI(ctrl)
	mockTransactionRepo := mockTransaction.NewMockTransactionRepositoryI(ctrl)

	uc := NewUserUsecase(
		mockPurchaseRepo,
		mockTransactionRepo,
		mockUserRepo,
		logrus.New(),
	)

	ctx := context.Background()
	userID := uint(1)
	testError := errors.New("test error")

	t.Run("successful response", func(t *testing.T) {
		user := &model.User{Coins: 100}
		productType := "t-shirt"
		purchaseGroup := purchaseEntity.PurchaseGroup{
			PurchaseTypeName: productType,
			Quantity:         1,
		}
		inventory := &purchaseEntity.Inventory{purchaseGroup}
		sentTrancationGroup := transactionEntity.SentTransactionGroup{
			ReceiverUsername: "user2",
			Amount:           500,
		}
		receivedTrancationGroup := transactionEntity.ReceivedTransactionGroup{
			SenderUsername: "user2",
			Amount:         300,
		}
		sent := &transactionEntity.SentHistory{sentTrancationGroup}
		received := &transactionEntity.ReceivedHistory{receivedTrancationGroup}

		mockUserRepo.EXPECT().GetById(ctx, userID).Return(user, nil)
		mockPurchaseRepo.EXPECT().GetPurchasesByUserId(ctx, userID).Return(*inventory, nil)
		mockTransactionRepo.EXPECT().GetSentByUserID(ctx, userID).Return(*sent, nil)
		mockTransactionRepo.EXPECT().GetReceivedByUserID(ctx, userID).Return(*received, nil)

		resp, err := uc.GetInfoById(ctx, userID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if resp.Coins != user.Coins {
			t.Errorf("expected coins %d, got %d", user.Coins, resp.Coins)
		}
		if len(resp.Inventory) != 1 || resp.Inventory[0].Quantity != 1 || resp.Inventory[0].PurchaseTypeName != productType {
			t.Error("invalid inventory items")
		}
		if len(resp.CoinHistory.SentHistory) != 1 {
			t.Error("invalid sent transactions count")
		}
		if len(resp.CoinHistory.ReceivedHistory) != 1 {
			t.Error("invalid received transactions count")
		}
	})

	t.Run("user repo error", func(t *testing.T) {
		mockUserRepo.EXPECT().GetById(ctx, userID).Return(nil, testError)
		mockPurchaseRepo.EXPECT().GetPurchasesByUserId(gomock.Any(), gomock.Any()).Times(0)
		mockTransactionRepo.EXPECT().GetSentByUserID(gomock.Any(), gomock.Any()).Times(0)
		mockTransactionRepo.EXPECT().GetReceivedByUserID(gomock.Any(), gomock.Any()).Times(0)

		_, err := uc.GetInfoById(ctx, userID)
		if !errors.Is(err, testError) {
			t.Errorf("expected error %v, got %v", testError, err)
		}
	})

	t.Run("purchase repo error", func(t *testing.T) {
		user := &model.User{Coins: 100}

		mockUserRepo.EXPECT().GetById(ctx, userID).Return(user, nil)
		mockPurchaseRepo.EXPECT().GetPurchasesByUserId(ctx, userID).Return(nil, testError)
		mockTransactionRepo.EXPECT().GetSentByUserID(gomock.Any(), gomock.Any()).Times(0)
		mockTransactionRepo.EXPECT().GetReceivedByUserID(gomock.Any(), gomock.Any()).Times(0)

		_, err := uc.GetInfoById(ctx, userID)
		if !errors.Is(err, testError) {
			t.Errorf("expected error %v, got %v", testError, err)
		}
	})

	t.Run("sent transactions error", func(t *testing.T) {
		user := &model.User{Coins: 100}
		inventory := &purchaseEntity.Inventory{}

		mockUserRepo.EXPECT().GetById(ctx, userID).Return(user, nil)
		mockPurchaseRepo.EXPECT().GetPurchasesByUserId(ctx, userID).Return(*inventory, nil)
		mockTransactionRepo.EXPECT().GetSentByUserID(ctx, userID).Return(nil, testError)
		mockTransactionRepo.EXPECT().GetReceivedByUserID(gomock.Any(), gomock.Any()).Times(0)

		_, err := uc.GetInfoById(ctx, userID)
		if !errors.Is(err, testError) {
			t.Errorf("expected error %v, got %v", testError, err)
		}
	})

	t.Run("received transactions error", func(t *testing.T) {
		user := &model.User{Coins: 100}
		inventory := &purchaseEntity.Inventory{}
		sent := &transactionEntity.SentHistory{}

		mockUserRepo.EXPECT().GetById(ctx, userID).Return(user, nil)
		mockPurchaseRepo.EXPECT().GetPurchasesByUserId(ctx, userID).Return(*inventory, nil)
		mockTransactionRepo.EXPECT().GetSentByUserID(ctx, userID).Return(*sent, nil)
		mockTransactionRepo.EXPECT().GetReceivedByUserID(ctx, userID).Return(nil, testError)

		_, err := uc.GetInfoById(ctx, userID)
		if !errors.Is(err, testError) {
			t.Errorf("expected error %v, got %v", testError, err)
		}
	})
}

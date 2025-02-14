package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/entity"
	transactionModel "github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/model"
	mockTransaction "github.com/artrsyf/avito-trainee-assignment/internal/transaction/repository/mock_repository"
	userModel "github.com/artrsyf/avito-trainee-assignment/internal/user/domain/model"
	mockUser "github.com/artrsyf/avito-trainee-assignment/internal/user/repository/mock_repository"
	"github.com/artrsyf/avito-trainee-assignment/pkg/uow"
	"github.com/artrsyf/avito-trainee-assignment/pkg/uow/mock_uow"
)

func TestTransactionUsecase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxRepo := mockTransaction.NewMockTransactionRepositoryI(ctrl)
	mockUserRepo := mockUser.NewMockUserRepositoryI(ctrl)
	mockUow := mock_uow.NewMockUnitOfWorkI(ctrl)

	uc := NewTransactionUsecase(mockTxRepo, mockUserRepo, mockUow)

	ctx := context.Background()
	testTransaction := &entity.Transaction{
		SenderUsername:   "sender",
		ReceiverUsername: "receiver",
		Amount:           100,
	}

	t.Run("successful transaction", func(t *testing.T) {
		sender := &userModel.User{ID: 1, Username: "sender", Coins: 200}
		receiver := &userModel.User{ID: 2, Username: "receiver", Coins: 50}

		mockUserRepo.EXPECT().GetByUsername(ctx, "sender").Return(sender, nil)
		mockUserRepo.EXPECT().GetByUsername(ctx, "receiver").Return(receiver, nil)
		mockUow.EXPECT().Begin(ctx).Return(nil)
		mockUserRepo.EXPECT().Update(ctx, mockUow, gomock.Any()).DoAndReturn(
			func(ctx context.Context, uow uow.UnitOfWorkI, user *userModel.User) error {
				switch user.ID {
				case 1:
					if user.Coins != 100 {
						t.Errorf("expected sender to have 100 coins, got %d", user.Coins)
					}
				case 2:
					if user.Coins != 150 {
						t.Errorf("expected receiver to have 150 coins, got %d", user.Coins)
					}
				default:
					t.Errorf("unexpected user ID: %d", user.ID)
				}
				return nil
			}).Times(2)
		mockTxRepo.EXPECT().Create(ctx, &transactionModel.Transaction{
			SenderUserID:   1,
			ReceiverUserID: 2,
			Amount:         100,
		}).Return(&transactionModel.Transaction{ID: 1}, nil)
		mockUow.EXPECT().Commit().Return(nil)

		err := uc.Create(ctx, testTransaction)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("insufficient balance", func(t *testing.T) {
		sender := &userModel.User{ID: 1, Username: "sender", Coins: 50}
		receiver := &userModel.User{ID: 2, Username: "receiver", Coins: 50}

		mockUserRepo.EXPECT().GetByUsername(ctx, "sender").Return(sender, nil)
		mockUserRepo.EXPECT().GetByUsername(ctx, "receiver").Return(receiver, nil)

		err := uc.Create(ctx, testTransaction)
		if !errors.Is(err, entity.ErrNotEnoughBalance) {
			t.Errorf("expected ErrNotEnoughBalance, got %v", err)
		}
	})

	t.Run("sender not found", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByUsername(ctx, "sender").Return(nil, errors.New("not found"))

		err := uc.Create(ctx, testTransaction)
		if err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("receiver not found", func(t *testing.T) {
		sender := &userModel.User{ID: 1, Username: "sender", Coins: 200}
		mockUserRepo.EXPECT().GetByUsername(ctx, "sender").Return(sender, nil)
		mockUserRepo.EXPECT().GetByUsername(ctx, "receiver").Return(nil, errors.New("not found"))

		err := uc.Create(ctx, testTransaction)
		if err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("update sender error", func(t *testing.T) {
		sender := &userModel.User{ID: 1, Username: "sender", Coins: 200}
		receiver := &userModel.User{ID: 2, Username: "receiver", Coins: 50}

		mockUserRepo.EXPECT().GetByUsername(ctx, "sender").Return(sender, nil)
		mockUserRepo.EXPECT().GetByUsername(ctx, "receiver").Return(receiver, nil)
		mockUow.EXPECT().Begin(ctx).Return(nil)
		mockUserRepo.EXPECT().Update(ctx, mockUow, gomock.Any()).Return(errors.New("update error"))
		mockUow.EXPECT().Rollback()

		err := uc.Create(ctx, testTransaction)
		if err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("transaction create error", func(t *testing.T) {
		sender := &userModel.User{ID: 1, Username: "sender", Coins: 200}
		receiver := &userModel.User{ID: 2, Username: "receiver", Coins: 50}

		mockUserRepo.EXPECT().GetByUsername(ctx, "sender").Return(sender, nil)
		mockUserRepo.EXPECT().GetByUsername(ctx, "receiver").Return(receiver, nil)
		mockUow.EXPECT().Begin(ctx).Return(nil)
		mockUserRepo.EXPECT().Update(ctx, mockUow, gomock.Any()).Times(2)
		mockTxRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil, errors.New("create error"))
		mockUow.EXPECT().Rollback()

		err := uc.Create(ctx, testTransaction)
		if err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("commit error", func(t *testing.T) {
		sender := &userModel.User{ID: 1, Username: "sender", Coins: 200}
		receiver := &userModel.User{ID: 2, Username: "receiver", Coins: 50}

		mockUserRepo.EXPECT().GetByUsername(ctx, "sender").Return(sender, nil)
		mockUserRepo.EXPECT().GetByUsername(ctx, "receiver").Return(receiver, nil)
		mockUow.EXPECT().Begin(ctx).Return(nil)
		mockUserRepo.EXPECT().Update(ctx, mockUow, gomock.Any()).Times(2)
		mockTxRepo.EXPECT().Create(ctx, gomock.Any()).Return(&transactionModel.Transaction{ID: 1}, nil)
		mockUow.EXPECT().Commit().Return(errors.New("commit error"))

		err := uc.Create(ctx, testTransaction)
		if err == nil {
			t.Error("expected error but got nil")
		}
	})
}

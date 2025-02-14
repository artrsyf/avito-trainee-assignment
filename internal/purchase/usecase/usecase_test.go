package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/dto"
	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/entity"
	purchaseModel "github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/model"
	mockPurchase "github.com/artrsyf/avito-trainee-assignment/internal/purchase/repository/mock_repository"
	userModel "github.com/artrsyf/avito-trainee-assignment/internal/user/domain/model"
	mockUser "github.com/artrsyf/avito-trainee-assignment/internal/user/repository/mock_repository"
	"github.com/artrsyf/avito-trainee-assignment/pkg/uow"
	"github.com/artrsyf/avito-trainee-assignment/pkg/uow/mock_uow"
)

func TestPurchaseUsecase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPurchaseRepo := mockPurchase.NewMockPurchaseRepositoryI(ctrl)
	mockUserRepo := mockUser.NewMockUserRepositoryI(ctrl)
	mockUow := mock_uow.NewMockUnitOfWorkI(ctrl)

	uc := NewPurchaseUsecase(mockPurchaseRepo, mockUserRepo, mockUow)

	ctx := context.Background()
	testRequest := &dto.PurchaseItemRequest{
		UserID:           1,
		PurchaseTypeName: "premium",
	}

	testPurchaseType := &purchaseModel.PurchaseType{
		Name: "premium",
		Cost: 100,
	}

	t.Run("successful purchase", func(t *testing.T) {
		user := &userModel.User{ID: 1, Coins: 200}

		mockUserRepo.EXPECT().GetById(ctx, uint(1)).Return(user, nil)
		mockPurchaseRepo.EXPECT().GetProductByType(ctx, "premium").Return(testPurchaseType, nil)
		mockUow.EXPECT().Begin(ctx).Return(nil)
		mockUserRepo.EXPECT().Update(ctx, mockUow, gomock.Any()).DoAndReturn(
			func(_ context.Context, _ uow.UnitOfWorkI, u *userModel.User) error {
				if u.Coins != 100 {
					t.Error("user coins not updated correctly")
				}
				return nil
			})
		mockPurchaseRepo.EXPECT().Create(ctx, &entity.Purchase{
			PurchaserId:      1,
			PurchaseTypeName: "premium",
		}).Return(&purchaseModel.Purchase{ID: 1}, nil)
		mockUow.EXPECT().Commit().Return(nil)

		err := uc.Create(ctx, testRequest)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("insufficient balance", func(t *testing.T) {
		user := &userModel.User{ID: 1, Coins: 50}

		mockUserRepo.EXPECT().GetById(ctx, uint(1)).Return(user, nil)
		mockPurchaseRepo.EXPECT().GetProductByType(ctx, "premium").Return(testPurchaseType, nil)

		err := uc.Create(ctx, testRequest)
		if !errors.Is(err, entity.ErrNotEnoughBalance) {
			t.Errorf("expected ErrNotEnoughBalance, got %v", err)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		mockUserRepo.EXPECT().GetById(ctx, uint(1)).Return(nil, errors.New("not found"))

		err := uc.Create(ctx, testRequest)
		if err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("product not found", func(t *testing.T) {
		user := &userModel.User{ID: 1, Coins: 200}

		mockUserRepo.EXPECT().GetById(ctx, uint(1)).Return(user, nil)
		mockPurchaseRepo.EXPECT().GetProductByType(ctx, "premium").Return(nil, errors.New("not found"))

		err := uc.Create(ctx, testRequest)
		if err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("begin transaction error", func(t *testing.T) {
		user := &userModel.User{ID: 1, Coins: 200}

		mockUserRepo.EXPECT().GetById(ctx, uint(1)).Return(user, nil)
		mockPurchaseRepo.EXPECT().GetProductByType(ctx, "premium").Return(testPurchaseType, nil)
		mockUow.EXPECT().Begin(ctx).Return(errors.New("tx error"))

		err := uc.Create(ctx, testRequest)
		if err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("user update error", func(t *testing.T) {
		user := &userModel.User{ID: 1, Coins: 200}

		mockUserRepo.EXPECT().GetById(ctx, uint(1)).Return(user, nil)
		mockPurchaseRepo.EXPECT().GetProductByType(ctx, "premium").Return(testPurchaseType, nil)
		mockUow.EXPECT().Begin(ctx).Return(nil)
		mockUserRepo.EXPECT().Update(ctx, mockUow, gomock.Any()).Return(errors.New("update error"))
		mockUow.EXPECT().Rollback()

		err := uc.Create(ctx, testRequest)
		if err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("purchase create error", func(t *testing.T) {
		user := &userModel.User{ID: 1, Coins: 200}

		mockUserRepo.EXPECT().GetById(ctx, uint(1)).Return(user, nil)
		mockPurchaseRepo.EXPECT().GetProductByType(ctx, "premium").Return(testPurchaseType, nil)
		mockUow.EXPECT().Begin(ctx).Return(nil)
		mockUserRepo.EXPECT().Update(ctx, mockUow, gomock.Any()).Return(nil)
		mockPurchaseRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil, errors.New("create error"))
		mockUow.EXPECT().Rollback()

		err := uc.Create(ctx, testRequest)
		if err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("commit error", func(t *testing.T) {
		user := &userModel.User{ID: 1, Coins: 200}

		mockUserRepo.EXPECT().GetById(ctx, uint(1)).Return(user, nil)
		mockPurchaseRepo.EXPECT().GetProductByType(ctx, "premium").Return(testPurchaseType, nil)
		mockUow.EXPECT().Begin(ctx).Return(nil)
		mockUserRepo.EXPECT().Update(ctx, mockUow, gomock.Any()).Return(nil)
		mockPurchaseRepo.EXPECT().Create(ctx, gomock.Any()).Return(&purchaseModel.Purchase{ID: 1}, nil)
		mockUow.EXPECT().Commit().Return(errors.New("commit error"))

		err := uc.Create(ctx, testRequest)
		if err == nil {
			t.Error("expected error but got nil")
		}
	})
}

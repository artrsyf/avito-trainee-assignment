package dto

import (
	sessionDTO "github.com/artrsyf/avito-trainee-assignment/internal/session/domain/dto"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/model"

	purchaseEntity "github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/entity"
	transactionEntity "github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/entity"

	"golang.org/x/crypto/bcrypt"
)

type SendCoinsRequest struct {
	ReceiverUsername string `json:"toUser"`
	Amount           uint   `json:"amount"`
}

type CoinHistory struct {
	ReceivedHistory transactionEntity.ReceivedHistory `json:"received"`
	SentHistory     transactionEntity.SentHistory     `json:"sent"`
}

type GetInfoResponse struct {
	Coins       uint                     `json:"coins"`
	Inventory   purchaseEntity.Inventory `json:"inventory"`
	CoinHistory CoinHistory              `json:"coinHistory"`
}

func CreateGetInfoResponse(userCoins uint, userInventory *purchaseEntity.Inventory, userSentTransactions *transactionEntity.SentHistory, userReceivedTransactions *transactionEntity.ReceivedHistory) *GetInfoResponse {
	return &GetInfoResponse{
		Coins:     userCoins,
		Inventory: *userInventory,
		CoinHistory: CoinHistory{
			ReceivedHistory: *userReceivedTransactions,
			SentHistory:     *userSentTransactions,
		},
	}
}

func AuthRequestToEntity(authRequest *sessionDTO.AuthRequest, coinsBalance uint) (*entity.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(authRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		Username:     authRequest.Username,
		Coins:        coinsBalance,
		PasswordHash: string(hashedPassword),
	}, nil
}

func ModelToEntity(user *model.User) *entity.User {
	return &entity.User{
		Username:     user.Username,
		Coins:        user.Coins,
		PasswordHash: user.PasswordHash,
	}
}

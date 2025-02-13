package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/dto"
	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/usecase"
	"github.com/artrsyf/avito-trainee-assignment/middleware"
)

type TransactionHandler struct {
	transactionUC usecase.TransactionUsecaseI
}

func NewTransactionHandler(transactionUsecase usecase.TransactionUsecaseI) *TransactionHandler {
	return &TransactionHandler{
		transactionUC: transactionUsecase,
	}
}

func (h *TransactionHandler) SendCoins(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	sendCoinsRequest := &dto.SendCoinsRequest{}
	err = json.Unmarshal(body, sendCoinsRequest)
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	/*TODO check senderUsername work*/
	senderUsername := r.Context().Value(middleware.UsernameContextKey).(string)
	transactionEntity := &entity.Transaction{
		SenderUsername:   senderUsername,
		ReceiverUsername: sendCoinsRequest.ReceiverUsername,
		Amount:           sendCoinsRequest.Amount,
	}

	err = h.transactionUC.Create(transactionEntity)
	if err != nil {
		/*Handle*/
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

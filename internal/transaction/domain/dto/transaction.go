package dto

type SendCoinsRequest struct {
	ReceiverUsername string `json:"toUser"`
	Amount           uint   `json:"amount"`
}

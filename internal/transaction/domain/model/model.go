package model

type Transaction struct {
	ID               uint   `db:"id"`
	SenderUsername   string `db:"sender_username"`
	ReceiverUsername string `db:"receiver_username"`
	Amount           uint   `db:"amount"`
}

package model

type Transaction struct {
	ID             uint `db:"id"`
	SenderUserID   uint `db:"sender_user_id"`
	ReceiverUserID uint `db:"receiver_user_id"`
	Amount         uint `db:"amount"`
}

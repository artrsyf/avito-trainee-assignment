package entity

type Transaction struct {
	SenderUsername   string
	ReceiverUsername string
	Amount           uint
}

type ReceivedTransactionGroup struct {
	SenderUsername string `json:"fromUser"`
	Amount         uint   `json:"amount"`
}

type SentTransactionGroup struct {
	ReceiverUsername string `json:"toUser"`
	Amount           uint   `json:"amount"`
}

type ReceivedHistory []ReceivedTransactionGroup

type SentHistory []SentTransactionGroup

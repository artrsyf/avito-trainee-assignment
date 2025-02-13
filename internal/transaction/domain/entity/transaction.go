package entity

type Transaction struct {
	SenderUsername   string
	ReceiverUsername string
	Amount           uint
}

type ReceivedTransactionGroup struct {
	SenderUsername string
	Amount         uint
}

type SentTransactionGroup struct {
	ReceiverUsername string
	Amount           uint
}

type ReceivedHistory []ReceivedTransactionGroup

type SentHistory []SentTransactionGroup

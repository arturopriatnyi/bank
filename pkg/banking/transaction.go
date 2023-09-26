package banking

import "github.com/google/uuid"

type Transaction struct {
	ID         string
	SenderID   int32
	ReceiverID int32
	Money      Money
}

func NewTransaction(senderID, receiverID int32, money Money) Transaction {
	return Transaction{
		ID:         uuid.NewString(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Money:      money,
	}
}

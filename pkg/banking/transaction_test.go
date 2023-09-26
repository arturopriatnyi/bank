package banking

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	money, err := NewMoney(5, 0, USD)
	assert.NoError(t, err)
	var senderID, receiverID int32 = 1, 2

	ts := NewTransaction(senderID, receiverID, money)

	_, err = uuid.Parse(string(ts.ID))
	assert.NoError(t, err)
	assert.Equal(t, senderID, ts.SenderID)
	assert.Equal(t, receiverID, ts.ReceiverID)
	assert.Equal(t, money, ts.Money)
}

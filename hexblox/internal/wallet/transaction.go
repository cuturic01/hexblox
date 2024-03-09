package wallet

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	id             string
	Input          *Input
	senderOutput   *Output
	receiverOutput *Output
}

func NewTransaction(senderWallet *Wallet, recipient string, amount float64) *Transaction {
	if amount > senderWallet.balance {
		fmt.Printf("Amount %f exceedes balance", amount)
		return nil
	}

	senderOutput := &Output{
		address: senderWallet.PublicKey,
		amount:  senderWallet.balance - amount,
	}
	receiverOutput := &Output{
		address: recipient,
		amount:  amount,
	}

	transaction := &Transaction{
		id:             uuid.NewString(),
		senderOutput:   senderOutput,
		receiverOutput: receiverOutput,
	}

	return SignTransaction(transaction, senderWallet)
}

func (transaction *Transaction) String() string {
	return fmt.Sprint(
		"-Transaction \n",
		"      Id:   ", transaction.id, "\n",
		"      Input:\n", transaction.Input.String(),
		"      Sender output:\n", transaction.senderOutput,
		"      Receiver output:\n", transaction.receiverOutput,
	)
}

func SignTransaction(transaction *Transaction, senderWallet *Wallet) *Transaction {
	transaction.Input = &Input{
		address:   senderWallet.PublicKey,
		timestamp: time.Now().UnixMilli(),
		amount:    senderWallet.balance,
		signature: GenerateHash(
			transaction.senderOutput.address,
			transaction.senderOutput.amount,
			transaction.receiverOutput.address,
			transaction.receiverOutput.amount),
	}
	return transaction
}

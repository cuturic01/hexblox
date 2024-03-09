package wallet

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	Id              string
	Input           *Input
	SenderOutput    *Output
	ReceiverOutputs []*Output
}

func NewTransaction(senderWallet *Wallet, recipient string, amount float64) *Transaction {
	if amount > senderWallet.balance {
		fmt.Printf("Amount %f exceedes balance", amount)
		return nil
	}

	senderOutput := &Output{
		Address: senderWallet.PublicKey,
		Amount:  senderWallet.balance - amount,
	}
	receiverOutput := &Output{
		Address: recipient,
		Amount:  amount,
	}

	transaction := &Transaction{
		Id:              uuid.NewString(),
		SenderOutput:    senderOutput,
		ReceiverOutputs: []*Output{receiverOutput},
	}

	SignTransaction(transaction, senderWallet)

	return transaction
}

func (transaction *Transaction) String() string {
	var outputsString string
	for _, output := range transaction.ReceiverOutputs {
		outputsString +=
			output.String() +
				"-----------------------------------------------------------------------------\n"
	}

	return fmt.Sprint(
		"-Transaction \n",
		"      Id:   ", transaction.Id, "\n",
		"      Input:\n", IndentString(transaction.Input.String(), "      "),
		"      Sender output:\n", IndentString(transaction.SenderOutput.String(), "      "),
		"      Receiver outputs:\n", IndentString(outputsString, "      "),
	)
}

func SignTransaction(transaction *Transaction, senderWallet *Wallet) {
	var outputsString string
	for _, output := range transaction.ReceiverOutputs {
		outputsString = fmt.Sprint(outputsString, output.String())
	}

	transaction.Input = &Input{
		Address:   senderWallet.PublicKey,
		Timestamp: time.Now().UnixMilli(),
		Amount:    senderWallet.balance,
		Signature: senderWallet.Sign(GenerateHash(outputsString)),
	}
}

func (transaction *Transaction) Update(senderWallet *Wallet, recipient string, amount float64) {
	if amount > transaction.SenderOutput.Amount {
		fmt.Printf("Amount %f exceedes balance", amount)
	}

	transaction.SenderOutput.Amount = transaction.SenderOutput.Amount - amount
	transaction.ReceiverOutputs = append(transaction.ReceiverOutputs, &Output{
		Address: recipient,
		Amount:  amount,
	})

	SignTransaction(transaction, senderWallet)
}

package wallet

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	Id              string
	Input           *Input
	senderOutput    *Output
	receiverOutputs []*Output
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
		Id:              uuid.NewString(),
		senderOutput:    senderOutput,
		receiverOutputs: []*Output{receiverOutput},
	}

	SignTransaction(transaction, senderWallet)

	return transaction
}

func (transaction *Transaction) String() string {
	var outputsString string
	for _, output := range transaction.receiverOutputs {
		outputsString +=
			output.String() +
				"-----------------------------------------------------------------------------\n"
	}

	return fmt.Sprint(
		"-Transaction \n",
		"      Id:   ", transaction.Id, "\n",
		"      Input:\n", IndentString(transaction.Input.String(), "      "),
		"      Sender output:\n", IndentString(transaction.senderOutput.String(), "      "),
		"      Receiver outputs:\n", IndentString(outputsString, "      "),
	)
}

func SignTransaction(transaction *Transaction, senderWallet *Wallet) {
	var outputsString string
	for _, output := range transaction.receiverOutputs {
		outputsString = fmt.Sprint(outputsString, output.String())
	}

	transaction.Input = &Input{
		Address:   senderWallet.PublicKey,
		timestamp: time.Now().UnixMilli(),
		amount:    senderWallet.balance,
		signature: senderWallet.Sign(GenerateHash(outputsString)),
	}
}

func (transaction *Transaction) Update(senderWallet *Wallet, recipient string, amount float64) {
	if amount > transaction.senderOutput.amount {
		fmt.Printf("Amount %f exceedes balance", amount)
	}

	transaction.senderOutput.amount = transaction.senderOutput.amount - amount
	transaction.receiverOutputs = append(transaction.receiverOutputs, &Output{
		address: recipient,
		amount:  amount,
	})

	SignTransaction(transaction, senderWallet)
}

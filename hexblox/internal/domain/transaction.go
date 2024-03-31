package domain

import (
	"fmt"
	"github.com/google/uuid"
	"hexblox/internal/config"
	"hexblox/internal/util"
	"time"
)

type Transaction struct {
	Id      string
	Input   *Input
	Outputs []*Output
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
		Id:      uuid.NewString(),
		Outputs: []*Output{senderOutput, receiverOutput},
	}

	SignTransaction(transaction, senderWallet)

	return transaction
}

func RewardTransaction(minerWallet *Wallet) *Transaction {
	input := &Input{
		Address:   "hexblox",
		Timestamp: time.Now().UnixMilli(),
		Amount:    100000,
		Signature: "reward-transaction",
	}
	output := &Output{
		Address: minerWallet.PublicKey,
		Amount:  config.MiningReword,
	}
	transaction := &Transaction{
		Id:      uuid.NewString(),
		Input:   input,
		Outputs: []*Output{output},
	}
	return transaction
}

func (transaction *Transaction) String() string {
	var outputsString string
	for _, output := range transaction.Outputs {
		outputsString +=
			output.String() +
				"-----------------------------------------------------------------------------\n"
	}

	return fmt.Sprint(
		"-Transaction \n",
		"      Id:   ", transaction.Id, "\n",
		"      Input:\n", util.IndentString(transaction.Input.String(), "      "),
		"      Outputs:\n", util.IndentString(outputsString, "      "),
	)
}

func (transaction *Transaction) Update(senderWallet *Wallet, recipient string, amount float64) {
	var senderOutput *Output
	for _, output := range transaction.Outputs {
		if output.Address == senderWallet.PublicKey {
			senderOutput = output
		}
	}
	if amount > senderOutput.Amount {
		fmt.Printf("Amount %f exceedes balance", amount)
	}

	senderOutput.Amount = senderOutput.Amount - amount
	transaction.Outputs = append(transaction.Outputs, &Output{
		Address: recipient,
		Amount:  amount,
	})

	SignTransaction(transaction, senderWallet)
}

func SignTransaction(transaction *Transaction, senderWallet *Wallet) {
	var outputsString string
	for _, output := range transaction.Outputs {
		outputsString = fmt.Sprint(outputsString, output.String())
	}
	transaction.Input = &Input{
		Address:   senderWallet.PublicKey,
		Timestamp: time.Now().UnixMilli(),
		Amount:    senderWallet.balance,
		Signature: senderWallet.Sign(util.GenerateHash(outputsString)),
	}
}

func Valid(transaction *Transaction) bool {
	var outputsString string
	for _, output := range transaction.Outputs {
		outputsString = fmt.Sprint(outputsString, output.String())
	}
	return util.VerifySignature(transaction.Input.Address, transaction.Input.Signature, util.GenerateHash(outputsString))
}

package domain

import (
	"fmt"
	"hexblox/internal/util"
)

type TransactionPool struct {
	Transactions []*Transaction
}

func NewTransactionPool() *TransactionPool {
	return &TransactionPool{
		Transactions: make([]*Transaction, 0),
	}
}

func (transactionPool *TransactionPool) AddTransaction(newTransaction *Transaction) {
	for _, transaction := range transactionPool.Transactions {
		if newTransaction.Id == transaction.Id {
			transaction = newTransaction
			return
		}
	}
	// fmt.Println(newTransaction)
	transactionPool.Transactions = append(transactionPool.Transactions, newTransaction)
}

func (transactionPool *TransactionPool) String() string {
	var transactionsString string
	for _, transaction := range transactionPool.Transactions {
		transactionsString +=
			transaction.String()
	}
	return fmt.Sprint(
		"-Transaction Pool \n",
		"      Transactions:\n", util.IndentString(transactionsString, "            "))
}

func (transactionPool *TransactionPool) ExistingTransaction(address string) *Transaction {
	for _, transaction := range transactionPool.Transactions {
		if transaction.Input.Address == address {
			return transaction
		}
	}

	return nil
}

func (transactionPool *TransactionPool) ValidTransactions() []*Transaction {
	validTransactions := make([]*Transaction, 0)
	for _, transaction := range transactionPool.Transactions {
		var totalAmount float64
		for _, output := range transaction.Outputs {
			totalAmount += output.Amount
		}
		if transaction.Input.Amount != totalAmount {
			fmt.Println("Input and output amounts don't match.")
			continue
		}

		if !Valid(transaction) {
			fmt.Println("Invalid transaction signature.")
			continue
		}

		validTransactions = append(validTransactions, transaction)

	}
	return validTransactions
}

func (transactionPool *TransactionPool) Clear() {
	transactionPool.Transactions = make([]*Transaction, 0)
	fmt.Println("Transaction pool cleared.")
}

package wallet

import (
	"fmt"
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
		"      Transactions:\n", IndentString(transactionsString, "            "))
}

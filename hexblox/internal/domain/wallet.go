package domain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"hexblox/internal/config"
	"hexblox/internal/util"
)

type Wallet struct {
	PublicKey string
	balance   float64
	keyPair   *ecdsa.PrivateKey
}

func NewWallet() *Wallet {
	keyPair := util.GenerateKeyPair()
	return &Wallet{
		PublicKey: util.EncodeKey(&keyPair.PublicKey),
		balance:   config.InitialBalance,
		keyPair:   keyPair,
	}
}

func (wallet *Wallet) String() string {
	return fmt.Sprint(
		"-Wallet \n",
		"      Public key: ", wallet.PublicKey, "\n",
		"      Balance:    ", wallet.balance, "\n",
	)
}

func (wallet *Wallet) Sign(hash string) string {
	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		panic(err)
	}
	signature, err := wallet.keyPair.Sign(rand.Reader, hashBytes, nil)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(signature)
}

func (wallet *Wallet) CreateTransaction(
	recipient string,
	amount float64,
	pool *TransactionPool,
	blockchain *Blockchain,
) *Transaction {
	wallet.CalculateBalance(blockchain)

	if amount > wallet.balance {
		fmt.Printf("Amount %f exceedes balance", amount)
		return nil
	}

	transaction := pool.ExistingTransaction(wallet.PublicKey)

	if transaction != nil {
		transaction.Update(wallet, recipient, amount)
		fmt.Println("Transaction updated.")
	} else {
		transaction = NewTransaction(wallet, recipient, amount)
		pool.AddTransaction(transaction)
	}

	return transaction
}

func (wallet *Wallet) CalculateBalance(blockchain *Blockchain) float64 {
	transactions := make([]*Transaction, 0)
	var mostRecentTimestamp int64 = 0
	var mostRecentTransaction *Transaction

	for _, block := range blockchain.Chain() {
		for _, transaction := range block.Data {
			transactions = append(transactions, transaction)
			if transaction.Input.Address != wallet.PublicKey {
				continue
			}

			if transaction.Input.Timestamp > mostRecentTimestamp {
				mostRecentTimestamp = transaction.Input.Timestamp
				mostRecentTransaction = transaction
			}
		}
	}

	if len(transactions) == 0 {
		fmt.Printf("No recent transactions.")
		return wallet.balance
	}

	if mostRecentTransaction == nil {
		fmt.Printf("No recent transactions from this wallet.")
		return wallet.balance
	}

	for _, output := range mostRecentTransaction.Outputs {
		if output.Address == wallet.PublicKey {
			wallet.balance = output.Amount
		}
	}

	for _, transaction := range transactions {
		if transaction.Input.Timestamp <= mostRecentTransaction.Input.Timestamp {
			continue
		}

		for _, output := range transaction.Outputs {
			if output.Address == wallet.PublicKey {
				wallet.balance += output.Amount
			}
		}
	}

	return wallet.balance
}

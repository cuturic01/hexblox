package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"hexblox/internal/config"
)

type Wallet struct {
	PublicKey string
	balance   float64
	keyPair   *ecdsa.PrivateKey
}

func NewWallet() *Wallet {
	keyPair := GenerateKeyPair()
	return &Wallet{
		PublicKey: EncodeKey(&keyPair.PublicKey),
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
	signature, err := wallet.keyPair.Sign(rand.Reader, []byte(hash), nil)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(signature)
}

func (wallet *Wallet) CreateTransaction(recipient string, amount float64, pool *TransactionPool) *Transaction {
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

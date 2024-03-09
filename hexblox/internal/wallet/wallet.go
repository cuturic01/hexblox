package wallet

import (
	"bytes"
	"crypto/rsa"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"hexblox/internal/config"
	"hexblox/internal/util"
)

type Wallet struct {
	PublicKey string
	balance   float64
	keyPair   *rsa.PrivateKey
}

func NewWallet() *Wallet {
	keyPair := util.GenerateKeyPair()
	publicKey := keyPair.PublicKey

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(publicKey.E); err != nil {
		fmt.Println(err)
	}
	// TODO: see if this is acceptable
	publicKeyHex := hex.EncodeToString(append(publicKey.N.Bytes(), buf.Bytes()...))

	return &Wallet{
		PublicKey: publicKeyHex,
		balance:   config.InitialBalance,
		keyPair:   keyPair,
	}
}

func (wallet *Wallet) String() string {
	return fmt.Sprint(
		"-Wallet \n",
		"      Public key: ", wallet.PublicKey[:10], "...\n",
		"      Balance:    ", wallet.balance, "\n",
	)
}

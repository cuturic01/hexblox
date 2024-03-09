package wallet

import "fmt"

type Input struct {
	Address   string
	timestamp int64
	amount    float64
	signature string
}

func (input *Input) String() string {
	return fmt.Sprint(
		"      Address:   ", input.Address, "\n",
		"      Timestamp: ", input.timestamp, "\n",
		"      Amount:    ", input.amount, "\n",
		"      Signature: ", input.signature, "\n",
	)
}

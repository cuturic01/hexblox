package wallet

import "fmt"

type Input struct {
	Address   string
	Timestamp int64
	Amount    float64
	Signature string
}

func (input *Input) String() string {
	if input == nil {
		return "  "
	}
	return fmt.Sprint(
		"      Address:   ", input.Address, "\n",
		"      Timestamp: ", input.Timestamp, "\n",
		"      Amount:    ", input.Amount, "\n",
		"      Signature: ", input.Signature, "\n",
	)
}

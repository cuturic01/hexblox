package wallet

import "fmt"

type Output struct {
	Address string
	Amount  float64
}

func (output *Output) String() string {
	return fmt.Sprint(
		"      Address:   ", output.Address, "\n",
		"      Amount:    ", output.Amount, "\n",
	)
}

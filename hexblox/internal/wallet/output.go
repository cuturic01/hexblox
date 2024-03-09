package wallet

import "fmt"

type Output struct {
	address string
	amount  float64
}

func (output *Output) String() string {
	return fmt.Sprint(
		"            Address:   ", output.address, "\n",
		"            Amount:    ", output.amount, "\n",
	)
}

package valueobject

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Money struct {
	Amount decimal.Decimal
}

func NewMoney(amount decimal.Decimal) (Money, error) {
	if amount.IsNegative() {
		return Money{}, fmt.Errorf("amount must be >= 0")
	}
	return Money{Amount: amount.Round(2)}, nil
}

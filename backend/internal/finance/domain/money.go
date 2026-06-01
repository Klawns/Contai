package domain

type Money int64

func NewMoney(cents int64) Money {
	return Money(cents)
}

func (money Money) Cents() int64 {
	return int64(money)
}

func (money Money) IsPositive() bool {
	return money > 0
}

func (money Money) IsZero() bool {
	return money == 0
}

func (money Money) Add(other Money) Money {
	return money + other
}

func (money Money) Sub(other Money) Money {
	return money - other
}

func (money Money) Neg() Money {
	return -money
}

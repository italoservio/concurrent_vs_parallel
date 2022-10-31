package structures

type Output struct {
	Entry  string
	Amount float32
}

func NewOutput(entry string, amount float32) *Output {
	return &Output{
		Entry:  entry,
		Amount: amount,
	}
}

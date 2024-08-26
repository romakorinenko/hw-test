package shape

type Square struct {
	side int
}

func NewSquare(side int) *Square {
	return &Square{side: side}
}

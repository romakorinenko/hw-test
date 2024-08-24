package shape

import (
	"fmt"
	"math"

	"github.com/shopspring/decimal"
)

type Circle struct {
	radius int64
}

func NewCircle(radius int64) *Circle {
	return &Circle{radius: radius}
}

func (c *Circle) Square() float64 {
	radius := decimal.NewFromInt(c.radius)

	return decimal.NewFromFloat(math.Pi).Mul(radius.Pow(decimal.NewFromFloat32(2))).InexactFloat64()
}

func (c *Circle) String(square float64) string {
	return fmt.Sprintf("Круг: радиус %d Площадь: %.4f", c.radius, square)
}

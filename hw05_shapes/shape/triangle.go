package shape

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Triangle struct {
	base   int64
	height int64
}

func NewTriangle(base int64, height int64) *Triangle {
	return &Triangle{base: base, height: height}
}

func (t *Triangle) Square() float64 {
	return decimal.NewFromFloat(0.5).Mul(decimal.NewFromInt(t.base * t.height)).InexactFloat64()
}

func (t *Triangle) Description(square float64) string {
	return fmt.Sprintf("Треугольник: основание %d, высота %d Площадь: %.0f", t.base, t.height, square)
}

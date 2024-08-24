package shape

import "fmt"

type Rectangle struct {
	length int
	width  int
}

func NewRectangle(length int, width int) *Rectangle {
	return &Rectangle{length: length, width: width}
}

func (r *Rectangle) Square() float64 {
	return float64(r.length * r.width)
}

// Description у прямоугольника все таки длина, а не высота.
func (r *Rectangle) Description(square float64) string {
	return fmt.Sprintf("Прямоугольник: ширина %d, высота %d Площадь: %.0f", r.width, r.length, square)
}

package shapes

import (
	"errors"

	"github.com/fixme_my_friend/hw06_testing/shapes/shape"
)

func calculateArea(s any) (float64, error) {
	switch shapeType := s.(type) {
	case shape.Shape:
		return shapeType.Square(), nil
	default:
		return 0, errors.New("ошибка: переданный объект не является фигурой")
	}
}

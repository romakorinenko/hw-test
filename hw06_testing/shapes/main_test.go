package main

import (
	"errors"
	"math"
	"testing"

	"github.com/fixme_my_friend/hw06_testing/shapes/shape"
	"github.com/stretchr/testify/require"
)

func Test_CalculateArea(t *testing.T) {
	testCases := []struct {
		name           string
		shape          any
		expectedSquare float64
		expectedError  error
	}{
		{
			name:           "Площадь для круга",
			shape:          shape.NewCircle(5),
			expectedSquare: math.Pi * math.Pow(5, 2),
			expectedError:  nil,
		},
		{
			name:           "Площадь для прямоугольника",
			shape:          shape.NewRectangle(5, 4),
			expectedSquare: 20,
			expectedError:  nil,
		},
		{
			name:           "Площадь для треугольника",
			shape:          shape.NewTriangle(5, 4),
			expectedSquare: 10,
			expectedError:  nil,
		},
		{
			name:           "Площадь для квадрата (не имплементирует интерфейс Shape)",
			shape:          shape.NewSquare(4),
			expectedSquare: 0,
			expectedError:  errors.New("ошибка: переданный объект не является фигурой"),
		},
		{
			name:           "shape = nil",
			shape:          nil,
			expectedSquare: 0,
			expectedError:  errors.New("ошибка: переданный объект не является фигурой"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			square, err := calculateArea(testCase.shape)

			require.Equal(t, testCase.expectedSquare, square)
			require.Equal(t, testCase.expectedError, err)
		})
	}
}

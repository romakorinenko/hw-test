package main

import (
	"fmt"

	"github.com/romakorinenko/hw-test/hw05_shapes/errors"
	"github.com/romakorinenko/hw-test/hw05_shapes/shape"
)

func main() {
	circle := shape.NewCircle(5)
	rectangle := shape.NewRectangle(5, 10)
	triangle := shape.NewTriangle(8, 6)
	square := shape.NewSquare(5)

	if circleSquare, err := calculateArea(circle); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(circle.String(circleSquare))
	}

	if rectangleSquare, err := calculateArea(rectangle); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(rectangle.String(rectangleSquare))
	}

	if triangleSquare, err := calculateArea(triangle); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(triangle.String(triangleSquare))
	}

	if squareSquare, err := calculateArea(square); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Unexpected result", squareSquare)
	}
}

func calculateArea(s any) (float64, error) {
	switch shapeType := s.(type) {
	case *shape.Circle:
		return shapeType.Square(), nil
	case *shape.Rectangle:
		return shapeType.Square(), nil
	case *shape.Triangle:
		return shapeType.Square(), nil
	default:
		return 0, errors.NoShapeError{}
	}
}

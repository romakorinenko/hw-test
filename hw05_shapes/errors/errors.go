package errors

type NoShapeError struct{}

func (e NoShapeError) Error() string {
	return "Ошибка: переданный объект не является фигурой."
}

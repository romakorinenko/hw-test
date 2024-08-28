package chessboard

import (
	"errors"
	"strings"
)

const (
	black = "#"
	white = " "
)

func Chessboard(size int) (string, error) {
	if size <= 0 {
		return "", errors.New("размер доски должен быть больше 1")
	}

	var res strings.Builder
	for str := 0; str < size; str++ {
		res.WriteString(chessString(size, str%2 == 0))
	}

	return res.String(), nil
}

func chessString(size int, isFirstBlack bool) string {
	var res strings.Builder

	for cell := 0; cell < size; cell++ {
		if cell%2 == 0 == isFirstBlack {
			res.WriteString(black)
		} else {
			res.WriteString(white)
		}
	}

	res.WriteString("\n")

	return res.String()
}

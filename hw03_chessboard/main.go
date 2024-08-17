package main

import (
	"fmt"
	"strings"
)

const (
	black = "#"
	white = " "
)

func main() {
	var size int
	fmt.Print("Enter chessboard size: ")
	if _, err := fmt.Scanf("%d", &size); err != nil {
		fmt.Printf("getting chessboard size error. Cause: %v", err)
	}

	fmt.Println(chessboard(size))
}

func chessboard(size int) string {
	var res strings.Builder
	for str := 0; str < size; str++ {
		res.WriteString(chessString(size, str%2 == 0))
	}

	return res.String()
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

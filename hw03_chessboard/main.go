package main

import "fmt"

const (
	black = "#"
	white = " "
)

func main() {
	var size int
	fmt.Print("Enter chessboard size: ")
	if _, err := fmt.Scanf("%d", &size); err != nil {
		fmt.Printf("getting chessboard size: %v", err)
	}

	chessboard(size)
}

func chessboard(size int) string {
	var res string
	for str := 0; str < size; str++ {
		res += chessString(size, str%2 == 0)
	}

	return res
}

func chessString(size int, isFirstBlack bool) string {
	var res string

	if isFirstBlack {
		for cell := 0; cell < size; cell++ {
			if cell%2 == 0 {
				res += black
			} else {
				res += white
			}
		}
	} else {
		for cell := 0; cell < size; cell++ {
			if cell%2 == 1 {
				res += black
			} else {
				res += white
			}
		}
	}

	return res + "\n"
}

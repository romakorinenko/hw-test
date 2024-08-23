package main

import (
	"fmt"

	"github.com/romakorinenko/hw-test/hw04_struct_comparator/comparator"
	"github.com/romakorinenko/hw-test/hw04_struct_comparator/model"
)

func main() {
	firstBook := model.NewBook(1, "book 1", "author 1", 1999, 546, 4.6)
	secondBook := model.NewBook(2, "book 2", "author 2", 2000, 300, 4.3)

	bookYearComparator := comparator.NewBookComparator(comparator.YearComparing)
	bookSizeComparator := comparator.NewBookComparator(comparator.SizeComparing)
	bookRateComparator := comparator.NewBookComparator(comparator.RateComparing)

	fmt.Println("is it first book older?", bookYearComparator.CompareByMode(firstBook, secondBook))
	fmt.Println("is it first book bigger?", bookSizeComparator.CompareByMode(firstBook, secondBook))
	fmt.Println("is it first book more popular?", bookRateComparator.CompareByMode(firstBook, secondBook))
}

package comparator

import (
	"fmt"

	"github.com/romakorinenko/hw-test/hw04_struct_comparator/model"
	"github.com/shopspring/decimal"
)

const (
	YearComparing = iota
	SizeComparing
	RateComparing
)

type BookComparator struct {
	comparatorMode int
}

func NewBookComparator(comparatorMode int) *BookComparator {
	return &BookComparator{comparatorMode: comparatorMode}
}

func (c *BookComparator) CompareByMode(firstBook, secondBook *model.Book) bool {
	switch c.comparatorMode {
	case YearComparing:
		return compareYears(firstBook.Year(), secondBook.Year())
	case SizeComparing:
		return compareSizes(firstBook.Size(), secondBook.Size())
	case RateComparing:
		return compareRates(firstBook.Rate(), secondBook.Rate())
	default:
		panic(fmt.Sprintf("comparator mode must have value from 0 to 2. received: %d", c.comparatorMode))
	}
}

func compareYears(firstBookYear, secondBookYear int) bool {
	return firstBookYear > secondBookYear
}

func compareSizes(firstBookSize, secondBookSize int) bool {
	return firstBookSize > secondBookSize
}

func compareRates(firstBookRate, secondBookRate float32) bool {
	return decimal.NewFromFloat32(firstBookRate).GreaterThan(decimal.NewFromFloat32(secondBookRate))
}

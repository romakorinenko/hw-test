package comparator

import (
	"errors"
	"fmt"

	"github.com/fixme_my_friend/hw06_testing/structComparator/model"
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

func (c *BookComparator) CompareByMode(firstBook, secondBook *model.Book) (bool, error) {
	if firstBook == nil || secondBook == nil {
		return false, errors.New("book cannot be null")
	}

	switch c.comparatorMode {
	case YearComparing:
		return compareYears(firstBook.Year(), secondBook.Year()), nil
	case SizeComparing:
		return compareSizes(firstBook.Size(), secondBook.Size()), nil
	case RateComparing:
		return compareRates(firstBook.Rate(), secondBook.Rate()), nil
	default:
		return false, fmt.Errorf("comparator mode must have value from 0 to 2. received: %d", c.comparatorMode)
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

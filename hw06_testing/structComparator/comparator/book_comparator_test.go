package comparator

import (
	"errors"
	"testing"

	"github.com/fixme_my_friend/hw06_testing/structComparator/model"
	"github.com/stretchr/testify/require"
)

func TestBookComparator_CompareByMode(t *testing.T) {
	testCases := []struct {
		name           string
		bookComparator *BookComparator
		firstBook      *model.Book
		secondBook     *model.Book
		expectedResult bool
		expectedError  error
	}{
		{
			name:           "книги одинаковые",
			bookComparator: NewBookComparator(YearComparing),
			firstBook:      model.NewBook(1, "book 1", "author 1", 1999, 546, 4.6),
			secondBook:     model.NewBook(1, "book 1", "author 1", 1999, 546, 4.6),
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:           "первая книга младше",
			bookComparator: NewBookComparator(YearComparing),
			firstBook:      model.NewBook(1, "book 1", "author 1", 2000, 546, 4.6),
			secondBook:     model.NewBook(1, "book 1", "author 1", 1999, 546, 4.6),
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:           "вторая книга младше",
			bookComparator: NewBookComparator(YearComparing),
			firstBook:      model.NewBook(1, "book 1", "author 1", 1999, 546, 4.6),
			secondBook:     model.NewBook(1, "book 1", "author 1", 2000, 546, 4.6),
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:           "первая книга толще",
			bookComparator: NewBookComparator(SizeComparing),
			firstBook:      model.NewBook(1, "book 1", "author 1", 1999, 600, 4.6),
			secondBook:     model.NewBook(1, "book 1", "author 1", 1999, 546, 4.6),
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:           "вторая книга толще",
			bookComparator: NewBookComparator(SizeComparing),
			firstBook:      model.NewBook(1, "book 1", "author 1", 1999, 546, 4.6),
			secondBook:     model.NewBook(1, "book 1", "author 1", 1999, 600, 4.6),
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:           "первая книга имеет более высокий рейтинг",
			bookComparator: NewBookComparator(RateComparing),
			firstBook:      model.NewBook(1, "book 1", "author 1", 1999, 546, 4.7),
			secondBook:     model.NewBook(1, "book 1", "author 1", 1999, 546, 4.6),
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:           "вторая книга имеет более высокий рейтинг",
			bookComparator: NewBookComparator(RateComparing),
			firstBook:      model.NewBook(1, "book 1", "author 1", 1999, 546, 4.7),
			secondBook:     model.NewBook(1, "book 1", "author 1", 1999, 546, 4.8),
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:           "вторая книга = nil",
			bookComparator: NewBookComparator(RateComparing),
			firstBook:      model.NewBook(1, "book 1", "author 1", 1999, 546, 4.7),
			secondBook:     nil,
			expectedResult: false,
			expectedError:  errors.New("book cannot be null"),
		},
		{
			name:           "указан неверный мод для сравнения книг",
			bookComparator: NewBookComparator(10),
			firstBook:      model.NewBook(1, "book 1", "author 1", 1999, 546, 4.7),
			secondBook:     model.NewBook(1, "book 1", "author 1", 1999, 546, 4.7),
			expectedResult: false,
			expectedError:  errors.New("comparator mode must have value from 0 to 2. received: 10"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mode, err := testCase.bookComparator.CompareByMode(testCase.firstBook, testCase.secondBook)

			require.Equal(t, testCase.expectedResult, mode)
			require.Equal(t, testCase.expectedError, err)
		})
	}
}

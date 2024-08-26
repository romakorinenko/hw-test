package chessboard

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChessboard(t *testing.T) {
	testCases := []struct {
		name              string
		size              int
		expectedBlackCell int
		expectedWhiteCell int
		expectedError     error
	}{
		{
			name:              "Размер доски 8",
			size:              8,
			expectedBlackCell: 32,
			expectedWhiteCell: 32,
			expectedError:     nil,
		},
		{
			name:              "Размер доски 7",
			size:              7,
			expectedBlackCell: 25,
			expectedWhiteCell: 24,
			expectedError:     nil,
		},
		{
			name:              "Размер доски 1",
			size:              1,
			expectedBlackCell: 1,
			expectedWhiteCell: 0,
			expectedError:     nil,
		},
		{
			name:              "Размер доски 0",
			size:              0,
			expectedBlackCell: 0,
			expectedWhiteCell: 0,
			expectedError:     errors.New("размер доски должен быть больше 1"),
		},
		{
			name:              "Размер доски -1",
			size:              0,
			expectedBlackCell: 0,
			expectedWhiteCell: 0,
			expectedError:     errors.New("размер доски должен быть больше 1"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			chessboard, err := Chessboard(testCase.size)
			require.Equal(t, testCase.expectedBlackCell, countCellsByColor(chessboard, rune(black[0])))
			require.Equal(t, testCase.expectedWhiteCell, countCellsByColor(chessboard, rune(white[0])))
			require.Equal(t, testCase.expectedError, err)
		})
	}
}

func countCellsByColor(chessboard string, color rune) int {
	count := 0

	for _, char := range chessboard {
		if char == color {
			count++
		}
	}

	return count
}

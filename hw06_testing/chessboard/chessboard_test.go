package chessboard

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChessboard(t *testing.T) {
	testCases := []struct {
		name               string
		size               int
		expectedChessboard string
		expectedError      error
	}{
		{
			name:               "Размер доски 8",
			size:               8,
			expectedChessboard: "# # # # \n # # # #\n# # # # \n # # # #\n# # # # \n # # # #\n# # # # \n # # # #\n",
			expectedError:      nil,
		},
		{
			name:               "Размер доски 7",
			size:               7,
			expectedChessboard: "# # # #\n # # # \n# # # #\n # # # \n# # # #\n # # # \n# # # #\n",
			expectedError:      nil,
		},
		{
			name:               "Размер доски 1",
			size:               1,
			expectedChessboard: "#\n",
			expectedError:      nil,
		},
		{
			name:               "Размер доски 0",
			size:               0,
			expectedChessboard: "",
			expectedError:      errors.New("размер доски должен быть больше 1"),
		},
		{
			name:               "Размер доски -1",
			size:               0,
			expectedChessboard: "",
			expectedError:      errors.New("размер доски должен быть больше 1"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			chessboard, err := Chessboard(testCase.size)
			require.Equal(t, testCase.expectedChessboard, chessboard)
			require.Equal(t, testCase.expectedError, err)
		})
	}
}

package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCountWithPool(t *testing.T) {
	testCases := []struct {
		name           string
		iterations     int
		expectedResult int
		expectedError  error
	}{
		{
			name:           "Количество итераций = 0",
			iterations:     0,
			expectedResult: 0,
			expectedError:  nil,
		},
		{
			name:           "Количество итераций = 1000",
			iterations:     1000,
			expectedResult: 1000,
			expectedError:  nil,
		},
		{
			name:           "Количество итераций = -1",
			iterations:     -1,
			expectedResult: 0,
			expectedError:  errors.New("iterations cannot be less than 0. iterations = -1"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := countWithPool(testCase.iterations)
			require.Equal(t, testCase.expectedError, err)
			require.Equal(t, testCase.expectedResult, result)
		})
	}
}

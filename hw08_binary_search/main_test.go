package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBinarySearchRecursive(t *testing.T) {
	testCases := []struct {
		name          string
		array         []int
		target        int
		expectedIndex int
	}{
		{
			name:          "получен пустой слайс",
			array:         []int{},
			target:        2,
			expectedIndex: -1,
		},
		{
			name:          "получен слайс элементов, значение target есть",
			array:         []int{1, 2, 3, 4, 5, 6},
			target:        4,
			expectedIndex: 3,
		},
		{
			name:          "получен слайс элементов, значение target нет",
			array:         []int{1, 2, 3, 4, 5, 6},
			target:        7,
			expectedIndex: -1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualIndex := binarySearchRecursive(testCase.array, 0, len(testCase.array)-1, testCase.target)

			require.Equal(t, testCase.expectedIndex, actualIndex)
		})
	}
}

func TestBinarySearchIterative(t *testing.T) {
	testCases := []struct {
		name          string
		array         []int
		target        int
		expectedIndex int
	}{
		{
			name:          "получен пустой слайс",
			array:         []int{},
			target:        2,
			expectedIndex: -1,
		},
		{
			name:          "получен слайс элементов, значение target есть",
			array:         []int{1, 2, 3, 4, 5, 6},
			target:        4,
			expectedIndex: 3,
		},
		{
			name:          "получен слайс элементов, значение target нет",
			array:         []int{1, 2, 3, 4, 5, 6},
			target:        7,
			expectedIndex: -1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualIndex := binarySearchIterative(testCase.array, testCase.target)

			require.Equal(t, testCase.expectedIndex, actualIndex)
		})
	}
}

package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCountWords(t *testing.T) {
	testCases := []struct {
		name           string
		text           string
		expectedResult map[string]int
	}{
		{
			name:           "В метод передана пустая строка",
			text:           "",
			expectedResult: map[string]int{},
		},
		{
			name:           "В метод передан текст на русском",
			text:           "рыжий кот рыжий кот",
			expectedResult: map[string]int{"рыжий": 2, "кот": 2},
		},
		{
			name:           "В метод передан текст на английском",
			text:           "ginger cat cat ginger",
			expectedResult: map[string]int{"ginger": 2, "cat": 2},
		},
		{
			name:           "В метод передан текст с разным регистром",
			text:           "gingER Cat GiNger cAt",
			expectedResult: map[string]int{"ginger": 2, "cat": 2},
		},
		{
			name:           "В метод передан текст, где помимо пробелов присутствуют другие знаки препинания",
			text:           "ginger,cat:ginger-cat",
			expectedResult: map[string]int{"ginger": 2, "cat": 2},
		},
		{
			name:           "В метод передан комбинированный текст, включающий все предыдущие кейсы",
			text:           "Рыжий;коТ-Ginger cat+кот",
			expectedResult: map[string]int{"ginger": 1, "cat": 1, "рыжий": 1, "кот": 2},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			wordToCountMap := countWords(testCase.text)
			require.Equal(t, testCase.expectedResult, wordToCountMap)
		})
	}
}

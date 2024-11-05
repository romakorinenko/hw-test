package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetParams_OnlyEnvFile(t *testing.T) {
	file, level, output, err := getParams()
	require.NoError(t, err)
	require.Equal(t, "log.log", file)
	require.Equal(t, "INFO", level)
	require.Equal(t, "result.txt", output)
}

func TestAnalyzeLogFile(t *testing.T) {
	testCases := []struct {
		name           string
		level          string
		fileName       string
		expectedResult string
		hasError       bool
	}{
		{
			name:           "Указан несуществующий log файл",
			level:          "INFO",
			fileName:       "log2.log",
			expectedResult: "",
			hasError:       true,
		},
		{
			name:           "Указан существующий log файл с уровнем INFO",
			level:          "INFO",
			fileName:       "log.log",
			expectedResult: "файл log.log содержит 12 строк с уровнем INFO",
			hasError:       false,
		},
		{
			name:           "Указан существующий log файл с уровнем WARN",
			level:          "WARN",
			fileName:       "log.log",
			expectedResult: "файл log.log содержит 3 строк с уровнем WARN",
			hasError:       false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := analyzeLogFile(testCase.fileName, testCase.level)
			require.Equal(t, testCase.hasError, err != nil)
			require.Equal(t, testCase.expectedResult, result)
		})
	}
}

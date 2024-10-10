package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCollectData(t *testing.T) {
	fromSensorChannel := make(chan int64)

	seconds := 1 * time.Second
	go collectData(fromSensorChannel, seconds)

	to := time.After(seconds)

	results := make([]int64, 0, 10)

	select {
	case <-to:
		break
	case i := <-fromSensorChannel:
		results = append(results, i)
	}

	require.True(t, len(results) > 0)
}

func TestAnalyseData_TenNumbersReceived(t *testing.T) {
	fromSensorChannel := make(chan int64)
	analysedDataChannel := make(chan float32)

	go analyseData(fromSensorChannel, analysedDataChannel)

	testData := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, number := range testData {
		fromSensorChannel <- number
	}
	close(fromSensorChannel)

	result := <-analysedDataChannel

	expected := float32(5.5)
	require.Equal(t, expected, result)
}

func TestAnalyseData_TwoNumbersReceived(t *testing.T) {
	fromSensorChannel := make(chan int64)
	analysedDataChannel := make(chan float32)

	go analyseData(fromSensorChannel, analysedDataChannel)

	testData := []int64{1, 2}
	for _, number := range testData {
		fromSensorChannel <- number
	}
	close(fromSensorChannel)

	select {
	case <-time.After(100 * time.Millisecond):
	default:
		_, ok := <-analysedDataChannel
		require.False(t, ok)
	}
}

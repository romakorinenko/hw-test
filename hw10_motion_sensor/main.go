package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"
)

func main() {
	fromSensorChannel := make(chan int64)
	analysedDataChannel := make(chan float32)

	go collectData(fromSensorChannel, 60*time.Minute)
	go analyseData(fromSensorChannel, analysedDataChannel)

	for averageData := range analysedDataChannel {
		fmt.Println(averageData)
	}
}

func collectData(fromSensorChannel chan<- int64, seconds time.Duration) {
	defer close(fromSensorChannel)

	to := time.After(seconds)

	for {
		select {
		case <-to:
			return
		default:
			randomNumber, err := rand.Int(rand.Reader, big.NewInt(10))
			if err != nil {
				log.Println("cannot receive random number")
			}
			fromSensorChannel <- randomNumber.Int64()
		}
	}
}

func analyseData(fromSensorChannel <-chan int64, analysedDataChannel chan<- float32) {
	defer close(analysedDataChannel)

	var data int64
	var count int

	for i := range fromSensorChannel {
		data += i
		count++
		if count == 10 {
			analysedDataChannel <- float32(data) / 10
			data = 0
			count = 0
		}
	}
}

package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"
)

func main() {
	fromSensorChannel := make(chan int64)
	analysedDataChannel := make(chan float32)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	go collectData(ctx, fromSensorChannel)
	go analyseData(fromSensorChannel, analysedDataChannel)

	for b := range analysedDataChannel {
		fmt.Println(b)
	}
}

func collectData(ctx context.Context, fromSensorChannel chan<- int64) {
	defer close(fromSensorChannel)

	for {
		select {
		case <-ctx.Done():
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

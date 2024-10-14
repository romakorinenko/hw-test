package main

import (
	"fmt"
	"sync"
)

func main() {
	iterations := 100
	res, err := countWithPool(iterations)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("result = %d after %d iterations", res, iterations)
}

func countWithPool(iterations int) (int, error) {
	if iterations < 0 {
		return 0, fmt.Errorf("iterations cannot be less than 0. iterations = %d", iterations)
	}

	mutex := sync.Mutex{}
	waitGroup := sync.WaitGroup{}

	var counter int

	for iteration := 0; iteration < iterations; iteration++ {
		waitGroup.Add(1)

		go func() {
			mutex.Lock()
			counter++
			fmt.Printf("goroutine increase counter value to %d\n", counter)
			mutex.Unlock()
			waitGroup.Done()
		}()
		waitGroup.Wait()
	}
	return counter, nil
}

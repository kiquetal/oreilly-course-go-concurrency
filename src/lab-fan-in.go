package main

import (
	"fmt"
	"sync"
)

func main() {
	var producerWg sync.WaitGroup

	fanInChan := make(chan int)
	for i := 0; i < 10; i++ {
		producerWg.Add(1)
		go func(workerID int) {
			defer func() {
				producerWg.Done()
				if workerID == 0 {
					producerWg.Wait()
					close(fanInChan)
				}
			}()
			for j := 0; j < 10; j++ {
				fanInChan <- j
			}
			fmt.Printf("producer %d finished\n", workerID)

		}(i)
	}

	var consumerWg sync.WaitGroup
	consumerWg.Add(1)
	go func() {
		defer consumerWg.Done()
		for msg := range fanInChan {
			fmt.Printf("received %d\n", msg)
		}
		fmt.Printf("consumer finished - end of all\n")
	}()
	producerWg.Wait()
	consumerWg.Wait()

}

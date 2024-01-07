package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	n := 100
	//limit to 10 concurrent calls
	ch := make(chan int)
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(workerId int) {
			for msgID := range ch {
				DoRPC(workerId, msgID) //blocking call
			}
			wg.Done()
		}(i)
	}

	for i := 0; i < n; i++ {
		ch <- i // sending blocks until one of the workers is ready to receive
	}
	close(ch)
	wg.Wait()

}

func DoRPC(workerId int, msgID int) {
	time.Sleep(150 * time.Millisecond)
	fmt.Printf("worker %d sending message %d\n", workerId, msgID)

}

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 350*time.Millisecond)
	defer cancel()
	var wg sync.WaitGroup
	n := 100
	//limit to 10 concurrent calls
	ch := make(chan int)
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(workerId int) {
			defer wg.Done()

			for {
				select {
				case msgID, ok := <-ch:
					if !ok {
						fmt.Printf("worker %d shutting down via channel\n", workerId)
						return
					}
					DoRPC(ctx, workerId, msgID)
				case <-ctx.Done():
					fmt.Printf("worker %d done via timeout\n", workerId)
					return // return to the block of select or for loop
				}
			}

		}(i)
	}
loop:
	for i := 0; i < n; i++ {
		select {
		case ch <- i:
		case <-ctx.Done():
			break loop
		}
	}
	close(ch)
	wg.Wait()
	fmt.Println("main done")
}

func DoRPC(ctx context.Context, workerId int, msgID int) {
	time.Sleep(150 * time.Millisecond)

	fmt.Printf("worker %d sending message %d\n", workerId, msgID)

}

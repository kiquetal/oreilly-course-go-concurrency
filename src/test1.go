package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	//use waitGroup
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go DoRPC(i)
	}
	wg.Wait()

}

func DoRPC(msgID int) {
	time.Sleep(1 * time.Second)
	fmt.Printf("sending message %d\n", msgID)

}

package main

import (
	"errors"
	"sync"
)

func main() {

	go func() {

	}()

}

type BufferedChannel struct {
	data []int
	head int // next empty slot
	tail int // next slot to read
	mut  *sync.Mutex
}

var ErrFull = errors.New("buffer is full")
var ErrEmpty = errors.New("buffer is empty")

func NewBufferedChannel(size int) *BufferedChannel {
	return &BufferedChannel{
		data: make([]int, size),
		mut:  new(sync.Mutex),
	}
}

func (bc *BufferedChannel) Send(msg int) error {
	bc.mut.Lock()
	defer bc.mut.Unlock()
	bc.data[bc.head] = msg
	bc.head = (bc.head + 1) % len(bc.data)
	return nil
}
func (bc *BufferedChannel) Receive() (int, error) {
	bc.mut.Lock()
	defer bc.mut.Unlock()
	msg := bc.data[bc.tail]
	bc.tail = (bc.tail + 1) % len(bc.data)
	return msg, nil

}

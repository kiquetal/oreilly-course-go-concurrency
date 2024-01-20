package main

import (
	"fmt"
	"strings"
	"sync"
)

func main() {
	lines := []string{
		"Loreum ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		"Do not go gentle into that good night, Old age should burn and rave at close of day; Rage, rage against the dying of the light.",
		"Rage, rage against the dying of the light.",
		"My name is Ozymandias, king of kings: Look on my works, ye Mighty, and despair!",
	}
	linesChan := make(chan string)
	wordsChan := make([]chan string, 26)
	for i := 0; i < 26; i++ {
		wordsChan[i] = make(chan string)
	}
	wordCountChan := make(chan map[string]int)
	numMappers := 3
	numReducers := 26

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, line := range lines {
			linesChan <- line
		}
		close(linesChan)
		fmt.Printf("producer finished\n")
	}()

	var wgMapper sync.WaitGroup
	for i := 0; i < numMappers; i++ {
		wgMapper.Add(1)
		go func(id int) {
			defer func() {
				wgMapper.Done()
				if id == 0 {
					wgMapper.Wait()
					for _, ch := range wordsChan {
						close(ch)
					}
				}
				fmt.Printf("mapper %d finished\n", id)
			}()
			for line := range linesChan {
				tokens := strings.Split(line, " ")
				for _, token := range tokens {
					token = strings.TrimSpace(token)
					println(token)
					idx := int(token[0] - 'a')
					wordsChan[idx] <- token
				}
			}
		}(i)
	}

	var wgReducer sync.WaitGroup
	for i := 0; i < numReducers; i++ {
		wgReducer.Add(1)
		go func(workerId int) {
			defer func() {
				wgReducer.Done()
				if workerId == 0 {
					wgReducer.Wait()
					close(wordCountChan)
				}
				fmt.Printf("reducer %d finished\n", workerId)
			}()
			count := make(map[string]int)
			for word := range wordsChan[workerId] {
				count[word]++
			}
			wordCountChan <- count
		}(i)
	}

	go func() {
		defer wg.Done()
		for counts := range wordCountChan {
			fmt.Printf("counts received: %v\n", counts)
		}
		fmt.Printf("consumer finished - end of all\n")
	}()
	wg.Wait()
	fmt.Println("main routine")

}

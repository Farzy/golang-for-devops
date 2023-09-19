package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	fmt.Printf("one\n")
	wg.Add(1)
	go func() {
		defer wg.Done()
		testFunction("Thread 1", 300*time.Millisecond)
	}()
	fmt.Printf("two\n")
	wg.Add(1)
	go func() {
		defer wg.Done()
		testFunction("Thread 2", 800*time.Millisecond)
	}()
	fmt.Printf("three\n")
	wg.Wait()
	fmt.Printf("We are finished!\n")
}

func testFunction(id string, pause time.Duration) {
	for i := 0; i < 5; i++ {
		fmt.Printf("#%s: checking…\n", id)
		time.Sleep(pause)
	}
}

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	var once sync.Once

	fmt.Printf("one\n")
	wg.Add(1)
	go func() {
		defer wg.Done()
		testFunction("Thread 1", 300*time.Millisecond, &once)
	}()
	fmt.Printf("two\n")
	wg.Add(1)
	go func() {
		defer wg.Done()
		testFunction("Thread 2", 800*time.Millisecond, &once)
	}()
	fmt.Printf("three\n")
	wg.Wait()
	fmt.Printf("We are finished!\n")
}

func testFunction(id string, pause time.Duration, once *sync.Once) {
	for i := 0; i < 5; i++ {
		once.Do(func() {
			fmt.Printf("#%s go the Once call\n", id)
		})

		fmt.Printf("#%s: checking…\n", id)
		time.Sleep(pause)
	}
}

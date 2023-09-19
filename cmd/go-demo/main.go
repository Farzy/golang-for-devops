package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Printf("one\n")
	c := make(chan string)
	c2 := make(chan string)
	go testFunction("Thread 1", 300*time.Millisecond, c)
	fmt.Printf("two\n")
	go testFunction("Thread 2", 800*time.Millisecond, c2)
	fmt.Printf("three\n")
	areWeFinished := []string{<-c, <-c2}
	fmt.Printf("areWeFinished: %v\n", areWeFinished)
}

func testFunction(id string, pause time.Duration, c chan string) {
	for i := 0; i < 5; i++ {
		fmt.Printf("#%s: checking…\n", id)
		time.Sleep(pause)
	}
	c <- fmt.Sprintf("#%s: we are finished", id)
}

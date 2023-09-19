package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Printf("one\n")
	c := make(chan string)
	go testFunction(c)
	fmt.Printf("two\n")
	areWeFinished := <-c
	fmt.Printf("areWeFinished: %v\n", areWeFinished)
}

func testFunction(c chan string) {
	for i := 0; i < 5; i++ {
		fmt.Printf("checking…\n")
		time.Sleep(1 * time.Second)
	}
	c <- "we are finished"
}

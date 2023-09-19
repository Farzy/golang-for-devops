package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Printf("one\n")
	go testFunction()
	fmt.Printf("two\n")
	time.Sleep(3 * time.Second)
}

func testFunction() {
	for {
		fmt.Printf("checking…\n")
		time.Sleep(1 * time.Second)
	}
}

package main

import (
	"fmt"
	"sync"
)

type mytype struct {
	counter int
	mu      sync.Mutex
}

func main() {
	threadCount := 100
	myTypeInstance := mytype{}
	finished := make(chan bool)
	for i := 0; i < threadCount; i++ {
		go func(myTypeInstance *mytype) {
			myTypeInstance.counter++
			finished <- true
		}(&myTypeInstance)
	}
	for i := 0; i < threadCount; i++ {
		<-finished
	}
	fmt.Printf("Counter: %d\n", myTypeInstance.counter)
}

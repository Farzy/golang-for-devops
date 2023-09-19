package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type mytype struct {
	counter int
	mu      sync.Mutex
}

func main() {
	threadCount := 10
	myTypeInstance := mytype{}
	finished := make(chan bool)
	for i := 0; i < threadCount; i++ {
		go func(myTypeInstance *mytype) {
			myTypeInstance.mu.Lock()
			fmt.Printf("input counter: %d\n", myTypeInstance.counter)
			myTypeInstance.counter++
			time.Sleep(time.Duration(rand.Intn(800)) * time.Millisecond)
			if myTypeInstance.counter == 5 {
				fmt.Printf("Found counter == 5\n")
			}
			fmt.Printf("output counter: %d\n", myTypeInstance.counter)
			finished <- true
			myTypeInstance.mu.Unlock()
		}(&myTypeInstance)
	}
	for i := 0; i < threadCount; i++ {
		<-finished
	}
	fmt.Printf("Counter: %d\n", myTypeInstance.counter)
}

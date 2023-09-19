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
			myTypeInstance.counter++
			time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
			if myTypeInstance.counter == 5 {
				fmt.Printf("Found counter == 5\n")
			}
			finished <- true
			myTypeInstance.mu.Unlock()
		}(&myTypeInstance)
	}
	for i := 0; i < threadCount; i++ {
		<-finished
	}
	fmt.Printf("Counter: %d\n", myTypeInstance.counter)
}

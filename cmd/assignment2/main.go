package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	for i := 0; i < 5; i++ {
		go callRatelimit()
	}
	time.Sleep(10 * time.Minute)
}

func callRatelimit() {
	for {
		response, err := http.Get("http://localhost:8080/ratelimit")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		body, err := io.ReadAll(response.Body)
		_ = response.Body.Close()
		if err != nil {
			fmt.Printf("Error readall: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s: %s", response.Status, body)
		time.Sleep(1 * time.Second)
	}
}

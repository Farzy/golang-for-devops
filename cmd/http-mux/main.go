package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var handlers = []struct {
	path string
	fn   func(w http.ResponseWriter, r *http.Request)
	cost int8
}{
	{
		"/v1/hello",
		HelloHandler,
		1,
	},
	{
		"/v1/time",
		CurrentTimeHandler,
		2,
	},
}

func HelloHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Hello, World!\n"))
}

func CurrentTimeHandler(w http.ResponseWriter, _ *http.Request) {
	curTime := time.Now().Format(time.RFC3339)
	_, _ = w.Write([]byte(fmt.Sprintf("The current time is %v\n", curTime)))
}

func main() {
	addr, exists := os.LookupEnv("ADDR")
	if !exists {
		addr = "localhost:8080"
	}

	mux := http.NewServeMux()
	for _, handler := range handlers {
		mux.HandleFunc(handler.path, handler.fn)
	}

	log.Printf("Server is listening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

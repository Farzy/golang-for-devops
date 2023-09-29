package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!\n"))
}

func CurrentTimeHandler(w http.ResponseWriter, r *http.Request) {
	curTime := time.Now().Format(time.RFC3339)
	w.Write([]byte(fmt.Sprintf("The current time is %v\n", curTime)))
}

func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = "localhost:8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/hello", HelloHandler)
	mux.HandleFunc("/v1/time", CurrentTimeHandler)

	log.Printf("Server is listening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

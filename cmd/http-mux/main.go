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
	fn   http.HandlerFunc
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

// Logger is a middleware handler that does request logging
type Logger struct {
	handler http.Handler
}

// ServeHTTP handles the request by passing it to the real
// handler and logging the request details
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	log.Printf("Logger: %s %s %v", r.Method, r.URL.Path, time.Since(start))
}

// NewLogger constructs a new Logger middleware handler
func NewLogger(next http.Handler) *Logger {
	return &Logger{next}
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
	wrappedMux := NewLogger(mux)

	log.Printf("Server is listening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, wrappedMux))
}

package main

import (
	"fmt"
	"io"
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
		"/v1/hello/",
		HelloHandler,
		1,
	},
	{
		"/v1/time",
		CurrentTimeHandler,
		2,
	},
	{
		"/v1/trailers",
		SendTrailersHandler,
		2,
	},
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(fmt.Sprintf(`Hello, World!
This is the path: %s
`, r.URL.Path)))
}

func CurrentTimeHandler(w http.ResponseWriter, _ *http.Request) {
	curTime := time.Now().Format(time.RFC3339)
	_, _ = w.Write([]byte(fmt.Sprintf("The current time is %v\n", curTime)))
}

// SendTrailersHandler add HTTP Trailers, they are a set of key/value pairs like headers that come after
// the HTTP response, instead of before.
func SendTrailersHandler(w http.ResponseWriter, _ *http.Request) {
	// Before any call to WriteHeader or Write, declare
	// the trailers you will set during the HTTP
	// response. These three headers are actually sent in
	// the trailer.
	w.Header().Set("Trailer", "AtEnd1, AtEnd2")
	w.Header().Add("Trailer", "AtEnd3")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)

	w.Header().Set("AtEnd1", "value 1")
	_, _ = io.WriteString(w, "This HTTP response has both headers before this text and trailers at the end.\n")
	w.Header().Set("AtEnd2", "value 2")
	w.Header().Set("AtEnd3", "value 3") // These will appear as trailers.
}

// Logger is a middleware handler that does request logging
type Logger struct {
	handler http.Handler
}

// ServeHTTP handles the request by passing it to the real
// handler and logging the request details
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Begin Logger")
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	log.Printf("End Logger: %s %s %v", r.Method, r.URL.Path, time.Since(start))
}

// NewLogger constructs a new Logger middleware handler
func NewLogger(next http.Handler) *Logger {
	return &Logger{next}
}

type ResponseHeader struct {
	handler     http.Handler
	headerName  string
	headerValue string
}

func (rh ResponseHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Begin ResponseHeader")
	w.Header().Add(rh.headerName, rh.headerValue)
	rh.handler.ServeHTTP(w, r)
	log.Printf("End ResponseHeader")
}

func NewResponseHeader(next http.Handler, headerName string, headerValue string) http.Handler {
	return &ResponseHeader{
		next,
		headerName,
		headerValue,
	}
}

func Middleware1(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Begin Middleware1")
			log.Printf("  Request: %+v", r)
			log.Printf("  Response: %+v", w)
			next.ServeHTTP(w, r)
			log.Printf("End Middleware1")
			log.Printf("  Response: %+v", w)
		},
	)
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
	wrappedMux :=
		Middleware1(
			NewLogger(
				NewResponseHeader(mux, "X-My-Header", "My header value")))

	log.Printf("Server is listening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, wrappedMux))
}

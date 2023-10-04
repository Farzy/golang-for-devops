package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Farzy/golang-for-devops/pkg/mux/helpers"
	"github.com/Farzy/golang-for-devops/pkg/mux/token-bucket"
)

type HTTPHandlerConfig struct {
	pathTrailer string
	fn          http.HandlerFunc
	maxTokens   int64
	rate        int64
}

var handlers = map[string]*HTTPHandlerConfig{
	"/v1/hello": {
		"/",
		HelloHandler,
		50,
		2,
	},
	"/v1/time": {
		"",
		CurrentTimeHandler,
		25,
		1,
	},
	"/v1/trailers": {
		"",
		SendTrailersHandler,
		10,
		1,
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
	w.Header().Add(rh.headerName, rh.headerValue)
	rh.handler.ServeHTTP(w, r)
}

func NewResponseHeader(next http.Handler, headerName string, headerValue string) http.Handler {
	return &ResponseHeader{
		next,
		headerName,
		headerValue,
	}
}

type CostHeader struct {
	handler       http.Handler
	pathBucketMap map[string]*token_bucket.TokenBucket
	m             sync.Mutex
}

func (ch *CostHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var pathBucket *token_bucket.TokenBucket

	// Keep only the first 2 components of the path, without the trailing '/'
	path := helpers.TruncateFromNthOccurrence(r.URL.Path, '/', 3)

	handlerConfig, found := handlers[path]
	if !found {
		handlerConfig = &HTTPHandlerConfig{
			pathTrailer: "",
			fn:          nil,
			maxTokens:   1,
			rate:        1,
		}
	}
	w.Header().Set(
		"X-Request-Cost",
		fmt.Sprintf("max %d / rate %d", handlerConfig.maxTokens, handlerConfig.rate))

	if pathBucket, found = ch.pathBucketMap[path]; !found {
		ch.m.Lock()
		pathBucket = token_bucket.NewTokenBucket(handlerConfig.rate, handlerConfig.maxTokens)
		ch.pathBucketMap[path] = pathBucket
		ch.m.Unlock()
	}
	if pathBucket.IsRequestAllowed(1) {
		ch.handler.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte("Too many requests!\n"))
		return
	}
}

func NewCostHeader(next http.Handler) http.Handler {
	ch := CostHeader{
		handler:       next,
		pathBucketMap: make(map[string]*token_bucket.TokenBucket),
	}

	return &ch
}

func Middleware1(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Middleware1 Request: %+v", r)
			next.ServeHTTP(w, r)
			log.Printf("Middleware1 Response: %+v", w)
		},
	)
}

func main() {
	addr, exists := os.LookupEnv("ADDR")
	if !exists {
		addr = "localhost:8080"
	}

	mux := http.NewServeMux()
	for path, handlerConfig := range handlers {
		mux.HandleFunc(path+handlerConfig.pathTrailer, handlerConfig.fn)
	}
	wrappedMux :=
		Middleware1(
			NewLogger(
				NewResponseHeader(
					NewCostHeader(mux), "X-My-Header", "My header value")))

	log.Printf("Server is listening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, wrappedMux))
}

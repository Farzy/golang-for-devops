package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func formatHeaders(header *http.Header) string {
	s := strings.Builder{}
	for k, v := range *header {
		_, err := fmt.Fprintf(&s, "- %s: %s\n", k, v)
		if err != nil {
			return "Error building string"
		}
	}
	return s.String()
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "It works.\n")
	fmt.Fprintf(w, "Method: %v\n", r.Method)
	fmt.Fprintf(w, "Headers:\n%s\n", formatHeaders(&r.Header))
}
func main() {
	port := "8080"
	http.HandleFunc("/", index)
	fmt.Printf("Starting server on port %s…\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}

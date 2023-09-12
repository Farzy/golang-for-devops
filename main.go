package main

import (
	"fmt"
	"io"
	"log"
)

// Words structure
// Example json: {"page":"words","input":"word1","words":["word3","word2","word1"]}
type Words struct {
	Page  string   `json:"page"`
	Input string   `json:"input"`
	Words []string `json:"words"`
}

type MySlowReader struct {
	Contents string
	pos      int
}

func (m *MySlowReader) Read(p []byte) (n int, err error) {
	if m.pos < len(m.Contents) {
		n := copy(p, m.Contents[m.pos:m.pos+1])
		m.pos++
		return n, nil
	}
	return 0, io.EOF
}

func main() {
	mySlowReaderInstance := &MySlowReader{
		Contents: "hello world!",
	}

	out, err := io.ReadAll(mySlowReaderInstance)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("output: %s\n", out)
}

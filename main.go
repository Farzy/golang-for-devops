package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"unicode"
)

// Word Create custom type, in order to have the Unmarshaler code capitalize the first letter
// when converting from JSON to our struct.
type Word string

func (w *Word) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if len(s) > 0 {
		runes := []rune(s)
		runes[0] = unicode.ToUpper(runes[0])
		s = string(runes)
	}
	*w = Word(s)

	return nil
}

// Words structure
// Example json: {"page":"words","input":"word1","words":["word3","word2","word1"]}
type Words struct {
	Page  string `json:"page"`
	Input string `json:"input"`
	Words []Word `json:"words"`
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Usage: ./http-get <url>\n")
		os.Exit(1)
	}
	myUrl := args[1]
	if _, err := url.ParseRequestURI(myUrl); err != nil {
		fmt.Printf("URL is in invalid format: %s\n", err)
		os.Exit(1)
	}

	response, err := http.Get(myUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error while closing Body stream:", err)
		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("HTTP Status code: %d\nBody: %s\n", response.StatusCode, body)
	// fmt.Printf("HTTP Status code: %d\nBody: %v\n", response.StatusCode, string(body)) // Explicit []byte to string conversion

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Invalid output (HTTP code: %d): %s\n", response.StatusCode, body)
		os.Exit(1)
	}

	var words Words

	err = json.Unmarshal(body, &words)
	if err != nil {
		log.Fatal(err)
	}

	var strWords = make([]string, 0, len(words.Words))
	for _, w := range words.Words {
		strWords = append(strWords, string(w))
	}
	fmt.Printf("JSON Parsed\nPage: %s\nWords: %v\n", words.Page, strings.Join(strWords, ", "))
}

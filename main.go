package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Response interface {
	GetResponse() string
}

type Page struct {
	Name string `json:"page"`
}

// Words structure
// Example json: {"page":"words","input":"word1","words":["word3","word2","word1"]}
type Words struct {
	Input string   `json:"input"`
	Words []string `json:"words"`
}

func (w Words) GetResponse() string {
	return fmt.Sprintf("%s", strings.Join(w.Words, ", "))
}

type Occurrence struct {
	Words map[string]int `json:"words"`
}

func (o Occurrence) GetResponse() string {
	var out []string
	for word, occurrence := range o.Words {
		out = append(out, fmt.Sprintf("%s (%d)", word, occurrence))
	}
	return fmt.Sprintf("%s", strings.Join(out, ", "))
}

func main() {
	var (
		requestURL string
		password   string
		parsedURL  *url.URL
		err        error
	)

	flag.StringVar(&requestURL, "url", "", "URL to access")
	flag.StringVar(&password, "password", "", "Use a password to access our API")
	flag.Parse()

	if parsedURL, err = url.ParseRequestURI(requestURL); err != nil {
		fmt.Printf("Validation error: URL '%s' is in invalid format: %s\n", requestURL, err)
		flag.Usage()
		os.Exit(1)
	}

	client := http.Client{}

	if password != "" {
		token, err := doLoginRequest(client, parsedURL.Scheme+"://"+parsedURL.Host+"/login", password)
		if err != nil {
			var requestErr RequestError
			if errors.As(err, &requestErr) {
				fmt.Printf("Error: %s (HTTP Code: %d, Body: %s)\n",
					requestErr.Err, requestErr.HTTPCode, requestErr.Body)
				os.Exit(1)

			}
		}
		client.Transport = MyJWTTransport{
			transport: http.DefaultTransport,
			token:     token,
		}
	}
	res, err := doRequest(client, parsedURL.String())
	if err != nil {
		var requestErr RequestError
		//if requestErr, ok := err.(RequestError); ok {
		if errors.As(err, &requestErr) {
			fmt.Printf("Error: %s (HTTP Code: %d, Body: %s)\n",
				requestErr.Err, requestErr.HTTPCode, requestErr.Body)
			os.Exit(1)

		}
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	if res == nil {
		fmt.Printf("No response\n")
		os.Exit(1)
	}

	fmt.Printf("Response: %s\n", res.GetResponse())
}

func doRequest(client http.Client, requestURL string) (Response, error) {
	response, err := client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP Get error: %s", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error while closing Body stream:", err)
		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("ReadAll error: %s", err)
	}

	//fmt.Printf("HTTP Status code: %d\nBody: %s\n", response.StatusCode, body)
	// fmt.Printf("HTTP Status code: %d\nBody: %v\n", response.StatusCode, string(body)) // Explicit []byte to string conversion

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid output (HTTP code: %d): %s\n", response.StatusCode, body)
	}

	if !json.Valid(body) {
		return nil, RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("No valid JSON returned"),
		}
	}

	var page Page

	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("Page unmasharl error: %s", err),
		}
	}

	switch page.Name {
	case "words":
		var words Words

		err = json.Unmarshal(body, &words)
		if err != nil {
			return nil, RequestError{
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Words unmasharl error: %s", err),
			}
		}

		return words, nil
	case "occurrence":
		var occurrence Occurrence

		err = json.Unmarshal(body, &occurrence)
		if err != nil {
			return nil, RequestError{
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Occurrence unmasharl error: %s", err),
			}
		}

		return occurrence, nil
	}

	return nil, nil
}

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

type WordsPage struct {
	Page
	Words
}

type ExtraSpecial struct {
	one   int
	two   int
	three string
}

// UnmarshalJSON Special treatment needed because of mixed types
// Reference: https://boldlygo.tech/posts/2019-12-19-go-json-tricks-array-as-struct/
func (e *ExtraSpecial) UnmarshalJSON(p []byte) error {
	var tmp []json.RawMessage
	if err := json.Unmarshal(p, &tmp); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[0], &e.one); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[1], &e.two); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[2], &e.three); err != nil {
		return err
	}
	return nil
}

type Assignment1 struct {
	Words        []string           `json:"words"`
	Percentages  map[string]float64 `json:"percentages"`
	Special      []string           `json:"special"`
	ExtraSpecial ExtraSpecial       `json:"extraSpecial"`
}

func (a Assignment1) GetResponse() string {
	var resp = make([]byte, 0, 50)
	resp = fmt.Appendf(resp, "Assignment1:\n")
	resp = fmt.Appendf(resp, "Words:\n")
	for _, w := range a.Words {
		resp = fmt.Appendf(resp, "- %q\n", w)
	}
	resp = fmt.Appendf(resp, "Percentages:\n")
	for _, p := range a.Percentages {
		resp = fmt.Appendf(resp, "- %v\n", p)
	}
	resp = fmt.Appendf(resp, "Special:\n")
	for _, w := range a.Special {
		resp = fmt.Appendf(resp, "- %q\n", w)
	}
	resp = fmt.Appendf(resp, "ExtraSpecial:\n")
	resp = fmt.Appendf(resp, "1: %d\n2: %d\n3: %q\n",
		a.ExtraSpecial.one, a.ExtraSpecial.two, a.ExtraSpecial.three)
	return string(resp)
}

func (a API) DoGetRequest(requestURL string) (Response, error) {
	response, err := a.Client.Get(requestURL)
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
	case "assignment1":
		var assignment1 Assignment1

		err = json.Unmarshal(body, &assignment1)
		if err != nil {
			return nil, RequestError{
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Assignment1 unmasharl error: %s", err),
			}
		}

		return assignment1, nil
	}

	return nil, nil
}

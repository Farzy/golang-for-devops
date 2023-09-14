package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type MockRoundTripper struct {
	RoundTripperOutput *http.Response
	IsAuthenticated    bool
}

const TOKEN = "abc123"

func (m MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.IsAuthenticated {
		if auth := req.Header.Get("Authorization"); auth != "Bearer "+TOKEN {
			return nil, fmt.Errorf("wrong authorization header: %s", auth)
		}
	}
	return m.RoundTripperOutput, nil
}

func TestMyJWTTransport_RoundTrip(t *testing.T) {
	myJWTTransport := MyJWTTransport{
		transport: MockRoundTripper{
			RoundTripperOutput: &http.Response{
				StatusCode: 200,
			},
			IsAuthenticated: false,
		},
	}
	req := &http.Request{}
	res, err := myJWTTransport.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip error: %s", err)
	}
	if res.StatusCode != 200 {
		t.Errorf("StatusCode is not 200, got: %d", res.StatusCode)
	}
}

func TestMyJWTTransport_RoundTripPassword(t *testing.T) {
	loginResponse := LoginResponse{
		Token: TOKEN,
	}
	loginResponseBytes, err := json.Marshal(loginResponse)
	if err != nil {
		t.Errorf("marshal error: %s", err)
	}
	myJWTTransport := MyJWTTransport{
		transport: MockRoundTripper{
			RoundTripperOutput: &http.Response{
				StatusCode: 200,
			},
			IsAuthenticated: true,
		},
		httpClient: MockClient{
			PostResponseOutput: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(loginResponseBytes)),
			},
		},
		password: "abc",
	}
	req := &http.Request{
		Header: make(http.Header),
	}
	res, err := myJWTTransport.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip error: %s", err)
	}
	if res.StatusCode != 200 {
		t.Errorf("StatusCode is not 200, got: %d", res.StatusCode)
	}
}

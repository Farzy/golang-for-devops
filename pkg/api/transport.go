package api

import (
	"net/http"
)

type MyJWTTransport struct {
	transport  http.RoundTripper
	token      string
	password   string
	loginURL   string
	httpClient ClientIface
}

func (m *MyJWTTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if m.token == "" {
		if m.password != "" {
			token, err := doLoginRequest(m.httpClient, m.loginURL, m.password)
			if err != nil {
				return nil, err
			}
			m.token = token
		}
	}
	if m.token != "" {
		request.Header.Add("Authorization", "Bearer "+m.token)
	}
	return m.transport.RoundTrip(request)
}

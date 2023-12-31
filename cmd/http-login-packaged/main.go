package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/Farzy/golang-for-devops/pkg/api"
)

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

	apiInstance := api.New(api.Options{
		Password: password,
		LoginURL: parsedURL.Scheme + "://" + parsedURL.Host + "/login",
	})

	res, err := apiInstance.DoGetRequest(parsedURL.String())
	if err != nil {
		var requestErr api.RequestError
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

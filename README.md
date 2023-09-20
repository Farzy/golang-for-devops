# Golang training

This training is based on the Udemy course [Golang for DevOps and Cloud Engineers](https://www.udemy.com/share/107N563@yDwZ8kiQ8Q0_E4TwrSv9vCsCJ-UA3XCcSWWcn-x_6x6EoFHIkzHflhaT0KitsTaNvw==/).

The course material is available at on [Udemy](https://github.com/wardviaene/golang-for-devops-course).

## Sample commands

Building binaries for different OS and architectures.

```shell
GOOS=linux GOARCH=amd64 go build -o assignement2-linux-amd64 cmd/assignment2/*.go
GOOS=linux GOARCH=arm64 go build -o assignement2-linux-arm64 cmd/assignment2/*.go
GOOS=darwin GOARCH=amd64 go build -o assignement2-darwin-amd64 cmd/assignment2/*.go
GOOS=darwin GOARCH=arm64 go build -o assignement2-darwin-arm64 cmd/assignment2/*.go
```

package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	fmt.Printf("Hello World!\nos.Args: %v\nArgument: %v\n", args, args[1:])
}

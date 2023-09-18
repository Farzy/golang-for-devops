package main

import (
	"fmt"
	"reflect"
)

func main() {
	var t1 string = "this is a string"
	discover(t1)
	discover2(t1)
	var t2 *string = &t1
	discover2(t2)
	var t3 int = 123
	discover2(t3)
	discover2(nil)
}

func discover(t any) {
	switch t.(type) {
	case string:
		// This will not work
		//t2 := t + "..."
		t2 := t.(string) + "..."
		fmt.Printf("String found: %s\n", t)
		fmt.Printf("New string: %s\n", t2)
	}
}

func discover2(t any) {
	switch v := t.(type) {
	case string:
		v += "..."
		fmt.Printf("String found: %s\n", t)
		fmt.Printf("New string: %s\n", v)
	case *string:
		fmt.Printf("Pointer string found: %v\n", v)
		fmt.Printf("Pointer string found: %s\n", *v)
	case int:
		fmt.Printf("We have an integer: %d\n", v)
	default:
		myType := reflect.TypeOf(t)
		if myType == nil {
			fmt.Printf("Type is nil\n")
		} else {
			fmt.Printf("Type not found: %s\n", reflect.TypeOf(t))
		}
	}
}

package main

import (
	"fmt"
	"os"
	"strconv"
)

// Print error message
func Stderr(msg ...interface{}) { //←[ Stderr msg does not escape]
	fmt.Fprintln(os.Stderr, msg...)
}

// Atoi that panics on bad input
func Atoi(a string) int { //←[ leaking param: a]
	i, err := strconv.Atoi(a)
	if err != nil {
		panic(err)
	}
	return i
}

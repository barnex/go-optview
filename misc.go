package main

import (
	"fmt"
	"os"
	"strconv"
)

// Print error message
func Stderr(msg ...interface{}) {
	fmt.Fprintln(os.Stderr, msg...)
}

// Atoi that panics on bad input
func Atoi(a string) int {
	i, err := strconv.Atoi(a)
	if err != nil {
		panic(err)
	}
	return i
}

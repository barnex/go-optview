package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	readFrom(os.Stdin)
}

func readFrom(in_ io.Reader) {
	in := bufio.NewReader(in_)
	for l, _, err := in.ReadLine(); err == nil; l, _, err = in.ReadLine() {
		line := string(l)
		parseLine(line)
	}
}

func parseLine(line string) {
	defer func() {
		err := recover()
		if err != nil {
			stderr(`optview:`, err, ` while parsing "`, line, `"`)
		}
	}()
	words := strings.SplitN(line, ":", 3)
	fmt.Println(words)
}

func stderr(msg ...interface{}) {
	fmt.Fprintln(os.Stderr, msg...)
}

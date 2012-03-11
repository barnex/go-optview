package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	files map[string]*SourceFile = make(map[string]*SourceFile)
)

func main() {
	flag.Parse()
	readFrom(os.Stdin)

	for f := range files {
		fmt.Println(`//`, f, ":")
	}
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
	filename := words[0]
	lineNo := Atoi(words[1])
	info := words[2]
	Source(filename).AddOpt(lineNo, info)
}

func stderr(msg ...interface{}) {
	fmt.Fprintln(os.Stderr, msg...)
}

func Atoi(a string) int {
	i, err := strconv.Atoi(a)
	if err != nil {
		panic(err)
	}
	return i
}

func Source(fileName string) *SourceFile {
	if file, ok := files[fileName]; ok {
		return file
	}
	file := NewSourceFile(fileName)
	files[fileName] = file
	return file
}

func NewSourceFile(fileName string) *SourceFile {
	return &SourceFile{fileName, make(map[int][]string)}
}

type SourceFile struct {
	Name string
	Opts map[int][]string
}

func (this *SourceFile) AddOpt(line int, opt string) {
	this.Opts[line] = append(this.Opts[line], opt)
}

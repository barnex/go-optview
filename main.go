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
	ReadCompilerOutput(os.Stdin)

	for n, f := range files {
		fmt.Fprintln(os.Stdout, "//", n, ":")
		f.FormatTo(os.Stdout)
	}
}

func ReadCompilerOutput(in_ io.Reader) {
	in := bufio.NewReader(in_)
	for l, prefix, err := in.ReadLine(); err == nil; l, prefix, err = in.ReadLine() {
		if prefix {
			panic("Enlarge your buffer!")
		}
		line := string(l)
		parseLine(line)
	}
}

func parseLine(line string) {
	defer func() {
		err := recover()
		if err != nil {
			Stderr(`optview:`, err, ` while parsing "`, line, `"`)
		}
	}()
	words := strings.SplitN(line, ":", 3)
	filename := words[0]
	lineNo := Atoi(words[1])
	info := words[2]
	GetSourceFile(filename).AddOpt(lineNo, info)
}

func Stderr(msg ...interface{}) {
	fmt.Fprintln(os.Stderr, msg...)
}

func Atoi(a string) int {
	i, err := strconv.Atoi(a)
	if err != nil {
		panic(err)
	}
	return i
}

func (f *SourceFile) FormatTo(out io.Writer) {
	in_, err := os.Open(f.Name)
	if err != nil {
		Stderr(err)
		return
	}
	in := bufio.NewReader(in_)

	lineNo := 1
	for l, prefix, err := in.ReadLine(); err == nil; l, prefix, err = in.ReadLine() {
		if prefix {
			panic("Enlarge your buffer!")
		}
		fmt.Fprint(out, string(l))
		if opts, ok := f.Opts[lineNo]; ok {
			fmt.Print("//", opts)
		}
		fmt.Fprintln(out)
		lineNo++
	}
}

func GetSourceFile(fileName string) *SourceFile {
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

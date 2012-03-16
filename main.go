package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

var (
	flag_version   *bool   = flag.Bool("V", false, "print version and exit")
	flag_writeback *bool   = flag.Bool("w", false, "write result to source files instead of stdout")
	flag_clean     *bool   = flag.Bool("c", false, "clean existing optimization messages from source")
	flag_prefix    *string = flag.String("prefix", "//"+"←", "prefix for optimization messages")
)

var files map[string]*SourceFile = make(map[string]*SourceFile)

func main() {
	flag.Parse()

	if *flag_version {
		fmt.Println("go-optview 0\nGo", runtime.Version())
		return
	}
	ReadCompilerOutput(os.Stdin)

	for name, f := range files {
		if *flag_writeback {
			buf := new(bytes.Buffer)
			f.WriteTo(buf)
			out, _ := os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
			out.Write(buf.Bytes())
		} else {
			fmt.Fprintln(os.Stdout, *flag_prefix, name, ":")
			f.WriteTo(os.Stdout)
		}
	}
}

// Read output from "gc -m". E.g.:
//	main.go:91: can inline NewSourceFile
//	main.go:80: inlining call to NewSourceFile
//	main.go:21: main ... argument does not escape
//	main.go:26: leaking param: in_
//	main.go:37: parseCompilerLine line does not escape
// Optimization comments are stored in "files" map
func ReadCompilerOutput(in_ io.Reader) {
	in := bufio.NewReader(in_)
	for l, prefix, err := in.ReadLine(); err == nil; l, prefix, err = in.ReadLine() {
		if prefix {
			panic("Enlarge your buffer!")
		}
		line := string(l)
		parseCompilerLine(line)
	}
}

// Parse single line of "gc -m" output
func parseCompilerLine(line string) {
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
	GetSourceFile(filename).AddMsg(lineNo, info)
}

// Print source code + optimization messages to out
func (f *SourceFile) WriteTo(out io.Writer) {
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

		// print source line
		Check(fmt.Fprint(out, cleanSourceLine(string(l))))

		// print messages
		if !*flag_clean {
			if opts, ok := f.Msg[lineNo]; ok {
				Check(fmt.Fprint(out, *flag_prefix, opts))
			}
		}
		fmt.Fprintln(out)
		lineNo++
	}
}

func Check(n int, err error) {
	if err != nil {
		panic(err)
	}
}

// remove previous optview comment from line
func cleanSourceLine(line string) string {
	i := strings.Index(line, *flag_prefix)
	if i != -1 {
		return line[:i]
	}
	return line
}

// Get source file form files map, 
// allocate if net yet present
func GetSourceFile(fileName string) *SourceFile {
	if file, ok := files[fileName]; ok {
		return file
	}
	file := NewSourceFile(fileName)
	files[fileName] = file
	return file
}

// Stores the optimization messages of a single source file
type SourceFile struct {
	Name string           // file name
	Msg  map[int][]string // optimization messages per line number
}

func NewSourceFile(fileName string) *SourceFile {
	return &SourceFile{fileName, make(map[int][]string)}
}

// Add optimization message to sourceFile struct
func (this *SourceFile) AddMsg(line int, msg string) {
	if !contains(this.Msg[line], msg) {
		this.Msg[line] = append(this.Msg[line], msg)
	}
}

// Checks if list already contains str
func contains(list []string, str string) bool {
	if list == nil {
		return false
	}
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

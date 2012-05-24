package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
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
	args := flag.Args() //←[ inlining call to flag.Args]

	if *flag_version {
		fmt.Println("go-optview 0.2\nGo", runtime.Version()) //←[ inlining call to runtime.Version  main ... argument does not escape]
		return
	}

	var in io.Reader

	if len(args) > 0 {
		var err error
		cmd := exec.Command("go", append([]string{"tool"}, args...)...) //←[ main []string literal does not escape]
		in, err = cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err) //←[ main ... argument does not escape]
		}
		err = cmd.Start()
		if err != nil {
			log.Fatal(err) //←[ main ... argument does not escape]
		}
	} else {
		in = os.Stdin
	}

	ReadCompilerOutput(in)

	for name, f := range files {
		if *flag_writeback {
			buf := new(bytes.Buffer) //←[ new(bytes.Buffer) escapes to heap]
			f.WriteTo(buf)           // TODO: return err, don't writeback on err!
			out, _ := os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
			out.Write(buf.Bytes()) //←[ inlining call to Bytes]
		} else {
			fmt.Fprintln(os.Stdout, *flag_prefix, name, ":") //←[ main ... argument does not escape]
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
func ReadCompilerOutput(in_ io.Reader) { //←[ leaking param: in_]
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
func parseCompilerLine(line string) { //←[ leaking param: line]
	defer func() { //←[ parseCompilerLine func literal does not escape]
		err := recover()
		if err != nil {
			Stderr(`optview:`, err, ` while parsing "`, line, `"`) //←[ leaking closure reference line  <S> ... argument does not escape  _func_001 ... argument does not escape]
		}
	}()
	words := strings.SplitN(line, ":", 3)
	filename := words[0]
	lineNo := Atoi(words[1])
	info := words[2]
	GetSourceFile(filename).AddMsg(lineNo, info)
}

// Print source code + optimization messages to out
func (f *SourceFile) WriteTo(out io.Writer) { //←[ leaking param: out  (*SourceFile).WriteTo f does not escape]
	in_, err := os.Open(f.Name)
	if err != nil {
		Stderr(err) //←[ (*SourceFile).WriteTo ... argument does not escape]
		return
	}
	in := bufio.NewReader(in_)

	lineNo := 1
	for l, prefix, err := in.ReadLine(); err == nil; l, prefix, err = in.ReadLine() {
		if prefix {
			panic("Enlarge your buffer!")
		}

		// print source line
		Check(fmt.Fprint(out, cleanSourceLine(string(l)))) //←[ (*SourceFile).WriteTo ... argument does not escape]

		// print messages
		if !*flag_clean {
			if opts, ok := f.Msg[lineNo]; ok {
				Check(fmt.Fprint(out, *flag_prefix, opts)) //←[ (*SourceFile).WriteTo ... argument does not escape]
			}
		}
		fmt.Fprintln(out)
		lineNo++
	}
}

func Check(n int, err error) { //←[ leaking param: err]
	if err != nil {
		panic(err)
	}
}

// remove previous optview comment from line
func cleanSourceLine(line string) string { //←[ leaking param: line]
	i := strings.Index(line, *flag_prefix)
	if i != -1 {
		return line[:i]
	}
	return line
}

// Get source file form files map, 
// allocate if net yet present
func GetSourceFile(fileName string) *SourceFile { //←[ leaking param: fileName]
	if file, ok := files[fileName]; ok {
		return file
	}
	file := NewSourceFile(fileName) //←[ inlining call to NewSourceFile  make(map[int][]string, 0) escapes to heap  &SourceFile literal escapes to heap]
	files[fileName] = file
	return file
}

// Stores the optimization messages of a single source file
type SourceFile struct {
	Name string           // file name
	Msg  map[int][]string // optimization messages per line number
}

func NewSourceFile(fileName string) *SourceFile { //←[ can inline NewSourceFile  leaking param: fileName]
	return &SourceFile{fileName, make(map[int][]string)} //←[ make(map[int][]string, 0) escapes to heap  &SourceFile literal escapes to heap]
}

// Add optimization message to sourceFile struct
func (this *SourceFile) AddMsg(line int, msg string) { //←[ leaking param: msg  (*SourceFile).AddMsg this does not escape]
	if !contains(this.Msg[line], msg) {
		this.Msg[line] = append(this.Msg[line], msg)
	}
}

// Checks if list already contains str
func contains(list []string, str string) bool { //←[ contains list does not escape  contains str does not escape]
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

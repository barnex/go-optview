/*
Go-optview takes the output of gc -m (compiler's optimization decissions) and presents it side-by-side with the source code.

Use:

	go-optview [flags] [compiler command]

Flags:

	-V: print version and exit
	-c: clean existing optimization messages from source
	-prefix: prefix for optimization messages, default "//←"
	-w: write result to source files instead of stdout

The compiler command is typically "6g -m", "go tool" is automatically prefixed. If none is specified, go-optview reads gc output from stdin.

Usage examples:

Write optimization messages back to source files:
	go-optview -w 6g -m *.go

Clean source files:
	go-optview -w -c 6g -m *.go

Write optimization messages to stdout:
	go-optview 6g -m *.go

Use of stdin:
	go tool 6g -m | go-optview -w

Example output:

	func MakeSlice(N int) Slice { //←[ can inline MakeSlice]
		return make(Slice, N) //←[ make(Slice, N) escapes to heap]
	}

	func (s Slice) N() int { //←[ can inline Slice.N  Slice.N s does not escape]
		return len(s)
	}
*/
package main
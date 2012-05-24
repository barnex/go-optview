all: *.go
	go install -v
	gofmt -w *.go


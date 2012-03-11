all: *.go
	gofmt -w *.go
	go tool 6g -o optview.6 *.go
	go tool 6l -o optview optview.6


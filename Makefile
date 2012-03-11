all: *.go
	gofmt -w *.go
	go tool 6g -m -o optview.6 *.go
	go tool 6l -m -o optview optview.6


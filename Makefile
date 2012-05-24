all: *.go
	go install -v
	gofmt -w *.go

opt: all *.go
	go-optview -w 6g -m *.go
	gofmt -w *.go

.PHONY: clean
clean:
	go clean
	go-optview -w -c 6g -m *.go
	gofmt -w *.go
	

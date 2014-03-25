goproxy: $(shell find . -name "*.go")
	GOPATH=`pwd` go build src/cmd/goproxy.go

run: goproxy
	./goproxy

.PHONY: run clean

clean:
	-rm goproxy
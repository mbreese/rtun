SOURCES := $(shell find . -name '*.go')

bin/rtun.darwin_amd64: bin/rtun.linux_amd64
	GOOS=darwin GOARCH=amd64 go build -o bin/rtun.darwin_amd64 main.go

bin/rtun.linux_amd64: go.mod go.sum $(SOURCES) 
	go build -o bin/rtun.linux_amd64 main.go

run:
	go run main.go

clean:
	rm bin/*

.PHONY: run clean test

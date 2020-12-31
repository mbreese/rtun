SOURCES := $(shell find . -name '*.go')

bin/rtun.macos: bin/rtun
	GOOS=darwin GOARCH=amd64 go build -o bin/rtun.macos main.go

bin/rtun: go.mod go.sum $(SOURCES) 
	go build -o bin/rtun main.go

run:
	go run main.go

clean:
	rm bin/*

.PHONY: run clean test


.PHONY: all build run clean

all: build

build:
	go build -o bin/mini_agent ./cmd/mini_agent

run: build
	./bin/mini_agent

clean:
	rm -f bin/mini_agent

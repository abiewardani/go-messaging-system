.PHONY: build run clean

build:
	go build -o bin/server ./cmd/server

run: build
	./bin/server

clean:
	go clean
	rm -rf bin/
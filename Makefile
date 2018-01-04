all: gocom

gocom: $(shell find . -name '*.go')
	go build

test:
	go test ./...

clean:
	rm -f gocom

.PHONY: clean test all

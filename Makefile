NAME			= gkeepassxreader
BINARY_PATH 	= $(NAME)

.PHONY: all install build test rebuild clean check

all: build test

build:
	go build -o $(BINARY_PATH)

test:
	ginkgo -race -cover -progress -keepGoing ./...

rebuild:
	go build -v -race -a -o $(BINARY_PATH)

check: test
	@echo Checking code is gofmted
	@bash -c 'if [ -n "$(gofmt -s -l .)" ]; then echo "Go code is not formatted:"; gofmt -s -d -e .; exit 1;fi'

clean:
	@rm -f $(BINARY_PATH)

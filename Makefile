NAME			= gkeepassxreader
BINARY_PATH 	= $(NAME)

.PHONY: all install build test rebuild clean

all: install build test

install:
	dep ensure

build:
	go build -o $(BINARY_PATH)

test:
	ginkgo -race -cover -progress -keepGoing ./...

rebuild:
	go build -v -race -a -o $(BINARY_PATH)

clean:
	@rm -f $(BINARY_PATH)

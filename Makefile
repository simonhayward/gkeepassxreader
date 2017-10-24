NAME			= gkeepassxreader
BINARY_PATH 	= $(NAME)
PACKAGES		= `glide novendor | grep -v "fakes" | grep -v -x "."`

.PHONY: all install build test rebuild clean

all: install test

install:
	glide install

build:
	go build -o $(BINARY_PATH)

test:
	ginkgo -race -cover --progress $(PACKAGES)

rebuild:
	go build -v -race -a -o $(BINARY_PATH)

clean:
	@rm -f $(BINARY_PATH)

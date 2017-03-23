.PHONY: all build lint clean

all: build

build:
	go build -v -o uoft-wifi-usage main.go

lint:
	${GOPATH}/bin/golint .

clean:
	rm uoft-wifi-usage

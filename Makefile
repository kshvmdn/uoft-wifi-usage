all: uoft-wifi-usage

uoft-wifi-usage:
	go build -o uoft-wifi-usage main.go

clean:
	rm uoft-wifi-usage

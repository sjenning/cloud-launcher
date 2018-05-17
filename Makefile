all: build

build:
	go build github.com/sjenning/cloud-launcher/cmd/cloud-launcher

clean:
	rm cloud-launcher

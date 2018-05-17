all: build

build:
	go build github.com/sjenning/cloud-launcher/cmd/aws-launcher

clean:
	rm cloud-launcher
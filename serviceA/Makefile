export tag=v1.0
root:
	export ROOT=github.com/cncamp/service2

build:
	echo "building serviceA binary"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .

release: build
	echo "building service2 container"
	docker build -t winkyi/servicea:${tag} .
export tag=v1.0
root:
	set ROOT=github.com/winkyi/CloudNative

build-linux:
	echo "building httpserver binary"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o deploy/amd64 .


build-windows:
	echo "building httpserver binary"
	go build -o deploy/windows .
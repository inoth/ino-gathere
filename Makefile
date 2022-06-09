
#all: build down setup clean log

start: clean build package
start-arm: clean build-arm package-arm

clean:
	- rm -vrf ./release/*

build:
	- go mod tidy
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./release/ino-gathere main.go
	cp setup.sh release/
	cp ./private_key.pem release/	
	cp ./public_key.pem release/
	cp Dockerfile release/
	- mkdir release/config
	cp config/*.yaml release/config/

build-arm:
	- go mod tidy
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./release/ino-gathere main.go
	cp setup.sh release/
	cp ./private_key.pem release/
	cp ./public_key.pem release/
	cp Dockerfile-arm release/
	- mkdir release/config
	cp config/*.yaml release/config/.

package:
	#sudo docker build -t ino-gathere:latest .
	sudo docker buildx build -t ino-gathere:latest --platform linux/amd64 .

package-arm:
	export DOCKER_CLI_EXPERIMENTAL=enabled
	sudo docker buildx build -t ino-gathere:latest -f Dockerfile-arm --platform linux/arm64 .

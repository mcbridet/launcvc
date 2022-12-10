build:
	go build -o bin/launstat_linux_amd64 main.go

build-all:
	./build-all.sh

run:
	go run main.go
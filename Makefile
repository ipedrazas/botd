ITTERATION := $(shell date +%s)

all: build

deps:
	go get github.com/gorilla/mux

test:
	go test api/api_test.go

build:
	CGO_ENABLED=0 GOOS=linux go build  -a -installsuffix cgo -o botd .

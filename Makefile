ITTERATION := $(shell date +%s)

all: build

deps:
	go get github.com/gorilla/mux
	go get gopkg.in/redis.v3

test:
	go test api/api_test.go

build:
	CGO_ENABLED=0 GOOS=linux go build  -a -installsuffix cgo -o botd .

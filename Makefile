ITTERATION := $(shell date +%s)

all: build

deps:
	go get github.com/gorilla/mux

test:
	@test -z "$(shell find . -name '*.go' | xargs gofmt -l)" || (echo "Need to run 'go fmt ./...'"; exit 1)
	go test -cover -short ./...

build:
	CGO_ENABLED=0 GOOS=linux go build  -a -installsuffix cgo -o botd .

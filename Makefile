ITTERATION := $(shell date +%s)

all: build

deps:
	go get github.com/gorilla/mux
	go get -t -v ./...

test:
	@test -z "$(shell find . -name '*.go' | xargs gofmt -l)" || (echo "Need to run 'go fmt ./...'"; exit 1)
	go test -cover -short ./...

build:
	go build CGO_ENABLED=0 GOOS=linux -a -installsuffix cgo -o botd .

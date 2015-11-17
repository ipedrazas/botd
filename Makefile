ITTERATION := $(shell date +%s)

all: build

deps:
	go get github.com/gorilla/mux
	go get gopkg.in/redis.v3
	go get -u github.com/jstemmer/go-junit-report

test:
	go test -v api/api_test.go | go-junit-report > report.xml
	cat report.xml

build:
	CGO_ENABLED=0 GOOS=linux go build  -a -installsuffix cgo -o botd .

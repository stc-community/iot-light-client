.PHONY: all build

PROG=bin/iotc
SRCS=.

# git commit hash
COMMIT_HASH=$(shell git rev-parse --short HEAD || echo "GitNotFound")
BUILD_DATE=$(shell date '+%Y-%m-%d %H:%M:%S')
CFLAGS = -ldflags "-s -w -X \"main.BuildVersion=${COMMIT_HASH}\" -X \"main.BuildDate=$(BUILD_DATE)\""

all:
	if [ ! -d "./bin/" ]; then \
	mkdir bin; \
	fi
	GOPROXY=$(GOPROXY) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  $(CFLAGS) -o $(PROG) $(SRCS)

build:
	go build -race -tags=jsoniter


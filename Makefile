# Modified Makefile from https://gist.github.com/Stratus3D/a5be23866810735d7413

.PHONY: build doc fmt lint run test install vet

# Prepend our _vendor directory to the system GOPATH
# so that import path resolution will prioritize
# our third party snapshots.
GOPATH := ${PWD}/_vendor:${GOPATH}
export GOPATH

default: build

build: vet
	go build -v -o ./bin/ftpbot ./src/

doc:
	godoc -http=:6060 -index

# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
fmt:
	go fmt ./src/...

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint:
	golint ./src

run: build
	./bin/ftpbot

test:
	go test ./src/...


install:
	cd src && go get && go install

vet:
	go vet ./src/...

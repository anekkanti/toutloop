all: build test prof

dependencies:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

build: dependencies
	go build ./...

test: build
	go test ./...

prof: build
	go test -bench=. -benchmem -cpuprofile profile.out

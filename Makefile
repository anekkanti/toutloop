all: build test prof

dependencies:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

build: dependencies
	go build ./...

tests: build
	go test -v -coverprofile=cover.out ./...

prof: build
	go test -v -bench=. -benchmem -cpuprofile profile.out

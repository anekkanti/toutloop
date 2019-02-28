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

documentation:
	go get github.com/robertkrimen/godocdown/godocdown
	godocdown github.com/anekkanti/toutloop > docs/README

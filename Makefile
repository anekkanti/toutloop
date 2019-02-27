all: build test prof

build: 
	go build ./...

test:
	go test ./...

prof:
	go test -bench=. -benchmem -cpuprofile profile.out

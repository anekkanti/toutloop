# Timeout Loop

[![CircleCI](https://circleci.com/gh/anekkanti/toutloop.svg?style=svg)](https://circleci.com/gh/anekkanti/toutloop)
[![GoDoc](https://godoc.org/github.com/anekkanti/toutloop?status.svg)](https://godoc.org/github.com/anekkanti/toutloop)
[![Go Report Card](https://goreportcard.com/badge/github.com/anekkanti/toutloop)](https://goreportcard.com/report/github.com/anekkanti/toutloop)
[![codecov](https://codecov.io/gh/anekkanti/toutloop/branch/master/graph/badge.svg)](https://codecov.io/gh/anekkanti/toutloop)

Useful when scheduling at scale. Uses a heap to track timeouts.

Example
```go

type tjob struct {
	name string
}

func main() {

        tloop := New()
	tloop.Run()

	j1 := &tjob{name: "j1"}
	err := tloop.Add(j1.name, j1, time.Millisecond*300)
	if err != nil {
		panic(err)
	}

	err = tloop.Reschedule(j1.name, time.Millisecond*400)
	if err != nil {
		panic(err)
	}

	for j := range tloop.C {
		if j.(*tjob) == j1 {
			break
		}
	}
	tloop.Stop()
}
```

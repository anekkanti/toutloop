# Timeout Loop

[![Codeship Status for anekkanti/toutloop](https://app.codeship.com/projects/68383590-1da4-0137-3f03-0ec35b27b473/status?branch=master)](https://app.codeship.com/projects/329066)
[![GoDoc](https://godoc.org/github.com/anekkanti/toutloop?status.svg)](https://godoc.org/github.com/anekkanti/toutloop)

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

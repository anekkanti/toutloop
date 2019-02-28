# Timeout Loop

[![GoDoc](https://godoc.org/github.com/anekkanti/toutloop?status.svg)](https://godoc.org/github.com/anekkanti/toutloop)

Useful to scheduling at scale. Uses a timeout heap to track timeouts.

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

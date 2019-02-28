# toutloop
--
    import "github.com/anekkanti/toutloop"


## Usage

#### type ToutLoop

```go
type ToutLoop struct {
	C chan interface{}
}
```

ToutLoop or the timeout loop. The loop uses a heap to track and dispatches
events when their timeout's expire Listen to C to recieve events

#### func  New

```go
func New() *ToutLoop
```
New returns a new timeout looop

#### func (*ToutLoop) Add

```go
func (e *ToutLoop) Add(id string, object interface{}, after time.Duration) error
```
Add object with given id to be returned after given time

#### func (*ToutLoop) Remove

```go
func (e *ToutLoop) Remove(id string) error
```
Remove the object with the given id from the loop

#### func (*ToutLoop) Reschedule

```go
func (e *ToutLoop) Reschedule(id string, after time.Duration) error
```
Reschedule the object with the given id

#### func (*ToutLoop) Run

```go
func (e *ToutLoop) Run()
```
Run the timeout loop

#### func (*ToutLoop) Stop

```go
func (e *ToutLoop) Stop()
```
Stop the event loop

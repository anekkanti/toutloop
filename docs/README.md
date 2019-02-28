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

ToutLoop or the timeout event loop

#### func  NewToutLoop

```go
func NewToutLoop() *ToutLoop
```
NewToutLoop for scheduling stuff

#### func (*ToutLoop) Add

```go
func (e *ToutLoop) Add(id string, object interface{}, after time.Duration) error
```
Add job to run after given time

#### func (*ToutLoop) Remove

```go
func (e *ToutLoop) Remove(id string) error
```
Remove job

#### func (*ToutLoop) Reschedule

```go
func (e *ToutLoop) Reschedule(id string, after time.Duration) error
```
Reschedule job to run after given time

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

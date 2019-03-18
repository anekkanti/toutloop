package toutloop

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

type operation int

const (
	addOp        operation = iota
	rescheduleOp operation = iota
	removeOp     operation = iota
)

type request struct {
	operation operation
	id        string
	object    interface{}
	runTime   time.Time
	reply     bool
}

type timeout struct {
	id      string
	object  interface{}
	runTime time.Time
	index   int
}
type toutHeap []*timeout

func (t toutHeap) Len() int           { return len(t) }
func (t toutHeap) Less(i, j int) bool { return t[i].runTime.Before(t[j].runTime) }
func (t toutHeap) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
	t[i].index = i
	t[j].index = j
}

func (t *toutHeap) Push(x interface{}) {
	n := len(*t)
	tout := x.(*timeout)
	tout.index = n
	*t = append(*t, tout)
}

func (t *toutHeap) Pop() interface{} {
	old := *t
	n := len(old)
	x := old[n-1]
	x.index = -1
	*t = old[0 : n-1]
	return x
}

// ToutLoop or the timeout loop.
// The loop uses a heap to track and dispatches events when their timeout's expire
// Listen to C to recieve events
type ToutLoop struct {
	heap     toutHeap
	requests chan *request
	reply    chan error
	C        chan interface{}
	wg       sync.WaitGroup
	store    map[string]*timeout
}

func (e *ToutLoop) handleRequest(req *request, dispatchID *string) {
	if req == nil {
		return
	}
	var err error
	switch req.operation {
	case addOp:
		// ignore request if the object already exists in store
		if _, ok := e.store[req.id]; !ok {
			tout := &timeout{
				id:      req.id,
				object:  req.object,
				runTime: req.runTime,
			}
			e.store[req.id] = tout
			heap.Push(&e.heap, tout)
		} else {
			err = fmt.Errorf("object with id=%s already exists", req.id)
		}
	case rescheduleOp:
		// ignore request if the object does not exists in store
		if tout, ok := e.store[req.id]; ok {
			tout.runTime = req.runTime
			if tout.index != -1 {
				heap.Fix(&e.heap, tout.index)
			} else {
				heap.Push(&e.heap, tout)
			}
		} else {
			err = fmt.Errorf("object with id=%s does not exists", req.id)
		}
	case removeOp:
		// ignore request if the object does not exist in store
		if tout, ok := e.store[req.id]; ok {
			if tout.index != -1 {
				heap.Remove(&e.heap, tout.index)
			}
			delete(e.store, req.id)
		} else {
			err = fmt.Errorf("object with id=%s does not exists", req.id)
		}
	}

	if *dispatchID == req.id {
		*dispatchID = ""
	}
	if req.reply {
		e.reply <- err
	}
}

// Run the timeout loop
func (e *ToutLoop) Run() {
	e.wg.Add(1)
	go func(e *ToutLoop) {
		defer e.wg.Done()
		var dispatchID string
		var toutTimer = time.NewTimer(time.Second)
		var ok = true
		for ok {
			var req *request
			if dispatchID != "" {
				select {
				case req, ok = <-e.requests:
					e.handleRequest(req, &dispatchID)
				case e.C <- e.store[dispatchID].object:
					dispatchID = ""
					delete(e.store, dispatchID)
				}
			} else if len(e.heap) > 0 {
				now := time.Now()
				if e.heap[0].runTime.After(now) {
					toutTimer.Reset(e.heap[0].runTime.Sub(now))
					select {
					case req, ok = <-e.requests:
						e.handleRequest(req, &dispatchID)
						continue
					case <-toutTimer.C:
					}
				}
				tout := heap.Pop(&e.heap).(*timeout)
				dispatchID = tout.id
			} else {
				select {
				case req, ok = <-e.requests:
					e.handleRequest(req, &dispatchID)
				}
			}
		}
		close(e.C)
		close(e.reply)
	}(e)
}

// New returns a new timeout looop
func New() *ToutLoop {
	e := &ToutLoop{
		requests: make(chan *request),
		reply:    make(chan error),
		C:        make(chan interface{}),
		wg:       sync.WaitGroup{},
		store:    make(map[string]*timeout),
	}
	heap.Init(&e.heap)
	return e
}

// Stop the event loop
func (e *ToutLoop) Stop() {
	close(e.requests)
	e.wg.Wait()
}

// Add object with given id to be returned after given time
func (e *ToutLoop) Add(id string, object interface{}, after time.Duration) error {
	e.requests <- &request{
		operation: addOp,
		id:        id,
		object:    object,
		runTime:   time.Now().Add(after),
		reply:     true,
	}
	return <-e.reply
}

// Reschedule the object with the given id
func (e *ToutLoop) Reschedule(id string, after time.Duration) error {
	e.requests <- &request{
		operation: rescheduleOp,
		id:        id,
		runTime:   time.Now().Add(after),
		reply:     true,
	}
	return <-e.reply
}

// Remove the object with the given id from the loop
func (e *ToutLoop) Remove(id string) error {
	e.requests <- &request{
		operation: removeOp,
		id:        id,
		reply:     true,
	}
	return <-e.reply
}

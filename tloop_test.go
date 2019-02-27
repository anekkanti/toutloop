package toutloop

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type tjob struct {
	timeout
	name string
	when time.Time
}

func TestToutLoop(t *testing.T) {
	assert := assert.New(t)
	tloop := NewToutLoop()
	tloop.Run()

	j1 := &tjob{name: "job-1", when: time.Now().Add(time.Millisecond * 300)}
	err := tloop.Add(j1.name, j1, time.Millisecond*300)
	assert.NoError(err)

	j2 := &tjob{name: "job-2", when: time.Now().Add(time.Millisecond * 200)}
	err = tloop.Add(j2.name, j2, time.Millisecond*200)
	assert.NoError(err)
	err = tloop.Add(j2.name, j2, time.Millisecond*200)
	assert.Error(err)

	j3 := &tjob{name: "job-3", when: time.Now().Add(time.Millisecond * 100)}
	err = tloop.Add(j3.name, j3, time.Millisecond*100)
	assert.NoError(err)

	j4 := &tjob{name: "job-4", when: time.Now().Add(time.Millisecond * 400)}
	err = tloop.Add(j4.name, j4, time.Millisecond*100)
	assert.NoError(err)
	err = tloop.Reschedule(j4.name, time.Millisecond*400)
	assert.NoError(err)
	err = tloop.Reschedule("non-existing", time.Millisecond*400)
	assert.Error(err)

	j5 := &tjob{name: "job-5", when: time.Now().Add(time.Millisecond * 520)}
	err = tloop.Add(j5.name, j5, time.Millisecond*10)
	assert.NoError(err)
	time.Sleep(time.Millisecond * 20)
	err = tloop.Reschedule(j5.name, time.Millisecond*500)
	assert.NoError(err)

	j6 := &tjob{name: "job-6", when: time.Now().Add(time.Millisecond * 100)}
	err = tloop.Add(j6.name, j6, time.Millisecond*100)
	assert.NoError(err)
	err = tloop.Remove(j6.name)
	assert.NoError(err)
	err = tloop.Remove("non-existing")
	assert.Error(err)

	j7 := &tjob{name: "job-7", when: time.Now().Add(time.Millisecond * 100)}
	err = tloop.Add(j7.name, j7, time.Millisecond*10)
	assert.NoError(err)
	time.Sleep(time.Millisecond * 20)
	err = tloop.Remove(j7.name)
	assert.NoError(err)

	j := (<-tloop.C).(*tjob)
	assert.Equal(j3, j)
	assert.WithinDuration(time.Now(), j.when, time.Millisecond*5)

	j = (<-tloop.C).(*tjob)
	assert.Equal(j2, j)
	assert.WithinDuration(time.Now(), j.when, time.Millisecond*5)

	j = (<-tloop.C).(*tjob)
	assert.Equal(j1, j)
	assert.WithinDuration(time.Now(), j.when, time.Millisecond*5)

	j = (<-tloop.C).(*tjob)
	assert.Equal(j4, j)
	assert.WithinDuration(time.Now(), j.when, time.Millisecond*5)

	j = (<-tloop.C).(*tjob)
	assert.Equal(j5, j)
	assert.WithinDuration(time.Now(), j.when, time.Millisecond*5)

	select {
	case <-tloop.C:
		assert.FailNow("should not have recieved any jobs")
	case <-time.Tick(500 * time.Millisecond):
		break
	}
	tloop.Stop()
}

func BenchmarkToutLoopNJobs(t *testing.B) {
	assert := assert.New(t)
	tloop := NewToutLoop()
	tloop.Run()

	numberOfJobs := 10000

	go func() {
		for i := 0; i < numberOfJobs; i++ {
			after := time.Duration(rand.Int()%500) * time.Millisecond
			j := &tjob{name: fmt.Sprintf("job-%d", i), when: time.Now().Add(after)}
			err := tloop.Add(j.name, j, after)
			assert.NoError(err)
		}
	}()

	count := 0
	var deltaSum time.Duration
	var deltaMax time.Duration
	for o := range tloop.C {
		j := o.(*tjob)
		delta := time.Now().Sub(j.when)
		deltaSum = deltaSum + delta
		if deltaMax < delta {
			deltaMax = delta
		}
		count++
		if count == numberOfJobs {
			tloop.Stop()
		}
	}
	deltaAvg := deltaSum / time.Duration(count)
	t.Logf("done: recieved %d events", count)
	t.Logf("avg delta: %s", deltaAvg)
	t.Logf("max delta: %s", deltaMax)

	assert.Equal(true, deltaAvg < time.Millisecond*1)
	assert.Equal(true, deltaMax < time.Millisecond*10)
	assert.Equal(numberOfJobs, count)
}

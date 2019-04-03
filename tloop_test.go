package toutloop

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type tjob struct {
	name string
	when time.Time
}

func TestExample(t *testing.T) {
	if !testing.Short() {
		t.Skip("skipping perf tests in perf mode")
	}
	tloop := New(0 /*recieveChanBufferSize*/)
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

func TestToutLoop(t *testing.T) {
	if !testing.Short() {
		t.Skip("skipping perf tests in perf mode")
	}
	assert := assert.New(t)
	tloop := New(0)
	tloop.Run()

	j1 := &tjob{name: "job-1", when: time.Now().Add(time.Millisecond * 300)}
	err := tloop.Add(j1.name, j1, time.Millisecond*300)
	assert.NoError(err)
	assert.True(tloop.Exists(j1.name))
	assert.False(tloop.Exists("non-existing"))

	j2 := &tjob{name: "job-2", when: time.Now().Add(time.Millisecond * 200)}
	err = tloop.Add(j2.name, j2, time.Millisecond*200)
	assert.NoError(err)
	assert.True(tloop.Exists(j2.name))
	err = tloop.Add(j2.name, j2, time.Millisecond*200)
	assert.Error(err)

	j3 := &tjob{name: "job-3", when: time.Now().Add(time.Millisecond * 100)}
	err = tloop.Add(j3.name, j3, time.Millisecond*100)
	assert.NoError(err)
	assert.True(tloop.Exists(j3.name))

	j4 := &tjob{name: "job-4", when: time.Now().Add(time.Millisecond * 400)}
	err = tloop.Add(j4.name, j4, time.Millisecond*100)
	assert.NoError(err)
	err = tloop.Reschedule(j4.name, time.Millisecond*400)
	assert.NoError(err)
	assert.True(tloop.Exists(j4.name))
	err = tloop.Reschedule("non-existing", time.Millisecond*400)
	assert.Error(err)

	j5 := &tjob{name: "job-5", when: time.Now().Add(time.Millisecond * 520)}
	err = tloop.Add(j5.name, j5, time.Millisecond*10)
	assert.NoError(err)
	time.Sleep(time.Millisecond * 20)
	err = tloop.Reschedule(j5.name, time.Millisecond*500)
	assert.NoError(err)
	assert.True(tloop.Exists(j5.name))

	j6 := &tjob{name: "job-6", when: time.Now().Add(time.Millisecond * 100)}
	err = tloop.Add(j6.name, j6, time.Millisecond*100)
	assert.True(tloop.Exists(j6.name))
	assert.NoError(err)
	err = tloop.Remove(j6.name)
	assert.NoError(err)
	assert.False(tloop.Exists(j6.name))
	err = tloop.Remove("non-existing")
	assert.Error(err)

	j7 := &tjob{name: "job-7", when: time.Now().Add(time.Millisecond * 100)}
	err = tloop.Add(j7.name, j7, time.Millisecond*10)
	assert.NoError(err)
	time.Sleep(time.Millisecond * 20)
	err = tloop.Remove(j7.name)
	assert.NoError(err)
	assert.False(tloop.Exists(j7.name))

	j := (<-tloop.C).(*tjob)
	assert.Equal(j3, j)
	assert.WithinDuration(time.Now(), j.when, time.Millisecond*5)
	assert.False(tloop.Exists(j3.name))
	assert.True(tloop.Exists(j1.name))

	j = (<-tloop.C).(*tjob)
	assert.Equal(j2, j)
	assert.WithinDuration(time.Now(), j.when, time.Millisecond*5)
	assert.False(tloop.Exists(j2.name))
	assert.True(tloop.Exists(j1.name))

	j = (<-tloop.C).(*tjob)
	assert.Equal(j1, j)
	assert.WithinDuration(time.Now(), j.when, time.Millisecond*5)
	assert.False(tloop.Exists(j1.name))

	j = (<-tloop.C).(*tjob)
	assert.Equal(j4, j)
	assert.WithinDuration(time.Now(), j.when, time.Millisecond*5)
	assert.False(tloop.Exists(j4.name))

	j = (<-tloop.C).(*tjob)
	assert.Equal(j5, j)
	assert.WithinDuration(time.Now(), j.when, time.Millisecond*5)
	assert.False(tloop.Exists(j5.name))

	select {
	case <-tloop.C:
		assert.FailNow("should not have received any jobs")
	case <-time.Tick(500 * time.Millisecond):
		break
	}
	tloop.Stop()
}

func runToutLoopWithNJobs(numberOfJobsPerSec int64, assert *assert.Assertions) (avg, max time.Duration) {
	tloop := New(10)
	tloop.Run()

	mult := int64(1)

	go func() {
		for i := int64(0); i < numberOfJobsPerSec*mult; i++ {
			after := time.Duration(rand.Int()%(int(mult)*1000)) * time.Millisecond
			j := &tjob{name: fmt.Sprintf("job-%d", i), when: time.Now().Add(after)}
			err := tloop.Add(j.name, j, after)
			assert.NoError(err)
		}
	}()

	count := int64(0)
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
		if count == numberOfJobsPerSec*mult {
			tloop.Stop()
		}
	}

	assert.Equal(numberOfJobsPerSec*mult, count)

	deltaAvg := deltaSum / time.Duration(count)
	return deltaAvg, deltaMax
}

func TestToutLoop1KJobs(t *testing.T) {
	if !testing.Short() {
		t.Skip("skipping perf tests in perf mode")
	}
	assert := assert.New(t)
	deltaAvg, deltaMax := runToutLoopWithNJobs(1000, assert)
	t.Logf("avg delta: %s", deltaAvg)
	t.Logf("max delta: %s", deltaMax)
}

func TestToutLoopPerf(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping perf tests in short mode")
	}
	assert := assert.New(t)

	f, err := os.Create("prof.csv")
	if err != nil {
		t.Logf("failed to prof.csv, err=%s", err)
	}

	if f != nil {
		f.WriteString("rate(events per second),avg delay(ms),max delay(ms)\n")
	}

	count := int64(1)
	for i := 0; i < 23; i++ {
		deltaAvg, deltaMax := runToutLoopWithNJobs(count, assert)
		t.Logf("count: %d", count)
		t.Logf("avg delta: %s", deltaAvg)
		t.Logf("max delta: %s", deltaMax)
		f.WriteString(fmt.Sprintf("%d,%f,%f\n", count, deltaAvg.Seconds()*1000, deltaMax.Seconds()*1000))
		count = count * 2
	}
}

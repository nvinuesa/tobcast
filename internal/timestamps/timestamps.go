package timestamps

import (
	"math"
	"sync"
)

type Timestamps struct {
	times map[int]int64
	mux   sync.Mutex
}

func New(ports []int) *Timestamps {
	init := make(map[int]int64)
	for _, port := range ports {
		init[port] = math.MaxInt64
	}
	return &Timestamps{times: init}
}

func (t *Timestamps) Get(port int) int64 {
	t.mux.Lock()
	res := t.times[port]
	defer t.mux.Unlock()
	return res
}

func (t *Timestamps) Set(port int, timestamp int64) {
	t.mux.Lock()
	t.times[port] = timestamp
	t.mux.Unlock()
}

func (t *Timestamps) Incr(port int) int64 {
	t.mux.Lock()
	res := t.times[port] + 1
	t.times[port] = res
	defer t.mux.Unlock()
	return res
}

func (t *Timestamps) IncrOrSet(port int, timestamp int64) {
	t.mux.Lock()
	if timestamp > t.times[port] {
		t.times[port] = timestamp + 1
	} else {
		t.Incr(port)
	}
	t.mux.Unlock()
}

func (t *Timestamps) Min() int64 {
	t.mux.Lock()
	keys := make([]int, 0, len(t.times))
	for k := range t.times {
		keys = append(keys, k)
	}
	min := t.times[keys[0]]
	for _, v := range t.times {
		if v < min {
			min = v
		}
	}
	defer t.mux.Unlock()
	return min
}

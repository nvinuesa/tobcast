package timestamps

type Timestamps struct {
	times map[int]int64
}

func New(ports []int) *Timestamps {
	init := make(map[int]int64)
	for _, port := range ports {
		init[port] = 0
	}
	return &Timestamps{times: init}
}

func (t *Timestamps) Incr(port int) {
	t.times[port] = t.times[port] + 1
}

func (t *Timestamps) Get(port int) int64 {
	return t.times[port]
}

func (t *Timestamps) Update(port int, timestamp int64) {
	if timestamp > t.times[port] {
		t.times[port] = timestamp + 1
	} else {
		t.Incr(port)
	}
}

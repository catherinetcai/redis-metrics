package decoder

import "strings"

type Counters []*Counter

// Finds or creates a counter to the key and increment
func (c Counters) FindOrCreateCounter(key string) {
	// If the counter exists, then increment
	for _, cnt := range c {
		if strings.Contains(key, cnt.Key) {
			cnt.Increment()
		} else {
			// Else, create counter and increment it
			ncnt := NewCounter(key)
			ncnt.Increment()
			c = append(c, ncnt)
		}
	}
}

type Counter struct {
	Key   string
	Count int
}

func NewCounter(key string) *Counter {
	return &Counter{Key: key}
}

func (c *Counter) Increment() {
	c.Count++
}

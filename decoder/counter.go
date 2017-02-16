package decoder

import (
	"fmt"
	"os"
	"strings"
)

// ** THIS IS NOW DEPRECRATED, WE CAN DUMP THIS PROBABLY **
type Counters []*Counter

// Finds or creates a counter to the key and increment
func (c *Counters) FindOrCreateCounter(key string, val []byte) {
	mem := len(val)
	if len(*c) < 1 {
		ncnt := NewCounter(key)
		ncnt.Increment()
		ncnt.AddMemory(mem)
		*c = append(*c, ncnt)
		return
	}
	// If the counter exists, then increment
	for i, _ := range *c {
		if strings.EqualFold(key, (*c)[i].Key) {
			(*c)[i].Increment()
			(*c)[i].AddMemory(mem)
		} else {
			// Else, create counter and increment it
			ncnt := NewCounter(key)
			ncnt.Increment()
			ncnt.AddMemory(mem)
			*c = append(*c, ncnt)
		}
	}
}

func (c *Counters) WriteCountersToFile(f *os.File) {
	for _, cnt := range *c {
		_, err := f.WriteString(fmt.Sprintf("%s: %d - %d", cnt.Key, cnt.Count, cnt.Memory))
		if err != nil {
			fmt.Sprintf("Error writing to file: %v\n", err)
		}
	}
}

type Counter struct {
	Key    string
	Count  int
	Memory int
}

func NewCounter(key string) *Counter {
	return &Counter{Key: key, Count: 0, Memory: 0}
}

func (c *Counter) Increment() {
	c.Count++
}

func (c *Counter) AddMemory(mem int) {
	c.Memory = c.Memory + mem
}

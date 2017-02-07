package decoder

import (
	"fmt"
	"regexp"

	metrics "github.com/rcrowley/go-metrics"
)

const (
	DefaultCounterInc    = 1
	DefaultReservoirSize = 1028
	keyRegex             = "(?P<key>\\D+):?(?P<rest>[0-9]+)"
	counter              = "counter"
	memCounter           = "memCounter"
	memHistogram         = "memHistogram"
)

var keyMatchRegexp = regexp.MustCompile(keyRegex)
var defaultSample = metrics.NewUniformSample(DefaultReservoirSize)

type MatchKeys struct {
	Keys Keys
}

type Keys []string

// if is match, then collect it?
func (k Keys) WriteToCountersIfMatch(key string, val []byte, c *Counters) {
	if len(k) < 1 {
		c.FindOrCreateCounter(key, val)
		return
	}
	mkey := isMatch(key)
	if mkey != "" {
		c.FindOrCreateCounter(mkey, val)
	}
}

func CollectMetricsIfMatch(key string, val []byte, r *metrics.Registry) {
	mkey := isMatch(key)
	if mkey != "" {
		c := metrics.GetOrRegisterCounter(counterKey(mkey), *r)
		c.Inc(DefaultCounterInc)
		if len(val) > 0 {
			collectMemoryMetrics(mkey, val, r)
		}
	}
}

func collectMemoryMetrics(key string, val []byte, r *metrics.Registry) {
	// Increment bytes by key
	c := metrics.GetOrRegisterCounter(memoryCounterKey(key), *r)
	c.Inc(int64(len(val)))
	h := metrics.GetOrRegisterHistogram(memoryHistogramKey(key), *r, defaultSample)
	h.Update(int64(len(val)))
}

// Namespace the key for key counts
func counterKey(key string) string {
	return fmt.Sprintf("%s::%s", key, counter)
}

// Namespace the keyf or memory count
func memoryCounterKey(key string) string {
	return fmt.Sprintf("%s::%s", key, memCounter)
}

func memoryHistogramKey(key string) string {
	return fmt.Sprintf("%s::%s", key, memHistogram)
}

func isMatch(key string) string {
	//	for _, mkey := range k {
	matches := keyMatchRegexp.FindStringSubmatch(key)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

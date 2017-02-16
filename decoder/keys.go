package decoder

import (
	"fmt"
	"regexp"

	metrics "github.com/rcrowley/go-metrics"
)

const (
	// keyRegex is the regex for attempting to pull out keys
	keyRegex = "(?P<key>\\D+):?(?P<rest>[0-9]+)"

	// Default counter constants
	DefaultCounterInc    = 1
	DefaultReservoirSize = 1028

	// counter keys
	counter      = "counter"
	memCounter   = "memCounter"
	memHistogram = "memHistogram"
)

// keyMatchRegexp initializes the key regex
var keyMatchRegexp = regexp.MustCompile(keyRegex)

// defaultSample just creates a default uniform sample for the metrics lib
var defaultSample = metrics.NewUniformSample(DefaultReservoirSize)

// MatchKeys is a struct that isn't yet implemented. Will be implemented to allow for specific keys to be registered as a metric
type MatchKeys struct {
	Keys Keys
}

// Keys is just a slice of strings
type Keys []string

// CollectMetricsIfMatch does a regex match on keys, and if it matches, it will increment a counter. If there's length to the val, it means that it should be incremented in byte usage
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

// collectMemoryMetrics just writes a counter for memory and adds it to a histogram
func collectMemoryMetrics(key string, val []byte, r *metrics.Registry) {
	// Increment bytes by key
	c := metrics.GetOrRegisterCounter(memoryCounterKey(key), *r)
	c.Inc(int64(len(val)))
	h := metrics.GetOrRegisterHistogram(memoryHistogramKey(key), *r, defaultSample)
	h.Update(int64(len(val)))
}

// counterKey namespaces the key for key counts
func counterKey(key string) string {
	return fmt.Sprintf("%s::%s", key, counter)
}

// memoryCounterKey namespaces the key for memory counts
func memoryCounterKey(key string) string {
	return fmt.Sprintf("%s::%s", key, memCounter)
}

// memoryHistogramKey namespaces the key for histograms
func memoryHistogramKey(key string) string {
	return fmt.Sprintf("%s::%s", key, memHistogram)
}

// isMatch just does a regexp match to pull out the key
func isMatch(key string) string {
	//	for _, mkey := range k {
	matches := keyMatchRegexp.FindStringSubmatch(key)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

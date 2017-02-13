package decoder

import (
	"fmt"
	"os"

	"github.com/cupcake/rdb/nopdecoder"
	"github.com/davecgh/go-spew/spew"
	metrics "github.com/rcrowley/go-metrics"
)

type Decoder struct {
	db int
	i  int
	nopdecoder.NopDecoder
	OutFile *os.File
	// Counters        *Counters
	// MatchKeys       Keys
	MetricsRegistry *metrics.Registry
}

func (p *Decoder) StartDatabase(n int) {
	p.db = n
}

func (p *Decoder) StartSet(key []byte, cardinality, expiry int64) {
	p.i = 0
	CollectMetricsIfMatch(string(key), []byte{}, p.MetricsRegistry)
}

func (p *Decoder) Sadd(key, member []byte) {
	p.i++
	CollectMetricsIfMatch(string(key), member, p.MetricsRegistry)
}

func (p *Decoder) Set(key, value []byte, expiry int64) {
	CollectMetricsIfMatch(string(key), []byte{}, p.MetricsRegistry)
}

func (p *Decoder) StartHash(key []byte, length, expiry int64) {
	p.i = 0
	CollectMetricsIfMatch(string(key), []byte{}, p.MetricsRegistry)
}

func (p *Decoder) Hset(key, field, value []byte) {
	p.i++
	CollectMetricsIfMatch(string(key), value, p.MetricsRegistry)
}

func (p *Decoder) EndHash(key []byte) {}

func (p *Decoder) StartList(key []byte, length, expiry int64) {
	p.i = 0
	CollectMetricsIfMatch(string(key), []byte{}, p.MetricsRegistry)
}

func (p *Decoder) Rpush(key, value []byte) {
	p.i++
	CollectMetricsIfMatch(string(key), value, p.MetricsRegistry)
}

func (p *Decoder) StartZSet(key []byte, cardinality, expiry int64) {
	p.i = 0
	CollectMetricsIfMatch(string(key), []byte{}, p.MetricsRegistry)
}

func (p *Decoder) Zadd(key []byte, score float64, member []byte) {
	p.i++
	CollectMetricsIfMatch(string(key), member, p.MetricsRegistry)
}

func (p *Decoder) EndRDB() {
	fmt.Println("Finished parsing through keys... writing to file now...")
	spew.Dump(p.MetricsRegistry)
}

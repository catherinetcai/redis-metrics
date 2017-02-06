package decoder

import (
	"fmt"
	"os"

	"github.com/cupcake/rdb/nopdecoder"
)

type Decoder struct {
	db int
	i  int
	nopdecoder.NopDecoder
	OutFile *os.File
}

func (p *Decoder) StartDatabase(n int) {
	p.db = n
}

func (p *Decoder) Set(key, value []byte, expiry int64) {
	fmt.Printf("db=%d %q -> %q\n", p.db, key, value)
}

func (p *Decoder) Hset(key, field, value []byte) {
	fmt.Printf("db=%d %q . %q -> %q\n", p.db, key, field, value)
}

func (p *Decoder) Sadd(key, member []byte) {
	fmt.Printf("db=%d %q { %q }\n", p.db, key, member)
}

func (p *Decoder) StartList(key []byte, length, expiry int64) {
	p.i = 0
}

func (p *Decoder) Rpush(key, value []byte) {
	fmt.Printf("db=%d %q[%d] -> %q\n", p.db, key, p.i, value)
	p.i++
}

func (p *Decoder) StartZSet(key []byte, cardinality, expiry int64) {
	p.i = 0
}

func (p *Decoder) Zadd(key []byte, score float64, member []byte) {
	fmt.Printf("db=%d %q[%d] -> {%q, score=%g}\n", p.db, key, p.i, member, score)
	p.i++
}

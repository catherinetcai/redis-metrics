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
	Outfile *os.File
}

func write(line string, f *os.File) {
	fmt.Print(line)
	f.WriteString(line)
}

func (p *Decoder) StartDatabase(n int) {
	p.db = n
}

func (p *Decoder) Set(key, value []byte, expiry int64) {
	line := fmt.Sprintf("db=%d %q -> %q\n", p.db, key, value)
	write(line, p.Outfile)
}

func (p *Decoder) Hset(key, field, value []byte) {
	line := fmt.Sprintf("db=%d %q . %q -> %q\n", p.db, key, field, value)
	write(line, p.Outfile)
}

func (p *Decoder) Sadd(key, member []byte) {
	line := fmt.Sprintf("db=%d %q { %q }\n", p.db, key, member)
	write(line, p.Outfile)
}

func (p *Decoder) StartList(key []byte, length, expiry int64) {
	p.i = 0
}

func (p *Decoder) Rpush(key, value []byte) {
	line := fmt.Sprintf("db=%d %q[%d] -> %q\n", p.db, key, p.i, value)
	write(line, p.Outfile)
	p.i++
}

func (p *Decoder) StartZSet(key []byte, cardinality, expiry int64) {
	p.i = 0
}

func (p *Decoder) Zadd(key []byte, score float64, member []byte) {
	line := fmt.Sprintf("db=%d %q[%d] -> {%q, score=%g}\n", p.db, key, p.i, member, score)
	write(line, p.Outfile)
	p.i++
}

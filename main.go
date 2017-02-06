// outputs a human readable diffable dump of the rdb file.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cupcake/rdb"
	"github.com/grindrllc/redis-metrics/decoder"
)

const (
	defaultOutFile = "rdbout"
	inFileUsage    = "Input file for redis dump"
	outFileUsage   = "Output file for redis dump"
)

var (
	ErrNoInputFile = errors.New("ERR: Need to specify an input file with the -input or -i flags.")
)

// -i in file, -o outfile, if not, default
func main() {
	// Deal with args
	in := flag.String("input", "", inFileUsage)
	out := flag.String("output", defaultOutFile, outFileUsage)
	flag.StringVar(in, "i", "", inFileUsage)
	flag.StringVar(out, "o", defaultOutFile, outFileUsage)
	flag.Parse()

	// If in is blank or empty, error early
	if *in == "" || in == nil {
		maybeFatal(ErrNoInputFile)
	}
	absOut, oerr := filepath.Abs(*out)
	maybeFatal(oerr)
	outFile, ferr := createFile(absOut)
	maybeFatal(ferr)
	f, err := os.Open(*in)
	maybeFatal(err)
	err = rdb.Decode(f, &decoder.Decoder{OutFile: outFile})
	maybeFatal(err)
}

func maybeFatal(err error) {
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err)
		os.Exit(1)
	}
}

func createFile(filename string) (*os.File, error) {
	file, err := os.Create(filename)
	if os.IsExist(err) {
		file, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0666)
	}
	return file, err
}

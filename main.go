// outputs a human readable diffable dump of the rdb file.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	"github.com/catherinetcai/redis-metrics/decoder"
	metrics "github.com/rcrowley/go-metrics"

	"github.com/cupcake/rdb"
)

const (
	defaultOutFile = "rdbout"
	defaultKeyFile = "keys.yml"
	inFileUsage    = "Input file for redis dump"
	outFileUsage   = "Output file for redis dump"
	keyFileUsage   = "File with keys matching for redis dump"
)

var (
	ErrNoInputFile = errors.New("ERR: Need to specify an input file with the -input or -i flags.")
)

// -i in file, -o outfile, if not, default
func main() {
	// Deal with args
	in := flag.String("input", "", inFileUsage)
	out := flag.String("output", defaultOutFile, outFileUsage)
	keys := flag.String("keys", defaultKeyFile, keyFileUsage)
	flag.StringVar(in, "i", "", inFileUsage)
	flag.StringVar(out, "o", defaultOutFile, outFileUsage)
	flag.StringVar(keys, "k", defaultKeyFile, keyFileUsage)
	flag.Parse()

	// Get in file
	if *in == "" {
		maybeFatal(ErrNoInputFile)
	}
	absOut, oerr := filepath.Abs(*out)
	maybeFatal(oerr)
	outFile, ferr := createFile(absOut)
	defer outFile.Close()
	maybeFatal(ferr)
	f, err := os.Open(*in)
	defer f.Close()
	maybeFatal(err)

	// Get keys to match
	// matchKeys, kerr := getMatchKeys(*keys)
	// maybeFatal(kerr)

	// Decode
	r := metrics.NewRegistry()
	err = rdb.Decode(f, &decoder.Decoder{OutFile: outFile, MetricsRegistry: &r})
	maybeFatal(err)
}

func maybeFatal(err error) {
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err)
		os.Exit(1)
	}
}

func createFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	return file, err
}

func getMatchKeys(filename string) (decoder.Keys, error) {
	if filename == "" {
		fmt.Printf("No match keys set!")
		return decoder.Keys{}, nil
	}
	absFile, _ := filepath.Abs(filename)
	file, err := ioutil.ReadFile(absFile)
	if err != nil {
		return nil, err
	}

	keys := decoder.MatchKeys{}
	err = yaml.Unmarshal(file, &keys)
	if err != nil {
		return nil, err
	}

	return keys.Keys, nil
}

// outputs a human readable diffable dump of the rdb file.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/catherinetcai/redis-metrics/decoder"
	metrics "github.com/rcrowley/go-metrics"

	"github.com/cupcake/rdb"
)

const (
	// Files and extension defaults
	defaultOutFile = "rdbout"
	defaultKeyFile = "keys.yml"
	RdbExt         = ".rdb"

	// Flag descriptions
	inFileUsage  = "Input file for redis dump"
	keyFileUsage = "File with keys matching for redis dump"
	outFileUsage = "Output file for redis dump"
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

	f := getInputFiles(filepath.Dir(*in))

	// Decode
	r := metrics.NewRegistry()
	err := rdb.Decode(f, &decoder.Decoder{OutFile: outFile, MetricsRegistry: &r})
	maybeFatal(err)
}

// getInputFiles grabs a directory and globs all the files
func getInputFiles(dir string) io.Reader {
	files := concatFiles(dir)
	return bytes.NewBuffer(files)
}

// concatFiles grabs all .rdb extensions and concats them
func concatFiles(dir string) []byte {
	var files []byte
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), RdbExt) {
			fmt.Printf("Concating path: %v\n", path)
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			files = append(files, b...)
		}
		return nil
	})
	return files
}

// mayabeFatal - Is it fatal? Who knows!
func maybeFatal(err error) {
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err)
		os.Exit(1)
	}
}

// createFile creates the file - if exists, it'll delete it
func createFile(filename string) (*os.File, error) {
	if _, oerr := os.Stat(filename); oerr != nil {
		fmt.Prinf("%s exists. Deleting and overwriting...\n", filename)
		os.Remove(filename)
	}
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	return file, err
}

// NOT YET IMPLEMENTED - this will be for getting the keys that we want to specifically pull out
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

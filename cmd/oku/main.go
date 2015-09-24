package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/contiamo/oku"
)

type Config struct {
	Detect   bool
	Encoding string
	Output   string
}

var config Config

func init() {
	flag.BoolVar(&config.Detect, "d", false, "detect encoding and exit")
	flag.StringVar(&config.Encoding, "f", "", "from encoding: specify encoding of input (no detection)")
	flag.StringVar(&config.Output, "o", "", "output file")
}

func main() {
	flag.Parse()

	var fileName string
	args := flag.Args()
	if len(args) > 0 {
		fileName = args[len(args)-1]
	}

	var b []byte
	var err error

	if fileName == "" {
		// can we read from stdin?
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			b, err = ioutil.ReadAll(os.Stdin)
		} else {
			flag.Usage()
			return
		}
	} else {
		b, err = ioutil.ReadFile(fileName)
	}

	if err != nil {
		panic(err)
	}

	var charset string
	if config.Encoding == "" {
		res, err := oku.DetectEncoding(b)
		if err != nil {
			// write err to stderr and exit
			panic(err)
		}

		out := fmt.Sprintf("oku detected: %s, confidence: %d%%\n", res.Charset, res.Confidence)
		if config.Detect {
			// write detection to stdout and exit
			fmt.Print(out)
			return
		} else {
			fmt.Fprint(os.Stderr, out)
		}

		charset = res.Charset
	} else {
		charset = config.Encoding
	}

	reader := bytes.NewReader(b)
	// write file to stdout
	utf8Reader, err := oku.NewUTF8ReadCloser(reader, charset)
	if err != nil {
		panic(err)
	}

	// TODO: make buffered writes
	text, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		panic(err)
	}

	if config.Output == "" {
		fmt.Printf("%s", text)
	} else {
		if err := ioutil.WriteFile(config.Output, text, os.FileMode(0664)); err != nil {
			panic(err)
		}
	}
}

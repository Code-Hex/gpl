package gpl

import (
	"bytes"
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

// Options struct for parse command line arguments
type Options struct {
	Help    bool `short:"h" long:"help"`
	Version bool `short:"v" long:"version"`
	Stdin   bool `short:"s" long:"stdin"`
	Trace   bool `long:"trace"`
}

func (opts *Options) parse(argv []string) ([]string, error) {
	p := flags.NewParser(opts, flags.PrintErrors)
	args, err := p.ParseArgs(argv)

	if err != nil {
		os.Stderr.Write(opts.usage())
		return nil, errors.Wrap(err, "invalid command line options")
	}

	return args, nil
}

func (opts Options) usage() []byte {
	buf := new(bytes.Buffer)

	fmt.Fprintf(buf, msg+
		`Usage: 
  $ gpl [options] /your/project1 /your/project2 ...
  $ [commands] | gpl [-s|--stdin]

  Options:
  -h,  --help               print usage and exit
  -v,  --version            display the version of gpl and exit
  -s,  --stdin              read target directories from stdin
  --trace                   display detail error messages
`)
	return buf.Bytes()
}

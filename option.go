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
	Help    bool `short:"h" long:"help" description:"print usage and exit"`
	Version bool `short:"v" long:"version" description:"display the version of gpl and exit"`
	Trace   bool `long:"trace" description:"display detail error messages"`
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
	buf := bytes.Buffer{}

	fmt.Fprintf(&buf, msg+
		`Usage: gpl [options] /path/to/user/project1 /path/to/user/project2
  Options:
  -h,  --help                   print usage and exit
  -v,  --version                display the version of gpl and exit
  --trace                       display detail error messages
`)
	return buf.Bytes()
}

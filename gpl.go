package gpl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
)

const (
	version = "0.0.1"
	msg     = "gpl v" + version + ", Update multiple local repositories with parallel\n"
)

// Gpl struct
type Gpl struct {
	Trace       bool
	CPU         int
	TargetPaths []string
	Args        []string
	Stdin       io.Reader
	Stdout      io.Writer
	Stderr      io.Writer
}

type ignore struct {
	err error
}

type cause interface {
	Cause() error
}

// New will return Gpl struct
func New() *Gpl {
	return &Gpl{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		CPU:    runtime.NumCPU(),
	}
}

// Run is executed gpl command
// This is main method
func (gpl *Gpl) Run() int {
	if err := gpl.Execute(); err != nil {
		if gpl.Trace {
			fmt.Fprintf(gpl.Stderr, "Error:\n%+v\n", err)
		} else {
			fmt.Fprintf(gpl.Stderr, "Error:\n  %v\n", errmsg(err))
		}
		return 1
	}
	return 0
}

// Execute will be run gpl command.
// At first, parse command line arguments.
// Next, will detecting kind of repositories to update local repositories.
// Finally, update local repositories using parallel.
func (gpl *Gpl) Execute() error {
	if err := gpl.Parse(); err != nil {
		return errmsg(err)
	}

	dict := gpl.DetectRepository()
	if len(dict) == 0 {
		return errors.Errorf("There was not found repository from path")
	}

	if err := gpl.UpdateRepository(dict); err != nil {
		return errmsg(err)
	}

	return nil
}

// Parse method will parsing for gpl command line arguments.
func (gpl *Gpl) Parse() error {
	var opts Options
	if err := gpl.parseOptions(&opts, os.Args[1:]); err != nil {
		return errors.Wrap(err, "failed to parse command line args")
	}

	if err := gpl.parseRepositoryPath(&opts); err != nil {
		return errors.Wrap(err, "failed to parse repository path")
	}

	return nil
}

// Due to local repositories.
func (gpl *Gpl) parseRepositoryPath(opts *Options) error {
	for _, path := range gpl.Args {
		if isPath, _ := govalidator.IsFilePath(path); isPath {
			gpl.TargetPaths = append(gpl.TargetPaths, path)
		}
	}

	// Try read from stdin if have not been set filepath on argv.
	if len(gpl.TargetPaths) == 0 {
		scanner := bufio.NewScanner(gpl.Stdin)
		for scanner.Scan() {
			filepath := scanner.Text()
			if isPath, _ := govalidator.IsFilePath(filepath); isPath {
				gpl.TargetPaths = append(gpl.TargetPaths, filepath)
			}
		}
	}

	removeDuplicates(&gpl.TargetPaths)

	// Finally, return usage massage if have not been set filepath from stdin.
	if len(gpl.TargetPaths) == 0 {
		gpl.Stdout.Write(opts.usage())
		return makeIgnoreErr()
	}
	return nil
}

// Due to gpl command line arguments.
func (gpl *Gpl) parseOptions(opts *Options, argv []string) error {

	o, err := opts.parse(argv)
	if err != nil {
		return errors.Wrap(err, "failed to parse command line options")
	}

	if opts.Help {
		gpl.Stdout.Write(opts.usage())
		return makeIgnoreErr()
	}

	if opts.Version {
		gpl.Stdout.Write([]byte(msg))
		return makeIgnoreErr()
	}

	if opts.Trace {
		gpl.Trace = opts.Trace
	}

	gpl.Args = o

	return nil
}

// Remove duplicates elements in slice.
func removeDuplicates(s *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *s {
		if !found[x] {
			found[x] = true
			(*s)[j] = (*s)[i]
			j++
		}
	}
	*s = (*s)[:j]
}

// errmsg method will get important message from wrapped error message
func errmsg(err error) error {
	for e := err; e != nil; {
		switch e.(type) {
		case ignore:
			return nil
		case cause:
			e = e.(cause).Cause()
		default:
			return e
		}
	}

	return nil
}

func makeIgnoreErr() ignore {
	return ignore{
		err: errors.New("this is ignore message"),
	}
}

// Error due to options: version, usage
func (i ignore) Error() string {
	return i.err.Error()
}

// Cause due to options: version, usage
func (i ignore) Cause() error {
	return i.err
}

package gpl

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	paths := []string{
		"/your/project/path1",
		"/your/project/path2",
		"/your/project/path3",
		"/your/project/path4",
	}
	args := []string{"gpl", "/your/project/path1", "/your/project/path2"}

	os.Args = args

	gpl1 := New()
	assert.Equal(t, gpl1.CPU, runtime.NumCPU(), "failed to get cpu count")

	// First, testing parse to command line arguments
	if err := gpl1.Parse(); err != nil {
		t.Errorf("failed to parse arguments: %s", err.Error())
	}

	assert.Equal(t, paths[0], gpl1.TargetPaths[0], "failed to parse argument of filepath at first")
	assert.Equal(t, paths[1], gpl1.TargetPaths[1], "failed to parse argument of filepath at second")

	// Second, testing read from stdin
	input := bytes.NewBufferString(fmt.Sprintf("%s\n%s\n", paths[2], paths[3]))
	gpl2 := &Gpl{
		CPU:    gpl1.CPU,
		Stdin:  input,
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	os.Args = os.Args[:1]

	if err := gpl2.Parse(); err != nil {
		t.Errorf("failed to parse arguments: %s", err.Error())
	}
	assert.Equal(t, paths[2], gpl2.TargetPaths[0], "failed to parse argument of filepath at third")
	assert.Equal(t, paths[3], gpl2.TargetPaths[1], "failed to parse argument of filepath at fourth")

	// Third, testing do not pass arguments
	os.Args = os.Args[:1]
	output1 := new(bytes.Buffer)
	gpl3 := &Gpl{
		CPU:    gpl1.CPU,
		Stdin:  os.Stdin,
		Stdout: output1,
		Stderr: output1,
	}
	if err := gpl3.Parse(); err != nil {
		assert.Equal(t, nil, errmsg(err), "failed to display usage")
	}

	var opts Options
	assert.Equal(t, string(opts.usage()), output1.String(), "failed to display usage")

	// Fourth, testing -h flag
	os.Args = []string{"gpl", "-h"}
	output2 := new(bytes.Buffer)
	gpl4 := &Gpl{
		CPU:    gpl1.CPU,
		Stdin:  os.Stdin,
		Stdout: output2,
		Stderr: output2,
	}

	if err := gpl4.Parse(); err != nil {
		assert.Equal(t, nil, errmsg(err), "failed to display usage")
	}

	assert.Equal(t, string(opts.usage()), output2.String(), "failed to display usage")

	// Fifth, testing -v flag
	os.Args = []string{"gpl", "-v"}
	output3 := new(bytes.Buffer)
	gpl5 := &Gpl{
		CPU:    gpl1.CPU,
		Stdin:  os.Stdin,
		Stdout: output3,
		Stderr: output3,
	}

	if err := gpl5.Parse(); err != nil {
		assert.Equal(t, nil, errmsg(err), "failed to display usage")
	}

	assert.Equal(t, msg, output3.String(), "failed to display usage")
}

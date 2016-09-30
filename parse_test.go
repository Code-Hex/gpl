package gpl

import (
	"io"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	paths := []string{"/your/project/path1", "/your/prohject/path2", "/your/prohject/path3"}
	args := []string{"gpl", "/your/project/path1", "/your/prohject/path2"}
	os.Args = append(os.Args, args...)

	gpl := New()
	assert.Equal(t, gpl.CPU, runtime.NumCPU(), "failed to get cpu count")

	// First, testing parse to command line arguments
	if err := gpl.Parse(); err != nil {
		t.Errorf("failed to parse arguments: %s", err.Error())
	}

	assert.Equal(t, paths[0], gpl.TargetPaths[0], "failed to parse argument of filepath at first")
	assert.Equal(t, paths[1], gpl.TargetPaths[1], "failed to parse argument of filepath at second")

	// Second, testing read from stdin
	io.WriteString(os.Stdin, "/your/prohject/path3\n")

	gpl.TargetPaths = []string{}
	if err := gpl.Parse(); err != nil {
		t.Errorf("failed to parse arguments: %s", err.Error())
	}
	assert.Equal(t, paths[2], gpl.TargetPaths[0], "failed to parse argument of filepath at third")

}

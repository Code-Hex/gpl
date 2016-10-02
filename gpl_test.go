package gpl

import (
	"bytes"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	CreateFakeRepositories()

	expected := 0
	os.Args = []string{
		"gpl",
		"testdata/project1",
		"testdata/project2",
		"testdata/project3",
		"testdata/project4",
		"testdata/project5",
	}
	assert.Equal(t, expected, DummyNew().Run(), "failed to gpl.Run")

	RemoveFakeRepositories()
}

func DummyNew() *Gpl {
	return &Gpl{
		Stdin:  os.Stdin,
		Stdout: new(bytes.Buffer),
		Stderr: new(bytes.Buffer),
		CPU:    runtime.NumCPU(),
		Exec: func(path, command string, args ...string) (string, error) {
			return "", nil
		},
	}
}

func CreateFakeRepositories() {
	repositories := []string{
		"testdata/project1/.git",
		"testdata/project2/.git/svn",
		"testdata/project3/.hg",
		"testdata/project4/.svn",
		"testdata/project5/_darcs",
	}
	for _, path := range repositories {
		os.MkdirAll(path, 0755)
	}
}

func RemoveFakeRepositories() {
	repositories := []string{
		"testdata/project1",
		"testdata/project2",
		"testdata/project3",
		"testdata/project4",
		"testdata/project5",
	}
	for _, path := range repositories {
		os.RemoveAll(path)
	}
}

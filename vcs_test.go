package gpl

import (
	"bytes"
	"os"
	"runtime"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestVCS(t *testing.T) {
	output := new(bytes.Buffer)
	gpl := &Gpl{
		Stdin:  os.Stdin,
		Stdout: output,
		Stderr: output,
		CPU:    runtime.NumCPU(),
	}

	CreateFakeRepositories()

	fakeRepos := FakeRepositories()
	for path := range fakeRepos {
		gpl.TargetPaths = append(gpl.TargetPaths, path)
	}

	// Testing DetectRepository()
	dict := gpl.DetectRepository()
	for path := range fakeRepos {
		expected := fakeRepos[path]
		assert.Equal(t, expected, dict[path], "failed to detect repository type: "+expected)
	}

	// Testing UpdateRepository()
	// First, when Exec has error
	commands := ReposCommands()

	gpl.Exec = func(path, command string, args ...string) (string, error) {
		// Error has a string of command need to be executed.
		return "", errors.Errorf(join(command, args, ' '))
	}

	for key := range dict {
		path, repo := key, dict[key]
		expected := commands[repo]
		err := DoUpdate[repo](*gpl, path)
		if err == nil {
			t.Errorf("failed to DoUpdate[repo](gpl, path) closure")
		}
		assert.Equal(t, expected, err.Error(), "failed to invoke command when "+repo)
	}

	// Second, when Exec has not error
	gpl.Exec = func(path, command string, args ...string) (string, error) {
		// Nothing error
		return "", nil
	}

	if err := gpl.UpdateRepository(dict); err != nil {
		t.Errorf("failed to UpdateRepository()")
	}

	RemoveFakeRepositories()
}

func FakeRepositories() map[string]string {
	return map[string]string{
		"testdata/project1": ".git",
		"testdata/project2": ".git/svn",
		"testdata/project3": ".hg",
		"testdata/project4": ".svn",
		"testdata/project5": "_darcs",
	}
}

func ReposCommands() map[string]string {
	return map[string]string{
		".git":     "git pull --ff-only",
		".git/svn": "git svn rebase",
		".hg":      "hg pull --update",
		".svn":     "svn update",
		"_darcs":   "darcs pull",
	}
}

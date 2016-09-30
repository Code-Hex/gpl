package gpl

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

var repositories = []string{".git/svn", ".git", ".svn", ".hg", "_darcs"}

var doUpdate = map[string]func(string) error{
	".git": func(path string) error {
		return do(path, "git", "pull", "--ff-only")
	},
	".svn": func(path string) error {
		return do(path, "svn", "update")
	},
	".git/svn": func(path string) error {
		return do(path, "git", "svn", "rebase")
	},
	".hg": func(path string) error {
		return do(path, "hg", "pull", "--update")
	},
	"_darcs": func(path string) error {
		return do(path, "darcs", "pull")
	},
}

// UpdateRepository method for update repositories
// after detecting the kind of repositories from filepath
func (gpl Gpl) UpdateRepository() error {
	var (
		fi  os.FileInfo
		err error
		wg  sync.WaitGroup
	)

	semaphore := make(chan bool, gpl.CPU)
	errCh := make(chan error)
	go listenCh(errCh, len(gpl.TargetPaths))

	// foreach with filepath for potential repositories
	for _, path := range gpl.TargetPaths {
		// foreach with repositories
		for _, repo := range repositories {
			fi, err = os.Stat(filepath.Join(path, repo))
			if err == nil && fi.IsDir() {
				wg.Add(1)
				go func(path, repo string) {
					defer wg.Done()
					p(semaphore)
					// Execute command for each repositories
					errCh <- doUpdate[repo](path)
					v(semaphore)
				}(path, repo)
			}
		}
	}
	wg.Wait()

	return <-errCh
}

func listenCh(errCh chan error, totalPaths int) {
	errCount := 0
	for i := 0; i < totalPaths; i++ {
		if err := <-errCh; err != nil {
			errCount++
		}
	}

	if errCount > 0 {
		var word string
		if errCount == 1 {
			word = "repository"
		} else {
			word = "repositories"
		}
		errCh <- errors.Errorf("There was an error in the %d %s update", errCount, word)
	} else {
		errCh <- nil
	}
}

// This function execute repository update commands on your target directory.
func do(path, command string, args ...string) error {
	var stderr bytes.Buffer
	cmd := exec.Command(command, args...)
	cmd.Stderr = &stderr
	cmd.Dir = path

	cmdStr := join(command, args, ' ')
	fmt.Fprintf(os.Stdout, "[%s] %s (%s)\n", color.GreenString("Update"), path, cmdStr)

	if err := cmd.Run(); err != nil {
		// when failed execute command
		fmt.Fprintf(os.Stderr, "[%s] %s\n%s\n", color.RedString("Failed"), path, stderr.String())
		return err
	}

	// when finished execute command
	fmt.Fprintf(os.Stdout, "[%s] %s (%s)\n", color.YellowString("Done"), path, cmdStr)

	return nil
}

// this function like strings.Join()
// join(command, args, ' ')
func join(command string, args []string, sep byte) string {
	limit := len(args) - 1

	var str = make([]byte, 0, 20) // ensure of capacity 20 bytes
	str = append(str, command...)
	str = append(str, sep)
	for idx, chars := range args {
		if limit <= idx {
			str = append(str, chars...)
		} else {
			str = append(str, chars...)
			str = append(str, sep)
		}
	}
	return string(str)
}

// lock
func p(semaphore chan bool) {
	semaphore <- true
}

// unlock
func v(semaphore chan bool) {
	<-semaphore
}

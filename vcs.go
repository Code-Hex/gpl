package gpl

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

var repositories = []string{".git/svn", ".git", ".svn", ".hg", "_darcs"}

var DoUpdate = map[string]func(Gpl, string) error{
	".git": func(gpl Gpl, path string) error {
		return gpl.do(path, "git", "pull", "--ff-only")
	},
	".svn": func(gpl Gpl, path string) error {
		return gpl.do(path, "svn", "update")
	},
	".git/svn": func(gpl Gpl, path string) error {
		return gpl.do(path, "git", "svn", "rebase")
	},
	".hg": func(gpl Gpl, path string) error {
		return gpl.do(path, "hg", "pull", "--update")
	},
	"_darcs": func(gpl Gpl, path string) error {
		return gpl.do(path, "darcs", "pull")
	},
}

// DetectRepository for detecting the kind of repositories from filepath
func (gpl Gpl) DetectRepository() map[string]string {
	var (
		fi  os.FileInfo
		err error
	)

	dict := make(map[string]string)

	// filepath for potential repositories
	for _, path := range gpl.TargetPaths {
		// repositories
		for _, repo := range repositories {
			fi, err = os.Stat(filepath.Join(path, repo))
			if err == nil && fi.IsDir() {
				// Register repository type to dictionary
				dict[path] = repo
				// break the repositories foreach loop
				break
			}
		}
	}

	return dict
}

// UpdateRepository for update repositories
func (gpl Gpl) UpdateRepository(dict map[string]string) error {

	errCh := make(chan error)
	// Counting semaphore
	semaphore := make(chan bool, gpl.CPU)

	for key := range dict {
		go func(path, repo string) {
			p(semaphore)
			// Execute command for each repositories
			errCh <- DoUpdate[repo](gpl, path)
			v(semaphore)
		}(key, dict[key])
	}

	err := wait(errCh, len(gpl.TargetPaths))

	close(semaphore)
	close(errCh)

	return err
}

func wait(errCh chan error, totalPaths int) error {
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
		return errors.Errorf("There was an error in the %d %s update", errCount, word)
	}

	return nil
}

// This function execute repository update commands on your target directory.
func (gpl Gpl) do(path, command string, args ...string) error {
	cmdStr := join(command, args, ' ')
	fmt.Fprintf(gpl.Stdout, "%s %s (%s)\n", color.GreenString("[Update]"), path, cmdStr)

	if reason, err := gpl.Exec(path, command, args...); err != nil {
		// when failed execute command
		fmt.Fprintf(gpl.Stderr, "%s %s\n%s\n", color.RedString("[Failed]"), path, reason)
		return err
	}

	// when finished execute command
	fmt.Fprintf(gpl.Stdout, "%s %s (%s)\n", color.YellowString("[Finish]"), path, cmdStr)

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

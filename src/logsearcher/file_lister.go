package logsearcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type invalidFilenameFormatErr string

func (e invalidFilenameFormatErr) Error() string {
	return string(e)
}

func ListFolders(dir string, ch chan string) {
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}
	dirLength := len(dir)

	var previousErr error

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			previousErr = err
			close(ch)
			return err
		}
		if info.IsDir() {
			ch <- path[dirLength:]
		}
		return nil
	})
	if previousErr == nil {
		close(ch)
	}
}

func ListFiles(dir string, ch chan string) error {
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}
	dirLength := len(dir)

	var previousErr error

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if previousErr != nil {
			return err
		}
		if err != nil {
			previousErr = err
			close(ch)
			return err
		}
		if !info.IsDir() {
			if !strings.HasPrefix(path, dir) {
				close(ch)
				return invalidFilenameFormatErr(fmt.Sprintf("Expected file %v to start with %v", path, dir))
			}
			ch <- path[dirLength:]
		}
		return nil
	})
	if previousErr == nil {
		close(ch)
	}
	return previousErr
}

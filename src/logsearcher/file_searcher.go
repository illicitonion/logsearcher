package logsearcher

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type Predicate func(string) bool

type FileEntry struct {
	Path     string
	Snippets chan string
}

func GetFileSnippets(dir string, predicate Predicate, files chan FileEntry) {
	rawFiles := make(chan string)
	go func() {
		for {
			filename, open := <-rawFiles
			if !open {
				close(files)
				break
			}

			lines := make(chan string)
			filteredLines := make(chan string)

			go Filter(lines, predicate, filteredLines)
			go ReadLines(dir+"/"+filename, lines)

			files <- FileEntry{filename, filteredLines}
		}
	}()
	ListFiles(dir, rawFiles)
}

func ReadLines(path string, ch chan string) (err error) {
	var file *os.File
	file, err = os.Open(path)
	if err != nil {
		close(ch)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		var line string
		line, err = reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return
		}
		trimmed := strings.Trim(line, "\n\r")
		if len(trimmed) > 0 {
			ch <- trimmed
		}
		if err == io.EOF {
			close(ch)
			return nil
		}
	}
	return
}

func Filter(in chan string, predicate Predicate, out chan string) {
	for {
		e, open := <-in
		if !open {
			close(out)
			break
		}
		if predicate(e) {
			out <- e
		}
	}
}

package logsearcher

import (
  "container/list"
  "fmt"
  "os"
  "path/filepath"
  "strings"
)

type invalidFilenameFormatErr string

func (e invalidFilenameFormatErr) Error() string {
  return string(e)
}

func ListFiles(dir string) (filenames []string, err error) {
  filenamesList := new(list.List)
  dirLength := len(dir)
  if !strings.HasSuffix(dir, "/") {
    dirLength++
  }
  err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }
    if !info.IsDir() {
      if !strings.HasPrefix(path, dir) {
        return invalidFilenameFormatErr(fmt.Sprintf("Expected file %v to start with %v", path, dir))
      }
      filenamesList.PushBack(path[dirLength:])
    }
    return nil
  })
  if err != nil {
    return []string{}, err
  }
  filenames = make([]string, filenamesList.Len())
  i := 0
  for e := filenamesList.Front() ; e != nil ; e = e.Next() {
    val := e
    filenames[i] = val.Value.(string)
    i++
  }
  return
}

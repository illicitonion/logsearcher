package logsearcher

import (
  "strconv"
  "testing"
)

func TestReadsLines(t *testing.T) {
  ch := make(chan string)
  go assertChanReturns(t, []string{"1","2","3","4","5"}, ch)
  err := ReadLines("../../testdata/files/2.txt", ch)
  checkNoError(t, err)
}

func TestFiltersLines(t *testing.T) {
  in := make(chan string)
  out := make(chan string)
  go func() {
    in <- "1"
    in <- "2"
    in <- "3"
    in <- "4"
    in <- "5"
    close(in)
  }()
  go assertChanReturns(t, []string{"2","4"}, out)

  isEven := func(item string) bool {
    num,_ := strconv.Atoi(item)
    return num % 2 == 0
  }

  Filter(in, isEven, out)
}

func assertFileEntryEquals(t *testing.T, expected FileEntry, actual FileEntry) {
  defer consume(expected.Snippets)
  defer consume(actual.Snippets)

  if expected.Path != actual.Path {
    t.Fatalf("Incorrect path: Expected %v but was %v", expected.Path, actual.Path)
  }
  for {
    expectedLine,expectedOk := <- expected.Snippets
    actualLine,actualOk := <- actual.Snippets
    if !expectedOk && !actualOk {
      break
    } else if expectedOk != actualOk {
      t.Errorf("Different numbers of entries in expected and actual snippets")
    } else if expectedOk && expectedLine != actualLine {
      t.Errorf("Expected snippet %v but was %v", expectedLine, actualLine)
    }
  }
}

func consumeFiles(ch chan FileEntry) {
  for {
    _,more := <- ch
    if !more {
      break
    }
  }
}

func stubChannel(values []string) (ch chan string) {
  ch = make(chan string)
  go func() {
    for i := range values {
      ch <- values[i]
    }
    close(ch)
  }()
  return
} 

func TestGetsFileSnippets(t *testing.T) {
  files := make(chan FileEntry)

  go GetFileSnippets("../../testdata/files", func(_ string) bool {
    return true
  }, files)

  expectedOne := FileEntry{"1.txt", stubChannel([]string{"1"})}
  expectedTwo := FileEntry{"2.txt", stubChannel([]string{"1","2","3","4","5"})}

  var file FileEntry
  var ok bool

  file,ok = <- files
  if !ok {
    t.Fatal("Expected file")
  }
  assertFileEntryEquals(t, expectedOne, file)

  file,ok = <- files
  if !ok {
    t.Fatal("Expected file")
  }
  assertFileEntryEquals(t, expectedTwo, file)

  file,ok = <- files
  if ok {
    t.Fatalf("Expected no more files but found %v", file.Path)
  }
}

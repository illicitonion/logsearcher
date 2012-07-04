package logsearcher

import (
  "testing"
)

func assertChanReturns(t *testing.T, expected []string, ch chan string) {
  for i := range expected {
    actual, open := <- ch
    if actual != expected[i] {
      t.Errorf("Mismatch at index %v: Expected %v but got: %v (already matched: %v)", i, expected[i], actual, expected[:i])
    }
    if !open && i < len(expected) {
      t.Fatalf("Expected %v values %v but only got %v %v", len(expected), expected, i, expected[:i])
    }
  }
  _,open := <- ch
  if open {
    t.Errorf("Channel still open after all matches")
  }
}
      
func checkNoError(t *testing.T, err error) {
  if err != nil {
    t.Fatalf("Expected error to be nil but was %v", err)
  }
}

func TestListFiles(t *testing.T) {
  expected_files := []string{"1.txt", "2.txt"}

  ch := make(chan string)
  go assertChanReturns(t, expected_files, ch)
  err := ListFiles("../../testdata/files", ch)
  checkNoError(t, err)
}

func TestListFilesIfDirHadTrailingSlash(t *testing.T) {
  expected_files := []string{"1.txt", "2.txt"}

  ch := make(chan string)
  go assertChanReturns(t, expected_files, ch)
  err := ListFiles("../../testdata/files/", ch)
  checkNoError(t, err)
}

func TestIncludesSubdirs(t *testing.T) {
  expected_files := []string{"1.txt", "child/2.txt", "child/grandchild/3.txt"}

  ch := make(chan string)
  go assertChanReturns(t, expected_files, ch)
  err := ListFiles("../../testdata/withchildren", ch)
  checkNoError(t, err)
}

func TestInvalidDirs(t *testing.T) {
  expected_files := []string{}

  invalid_dirs := []string{"!", "doesnotexist"}
  for i := range invalid_dirs {
    ch := make(chan string)
    go assertChanReturns(t, expected_files, ch)
    err := ListFiles(invalid_dirs[i], ch)

    if err == nil {
      t.Fatalf("Expected error but was nil")
    }
  }
}

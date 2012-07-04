package logsearcher

import (
  "testing"
)

func assertSlicesEqual(t *testing.T, expected_slice, actual_slice []string) {
  expected_len := len(expected_slice)
  actual_len := len(actual_slice)
  if expected_len != actual_len {
    t.Fatalf("Expected %v files %v but got: %v %v", expected_len, expected_slice, actual_len, actual_slice)
  }
  
  for i := range expected_slice {
    expected := expected_slice[i]
    actual := actual_slice[i]
    if expected != actual {
      t.Errorf("Mismatch at index %v: Expected %v but got: %v", i, expected, actual)
    }
  }
}

func checkNoError(t *testing.T, err error) {
  if err != nil {
    t.Fatalf("Expected error to be nil but was %v", err)
  }
}

func TestListFiles(t *testing.T) {
  expected_files := []string{"1.txt", "2.txt"}

  found_files, err := ListFiles("../../testdata/files")
  checkNoError(t, err)
  assertSlicesEqual(t, expected_files, found_files)
}

func TestListFilesIfDirHadTrailingSlash(t *testing.T) {
  expected_files := []string{"1.txt", "2.txt"}

  found_files, err := ListFiles("../../testdata/files/")
  checkNoError(t, err)
  assertSlicesEqual(t, expected_files, found_files)
}

func TestIncludesSubdirs(t *testing.T) {
  expected_files := []string{"1.txt", "child/2.txt", "child/grandchild/3.txt"}

  found_files, err := ListFiles("../../testdata/withchildren")
  checkNoError(t, err)
  assertSlicesEqual(t, expected_files, found_files)
}

func TestInvalidDirs(t *testing.T) {
  expected_files := []string{}

  invalid_dirs := []string{"!", "doesnotexist"}
  for i := range invalid_dirs {
    found_files, err := ListFiles(invalid_dirs[i])
    if err == nil {
      t.Fatalf("Expected error but was nil")
    }
    assertSlicesEqual(t, expected_files, found_files)
  }
}

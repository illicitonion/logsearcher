package logsearcher

import (
	"testing"
)

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

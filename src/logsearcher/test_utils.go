package logsearcher

import(
  "testing"
)

func assertChanReturns(t *testing.T, expected []string, ch chan string) {
  for i := range expected {
    actual, open := <- ch
    if actual != expected[i] {
      t.Fatalf("Mismatch at index %v: Expected %v but got: %v (already matched: %v)", i, expected[i], actual, expected[:i])
    }
    if !open && i < len(expected) {
      t.Fatalf("Expected %v values %v but only got %v %v", len(expected), expected, i, expected[:i])
    }
  }
  val,open := <- ch
  if open {
    consume(ch)
    close(ch)
    t.Fatalf("Channel still open after all matches, value returned: [%v]", val)
  }
}

func consume(ch chan string) {
  for {
    _,more := <- ch
    if !more {
      break
    }
  }
}

func checkNoError(t *testing.T, err error) {
  if err != nil {
    t.Fatalf("Expected error to be nil but was %v", err)
  }
}

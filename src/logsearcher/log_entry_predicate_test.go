package logsearcher

import (
	"testing"
)

func assertMatches(t *testing.T, predicate Predicate, s string, expected bool) {
	if predicate(s) != expected {
		t.Errorf("Expected predicate to be %v for [%v]\n", expected, s)
	}
}

func TestNameAndMessage(t *testing.T) {
	predicate := LogEntryPredicate("dawagner", "message")
	truths := []string{
		"[23:15:33] dawagner: some message",
		"[00:00:00] dawagner: Some other some message",
	}
	for i := range truths {
		assertMatches(t, predicate, truths[i], true)
	}
	falsehoods := []string{
		"[23:15:33] simonstewart: some message",
		"[00:00:00] dawagner: Something else",
	}
	for i := range falsehoods {
		assertMatches(t, predicate, falsehoods[i], false)
	}
}

func TestJustName(t *testing.T) {
	predicate := LogEntryPredicate("dawagner", "")
	truths := []string{
		"[23:15:33] dawagner: something",
		"[00:00:00] dawagner:",
	}
	for i := range truths {
		assertMatches(t, predicate, truths[i], true)
	}
	falsehoods := []string{
		"[23:15:33] simonstewart: something",
		"[00:00:00] PenguinOfDeath:",
	}
	for i := range falsehoods {
		assertMatches(t, predicate, falsehoods[i], false)
	}
}

func TestJustMessage(t *testing.T) {
	predicate := LogEntryPredicate("", "some message")
	truths := []string{
		"[23:15:33] dawagner: some message with trailing",
		"[23:15:34] dawagner: some message",
		"[00:00:00] simonstewart: leading some message",
	}
	for i := range truths {
		assertMatches(t, predicate, truths[i], true)
	}
	falsehoods := []string{
		"[23:15:33] dawagner: somemessage with trailing",
		"[00:00:00] simonstewart: leading somemessage",
	}
	for i := range falsehoods {
		assertMatches(t, predicate, falsehoods[i], false)
	}
}

func TestCaseInsensitivity(t *testing.T) {
	assertMatches(t, LogEntryPredicate("daWagner", "some Message"), "[23:15:33] Dawagner: Some message with trailing", true)
}

package logsearcher

import(
  "strings"
)

func LogEntryPredicate(usernameMixedCase string, messageMixedCase string) Predicate {
  username := strings.ToLower(usernameMixedCase)
  message := strings.ToLower(messageMixedCase)

  return func(strMixedCase string) bool {
    str := strings.ToLower(strMixedCase)
    fields := strings.Fields(str)

    return len(fields) > 1 &&
      (len(username) == 0 || username + ":" == fields[1]) &&
      strings.Contains(str, message)
  }
}

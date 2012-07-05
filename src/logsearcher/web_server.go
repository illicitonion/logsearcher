package logsearcher

import (
  "fmt"
  "net/http"
  "net/url"
)

type Escaper func(string) string

type LogSearchServer struct{
  LogPath string
  LogHttpRoot string
  Interface string
  Escape Escaper
}

const header = `<!DOCTYPE html>
<html>
<head><title>Selenium log search</title></head>
<body>`

const footer = `</form>
</body>
</html>`

const form string = header +
`<p>Both fields are optional, but it would be really nice if you supplied at least one (otherwise a 20MB+ page will be served).</p>
<p>Search is case-insensitive.</p>
<p>Username is exact-equality check.  Message is contains check.</p>
<form action="/" method="GET">
  <div>
    <label for="username">Username</label><input type="text" id="username" name="username" />
    <label for="message">Message</label><input type="text" id="message" name="message" size="50" />
    <input type="submit" value="Search" />
  </div>` + footer

func (server LogSearchServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  queryValues,_ := url.ParseQuery(r.URL.RawQuery)
  usernames,uok := queryValues["username"]
  messages,mok := queryValues["message"]
  if !uok || !mok {
    fmt.Fprint(w, form)
    return
  }

  username := usernames[0]
  message := messages[0]

  files := make(chan FileEntry)
  go func() {
    fmt.Fprint(w, header)
    for {
      file,ok := <- files
      if !ok {
        break
      }
      printed := false
      for {
        snippet,err := <- file.Snippets
        if !err {
          fmt.Fprint(w, "</ul>")
          break
        }
        if !printed {
          fmt.Fprintf(w, "<h3><a href='%v%v'>%v</a></h3><ul>", server.LogHttpRoot, file.Path, file.Path)
          printed = true
        }
        fmt.Fprintf(w, "<li>%v</li>", server.Escape(snippet))
      }
    }
    fmt.Fprint(w, footer)
  }()

  GetFileSnippets(server.LogPath, LogEntryPredicate(username, message), files)
}


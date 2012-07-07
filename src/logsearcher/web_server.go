package logsearcher

import (
  "fmt"
  "net/http"
  "net/url"
  "runtime"
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

const formBody string = `<p>Both fields are optional, but it would be really nice if you supplied at least one (otherwise a 20MB+ page will be served).</p>
<p>Search is case-insensitive.</p>
<p>Username is exact-equality check.  Message is contains check.</p>
<form action="/" method="GET">
  <div>
    <label for="username">Username</label><input type="text" id="username" name="username" />
    <label for="message">Message</label><input type="text" id="message" name="message" size="50" />
    <input type="submit" value="Search" />
  </div>`

func (server LogSearchServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  w.Header()["Connection"] = []string{"close"}
  _,printerr := fmt.Fprint(w, header)
  if printerr != nil {
    return
  }

  closed := false
  queryValues,_ := url.ParseQuery(r.URL.RawQuery)
  usernames,uok := queryValues["username"]
  messages,mok := queryValues["message"]
  if !uok || !mok {
    _,printerr := fmt.Fprint(w, formBody)
    if printerr != nil {
      return
    }
    folders := make(chan string)
    go func() {
      for {
        folder,ok := <- folders
        if !ok {
          _,printerr := fmt.Fprint(w, footer)
          closed = true
          if printerr != nil {
            return
          }
          break
        }
        _,printerr := fmt.Fprintf(w, "<p><a href='%v%v'>%v</a></p>", server.LogHttpRoot, folder, folder)
        if printerr != nil {
          return
        }
      }
    }()
    ListFolders(server.LogPath, folders)
    for !closed {
      runtime.Gosched()
    }
    return
  }

  username := usernames[0]
  message := messages[0]

  files := make(chan FileEntry)
  go func() {
    for {
      file,ok := <- files
      if !ok {
        break
      }
      printed := false
      for {
        snippet,err := <- file.Snippets
        if !err {
          _,printerr := fmt.Fprint(w, "</ul>")
          if printerr != nil {
            return
          }
          break
        }
        if !printed {
          _,printerr := fmt.Fprintf(w, "<h3><a href='%v%v'>%v</a></h3><ul>", server.LogHttpRoot, file.Path, file.Path)
          if printerr != nil {
            return;
          }
          printed = true
        }
        _,printerr := fmt.Fprintf(w, "<li>%v</li>", server.Escape(snippet))
        if printerr != nil {
          return
        }
      }
    }
    fmt.Fprint(w, footer)
    closed = true
  }()

  GetFileSnippets(server.LogPath, LogEntryPredicate(username, message), files)
  for !closed {
    runtime.Gosched()
  }
}


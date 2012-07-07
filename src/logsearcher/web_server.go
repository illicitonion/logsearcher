package logsearcher

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime"
)

type Escaper func(string) string

type LogSearchServer struct {
	LogPath     string
	LogHttpRoot string
	Interface   string
	Escape      Escaper
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

func (server LogSearchServer) writeFolders(w http.ResponseWriter, folders chan string, closed *bool) {
	for {
		folder, ok := <-folders
		if !ok {
			*closed = true
			return
		}
		fmt.Fprintf(w, "<p><a href='%v%v'>%v</a></p>", server.LogHttpRoot, folder, folder)
	}
}

func (server LogSearchServer) writeContent(w http.ResponseWriter, files chan FileEntry, closed *bool) {
	for {
		file, ok := <-files
		if !ok {
			break
		}
		printedAny := false
		for {
			snippet, ok := <-file.Snippets
			if !ok {
				if printedAny {
					fmt.Fprint(w, "</ul>")
				}
				break
			}
			if !printedAny {
				fmt.Fprintf(w, "<h3><a href='%v%v'>%v</a></h3><ul>", server.LogHttpRoot, file.Path, file.Path)
				printedAny = true
			}
			fmt.Fprintf(w, "<li>%v</li>", server.Escape(snippet))
		}
	}
	*closed = true
}

func (server LogSearchServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header()["Connection"] = []string{"close"}
	fmt.Fprint(w, header)

	closed := false
	queryValues, _ := url.ParseQuery(r.URL.RawQuery)
	usernames, uok := queryValues["username"]
	messages, mok := queryValues["message"]
	if !uok || !mok {
		fmt.Fprint(w, formBody)
		folders := make(chan string)
		go server.writeFolders(w, folders, &closed)
		ListFolders(server.LogPath, folders)
	} else {
		files := make(chan FileEntry)
		go server.writeContent(w, files, &closed)
		GetFileSnippets(server.LogPath, LogEntryPredicate(usernames[0], messages[0]), files)
	}

	for !closed {
		runtime.Gosched()
	}
	fmt.Fprint(w, footer)
}

package main

import(
  "flag"
  "fmt"
  "net/http"
  "os"
  "logsearcher"
  "strings"
)

func main() {
  logPath := flag.String("log_path", "/doesnotexist", "path to logs")
  httpPrefix := flag.String("http_prefix", "http://illicitonion.com/selogs/selenium/", "Prefix for links")
  iface := flag.String("if", "localhost:8983", "interface to bind to, e.g. localhost:8983")
  flag.Parse()

  file,err := os.Open(*logPath)
  if err != nil {
    fmt.Printf("File %v did not exist\n", *logPath)
    os.Exit(1)
  }
  file.Close()

  escapingReplacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "\"", "&quot;", "\\", "&#x27;", "/", "&#x2F;")

  escaper := func (str string) string {
    return escapingReplacer.Replace(str)
  }

  var server logsearcher.LogSearchServer = logsearcher.LogSearchServer{*logPath, *httpPrefix, *iface, escaper}
  http.ListenAndServe(server.Interface, server)
}

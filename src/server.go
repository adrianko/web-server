package main

import (
    "io"
    "io/ioutil"
    "log"
    "net/http"
	"strings"
)

var configuration map[string]string = make(map[string]string)

var handlers map[string]func(http.ResponseWriter, *http.Request) = make(map[string]func(http.ResponseWriter,
    *http.Request))

type Handler struct {}

func (*Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    
}

func main() {
    data, err := ioutil.ReadFile("../conf/config")
    check(err)
    parse_config(string(data))
    log.Printf("Running server %s:%s\n", configuration["interface"], configuration["port"])
    server := http.Server{Addr: configuration["interface"] + ":" + configuration["port"], Handler: &Handler{}}
    server.ListenAndServe()
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func parse_config(config string) {
	for _, c := range strings.Split(config, "\n") {
		line := strings.Split(c, "=")
		configuration[strings.Trim(line[0], " ")] = strings.Trim(line[1], " ")
	}
}

func log_request(r *http.Request, w http.ResponseWriter, status int) {
    log.Printf("%d %s: %s", status, r.Method, r.URL.String())
}

func send(r *http.Request, w http.ResponseWriter, status int, contentType string, content string) {
    set_content_type(w, contentType)
    w.WriteHeader(status)
    log_request(r, w, status)
    io.WriteString(w, content)
}

func set_content_type(w http.ResponseWriter, contentType string) {
    var cType string

    switch contentType {
    case "json":
        cType = "application/json"
        break
    case "html":
        cType = "text/html"
        break
    case "plain":
    default:
        cType = "text/plain"
        break
    }

    w.Header().Set("Content-Type", cType)
}

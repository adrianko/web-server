package main

import (
    "net/http"
    "io"
    "strconv"
    "log"
)

const PORT int = 8000

var handlers map[string]func(http.ResponseWriter, *http.Request) = make(map[string]func(http.ResponseWriter, 
    *http.Request))

type Handler struct {}

func (*Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if handle, ok := handlers[r.URL.String()]; ok {
        handle(w, r)
        return
    }
    
    error(w, r)
}

func main() {
    log.Printf("Running server 127.0.0.1:%d\n", PORT)
    server := http.Server{Addr: ":" + strconv.Itoa(PORT), Handler: &Handler{}}
    
    handlers["/"] = hello
    handlers["/html"] = helloHTML
    server.ListenAndServe()
}

func logRequest(r *http.Request) {
    log.Printf("%s: %s\n", r.Method, r.URL.String())
}

//handlers
func hello(w http.ResponseWriter, r *http.Request) {
    logRequest(r)
    io.WriteString(w, "Hello world")
}

func helloHTML(w http.ResponseWriter, r *http.Request) {
    logRequest(r)
    io.WriteString(w, "<h1>Hello world</h1>")
}

func error(w http.ResponseWriter, r *http.Request) {
    logRequest(r)
    io.WriteString(w, "<h1>Not found</h1>")
}
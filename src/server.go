package main

import (
    "fmt"
    "net/http"
    "io"
    "log"
)

const PORT = 8000

var handlers map[string]func(http.ResponseWriter, *http.Request) = make(map[string]func(http.ResponseWriter, 
    *http.Request))

type Handler struct {}

func (*Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if handle, ok := handlers[r.URL.String()]; ok {
        handle(w, r)
        return
    }
}

func main() {
    log.Printf("Running server 127.0.0.1:%d", PORT)
    server := http.Server{Addr: ":" + string(PORT), Handler: &Handler{}}
    
    handlers["/"] = hello
    handlers["/html"] = helloHTML
    server.ListenAndServe()
}

func logRequest(r *http.Request) {
    fmt.Printf("%s: %s\n", r.Method, r.URL.String())
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
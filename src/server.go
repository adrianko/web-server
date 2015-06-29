package main

import (
    "fmt"
    "net/http"
    "io"
)

var mux map[string]func(http.ResponseWriter, *http.Request) = make(map[string]func(http.ResponseWriter, *http.Request))

type Handler struct {}

func (*Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if handle, ok := mux[r.URL.String()]; ok {
        handle(w, r)
        return
    }
    
    io.WriteString(w, "Server: " + r.URL.String())
}

func main() {
    fmt.Println("Running server 127.0.0.1:8000")
    server := http.Server{Addr: ":8000", Handler: &Handler{}}
    
    mux["/"] = hello
    mux["/html"] = helloHTML
    server.ListenAndServe()
}

func logRequest(r *http.Request) {
    fmt.Printf("%s: %s", r.Method, r.URL.String())
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
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
    
    mux = make(map[string]func(http.ResponseWriter, *http.Request))
    mux["/"] = hello
    server.ListenAndServe()
}

//handlers
func hello(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "Hello world")
}
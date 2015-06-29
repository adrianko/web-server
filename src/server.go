package main

import (
    "fmt"
    "net/http"
    "io"
)

var mux map[string]func(http.ResponseWriter, *http.Request)

type Handler struct {}

func (*Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if h, ok := mux[r.URL.String()]; ok {
        h(w, r)
        return
    }
    
    io.WriteString(w, "Server: " + r.URL.String())
}

func main() {
    fmt.Println("Running server 127.0.0.1:8000")
    server := http.Server{Addr: ":8000", Handler: &Handler{}}
    
    mux = make(map[string]func(http.ResponseWriter, *http.Request))
    mux["/"] = nil
    server.ListenAndServe()
}
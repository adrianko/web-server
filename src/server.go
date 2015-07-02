package main

import (
    "net/http"
    "io"
    "strconv"
    "log"
)

const NIC string = "0.0.0.0"
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
    log.Printf("Running server %s:%d\n", NIC, PORT)
    server := http.Server{Addr: NIC + ":" + strconv.Itoa(PORT), Handler: &Handler{}}
    
    handlers["/"] = hello
    handlers["/html"] = helloHTML
    handlers["/json"] = helloJSON
    server.ListenAndServe()
}

func logRequest(r *http.Request, w http.ResponseWriter, status int) {
    log.Printf("%d %s: %s", status, r.Method, r.URL.String())
}

func sendOK(r *http.Request, w http.ResponseWriter, content string) {
    send(r, w, http.StatusOK, content)
}

func send(r *http.Request, w http.ResponseWriter, status int, content string) {
    logRequest(r, w, status)
    io.WriteString(w, content)
}

//handlers
func hello(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    sendOK(r, w, "Hello world")
}

func helloHTML(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    sendOK(r, w, "<h1>Hello world</h1>")
}

func helloJSON(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    sendOK(r, w, "{\"hello\": \"world\"}")
}

func error(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    send(r, w, http.StatusNotFound, "<h1>Not found</h1>")
}
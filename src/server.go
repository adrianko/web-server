package main

import (
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "strings"
    "fmt"
    "os"
)

var configuration map[string]string = make(map[string]string)

type Handler struct {}

func (*Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    files, err := ioutil.ReadDir(configuration["root"])
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(r.URL.String())
    
    for _, file := range files {
        fmt.Println(file.Name())
    }
    
    w.Header().Set("Content-Type", "")
    send(r, w, 404, "Hello, world")
}

func parse_config(config string) {
    for _, c := range strings.Split(config, "\n") {
        line := strings.Split(c, "=")
        configuration[strings.Trim(line[0], " ")] = strings.Trim(line[1], " ")
    }
    
    if root, ok := configuration["root"]; ok {
        if _, err := os.Stat(root); os.IsNotExist(err) {
            log.Fatal("Root path does not exist.")
        }
    } else {
        configuration["root"] = "/var/www" // default value
    }
}

func log_request(r *http.Request, w http.ResponseWriter, status int) {
    log.Printf("%d %s: %s", status, r.Method, r.URL.String())
}

func send(r *http.Request, w http.ResponseWriter, status int, content string) {
    w.WriteHeader(status)
    log_request(r, w, status)
    io.WriteString(w, content)
}

func start_server() {
    server := http.Server{Addr: configuration["interface"] + ":" + configuration["port"], Handler: &Handler{}}
    log.Printf("Running server %s:%s\n", configuration["interface"], configuration["port"])
    log.Fatal(server.ListenAndServe())
}

func main() {
    data, err := ioutil.ReadFile("../conf/config")
    
    if err != nil {
        log.Fatal(err)
    }
    
    parse_config(string(data))
    start_server()
}
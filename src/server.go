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
    if root, ok := configuration["root"]; ok {
        files, err := ioutil.ReadDir(root)
        check(err)
        fmt.Println(r.URL.String())
        
        for _, file := range files {
            fmt.Println(file.Name())
        }
    }
    
    w.Header().Set("Content-Type", "")
    send(r, w, 404, "Hello, world")
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
    
    if root, ok := configuration["root"]; ok {
        if _, err := os.Stat(root); os.IsNotExist(err) {
            log.Printf("Root path does not exist.")
            os.Exit(1)
        }
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

func main() {
    data, err := ioutil.ReadFile("../conf/config")
    check(err)
    parse_config(string(data))
    log.Printf("Running server %s:%s\n", configuration["interface"], configuration["port"])
    server := http.Server{Addr: configuration["interface"] + ":" + configuration["port"], Handler: &Handler{}}
    server.ListenAndServe()
}
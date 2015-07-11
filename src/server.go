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

var configuration_default = map[string]string{
    "root": "/var/www",
    "port": "80",
    "interface": "0.0.0.0",
}

var config_file string = "/etc/maester-http"

var configuration map[string]string = make(map[string]string)

type Handler struct {}

func (*Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    files, err := ioutil.ReadDir(configuration["root"])
    check(err)
    fmt.Println(r.URL.String())
    
    for _, file := range files {
        fmt.Println(file.Name())
    }
    
    send(r, w, 404, "", "Hello, world")
}

func parse_config(config string) {
    for _, c := range strings.Split(config, "\n") {
        line := strings.Split(c, "=")
        configuration[strings.Trim(line[0], " ")] = strings.Trim(line[1], " ")
    }
    
    for p, _ := range configuration_default {
        check_default(p)
    }
    
    if _, err := os.Stat(configuration["root"]); os.IsNotExist(err) {
        log.Fatal("Root path does not exist.")
    }
}

func check_default(property string) {
    if _, ok := configuration[property]; !ok {
        configuration[property] = configuration_default[property]
    }
}

func log_request(r *http.Request, w http.ResponseWriter, status int) {
    log.Printf("%d %s: %s", status, r.Method, r.URL.String())
}

func send(r *http.Request, w http.ResponseWriter, status int, content_type string, content string) {
    w.Header().Set("Content-Type", content_type)
    w.WriteHeader(status)
    log_request(r, w, status)
    io.WriteString(w, content)
}

func start_server() {
    server := http.Server{Addr: configuration["interface"] + ":" + configuration["port"], Handler: &Handler{}}
    log.Printf("Running server %s:%s\n", configuration["interface"], configuration["port"])
    log.Fatal(server.ListenAndServe())
}

func check(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    if len(os.Args) > 1 {
        config_file = os.Args[1]
    }
    
    data, err := ioutil.ReadFile(config_file)
    check(err)
    parse_config(string(data))
    start_server()
}
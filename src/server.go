package main

import (
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "strings"
    "os"
    "mime"
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
    static_file := configuration["root"] + r.URL.String()
    
    if strings.HasSuffix(static_file, "/") {
        static_file += "index.html"
    }
    
    if _, err := os.Stat(static_file); os.IsNotExist(err) {
        //file doesn't exist
        send(r, w, 404, "text/plain", "404: Not found")
        return
    }
    
    data, err := ioutil.ReadFile(static_file)
    
    if err != nil {
        log.Printf("Could not read file: " + static_file)
        return
    }
    
    fileName := strings.Split(static_file, ".")
    ext := "." + fileName[len(fileName) - 1]
    send(r, w, 200, mime.TypeByExtension(ext), string(data))
}

func load_config() {
    data, err := ioutil.ReadFile(config_file)
    check(err)
    parse_config(string(data))
}

func parse_config(config string) {
    for _, c := range strings.Split(config, "\n") {
        if strings.HasPrefix(c, "#") || strings.HasPrefix(c, ";") {
            continue
        }
        
        if strings.Contains(c, "=") {
            line := strings.Split(c, "=")
            configuration[strings.TrimSpace(line[0])] = strings.TrimSpace(line[1])
        }
    }
    
    for p, _ := range configuration_default {
        if _, ok := configuration[p]; !ok {
            configuration[p] = configuration_default[p]
        }
    }
    
    validate_config()
}

func validate_config() {
    // validate web root
    if _, err := os.Stat(configuration["root"]); os.IsNotExist(err) {
        log.Fatal("Root path does not exist.")
    }
}

func send(r *http.Request, w http.ResponseWriter, status int, content_type string, content string) {
    w.Header().Set("Content-Type", content_type)
    w.WriteHeader(status)
    log.Printf("%d %s: %s", status, r.Method, r.URL.String())
    io.WriteString(w, content)
}

func start_server() {
    server := http.Server{Addr: configuration["interface"] + ":" + configuration["port"], Handler: &Handler{}}
    log.Printf("Running server %s:%s\n", configuration["interface"], configuration["port"])
    check(server.ListenAndServe())
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
    
    load_config()
    start_server()
}
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

var config_file string = "/etc/maester-http"

var configuration map[string]string = map[string]string{
    "root": "/var/www",
    "port": "80",
    "interface": "0.0.0.0",
    "index": "index.html",
}

var index_files []string = []string{}

type Handler struct {}

func (*Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if valid, path := valid_path(configuration["root"] + r.URL.String()); valid {
        send(r, w, path)
    } else {
        send_response(r, w, 404, "text/plain", "404: Not found")
    }
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
    
    validate_config()
}

func validate_config() {
    // validate web root
    if _, err := os.Stat(configuration["root"]); os.IsNotExist(err) {
        log.Fatal("Root path does not exist.")
    }
    
    for _, in := range strings.Split(configuration["index"], " ") {
        index_files = append(index_files, strings.TrimSpace(in))
    }
}

func valid_path(path string) (bool, string) {
    if valid_file(path) {
        return true, path
    }
    
    if strings.HasSuffix(path, "/") {
        return valid_index(path)
    }
    
    if !strings.HasSuffix(path, "/") && !strings.HasSuffix(path, "index.html") {
        return valid_path(path + "/")
    }
    
    return false, path
}

func valid_index(path string) (bool, string) {
    if strings.HasSuffix(path, "/") {
        for _, in := range index_files {
            if valid_file(path + in) {
                return true, path + in
            }
        }
        
        return false, path
    }
        
    return valid_index(path + "/")
}

func valid_file(path string) bool {
    if info, err := os.Stat(path); err == nil && !info.IsDir() {
        return true
    }
    
    return false
}

func send(r *http.Request, w http.ResponseWriter, static_file string) {
    data, err := ioutil.ReadFile(static_file)
    
    if err != nil {
        log.Printf("Could not read file: " + static_file)
        return
    }
    
    fileName := strings.Split(static_file, ".")
    ext := "." + fileName[len(fileName) - 1]
    send_response(r, w, 200, mime.TypeByExtension(ext), string(data))
}

func send_response(r *http.Request, w http.ResponseWriter, status int, content_type string, content string) {
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
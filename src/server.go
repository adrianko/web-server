package main

import (
    "io"
    "io/ioutil"
    "log"
    "mime"
    "net/http"
    "os"
    "strings"
)
/**
 * TODO File index on directory
 * TODO Consolidate send function
 */
var config_file string = "/etc/maester-http"

var configuration map[string]string = map[string]string{
    "root":      "/var/www",
    "port":      "80",
    "interface": "0.0.0.0",
    "index":     "index.html",
    "error404":  "/error/404.html",
}

var index_files []string = []string{}

type Handler struct{}

func (*Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if valid, path := valid_path(configuration["root"] + r.URL.String()); valid {
        send(r, w, path)
    } else {
        send_not_found(r, w)
    }
}

func read_args() {
    if len(os.Args) > 1 {
        config_file = os.Args[1]
    }
}

func load_config() {
    data, err := ioutil.ReadFile(config_file)
    check(err)
    parse_config(string(data))
}

func parse_config(config string) {
    for _, conf_line := range strings.Split(config, "\n") {
        if strings.HasPrefix(conf_line, "#") || strings.HasPrefix(conf_line, ";") {
            continue
        }

        if strings.Contains(conf_line, " ") {
            property := strings.TrimSpace(string(conf_line[0:strings.Index(conf_line, " ")]))
            value := strings.TrimSpace(string(conf_line[strings.Index(conf_line, " "):len(conf_line)]))
            configuration[property] = value
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
    filePath := strings.Split(path, "/")

    if strings.HasPrefix(filePath[len(filePath) - 1], ".") {
        return false;
    }
    
    if info, err := os.Stat(path); err == nil && !info.IsDir() {
        return true
    }

    return false
}

func get_extension(file string) string {
    if strings.Contains(file, "/") {
        file = get_file(file)
    }
    
    if !valid_file(file) {
        return ""
    }
    
    fileName := strings.Split(filePath[len(filePath) - 1], ".")

    return "." + fileName[len(fileName) - 1]
}

func get_file(path string) string {
    
}

func send(r *http.Request, w http.ResponseWriter, static_file string) {
    data, err := ioutil.ReadFile(static_file)

    if err != nil {
        log.Printf("Could not read file: " + static_file)
        send_locked(r, w)
        return
    }

    filePath := strings.Split(static_file, "/")

    if !strings.HasPrefix(filePath[len(filePath) - 1], ".") {
        fileName := strings.Split(filePath[len(filePath) - 1], ".")
        ext := "." + fileName[len(fileName) - 1]
        send_response(r, w, 200, mime.TypeByExtension(ext), string(data))
    } else {
        send_not_found(r, w)
    }
}

func send_not_found(r *http.Request, w http.ResponseWriter) {
    file404 := configuration["root"] + configuration["error404"]
    
    if valid_file(file404) {
        data, err := ioutil.ReadFile(file404)
            
        if err != nil {
            log.Printf("Could not read file: " + file404)
            send_locked(r, w)
            return
        }

        filePath := strings.Split(file404, "/")
        fileName := strings.Split(filePath[len(filePath) - 1], ".")
        ext := "." + fileName[len(fileName) - 1]
        send_response(r, w, 404, mime.TypeByExtension(ext), string(data))
    } else {
        send_response(r, w, 404, "text/plain", "404: Not found")
    }
}

func send_locked(r *http.Request, w http.ResponseWriter) {
    send_response(r, w, 423, "", "")
}

func send_response(r *http.Request, w http.ResponseWriter, status int, content_type string, content string) {
    if content_type != "" {
        w.Header().Set("Content-Type", content_type)
    }

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
    read_args()
    load_config()
    start_server()
}

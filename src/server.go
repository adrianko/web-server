package main

import (
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strings"
    "gopkg.in/fsnotify.v1"
)
/**
 * TODO File index on directory
 * TODO File system watcher to invalidate cache for updated files
 */

const SERVER_NAME string = "Maester"
const VERSION string = "0.3"

var config_file string = "/etc/maester-http"

var configuration map[string]string = map[string]string{
    "root":      "/var/www",
    "port":      "80",
    "interface": "0.0.0.0",
    "index":     "index.html",
    "error404":  "error/error404.html",
    "showfiles": "off",
}

var index_files []string = []string{}

var file_cache map[string]CacheFile = make(map[string]CacheFile)

var file_watcher *fsnotify.Watcher

type CacheFile struct {
    content string
    content_type string
}

type Handler struct {}

func (*Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if valid, path := valid_path(configuration["root"] + r.URL.String()); valid {
        send_file(r, w, 200, path)
    } else if configuration["showfiles"] == "on" {
        log.Printf("Check directory: %s", r.URL.String())
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

    // create list of possible index files
    for _, in := range strings.Split(configuration["index"], " ") {
        index_files = append(index_files, strings.TrimSpace(in))
    }
    
    // validate error file on load instead of per request
    err_file := configuration["error404"]
    configuration["error404"] = ""
    
    if !strings.HasPrefix(err_file, "/") {
        err_file = configuration["root"] + "/" + err_file
    }
    
    if valid_file(err_file) {    
        _, err := ioutil.ReadFile(err_file)
            
        if err == nil {
            configuration["error404"] = err_file
        }
    }
}

func valid_path(path string) (bool, string) {
    if valid_file(path) {
        return true, path
    }

    if info, err := os.Stat(path); err == nil && info.IsDir() {
        return valid_index(path)
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
    
    fileName := strings.Split(file, ".")

    return "." + fileName[len(fileName) - 1]
}

func get_file(path string) string {
    filePath := strings.Split(path, "/")
    
    return filePath[len(filePath) - 1]
}

func get_mime_type(data []byte) string {
    return http.DetectContentType(data)
}

func load_file_watcher() {
    watcher, err := fsnotify.NewWatcher()
    check(err)
    file_watcher = watcher
    done := make(chan bool)

    go func() {
        for {
            select {
            case event := <-file_watcher.Events:
                if event.Op&fsnotify.Write == fsnotify.Write {
                    delete(file_cache, event.Name)
                    file_watcher.Remove(event.Name)
                }
            case err := <-file_watcher.Errors:
                log.Println("File watcher error: ", err)
            }
        }
    }()

    <-done
}

func send_file(r *http.Request, w http.ResponseWriter, status int, static_file string) {
    var data string
    var mime_type string

    if value, ok := file_cache[static_file]; ok {
        data = value.content  
        mime_type = value.content_type
    } else {
        file_data, err := ioutil.ReadFile(static_file)

        if err != nil {
            log.Println("Could not read file: " + static_file)
            send_locked(r, w)
            return
        }

        data = string(file_data)
        mime_type = get_mime_type(file_data)
        file_cache[static_file] = CacheFile{data, mime_type}
        file_watcher.Add(static_file)
    }

    if valid_file(static_file) {
        send_response(r, w, status, mime_type, data)
    } else {
        send_not_found(r, w)
    }
}

func send_not_found(r *http.Request, w http.ResponseWriter) {
    if configuration["error404"] != "" {
        send_file(r, w, 404, configuration["error404"])
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
    
    w.Header().Set("Server", SERVER_NAME + "/" + VERSION)
    w.WriteHeader(status)
    //log.Printf("%d %s: %s", status, r.Method, r.URL.String())
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
    go load_file_watcher()
    start_server()
}

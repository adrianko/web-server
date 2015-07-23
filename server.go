package main

import (
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strings"
    "strconv"
    "math"
    "fmt"
    "gopkg.in/fsnotify.v1"
)
/**
 * TODO Add icons to different folder and file types in file index
 */

 // Name of the server sent in the HTTP response header
const SERVER_NAME string = "Maester"

// Current version of the server
const VERSION string = "0.4"

// Number of bytes to use per Kb/Mb/Gb.
// Can also use 1000 for SI
const BYTES_PER_KB int64 = 1024

// Float of BYTES_PER_KB.
// Cached here to prevent frequent conversions of int64 to float64
const BYTES_PER_KB_FL float64 = float64(BYTES_PER_KB)

// Base64 encoded images sent back in file lists
var file_icons map[string]string = map[string]string{
    "back":     "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABQAAAAWCAMAAAD3n0w0AAAAElBMVEX////M//+ZmZlmZmYzMzMAA" +
                "ACei5rnAAAAAnRSTlP/AOW3MEoAAABVSURBVHgBbdFBCsBACENR45j7X7kQtC0T//KRjRhYevGgyjBL+VLZUtlS2VItS1AI1QQO" +
                "NgNZHCSUZJAc+ZB3sViFGzPcDmxZqdsvgRB/aJRu73D0HuO2BJfFYAozAAAAAElFTkSuQmCC",
    "file":     "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABQAAAAWCAMAAAD3n0w0AAAAD1BMVEX////M//+ZmZkzMzMAAABVs" +
                "TOVAAAAAnRSTlP/AOW3MEoAAABXSURBVHgBpcpBDsQwCENRY+f+Zx55QKShlbrozyrPQNcig9AJekJoI7mcUGo0FVobS/8v0X/u" +
                "aSNqIxMrDkxyQGMbD2wbaLojJnbz8gO6VxSPZIxYo4gfuU0C6reH1fMAAAAASUVORK5CYII=",
    "folder":   "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABQAAAAWCAMAAAD3n0w0AAAAElBMVEX/////zJnM//+ZZjMzMzMAA" +
                "ADCEvqoAAAAA3RSTlP//wDXyg1BAAAASElEQVR42s3KQQ6AQAhDUaXt/a/sQDrRJu7c+NmQB0e99B3lnqjT6cYx6zSIbV40n3D7" +
                "psYMoBoz4w8/EdNYQsbGEjNxYSljXTEsA9O1pLTvAAAAAElFTkSuQmCC",
    "image":    "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABQAAAAWCAMAAAD3n0w0AAAAJFBMVEX/////MzPM///MzMyZmZlmZ" +
                "mZmAAAzMzMAmcwAmTMAM2YAAADMt1kEAAAAA3RSTlP//wDXyg1BAAAAbUlEQVR42m3JQRKDMAxDUTWGKIH73xdNBB275Wv3BOaa" +
                "A/GN4BuCbyg1flz4Z8GOk2s/6CoCm0o4VQAFJxBLcWPIRIq6Hux9xjTGcaOajjESctnOkTFoKyhVo6B9VLT9YxwuYcrIrUQhWjt" +
                "rwgujCAczWx6q5QAAAABJRU5ErkJggg==",
    "text":     "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABQAAAAWCAMAAAD3n0w0AAAAD1BMVEX////M//+ZmZkzMzMAAABVs" +
                "TOVAAAAAnRSTlP/AOW3MEoAAABISURBVHjatcrRCgAgCENRbf7/N7dKomGvngjhMsPLD4NdMPwia438NRIyxsaL/XQZhyxpkC6z" +
                "yjLXGVXnkhqWJWIIrOgeinECLlUCjBCqNQoAAAAASUVORK5CYII=",
}

// The default configuration file
var config_file string = "/etc/maester-http"

// The default configuration settings used if not provided in a custom configuration
var configuration map[string]string = map[string]string{
    "root":      "/var/www",
    "port":      "80",
    "interface": "0.0.0.0",
    "index":     "index.html",
    "error404":  "error/error404.html",
    "showfiles": "off",
}

// Default index files to search for if not explicitly stated in the URL
var index_files []string = []string{}

// Cache of static files sent back in HTTP responses to prevent frequent file system IO
var file_cache map[string]CacheFile = make(map[string]CacheFile)

// File watcher to invalidate cached files
var file_watcher *fsnotify.Watcher

// CacheFile type storing both a file's contents and it's type
type CacheFile struct {
    content string
    content_type string
}

// Default server handler
type Handler struct {}

// All requests are parse and passed through ServeHTTP
// Acts as a controller
func (*Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if valid, path := valid_path(configuration["root"] + r.URL.String()); valid {
        send_file(r, w, 200, path)
    } else if configuration["showfiles"] == "on" {
        send_file_list(r, w, r.URL.String())
    } else {
        send_not_found(r, w)
    }
}

// Check if custom configuration file passed as argument
func read_args() {
    if len(os.Args) > 1 {
        config_file = os.Args[1]
    }
}

// Load the configuration file and attempt to read contents.
// FATAL: If cannot read file
func load_config() {
    data, err := ioutil.ReadFile(config_file)
    check(err)
    parse_config(string(data))
}

// Parse the loaded config line by line ignoring commented out lines
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

// Validate the configuration loaded and check each property contains valid value
// FATAL: If web root is invalid path
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

// Checks if a path is valid in the file system
// If is a valid file, returns true
// If is a valid directory, checks for an index file
// Otherwise returns false
func valid_path(path string) (bool, string) {
    if valid_file(path) {
        return true, path
    }

    if valid_dir(path) {
        return valid_index(path)
    }

    return false, path
}

// Checks if given directory has a valid index file
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

// Checks if a given file is valid and not a directory
// Reports dotfiles as invalid as they should not be sent back over HTTP e.g. .htaccess / .htpasswd
func valid_file(path string) bool {
    filePath := strings.Split(path, "/")

    if strings.HasPrefix(filePath[len(filePath) - 1], ".") {
        return false;
    }

    info, err := os.Stat(path); 
    
    return err == nil && !info.IsDir()
}

// Checks if given [ath is a valid directory
func valid_dir(path string) bool {
    info, err := os.Stat(path)
    
    return err == nil && info.IsDir()
}

// Returns the extension of a file or file in a path
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

// Returns the filename and extension from a path
func get_file(path string) string {
    filePath := strings.Split(path, "/")
    
    return filePath[len(filePath) - 1]
}

// Returns the MIME type of a given byte array
func get_mime_type(data []byte) string {
    return http.DetectContentType(data)
}

// Formats an int64 into bytes/Kb/Mb/Gb/Tb/Pb/Eb
func format_bytes(bytes int64) string {
    if bytes < BYTES_PER_KB {
        return strconv.FormatInt(bytes, 10) + " B"
    }

    bytes_fl := float64(bytes)
    exp := int(math.Log(bytes_fl) / math.Log(BYTES_PER_KB_FL))
    pre := []string{"K", "M", "G", "T", "P", "E"}[exp - 1]

    return fmt.Sprintf("%.1f %sB", bytes_fl / math.Pow(BYTES_PER_KB_FL, float64(exp)), pre)

}

// Loads the file watcher for cachced files
// Watched files are removed from cache and then the watcher if modified on the file system
// FATAL: If cannot start file watcher
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

// Returns the HTML img tag for a requested icon
func get_icon(icon string) string {
    return "<img src=\"" + file_icons[icon] + "\" alt=\"" + icon + " icon\" />"
}

// Return an HTML img tag for a given MIME type
func get_icon_by_mime(mime string) string {
    return ""
}

// Return either a file icon or a folder icon 
func file_folder_icon(is_directory bool) string {
    if is_directory {
        return get_icon("folder")
    }

    return get_icon("file")
}

// Send a file with given status and file path
// If the file cannot be read, HTTP code 423 is sent
// If the file is cached, send the file contents from cache
// If the file is not cached, attempt to read it, cache it and add it to the file watcher
// If the file is not valid, send HTTP code 404
// Other wise send HTTP code 200 with the file contents
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

// If the showfiles config setting is on, the client has request a directory and the directory does not have a valid
// index, send a list of files in the directory
func send_file_list(r *http.Request, w http.ResponseWriter, url string) {
    files, _ := ioutil.ReadDir(configuration["root"] + url)
    file_list := "<html>"
    file_list += "<head>"
    file_list += "<title>Index of: " + url + "</title>"
    file_list += "<style>"
    file_list += "body{font-family:Cambria}"
    file_list += "tr:first-child{font-weight:bold}"
    file_list += "td{padding:0 5px}"
    file_list += "img{vertical-align:middle}"
    file_list += "</style>"
    file_list += "<body>"
    file_list += "<h1>Directory index: <em>" + url + "</em></h1>"
    file_list += "<table>"
    file_list += "<tr><td>Name</td><td>Last modified</td><td>Size</td></tr>"
    file_list += "<tr><td>" + get_icon("back") + " <a href=\"../\">Parent directory</a></td><td></td><td></td></tr>"

    if !strings.HasSuffix(url, "/") {
        url += "/"
    }

    for _, f := range files {
        info, _ := os.Stat(configuration["root"] + url + f.Name())
        file_list += "<tr>"
        file_list += "<td>" + file_folder_icon(info.IsDir()) + " <a href=\"" + url + f.Name() + "\">" + f.Name() +
            "</a></td>"
        file_list += "<td>" + info.ModTime().String() + "</td>"
        file_list += "<td>"

        if !info.IsDir() {
            file_list += format_bytes(info.Size())
        } else {
            file_list += "-"
        }

        file_list += "</td>"
        file_list += "</tr>"
    }

    file_list += "</table>"
    file_list += "</body>"
    file_list += "</html>"

    send_response(r, w, 200, "text/html", file_list)
}

// Send a not found HTTP response
// If a custom 404 page is availabe and can be read, send it
// Otherwise send a plain 404 message
func send_not_found(r *http.Request, w http.ResponseWriter) {
    if configuration["error404"] != "" {
        send_file(r, w, 404, configuration["error404"])
    } else {
        send_response(r, w, 404, "text/plain", "404: Not found")
    }
}

// Send an HTTP resourse locked code: 423
func send_locked(r *http.Request, w http.ResponseWriter) {
    send_response(r, w, 423, "", "")
}

// Send an HTTP response with a given code, content type and contents
// Set the content type if provided
func send_response(r *http.Request, w http.ResponseWriter, status int, content_type string, content string) {
    if content_type != "" {
        w.Header().Set("Content-Type", content_type)
    }
    
    w.Header().Set("Server", SERVER_NAME + "/" + VERSION)
    w.WriteHeader(status)
    //log.Printf("%d %s: %s", status, r.Method, r.URL.String())
    io.WriteString(w, content)
}

// Set the server interface and port to the configuration set value and start the server
// FATAL: If the server cannot be started
func start_server() {
    server := http.Server{Addr: configuration["interface"] + ":" + configuration["port"], Handler: &Handler{}}
    log.Printf("Running server %s:%s\n", configuration["interface"], configuration["port"])
    check(server.ListenAndServe())
}

// Check if the error exists
// If so log a Fatal error
func check(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

// Load the args, load the config, start the file watcher and then start the server
func main() {
    read_args()
    load_config()
    go load_file_watcher()
    start_server()
}
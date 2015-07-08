package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const NIC string = "0.0.0.0"
const PORT int = 8000

var handlers map[string]func(http.ResponseWriter, *http.Request) = make(map[string]func(http.ResponseWriter,
	*http.Request))

type Handler struct{}

func (*Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handle, ok := handlers[r.URL.String()]; ok {
		handle(w, r)

		return
	}

	reqError(w, r)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	data, err := ioutil.ReadFile("../conf/config")
	check(err)
	log.Print(string(data))
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

func sendOK(r *http.Request, w http.ResponseWriter, contentType string, content string) {
	send(r, w, http.StatusOK, contentType, content)
}

func send(r *http.Request, w http.ResponseWriter, status int, contentType string, content string) {
	setContentType(w, contentType)
	w.WriteHeader(status)
	logRequest(r, w, status)
	io.WriteString(w, content)
}

func setContentType(w http.ResponseWriter, contentType string) {
	var cType string

	switch contentType {
	case "json":
		cType = "application/json"
		break
	case "html":
		cType = "text/html"
		break
	case "plain":
	default:
		cType = "text/plain"
		break
	}

	w.Header().Set("Content-Type", cType)
}

//handlers
func hello(w http.ResponseWriter, r *http.Request) {
	sendOK(r, w, "plain", "Hello world")
}

func helloHTML(w http.ResponseWriter, r *http.Request) {
	sendOK(r, w, "html", "<h1>Hello world</h1>")
}

func helloJSON(w http.ResponseWriter, r *http.Request) {
	sendOK(r, w, "json", "{\"hello\": \"world\"}")
}

func reqError(w http.ResponseWriter, r *http.Request) {
	send(r, w, http.StatusNotFound, "html", "<h1>Not found</h1>")
}

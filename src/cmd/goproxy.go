package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func basePath(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi World!")
}

var reverse *httputil.ReverseProxy

type ErrorCatcher struct {
	http.ResponseWriter
}

func (e *ErrorCatcher) Header() http.Header {
	return e.ResponseWriter.Header()
}
func (e *ErrorCatcher) Write(b []byte) (int, error) {
	return e.ResponseWriter.Write(b)
}
func (e *ErrorCatcher) WriteHeader(code int) {
	if 200 <= code && 400 > code {
		e.ResponseWriter.WriteHeader(code)
		return
	}
	panic(code)
}

func routeProxy(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if icode := recover(); nil != icode {
			switch icode.(type) {
			case int:
				t, err := template.ParseFiles("templates/500.html")
				if nil != err {
					log.Fatal(err.Error())
				}
				t.Execute(w, nil)
			}
		}
	}()
	reverse.ServeHTTP(&ErrorCatcher{w}, r)
}

func main() {
	port := flag.String("http", ":14031", "Bind address and port for webserver")
	dest := flag.String("dest", "http://localhost:9292", "Destination for reverse proxying")
	flag.Parse()

	http.Handle("/_static/", http.StripPrefix("/_static", http.FileServer(http.Dir("static"))))

	// http.HandleFunc("/", basePath)
	destUrl, _ := url.Parse(*dest)
	reverse = httputil.NewSingleHostReverseProxy(destUrl)
	http.HandleFunc("/", routeProxy)
	log.Fatal(http.ListenAndServe(*port, nil))
}

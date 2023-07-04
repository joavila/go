package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
	log.Printf("Healthcheck invoked from: %s", r.RemoteAddr )
        fmt.Fprint(w, "OK")
    })
    http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
	log.Printf("Invoked from: %s", r.RemoteAddr)
        fmt.Fprint(w, "hello world\n")
    })
    log.Fatal(http.ListenAndServe(":80", nil))
}

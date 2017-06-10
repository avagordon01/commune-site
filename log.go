package main

import (
    "log"
    "net/http"
)

func log_req(fn http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.Method, r.Proto, r.URL, r.Header, r.Body, r.RemoteAddr)
        fn(w, r)
    }
}

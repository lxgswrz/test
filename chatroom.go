package main

import (
    "io"
    "net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "pong~")
}

func main() {
    http.HandleFunc("/ping", Ping)
    http.ListenAndServe(":8080", nil)
}
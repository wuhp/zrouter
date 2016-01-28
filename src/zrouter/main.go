package main

import (
    "log"
    "net/http"
)

func startStatusServer() {
    log.Printf("Starting status server on port 10002 ...\n")
    log.Fatalln(http.ListenAndServe(":10002", NewRouter()))
}

func main() {
    go startStatusServer()
    log.Printf("Starting proxy server on port 10001 ...\n")
    log.Fatalln(http.ListenAndServe(":10001", proxy))
}

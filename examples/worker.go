package main

import (
    "fmt"
    "net/http"
    "os"
    "time"
)

func sleep1(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Sleeping 1s ...")
    time.Sleep(1 * time.Second)
}

func sleep5(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Sleeping 5s ...")
    time.Sleep(5 * time.Second)
}

func sleep10(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Sleeping 10s ...")
    time.Sleep(10 * time.Second)
}

func main() {
    if len(os.Args) != 2 {
        fmt.Printf("Usage: %s <port>\n", os.Args[0])
        return
    }

    http.HandleFunc("/sleep1", sleep1)
    http.HandleFunc("/sleep5", sleep5)
    http.HandleFunc("/sleep10", sleep10)

    fmt.Printf("Starting server on 0.0.0.0:%s ...\n", os.Args[1])
    http.ListenAndServe(fmt.Sprintf(":%s", os.Args[1]), nil)
}

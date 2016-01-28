package main

import (
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
)

func Ping(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "pong")
}

func ListService(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(proxy.ListService())
}

func PostService(w http.ResponseWriter, r *http.Request) {
    in := new(Service)
    if err := json.NewDecoder(r.Body).Decode(in); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    if len(in.Name) == 0 {
        http.Error(w, "empty service name", http.StatusBadRequest)
        return
    }

    if len(in.Url) == 0 {
        in.Url = "/"
    }

    if err := proxy.PostService(in); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
}

func GetService(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    sname := vars["service"]

    service, err := proxy.GetService(sname)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    json.NewEncoder(w).Encode(service)
}

func PutService(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    sname := vars["service"]

    in := new(Service)
    if err := json.NewDecoder(r.Body).Decode(in); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    if len(in.Url) == 0 {
        in.Url = "/"
    }

    in.Name = sname

    err := proxy.PutService(in)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
}

func DeleteService(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    sname := vars["service"]

    err := proxy.DeleteService(sname)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
}

func ListServicePool(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    sname := vars["service"]

    pools, err := proxy.ListServicePool(sname)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    json.NewEncoder(w).Encode(pools)
}

func GetServicePool(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    sname := vars["service"]
    pname := vars["pool"]

    pool, err := proxy.GetServicePool(sname, pname)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    json.NewEncoder(w).Encode(pool)
}

func PutServicePool(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    sname := vars["service"]
    pname := vars["pool"]

    in := new(Pool)
    if err := json.NewDecoder(r.Body).Decode(in); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if len(in.LBPolicy) == 0 {
        in.LBPolicy = "random"
    }

    if err := proxy.PutServicePool(sname, pname, in); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
}

func DeleteServicePool(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    sname := vars["service"]
    pname := vars["pool"]

    if err := proxy.DeleteServicePool(sname, pname); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
}

func ListServicePoolNode(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    sname := vars["service"]
    pname := vars["pool"]

    nodes, err := proxy.ListServicePoolNode(sname, pname)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    json.NewEncoder(w).Encode(nodes)
}

func PostServicePoolNode(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    sname := vars["service"]
    pname := vars["pool"]

    in := new(Node)
    if err := json.NewDecoder(r.Body).Decode(in); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if len(in.Name) == 0 {
        http.Error(w, "empty node name", http.StatusBadRequest)
        return
    }

    if in.Status == "unloading" {
        in.Status = "off"
    }

    in.ConnNum = 0

    if err := proxy.PostServicePoolNode(sname, pname, in); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
}

func GetServicePoolNode(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    sname := vars["service"]
    pname := vars["pool"]
    nname := vars["node"]

    node, err := proxy.GetServicePoolNode(sname, pname, nname)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    json.NewEncoder(w).Encode(node)
}

func PutServicePoolNode(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    sname := vars["service"]
    pname := vars["pool"]
    nname := vars["node"]

    in := new(Node)
    if err := json.NewDecoder(r.Body).Decode(in); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    in.Name = nname

    if err := proxy.PutServicePoolNode(sname, pname, in); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
}

func DeleteServicePoolNode(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    sname := vars["service"]
    pname := vars["pool"]
    nname := vars["node"]

    if err := proxy.DeleteServicePoolNode(sname, pname, nname); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
}

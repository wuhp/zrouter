package main

import (
    "math/rand"
    "net/http"
    "sort"
)

type Node struct {
    Name    string `json:"name"`
    Host    string `json:"host"`      // `host:port`
    Status  string `json:"status"`    // `on/off/unloading`
    Weight  int    `json:"weight"`    // `1 ~ 10`
    ConnNum int    `json:"conn_num"`
}

type Pattern struct {
    Type  string `json:"type"`       // `ip/header`, now only support header
    Value string `json:"value"`      // `MyHeader`, now only support checking some header exists
}

type Pool struct {
    Pattern  *Pattern `json:"pattern"`   // request will fall into the pool if it matching the pattern of the pool, always match if Pattern is nil
    LBPolicy string   `json:"lb_policy"` // load balance policy, now only support `random`
    Nodes    []*Node  `json:"-"`     // request will go to one of Nodes according to the LBPolicy
}

type Service struct {
    Name      string `json:"name"`
    Host      string `json:"host"`       // Host and Url represents a service
    Url       string `json:"url"`
    ProdPool  *Pool  `json:"-"`
    GrayPool  *Pool  `json:"-"`
    DebugPool *Pool  `json:"-"`
}

////////////////////////////////////////////////////////////////////////////////

func (p *Pattern) Match(req *http.Request) bool {
    if p.Type == "header" {
        for k, _ := range req.Header {
            if req.Header.Get(k) == p.Value {
                return true
            }
        }
    }

    return false
}

func (p *Pool) Pick() *Node {
    available := make([]*Node, 0)
    for _, node := range p.Nodes {
        if node.Status == "on" {
            available = append(available, node)
        }
    }

    if len(available) == 0 {
        return nil
    }

    if p.LBPolicy == "random" {
        return available[rand.Intn(len(available))]
    }

    return nil
}

////////////////////////////////////////////////////////////////////////////////

type Services []*Service

func (ss Services) Len() int {
    return len(ss)
}

func (ss Services) Swap(i, j int) {
    ss[i], ss[j] = ss[j], ss[i]
}

type ServiceReverseByUrl struct {
    Services
}

func (s ServiceReverseByUrl) Less(i, j int) bool {
    return s.Services[i].Url > s.Services[j].Url
}

////////////////////////////////////////////////////////////////////////////////

func findServiceByHost(services []*Service, host string) []*Service {
    result := make([]*Service, 0)
    for _, s := range services {
        if s.Host == host {
            result = append(result, s)
        }
    }

    return result
}

func sortServiceByUrlDsc(ss []*Service) []*Service {
    sort.Sort(ServiceReverseByUrl{ss})
    return ss
}

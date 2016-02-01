package main

import (
    "fmt"
    "io"
    "log"
    "net"
    "net/http"
    "strings"
    "sync"
)

type Proxy struct {
    sync.Mutex
    Services []*Service
}

var proxy *Proxy

func init() {
    proxy = new(Proxy)
    proxy.Services = make([]*Service, 0)
}

////////////////////////////////////////////////////////////////////////////////

func (p *Proxy) getService(name string) (*Service, error) {
    for _, s := range p.Services {
        if s.Name == name {
            return s, nil
        }
    }

    return nil, fmt.Errorf("service %s not found", name)
}

func (p *Proxy) getServicePool(sname, pname string) (*Pool, error) {
    service, err := p.getService(sname)
    if err != nil {
        return nil, err
    }

    switch pname {
    case "prod":
        return service.ProdPool, nil
    case "gray":
        return service.GrayPool, nil
    case "debug":
        return service.DebugPool, nil
    }

    return nil, fmt.Errorf("pool %s not found", pname)
}

func (p *Proxy) getServicePoolNode(sname, pname, nname string) (*Node, error) {
    pool, err := p.getServicePool(sname, pname)
    if err != nil {
        return nil, err
    }

    for _, node := range pool.Nodes {
        if node.Name == nname {
            return node, nil
        }
    }

    return nil, fmt.Errorf("node %s not found", nname)
}

func (p *Proxy) ListService() []*Service {
    p.Lock()
    defer p.Unlock()

    return p.Services
}

func (p *Proxy) PostService(service *Service) error {
    p.Lock()
    defer p.Unlock()

    if _, err := p.getService(service.Name); err == nil {
        return fmt.Errorf("duplicate service %s", service.Name)
    }

    service.ProdPool = new(Pool)
    service.ProdPool.LBPolicy = "random"
    service.ProdPool.Nodes = make([]*Node, 0)
    service.GrayPool = new(Pool)
    service.GrayPool.LBPolicy = "random"
    service.GrayPool.Nodes = make([]*Node, 0)
    service.DebugPool = new(Pool)
    service.DebugPool.LBPolicy = "random"
    service.DebugPool.Nodes = make([]*Node, 0)

    p.Services = append(p.Services, service)
    return nil
}

func (p *Proxy) GetService(name string) (*Service, error) {
    p.Lock()
    defer p.Unlock()

    return p.getService(name)
}

func (p *Proxy) PutService(service *Service) error {
    p.Lock()
    defer p.Unlock()

    s, err := p.getService(service.Name)
    if err != nil {
        return err
    }

    s.Host = service.Host
    s.Url = service.Url
    return nil
}

func (p *Proxy) DeleteService(name string) error {
    p.Lock()
    defer p.Unlock()

    var i int
    for i = 0; i < len(p.Services); i++ {
        if p.Services[i].Name == name {
            break
        }
    }

    if i == len(p.Services) {
        return fmt.Errorf("service %s not found", name)
    }

    p.Services = append(p.Services[:i], p.Services[i+1:]...)
    return nil
}

func (p *Proxy) ListServicePool(sname string) ([]string, error) {
    p.Lock()
    defer p.Unlock()

    if _, err := p.getService(sname); err != nil {
        return nil, err
    }

    return []string{"prod", "gray", "debug"}, nil
}

func (p *Proxy) GetServicePool(sname, pname string) (*Pool, error) {
    p.Lock()
    defer p.Unlock()

    return p.getServicePool(sname, pname)
}

func (p *Proxy) PutServicePool(sname, pname string, pl *Pool) error {
    p.Lock()
    defer p.Unlock()

    pool, err := p.getServicePool(sname, pname)
    if err != nil {
        return err
    }

    pool.Pattern = pl.Pattern
    pool.LBPolicy = pl.LBPolicy

    return nil
}

func (p *Proxy) DeleteServicePool(sname, pname string) error {
    p.Lock()
    defer p.Unlock()

    service, err := p.getService(sname)
    if err != nil {
        return err
    }

    switch pname {
    case "prod":
        service.ProdPool = nil
    case "gray":
        service.GrayPool = nil
    case "debug":
        service.DebugPool = nil
    default:
        return fmt.Errorf("pool %s not found", pname)
    }

    return nil
}

func (p *Proxy) ListServicePoolNode(sname, pname string) ([]*Node, error) {
    p.Lock()
    defer p.Unlock()

    pool, err := p.getServicePool(sname, pname)
    if err != nil {
        return nil, err
    }

    return pool.Nodes, nil
}

func (p *Proxy) PostServicePoolNode(sname, pname string, n *Node) error {
    p.Lock()
    defer p.Unlock()

    pool, err := p.getServicePool(sname, pname)
    if err != nil {
        return err
    }

    if _, err := p.getServicePoolNode(sname, pname, n.Name); err == nil {
        return fmt.Errorf("duplicate node %s", n.Name)
    }

    if pool.Nodes == nil {
        pool.Nodes = make([]*Node, 0)
    }

    pool.Nodes = append(pool.Nodes, n)
    return nil
}

func (p *Proxy) GetServicePoolNode(sname, pname, nname string) (*Node, error) {
    p.Lock()
    defer p.Unlock()

    return p.getServicePoolNode(sname, pname, nname)
}

func (p *Proxy) PutServicePoolNode(sname, pname string, n *Node) error {
    p.Lock()
    defer p.Unlock()

    node, err := p.getServicePoolNode(sname, pname, n.Name)
    if err != nil {
        return err
    }

    node.Weight = n.Weight
    node.Status = n.Status

    if node.Status == "unloading" && node.ConnNum == 0 {
        node.Status = "off"
    }

    return nil
}

func (p *Proxy) DeleteServicePoolNode(sname, pname, nname string) error {
    p.Lock()
    defer p.Unlock()

    pool, err := p.getServicePool(sname, pname)
    if err != nil {
        return err
    }

    var i int
    for i = 0; i < len(pool.Nodes); i++ {
        if pool.Nodes[i].Name == nname {
            break
        }
    }

    if i == len(pool.Nodes) {
        return fmt.Errorf("node %s not found", nname)
    }

    pool.Nodes = append(pool.Nodes[:i], pool.Nodes[i+1:]...)
    return nil
}

////////////////////////////////////////////////////////////////////////////////

func (p *Proxy) lookupService(req *http.Request) *Service {
    host := req.Host
    if i := strings.Index(host, ":"); i >= 0 {
        host = host[0:i]
    }

    path := req.URL.Path
    ss := findServiceByHost(p.Services, host)
    if len(ss) == 0 {
        ss = findServiceByHost(p.Services, "")
        if len(ss) == 0 {
            return nil
        }
    }

    for _, s := range sortServiceByUrlDsc(ss) {
        if strings.HasPrefix(path, s.Url) {
            return s
        }
    }

    return nil
}

func (p *Proxy) lookupNode(s *Service, req *http.Request) *Node {
    for _, pool := range []*Pool{s.DebugPool, s.GrayPool} {
        if pool.Pattern != nil && pool.Pattern.Match(req) {
            return pool.Pick()
        }
    }

    return s.ProdPool.Pick()
}

func (p *Proxy) lookup(req *http.Request) *Node {
    p.Lock()
    defer p.Unlock()

    s := p.lookupService(req)
    if s == nil {
        return nil
    }

    return p.lookupNode(s, req)
}

func (p *Proxy) increaseConn(node *Node) {
    p.Lock()
    defer p.Unlock()

    node.ConnNum++
}

func (p *Proxy) decreaseConn(node *Node) {
    p.Lock()
    defer p.Unlock()

    node.ConnNum--

    if node.ConnNum == 0 && node.Status == "unloading" {
        node.Status = "off"
    }
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
    node := p.lookup(req)
    if node == nil {
        rw.WriteHeader(http.StatusNotFound)
        return
    }

    p.increaseConn(node)
    defer p.decreaseConn(node)

    outreq := new(http.Request)
    *outreq = *req

    outreq.URL.Scheme = "http"
    outreq.URL.Host = node.Host

    outreq.Proto = "HTTP/1.1"
    outreq.ProtoMajor = 1
    outreq.ProtoMinor = 1
    outreq.Close = false

    if outreq.Header.Get("Connection") != "" {
        outreq.Header = make(http.Header)
        copyHeader(outreq.Header, req.Header)
        outreq.Header.Del("Connection")
    }

    if clientIp, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
        outreq.Header.Set("X-Forwarded-For", clientIp)
    }

    res, err := http.DefaultTransport.RoundTrip(outreq)
    if err != nil {
        log.Printf("proxy round trip error: %v", err)
        rw.WriteHeader(http.StatusInternalServerError)
        return
    }

    copyHeader(rw.Header(), res.Header)

    rw.WriteHeader(res.StatusCode)

    if res.Body != nil {
        var dst io.Writer = rw
        io.Copy(dst, res.Body)
    }
}

////////////////////////////////////////////////////////////////////////////////

func copyHeader(dst, src http.Header) {
    for k, vv := range src {
        for _, v := range vv {
            dst.Add(k, v)
        }
    }
}

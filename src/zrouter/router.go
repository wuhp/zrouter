package main

import (
    "fmt"
    "log"
    "net/http"
    "runtime/debug"
    "time"

    "github.com/gorilla/mux"
)

type Route struct {
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

var routes = []Route{
    Route{"GET", "/api/ping", Ping},

    Route{"GET",    "/api/services",           ListService  },
    Route{"POST",   "/api/services",           PostService  },
    Route{"GET",    "/api/services/{service}", GetService   },
    Route{"PUT",    "/api/services/{service}", PutService   },
    Route{"DELETE", "/api/services/{service}", DeleteService},

    Route{"GET",    "/api/services/{service}/pools",        ListServicePool  },
    Route{"GET",    "/api/services/{service}/pools/{pool}", GetServicePool   },
    Route{"PUT",    "/api/services/{service}/pools/{pool}", PutServicePool   },
    Route{"DELETE", "/api/services/{service}/pools/{pool}", DeleteServicePool},

    Route{"GET",    "/api/services/{service}/pools/{pool}/nodes",        ListServicePoolNode  },
    Route{"POST",   "/api/services/{service}/pools/{pool}/nodes",        PostServicePoolNode  },
    Route{"GET",    "/api/services/{service}/pools/{pool}/nodes/{node}", GetServicePoolNode   },
    Route{"PUT",    "/api/services/{service}/pools/{pool}/nodes/{node}", PutServicePoolNode   },
    Route{"DELETE", "/api/services/{service}/pools/{pool}/nodes/{node}", DeleteServicePoolNode},
}

type InnerResponseWriter struct {
    StatusCode int
    isSet      bool
    http.ResponseWriter
}

func (i *InnerResponseWriter) WriteHeader(status int) {
    if !i.isSet {
        i.StatusCode = status
        i.isSet = true
    }

    i.ResponseWriter.WriteHeader(status)
}

func (i *InnerResponseWriter) Write(b []byte) (int, error) {
    i.isSet = true
    return i.ResponseWriter.Write(b)
}

func wrapper(inner http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        s := time.Now()
        wr := &InnerResponseWriter{
            StatusCode:     200,
            isSet:          false,
            ResponseWriter: w,
        }

        defer func() {
            if err := recover(); err != nil {
                debug.PrintStack()
                wr.WriteHeader(http.StatusInternalServerError)
                log.Printf("Panic: %v\n", err)
                fmt.Fprintf(w, fmt.Sprintln(err))
            }

            d := time.Now().Sub(s)
            log.Printf("%s %s %d %s\n", r.Method, r.RequestURI, wr.StatusCode, d.String())
        }()

        inner.ServeHTTP(wr, r)
    })
}

func NewRouter() *mux.Router {
    router := mux.NewRouter()
    for _, route := range routes {
        router.Methods(route.Method).Path(route.Pattern).HandlerFunc(wrapper(route.HandlerFunc))
    }

    return router
}


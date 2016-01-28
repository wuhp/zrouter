# zrouter

A light weight http proxy, written in go.

## suitable situation
    reverse proxy
    load balance
    gray test
    hot deployment(upgrad without down time)

## build
    go install zrouter

## start
    ./bin/zrouter

## example
Start 2 worker instance
    cd example
    go run worker.go 20001   ## Terminal 1
    go run worker.go 20002   ## Terminal 2

Setup rules
    curl -i -X POST http://localhost:10002/api/services -d '{"name":"ping_server"}'
    curl -i -X POST http://localhost:10002/api/services/ping_server/pools/prod/nodes -d '{"name":"prod_001", "host":"127.0.0.1:20001", "status":"on"}'
    curl -i -X POST http://localhost:10002/api/services/ping_server/pools/prod/nodes -d '{"name":"prod_002", "host":"127.0.0.1:20002", "status":"on"}'
    curl -i -X GET http://localhost:10002/api/services/ping_server/pools/prod/nodes

Test
    curl -i http://localhost:10001/sleep1
    curl -i http://localhost:10001/sleep10

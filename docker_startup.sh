#!/bin/sh

SERVICE=$1
ENV=$2

if [ "$ENV" == "dev" ]; then
    case $SERVICE in
        frontend)        
        npm run dev
        ;;
        backend_api)
        reflex -r '\.go$' -s -- sh -c 'go run api/api.go'
        ;;
        zmq_router)
        reflex -r '\.go$' -s -- sh -c 'go run router/router.go'
        ;;
        zmq_worker)
        reflex -r '\.go$' -s -- sh -c 'go run worker/worker.go'
        ;;
    esac
else  
    case $SERVICE in
        frontend)
        npm run build
        npm run start
        ;;
        backend_api)
        go run api/api.go
        ;;
        zmq_router)
        go run router/router.go
        ;;
        zmq_worker)
        go run worker/worker.go
        ;;
    esac
fi
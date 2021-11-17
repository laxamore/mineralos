#!/bin/bash

SERVICE=$1
ENV=$2

if [ "$ENV" == "dev" ]
then
    case $SERVICE in
        frontend)        
        npm run dev
        ;;
        backend_api)
        reflex -r '\.go$' -s -- sh -c 'go run main.go'
        ;;
        zmq_router)
        reflex -r '\.go$' -s -- sh -c 'go run zmq/server/router/main.go'
        ;;
        zmq_worker)
        reflex -r '\.go$' -s -- sh -c 'go run zmq/server/worker/main.go'
        ;;
    esac
else  
    case $SERVICE in
        frontend)
        npm run build
        npm run start
        ;;
        backend_api)
        go run main.go
        ;;
        zmq_router)
        go run zmq/server/router/main.go
        ;;
        zmq_worker)
        go run zmq/server/router/main.go
        ;;
    esac
fi
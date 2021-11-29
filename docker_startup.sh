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
    esac
fi
FROM golang:1.17.3-alpine3.13

ENV DOCKER true

# Create app directory
RUN mkdir -p /go/src/github.com/laxamore/mineralos
WORKDIR /go/src/github.com/laxamore/mineralos

# Installing dependencies
COPY . /go/src/github.com/laxamore/mineralos
RUN go get -d -v ./...
RUN go install -v ./...
RUN go install github.com/cespare/reflex@latest

EXPOSE 5000
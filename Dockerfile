FROM golang:1.17

ENV DOCKER true

# Create app directory
RUN mkdir -p /go/src/github.com/laxamore/mineralos
WORKDIR /go/src/github.com/laxamore/mineralos

# Installing zmq
RUN apt install pkg-config -y
RUN echo "deb http://download.opensuse.org/repositories/network:/messaging:/zeromq:/release-stable/Debian_9.0/ ./" >> /etc/apt/sources.list
RUN wget https://download.opensuse.org/repositories/network:/messaging:/zeromq:/release-stable/Debian_9.0/Release.key -O- | apt-key add
RUN apt update -y
RUN apt install libzmq3-dev -y

# Installing dependencies
COPY . /go/src/github.com/laxamore/mineralos
RUN go get -d -v ./...
RUN go install -v ./...
RUN go install github.com/cespare/reflex@latest

EXPOSE 5000
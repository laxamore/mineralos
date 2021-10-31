package main

import (
	"fmt"
	"log"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	client_public, client_secret, err := zmq.NewCurveKeypair()
	if err != nil {
		log.Panic("NewCurveKeypair:", err)
	}
	server_public, server_secret, err := zmq.NewCurveKeypair()
	if err != nil {
		log.Panic("NewCurveKeypair:", err)
	}

	fmt.Printf("Client Public: %s\nClient Secret: %s\n", client_public, client_secret)
	fmt.Printf("Server Public: %s\nServer Secret: %s\n", server_public, server_secret)

	zmq.AuthCurveAdd("*", client_public)

	//  Create and bind server socket
	server, err := zmq.NewSocket(zmq.DEALER)
	if err != nil {
		log.Panic("NewSocket:", err)
	}
	defer func() {
		server.SetLinger(0)
		server.Close()
	}()
	server.SetIdentity("Server1")
	server.ServerAuthCurve("*", server_secret)
	err = server.Bind("tcp://*:9000")
	if err != nil {
		log.Panic("server.Bind:", err)
	}

	//  Create and connect client socket
	client, err := zmq.NewSocket(zmq.DEALER)
	if err != nil {
		log.Panic("NewSocket:", err)
	}
	defer func() {
		client.SetLinger(0)
		client.Close()
	}()
	server.SetIdentity("Client1")
	client.ClientAuthCurve(server_public, client_public, client_secret)
	err = client.Connect("tcp://127.0.0.1:9000")
	if err != nil {
		log.Panic("client.Connect:", err)
	}

	//  Send a message from client to server
	msg := []string{"Greetings", "Earthlings!"}
	_, err = client.SendMessage(msg[0], msg[1])
	if err != nil {
		log.Panic("client.SendMessage:", err)
	}

	// Receive message on the server
	message, err := server.RecvMessage(0)
	log.Print(message)
}

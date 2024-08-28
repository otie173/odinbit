package main

import (
	"fmt"

	"github.com/sacOO7/gowebsocket"
)

var (
	socket            gowebsocket.Socket
	connectedToServer bool
)

func connectServer(url string) {
	socket = gowebsocket.New(url)

	socket.OnConnected = func(s gowebsocket.Socket) {
		fmt.Println("Connected to server")
		connectedToServer = true
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		fmt.Println("Received connect error: ", err)
	}

	socket.OnDisconnected = func(err error, s gowebsocket.Socket) {
		fmt.Println("Disconnected from server: ", err)
		connectedToServer = false
	}

	socket.Connect()
}

func readServer() {
	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		fmt.Println("Received: ", message)
	}
}

func writeServer(message string) {
	socket.SendText(message)
}

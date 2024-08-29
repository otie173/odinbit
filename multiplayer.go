package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sacOO7/gowebsocket"
)

var (
	socket            gowebsocket.Socket
	connectedToServer bool
)

// Opcodes for responses from server
const (
	SEND_WORLD byte = iota
	SEND_DATA
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

func handleOpcode(opcode byte) {
	switch opcode {
	case SEND_WORLD:
		sendWorld()
	}
}

func readServer() {
	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		fmt.Println("Received: ", message)
	}

	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		handleOpcode(data[0])
	}
}

func sendWorld() {
	generateWorld()
	saveWorldFile()

	odinbitPath := getOdinbitPath()
	worldPath := filepath.Join(odinbitPath, "world_server.odn")
	worldData, err := os.ReadFile(worldPath)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	socket.SendBinary(worldData)
	fmt.Println("World was sended")
	err = os.Remove(worldPath)

	if err != nil {
		fmt.Println("Error: ", err)
	}
}

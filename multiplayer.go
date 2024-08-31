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
	RECEIVE_WORLD
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

	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		opcode := data[0]
		var messageData []byte

		if len(data) > 1 {
			messageData = data[1:]
			handleData(opcode, messageData)
		} else {
			handleData(opcode, nil)
		}

	}
}

func handleData(opcode byte, data []byte) {
	switch opcode {
	case SEND_WORLD:
		sendWorld()
	case RECEIVE_WORLD:
		receiveWorld(data)
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
	loadWorldFile()
	err = os.Remove(worldPath)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	currentScene = GAME
}

func receiveWorld(worldData []byte) {
	odinbitPath := getOdinbitPath()
	worldPath := filepath.Join(odinbitPath, "world_server.odn")

	if err := os.WriteFile(worldPath, worldData, 0644); err != nil {
		fmt.Println("Error with receive world from server: ", err)
	}

	world = loadWorldFile()
	if err := os.Remove(worldPath); err != nil {
		fmt.Println("Error with remove file: ", err)
	}
	currentScene = GAME
}

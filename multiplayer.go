package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sacOO7/gowebsocket"
)

var (
	socket            gowebsocket.Socket
	connectedToServer int32
	nickname          string
	password          string
	activeInput       int // 0 - nickname, 1 - password, 2 - ipAddress
	needSendWorld     int32
	needReceiveWorld  int32
	worldData         []byte
	worldDataMutex    sync.Mutex
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
		atomic.StoreInt32(&connectedToServer, 1)
		connectedToServer = 1
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		fmt.Println("Received connect error: ", err)
	}

	socket.OnDisconnected = func(err error, s gowebsocket.Socket) {
		fmt.Println("Disconnected from server: ", err)
		atomic.StoreInt32(&connectedToServer, 0)
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		fmt.Println("Received: ", message)
	}

	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		fmt.Println("Получено бинарное сообщение")
		opcode := data[0]
		var messageData []byte

		if len(data) > 1 {
			messageData = data[1:]
			handleData(opcode, messageData)
		} else {
			handleData(opcode, nil)
		}

	}
	socket.Connect()
}

func handleData(opcode byte, data []byte) {
	switch opcode {
	case SEND_WORLD:
		atomic.StoreInt32(&needSendWorld, 1)
		//sendWorld()
	case RECEIVE_WORLD:
		worldDataMutex.Lock()
		worldData = data
		worldDataMutex.Unlock()
		atomic.StoreInt32(&needReceiveWorld, 1)
		//receiveWorld(data)
	}
}

func sendWorld() {
	startTime := time.Now()
	log.Println("Я тут 1")
	//generateWorld()
	log.Println("Я тут 2")
	saveWorldFile()
	log.Println("Я тут 3")

	odinbitPath := getOdinbitPath()
	worldPath := filepath.Join(odinbitPath, "world_server.odn")
	worldData, err := os.ReadFile(worldPath)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	log.Println("Я тут 4")
	socket.SendBinary(worldData)
	fmt.Println("World was sended")
	loadWorldFile()
	err = os.Remove(worldPath)

	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("World sended ", time.Since(startTime))
	currentScene = GAME
}

func receiveWorld(worldData []byte) {
	startTime := time.Now()
	odinbitPath := getOdinbitPath()
	worldPath := filepath.Join(odinbitPath, "world_server.odn")

	if err := os.WriteFile(worldPath, worldData, 0644); err != nil {
		fmt.Println("Error with receive world from server: ", err)
	}

	world = loadWorldFile()
	if err := os.Remove(worldPath); err != nil {
		fmt.Println("Error with remove file: ", err)
	}

	atomic.StoreInt32(&needReceiveWorld, 1)
	fmt.Println("World received ", time.Since(startTime))
	currentScene = GAME
}

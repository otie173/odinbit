package main

import (
	"fmt"
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
	worldType         int32 // send(0) or receive(1)
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

	odinbitPath := getOdinbitPath()
	worldPath := filepath.Join(odinbitPath, "world_send.odn")
	worldData, err := os.ReadFile(worldPath)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	socket.SendBinary(worldData)
	loadWorldFile()
	//err = os.Remove(worldPath)

	//if err != nil {
	//	fmt.Println("Error: ", err)
	//}
	fmt.Println("World was sended in ", time.Since(startTime))
	currentScene = GAME
}

func receiveWorld(worldData []byte) {
	atomic.StoreInt32(&worldType, 1)

	startTime := time.Now()
	odinbitPath := getOdinbitPath()
	worldPath := filepath.Join(odinbitPath, "world_receive.odn")

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

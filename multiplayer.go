package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/sacOO7/gowebsocket"
	"github.com/vmihailenco/msgpack/v5"
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

	SEND_ID
	RECEIVE_ID

	BLOCK_PACKET
	ADD_BLOCK
	REMOVE_BLOCK
)

type BlockPacket struct {
	Action   byte
	Texture  uint32
	X        float32
	Y        float32
	Passable bool
}

func authPlayer() bool {
	posturl := fmt.Sprintf("http://%s/api/auth", ipAddress)
	body := []byte(fmt.Sprintf(`{"nickname":"%s","password":"%s"}`, nickname, password))
	resp, err := http.Post(posturl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error with auth: %v\n", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	resp.Body.Close()

	return string(respBody) == "OK"
}

func connectServer(url string) {
	socket = gowebsocket.New(url)
	socket.RequestHeader.Set("Session-Nickname", nickname)

	socket.OnConnected = func(s gowebsocket.Socket) {
		atomic.StoreInt32(&connectedToServer, 1)
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
	}

	socket.OnDisconnected = func(err error, s gowebsocket.Socket) {
		atomic.StoreInt32(&connectedToServer, 0)
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
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
	socket.Connect()
}

func handleData(opcode byte, data []byte) {
	switch opcode {
	case SEND_WORLD:
		atomic.StoreInt32(&needSendWorld, 1)
	case RECEIVE_WORLD:
		worldDataMutex.Lock()
		worldData = data
		worldDataMutex.Unlock()
		atomic.StoreInt32(&needReceiveWorld, 1)
	case SEND_ID:
		sendID()
	}
}

func sendWorld() {
	odinbitPath := getOdinbitPath()
	worldPath := filepath.Join(odinbitPath, "world_send.odn")
	worldData, err := os.ReadFile(worldPath)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	data2Send := append([]byte{RECEIVE_WORLD}, worldData...)

	socket.SendBinary(data2Send)
	loadWorldFile()
	currentScene = GAME
}

func receiveWorld(worldData []byte) {
	atomic.StoreInt32(&worldType, 1)

	odinbitPath := getOdinbitPath()
	worldPath := filepath.Join(odinbitPath, "world_receive.odn")

	if err := os.WriteFile(worldPath, worldData, 0644); err != nil {
		fmt.Println("Error with receive world from server: ", err)
	}

	world = loadWorldFile()
	currentScene = GAME
}

func sendID() {
	idData, err := msgpack.Marshal(&id)
	if err != nil {
		log.Println(err)
	}

	data2Send := append([]byte{RECEIVE_ID}, idData...)
	socket.SendBinary(data2Send)
}

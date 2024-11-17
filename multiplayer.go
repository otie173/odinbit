package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync/atomic"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/gowebsocket"
	"github.com/vmihailenco/msgpack/v5"
)

var (
	socket            *gowebsocket.Socket
	connectedToServer int32
	nickname          string
	password          string
	activeInput       int // 0 - nickname, 1 - password, 2 - ipAddress
	needLoadWorld     int32
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

	PLAYER_PACKET
	RECEIVE_PLAYER_DATA
)

type ServerStatus struct {
	Address          string `json:"address"`
	IdWaiting        bool   `json:"id_waiting"`
	WorldWaiting     bool   `json:"world_waiting"`
	PlayersConnected int    `json:"players_connected"`
	MaxPlayers       int    `json:"max_players"`
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

	socket.OnConnected = func(s *gowebsocket.Socket) {
		atomic.StoreInt32(&connectedToServer, 1)
	}

	socket.OnConnectError = func(err error, socket *gowebsocket.Socket) {
	}

	socket.OnDisconnected = func(err error, s *gowebsocket.Socket) {
		atomic.StoreInt32(&connectedToServer, 0)
	}

	socket.OnTextMessage = func(message string, socket *gowebsocket.Socket) {
	}

	socket.OnBinaryMessage = func(data []byte, socket *gowebsocket.Socket) {
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

func disconnectServer() {
	socket.Close()
}

func handleData(opcode byte, data []byte) {
	switch opcode {
	case BLOCK_PACKET:
		var packet map[string]interface{}
		if err := msgpack.Unmarshal(data, &packet); err != nil {
			log.Println(err)
		}
		log.Println(packet)
		switch GetByte(packet["Action"]) {
		case ADD_BLOCK:
			addBlockNetwork(rl.Texture2D{ID: GetUint32(packet["Texture"]), Width: 10, Height: 10, Mipmaps: 1, Format: 7}, GetFloat32(packet["X"]), GetFloat32(packet["Y"]), false)
			log.Println("Сетевой игрок поставил новый блок")
		case REMOVE_BLOCK:
			removeBlockNetwork(GetFloat32(packet["X"]), GetFloat32(packet["Y"]))
			log.Println("Сетевой игрок удалил блок")
		}
	}
}

func checkStatusRest() {
	var status ServerStatus

	r, err := http.Get(fmt.Sprintf("http://%s/api/status", ipAddress))
	if err != nil {
		log.Println(err)
	}
	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(&status)

	if status.IdWaiting {
		sendIdRest()
	}
	if status.WorldWaiting {
		sendWorldRest()
	}
}

func sendIdRest() {
	idData, err := msgpack.Marshal(&id)
	if err != nil {
		log.Println(err)
	}

	if _, err := http.Post(fmt.Sprintf("http://%s/api/loadid", ipAddress), "binary", bytes.NewBuffer(idData)); err != nil {
		log.Println(err)
	}
}

func sendWorldRest() {
	odinbitPath := getOdinbitPath()
	worldPath := filepath.Join(odinbitPath, "world_send.odn")
	worldData, err := os.ReadFile(worldPath)
	if err != nil {
		log.Println(err)
	}

	if _, err := http.Post(fmt.Sprintf("http://%s/api/loadworld", ipAddress), "binary", bytes.NewBuffer(worldData)); err != nil {
		log.Fatal(err)
	}
}

func loadWorldRest() {
	atomic.StoreInt32(&worldType, 1)

	resp, err := http.Get(fmt.Sprintf("http://%s/api/getworld", ipAddress))
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	odinbitPath := getOdinbitPath()
	worldPath := filepath.Join(odinbitPath, "world_receive.odn")
	worldData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	if err := os.WriteFile(worldPath, worldData, 0644); err != nil {
		fmt.Println("Error with receive world from server: ", err)
	}
	world = loadWorldFile()
	currentScene = GAME
}

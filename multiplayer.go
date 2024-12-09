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
	needSwitchScene   int32
)

const (
	blockPacket byte = iota
	blockAdd
	blockRemove
)

type ServerStatus struct {
	Address          string `json:"address"`
	IdWaiting        bool   `json:"id_waiting"`
	WorldWaiting     bool   `json:"world_waiting"`
	PlayersConnected int    `json:"players_connected"`
	MaxPlayers       int    `json:"max_players"`
}

type PlayerMultiplayer struct {
	Nickname string
	X        float32
	Y        float32
	TargetX  float32
	TargetY  float32
	Flipped  bool
}

func authPlayer() bool {
	posturl := fmt.Sprintf("http://%s/api/v1/player/auth", ipAddress)
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
		loadPlayerRest()
	}

	socket.OnConnectError = func(err error, socket *gowebsocket.Socket) {
		log.Println(err)
	}

	socket.OnDisconnected = func(err error, s *gowebsocket.Socket) {
		atomic.StoreInt32(&needSwitchScene, 1)
		atomic.StoreInt32(&connectedToServer, 0)
		savePlayerRest()
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
	case blockPacket:
		var packet map[string]interface{}
		if err := msgpack.Unmarshal(data, &packet); err != nil {
			log.Println(err)
		}
		log.Println(packet)
		switch GetByte(packet["Action"]) {
		case blockAdd:
			addBlockNetwork(rl.Texture2D{ID: GetUint32(packet["Texture"]), Width: 10, Height: 10, Mipmaps: 1, Format: 7}, GetFloat32(packet["X"]), GetFloat32(packet["Y"]), false)
		case blockRemove:
			removeBlockNetwork(GetFloat32(packet["X"]), GetFloat32(packet["Y"]))
		}
	}
}

func checkStatusRest() {
	var status ServerStatus

	r, err := http.Get(fmt.Sprintf("http://%s/api/v1/server/status", ipAddress))
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

	if _, err := http.Post(fmt.Sprintf("http://%s/api/v1/world/loadid", ipAddress), "binary", bytes.NewBuffer(idData)); err != nil {
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

	if _, err := http.Post(fmt.Sprintf("http://%s/api/v1/world/loadworld", ipAddress), "binary", bytes.NewBuffer(worldData)); err != nil {
		log.Fatal(err)
	}
}

func loadWorldRest() {
	atomic.StoreInt32(&worldType, 1)

	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/world/getworld", ipAddress))
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

func savePlayerRest() {
	playerMu.RLock()
	playerData := Player{
		X: playerPosition.X, Y: playerPosition.Y, TargetX: targetPosition.X, TargetY: targetPosition.Y,
		Health: playerHealth, WoodCount: woodCount, StoneCount: stoneCount, MetalCount: metalCount,
		PickaxeOpen: pickaxeIsOpen, AxeOpen: axeIsOpen, ShovelOpen: shovelIsOpen,
		WallOpen: wallIsOpen, WallWindowOpen: wallWindowIsOpen, FloorOpen: floorIsOpen,
		DoorOpen: doorIsOpen, DoorOpenOpen: doorOpenIsOpen, ChestOpen: chestIsOpen,
		WallCount: wallCount, WallWindowCount: wallWindowCount, FloorCount: floorCount,
		DoorCount: doorCount, ChestCount: chestCount, DoorOpenCount: doorOpenCount,
		BigBarrelOpen: bigBarrelIsOpen, BookshelfOpen: bookshelfIsOpen, ChairOpen: chairIsOpen,
		ClosetOpen: closetIsOpen, Fence1Open: fence1IsOpen, Fence2Open: fence2IsOpen,
		Floor2Open: floor2IsOpen, Floor4Open: floor4IsOpen, LampOpen: lampIsOpen,
		ShelfOpen: shelfIsOpen, SignOpen: signIsOpen, SmallBarrelOpen: smallBarrelIsOpen,
		TableOpen: tableIsOpen, TrashOpen: trashIsOpen, LootboxOpen: lootboxIsOpen, TombstoneOpen: tombstoneIsOpen, SaplingOpen: saplingIsOpen, SeedOpen: seedIsOpen, CabbageOpen: cabbageIsOpen,
		BigBarrelCount: bigBarrelCount, BookshelfCount: bookshelfCount, ChairCount: chairCount,
		ClosetCount: closetCount, Fence1Count: fence1Count, Fence2Count: fence2Count,
		Floor2Count: floor2Count, Floor4Count: floor4Count, LampCount: lampCount,
		ShelfCount: shelfCount, SignCount: signCount, SmallBarrelCount: smallBarrelCount,
		TableCount: tableCount, TrashCount: trashCount, LootboxCount: lootboxCount, TombstoneCount: tombstoneCount, SaplingCount: saplingCount, SeedCount: seedCount, CabaggeCount: cabbageCount,
	}
	playerMu.RUnlock()

	playerDataBinary, err := msgpack.Marshal(&playerData)
	if err != nil {
		log.Println("Error with marshal player data: ", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/api/v1/player/loadpdata", ipAddress), bytes.NewBuffer(playerDataBinary))
	req.Header.Add("Session-Nickname", nickname)
	if err != nil {
		log.Println("Error with create new request to server: ", err)
	}

	_, err = client.Do(req)
	if err != nil {
		log.Println("Error with request to server: ", err)
	}
}

func loadPlayerRest() {
	var playerData Player

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/api/v1/player/getpdata", ipAddress), nil)
	req.Header.Add("Session-Nickname", nickname)
	if err != nil {
		log.Println("Error with create new request to server: ", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("Error with request to server: ", err)
	}

	bodyData, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error with read response body: ", err)
	}

	if res.StatusCode == http.StatusOK {
		if err := msgpack.Unmarshal(bodyData, &playerData); err != nil {
			log.Println("Error with unmarshal player data: ", err)
		}

		playerPosition = rl.NewVector2(playerData.X, playerData.Y)
		targetPosition = rl.NewVector2(playerData.TargetX, playerData.TargetY)
		cam.Target = playerPosition
		playerHealth = playerData.Health

		woodCount = playerData.WoodCount
		stoneCount = playerData.StoneCount
		metalCount = playerData.MetalCount

		pickaxeIsOpen = playerData.PickaxeOpen
		axeIsOpen = playerData.AxeOpen
		shovelIsOpen = playerData.ShovelOpen

		wallIsOpen = playerData.WallOpen
		wallWindowIsOpen = playerData.WallWindowOpen
		floorIsOpen = playerData.FloorOpen
		doorIsOpen = playerData.DoorOpen
		chestIsOpen = playerData.ChestOpen
		doorOpenIsOpen = playerData.DoorOpen

		wallCount = playerData.WallCount
		wallWindowCount = playerData.WallWindowCount
		floorCount = playerData.FloorCount
		doorCount = playerData.DoorCount
		doorOpenCount = playerData.DoorOpenCount
		chestCount = playerData.ChestCount

		bigBarrelIsOpen = playerData.BigBarrelOpen
		bookshelfIsOpen = playerData.BookshelfOpen
		chairIsOpen = playerData.ChairOpen
		closetIsOpen = playerData.ClosetOpen
		fence1IsOpen = playerData.Fence1Open
		fence2IsOpen = playerData.Fence2Open
		floor2IsOpen = playerData.Floor2Open
		floor4IsOpen = playerData.Floor4Open
		lampIsOpen = playerData.LampOpen
		shelfIsOpen = playerData.ShelfOpen
		signIsOpen = playerData.SignOpen
		smallBarrelIsOpen = playerData.SmallBarrelOpen
		tableIsOpen = playerData.TableOpen
		trashIsOpen = playerData.TrashOpen
		lootboxIsOpen = playerData.LootboxOpen
		tombstoneIsOpen = playerData.TombstoneOpen
		saplingIsOpen = playerData.SaplingOpen
		seedIsOpen = playerData.SeedOpen
		cabbageIsOpen = playerData.CabbageOpen

		bigBarrelCount = playerData.BigBarrelCount
		bookshelfCount = playerData.BookshelfCount
		chairCount = playerData.ChairCount
		closetCount = playerData.ClosetCount
		fence1Count = playerData.Fence1Count
		fence2Count = playerData.Fence2Count
		floor2Count = playerData.Floor2Count
		floor4Count = playerData.Floor4Count
		lampCount = playerData.LampCount
		shelfCount = playerData.ShelfCount
		signCount = playerData.SignCount
		smallBarrelCount = playerData.SmallBarrelCount
		tableCount = playerData.TableCount
		trashCount = playerData.TrashCount
		lootboxCount = playerData.LootboxCount
		tombstoneCount = playerData.TombstoneCount
		saplingCount = playerData.SaplingCount
		seedCount = playerData.SeedCount
		cabbageCount = playerData.CabaggeCount
	} else if res.StatusCode == http.StatusNotFound {
		playerPosition = rl.NewVector2(0, 0)
		targetPosition = rl.NewVector2(0, 0)
		cam.Target = playerPosition
		playerHealth = 3

		woodCount = 0
		stoneCount = 0
		metalCount = 0

		pickaxeIsOpen = false
		axeIsOpen = false
		shovelIsOpen = false

		wallIsOpen = false
		wallWindowIsOpen = false
		floorIsOpen = false
		doorIsOpen = false
		chestIsOpen = false
		doorOpenIsOpen = false

		wallCount = 0
		wallWindowCount = 0
		floorCount = 0
		doorCount = 0
		doorOpenCount = 0
		chestCount = 0

		bigBarrelIsOpen = false
		bookshelfIsOpen = false
		chairIsOpen = false
		closetIsOpen = false
		fence1IsOpen = false
		fence2IsOpen = false
		floor2IsOpen = false
		floor4IsOpen = false
		lampIsOpen = false
		shelfIsOpen = false
		signIsOpen = false
		smallBarrelIsOpen = false
		tableIsOpen = false
		trashIsOpen = false
		lootboxIsOpen = false
		tombstoneIsOpen = false
		saplingIsOpen = false
		seedIsOpen = false
		cabbageIsOpen = false

		bigBarrelCount = 0
		bookshelfCount = 0
		chairCount = 0
		closetCount = 0
		fence1Count = 0
		fence2Count = 0
		floor2Count = 0
		floor4Count = 0
		lampCount = 0
		shelfCount = 0
		signCount = 0
		smallBarrelCount = 0
		tableCount = 0
		trashCount = 0
		lootboxCount = 0
		tombstoneCount = 0
		saplingCount = 0
		seedCount = 0
		cabbageCount = 0
	}
}

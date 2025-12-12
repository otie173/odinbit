package net

import (
	"log"
	"net"
	"time"

	"github.com/kelindar/binary"
	"github.com/minio/minlz"
	"github.com/otie173/odinbit/internal/client/player"
	"github.com/otie173/odinbit/internal/protocol/packet"
)

type Module struct {
	running    bool
	connection net.Conn
	connected  bool
	ready      bool
	dispatcher *Dispatcher
	loader     *Loader
}

func New(dispatcher *Dispatcher, loader *Loader) *Module {
	return &Module{
		dispatcher: dispatcher,
		loader:     loader,
	}
}

func (m *Module) Run() {
	for m.connection == nil && !m.connected {
		log.Println(m.connection, m.connected, m.ready)
		time.Sleep(100 * time.Millisecond)
	}

	log.Println(m.connection, m.connected, m.ready)
	if !m.ready {
		m.ready = true
		log.Println("Поставил на Ready")
	}

	if err := m.listen(); err != nil {
		log.Printf("Error with read buffer from server: %v\n", err)
		m.connection.Close()
		m.connected = false
	} else {
		m.running = true
	}
}

func (m *Module) IsRunning() bool {
	return m.running
}

func (m *Module) UpdateServerPos() {
	player.PlayerMu.Lock()
	pktStructure := packet.PlayerMove{
		CurrentX: player.GamePlayer.CurrentX,
		CurrentY: player.GamePlayer.CurrentY,
		Flipped:  player.GamePlayer.Flipped,
	}
	player.PlayerMu.Unlock()

	binaryStructure, err := binary.Marshal(&pktStructure)
	if err != nil {
		log.Printf("Error! Cant marshal player move structure: %v\n", err)
	}

	pkt := packet.Packet{
		Category: packet.CategoryPlayer,
		Opcode:   packet.OpcodePlayerMove,
		Payload:  binaryStructure,
	}

	data, err := binary.Marshal(&pkt)
	if err != nil {
		log.Printf("Error! Cant marshal player move packet: %v\n", err)
	}

	compressedPkt, err := CompressPkt(data)
	if err != nil {
		log.Printf("Error! Cant compress binary player move packet: %v\n", err)
	}

	if err := m.SendData(compressedPkt); err != nil {
		log.Printf("Error! Cant write player move packet data to server: %v\n", err)
	}
}

func (m *Module) LoadTextures(addr string) ([]byte, error) {
	data, err := m.loader.LoadTextures(addr)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m *Module) IsReady() bool {
	return m.ready
}

func (m *Module) IsConnected() bool {
	return m.connected
}

func (m *Module) Connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	m.connection = conn
	m.connected = true
	go m.Run()
	log.Println(m.connection, m.connected)
	return nil
}

func (m *Module) Disconnect() {
	m.running = false
	m.connection.Close()
	m.connection = nil
	m.connected = false
	m.ready = false
}

func decompressPacket(compressedPkt []byte) ([]byte, error) {
	decompressedData, err := minlz.Decode(nil, compressedPkt)
	if err != nil {
		return nil, err
	}
	return decompressedData, nil
}

func (m *Module) Dispatch(conn *net.Conn, pktCategory packet.PacketCategory, pktOpcode packet.PacketOpcode, data []byte) {
	m.dispatcher.Dispatch(conn, pktCategory, pktOpcode, data)
}

func (m *Module) listen() error {
	buffer := make([]byte, 1024*1024) // 1MB buffer
	for m.connected && m.connection != nil {
		n, err := m.connection.Read(buffer)
		if err != nil {
			return err
		}

		compressedPkt, err := decompressPacket(buffer[:n])
		if err != nil {
			log.Printf("Error! Cant decompress packet: %v\n", err)
		}

		pkt, err := parseBinaryPacket(compressedPkt)
		if err != nil {
			log.Printf("Error with parse packet from server: %v\n", err)
		}
		m.dispatcher.Dispatch(&m.connection, pkt.Category, pkt.Opcode, pkt.Payload)
	}
	return nil
}

func (m *Module) SendData(data []byte) error {
	if m.connected == true && m.connection != nil {
		if _, err := m.connection.Write(data); err != nil {
			return err
		}
	}
	return nil
}

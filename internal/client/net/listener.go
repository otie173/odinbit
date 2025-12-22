package net

import (
	"log"
	"net"

	"github.com/kelindar/binary"
	"github.com/minio/minlz"
	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/vmihailenco/msgpack/v5"
)

type Listener struct {
	connection net.Conn
	connected  bool
	dispatcher *Dispatcher
	loader     *Loader
}

func NewListener(dispatcher *Dispatcher, loader *Loader) *Listener {
	return &Listener{
		dispatcher: dispatcher,
		loader:     loader,
	}
}

func (l *Listener) IsConnected() bool {
	return l.connected
}

func (l *Listener) Connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return err
	}

	l.connection = conn
	l.connected = true
	return nil
}

func (l *Listener) Disconnect() {
	l.connection.Close()
}

func (l *Listener) ConvertPacket(pktCategory packet.PacketCategory, pktOpcode packet.PacketOpcode, pktData any) ([]byte, error) {
	data, err := msgpack.Marshal(&pktData)
	if err != nil {
		return nil, err
	}
	packet := packet.Packet{Category: pktCategory, Opcode: pktOpcode, Payload: data}
	binaryPacket, err := msgpack.Marshal(&packet)
	if err != nil {
		return nil, err
	}
	return binaryPacket, nil
}

func parsePacket(buffer []byte) (packet.Packet, error) {
	pkt := packet.Packet{}
	if err := msgpack.Unmarshal(buffer, &pkt); err != nil {
		return packet.Packet{}, err
	}
	return pkt, nil
}

func parseBinaryPacket(buffer []byte) (packet.Packet, error) {
	pkt := packet.Packet{}
	if err := binary.Unmarshal(buffer, &pkt); err != nil {
		return packet.Packet{}, err
	}

	return pkt, nil
}

func (l *Listener) LoadTextures(addr string) ([]byte, error) {
	data, err := l.loader.LoadTextures(addr)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (l *Listener) LoadWorld(addr string) ([]byte, error) {
	data, err := l.loader.LoadWorld(addr)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (l *Listener) Dispatch(conn *net.Conn, pktCategory packet.PacketCategory, pktOpcode packet.PacketOpcode, data []byte) {
	l.dispatcher.Dispatch(conn, pktCategory, pktOpcode, data)
}

func (l *Listener) decompressPacket(compressedPkt []byte) ([]byte, error) {
	//log.Printf("Info! Compressed packet lenght: %d\n", len(compressedPkt))
	decompressedData, err := minlz.Decode(nil, compressedPkt)
	if err != nil {
		return nil, err
	}
	return decompressedData, nil
}

func (l *Listener) Handle() {
	log.Println("Началась обработка соединения")
	buffer := make([]byte, 1024*1024) // 1MB buffer

	for {
		n, err := l.connection.Read(buffer)
		if err != nil {
			log.Printf("Error with read buffer from server: %v\n", err)
			l.connection.Close()
			l.connected = false
			return
		}

		compressedPkt, err := l.decompressPacket(buffer[:n])
		if err != nil {
			log.Printf("Error! Cant decompress packet: %v\n", err)
		}

		pkt, err := parseBinaryPacket(compressedPkt)
		if err != nil {
			log.Printf("Error with parse packet from server: %v\n", err)
		}
		l.dispatcher.Dispatch(&l.connection, pkt.Category, pkt.Opcode, pkt.Payload)
	}
}

func (l *Listener) Write(data []byte) error {
	if _, err := l.connection.Write(data); err != nil {
		return err
	}
	return nil
}

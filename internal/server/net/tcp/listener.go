package tcp

import (
	"log"
	"net"

	"github.com/kelindar/binary"
	"github.com/minio/minlz"
	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/otie173/odinbit/internal/server/game/player"
)

type Listener struct {
	players    player.Storage
	dispatcher *Dispatcher
}

func NewListener(players player.Storage, dispatcher *Dispatcher) *Listener {
	return &Listener{
		players:    players,
		dispatcher: dispatcher,
	}
}

func (l *Listener) decompressPacket(compressedPkt []byte) ([]byte, error) {
	decompressedPkt, err := minlz.Decode(nil, compressedPkt)
	if err != nil {
		return nil, err
	}
	return decompressedPkt, nil
}

func (l *Listener) parsePacket(buffer []byte) (packet.Packet, error) {
	pkt := packet.Packet{}
	if err := binary.Unmarshal(buffer, &pkt); err != nil {
		return packet.Packet{}, err
	}
	return pkt, nil
}

func (l *Listener) listen(conn net.Conn) {
	defer conn.Close()

	log.Printf("New connection listening : %s\n", conn.RemoteAddr().String())
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Error! Cant read buffer from %s: %v\n", conn.RemoteAddr().String(), err)
			log.Printf("Info! Closing this connection and removing player\n")
			l.players.RemovePlayer(conn)
			conn.Close()
			return
		}

		decompressedPkt, err := l.decompressPacket(buffer[:n])
		if err != nil {
			log.Printf("Error! Cant decompress packet from %s: %v\n", conn.RemoteAddr().String(), err)
		}

		pkt, err := l.parsePacket(decompressedPkt)
		if err != nil {
			log.Printf("Error! Cant parse packet from %s: %v\n", conn.RemoteAddr().String(), err)
			conn.Close()
			return
		}
		l.dispatcher.Dispatch(conn, pkt.Category, pkt.Opcode, pkt.Payload)
	}
}

func (l *Listener) Run(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error! Cant accept new connection: %v\n", err)
			break
		}
		log.Printf("New connection accepted: %s\n", conn.RemoteAddr().String())
		go l.listen(conn)
	}
	return nil
}

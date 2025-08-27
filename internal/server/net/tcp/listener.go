package tcp

import (
	"log"
	"net"

	"github.com/otie173/odinbit/internal/protocol/packet"
	"github.com/vmihailenco/msgpack/v5"
)

type Listener struct {
	dispatcher *Dispatcher
}

func NewListener(dispatcher *Dispatcher) *Listener {
	return &Listener{
		dispatcher: dispatcher,
	}
}

func parsePacket(buffer []byte) (packet.Packet, error) {
	pkt := packet.Packet{}
	if err := msgpack.Unmarshal(buffer, &pkt); err != nil {
		return packet.Packet{}, err
	}
	return pkt, nil
}

func (h *Listener) listen(conn net.Conn) {
	defer conn.Close()

	log.Printf("New connection listening : %s\n", conn.RemoteAddr().String())
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Error! Cant read buffer from %s: %v\n", conn.RemoteAddr().String(), err)
			conn.Close()
			return
		}

		pkt, err := parsePacket(buffer[:n])
		if err != nil {
			log.Printf("Error! Cant parse packet from %s: %v\n", conn.RemoteAddr().String(), err)
			conn.Close()
			return
		}
		h.dispatcher.Dispatch(conn, pkt.Type, pkt.Payload)
	}
}

func (h *Listener) Run(addr string) error {
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
		go h.listen(conn)
	}
	return nil
}

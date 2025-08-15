package connection

import (
	"log"
	"net"
)

var (
	connection net.Conn
	connected  bool
)

func IsConnected() bool {
	return connected
}

func Connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return err
	}

	connection = conn
	connected = true
	return nil
}

func Disconnect() {
	connection.Close()
}

func Write(data []byte) {
	if _, err := connection.Write(data); err != nil {
		log.Println("Error: ", err)
	}
}

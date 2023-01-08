package core

import (
	"fmt"
	"net"
	"smq/packets"
	"smq/utilities"
)

const (
	fixedHeaderLen = 5
	connect        = 1
	publish        = 2
	subscribe      = 3
	disconnect     = 4
)

type SMQServer struct {
	Addr string
	Port string
	Ln   net.Listener
}

type fixedHeader struct {
	packetType   uint8
	remainingLen uint
}

func NewSMQServer(addr, port string) (*SMQServer, error) {
	fmt.Printf("Starting SMQ Server on '%s:%s'\n", addr, port)
	ln, err := net.Listen("tcp", addr+":"+port)

	if err != nil {
		return nil, err
	}

	server := SMQServer{
		Addr: addr,
		Port: port,
		Ln:   ln,
	}

	return &server, nil
}

func (server *SMQServer) ConnectionHandler(conn net.Conn) {

	fixedBuffer := utilities.NewBuffer(fixedHeaderLen)
	nread, err := conn.Read(fixedBuffer.Bytes)

	if err != nil {
		fmt.Println(err)
		return
	}

	if nread != fixedHeaderLen {
		fmt.Println("invalid read size for connect header")
		return
	}

	packetType, err := fixedBuffer.UnpackUint8()

	if err != nil {
		fmt.Println(err)
		return
	}

	if packetType != connect {
		fmt.Println("Invalid packet type for connection")
		return
	}

	remaining, err := fixedBuffer.UnpackUint32()

	if err != nil {
		fmt.Println(err)
		return
	}

	connectBuffer := utilities.NewBuffer(uint(remaining))
	nread, err = conn.Read(connectBuffer.Bytes)

	if err != nil {
		fmt.Println(err)
		return
	}

	if nread != int(remaining) {
		fmt.Println("invalid read size for connect header")
		return
	}

	var connect packets.Connect
	connect.Parse(*connectBuffer)
	connect.Display()
}

// END OF SOURCE

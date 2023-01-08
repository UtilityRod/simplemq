package packets

import (
	"errors"
	"fmt"
	"net"
	"smq/utilities"
)

type Connect struct {
	ProtoName    string
	ProtoVersion uint8
	ClientName   string
	Username     string
	Password     string
}

func ReadConnect(conn net.Conn, size uint32) (*Connect, error) {
	var header Connect
	var err error
	buffer := utilities.NewBuffer(size)
	nread, err := conn.Read(buffer.Bytes)

	if err != nil {
		return nil, err
	}

	if nread != int(size) {
		return nil, errors.New("invalid read size for connect header")
	}

	header.ProtoName, err = buffer.UnpackString()

	if err != nil {
		return nil, err
	}

	header.ProtoVersion, err = buffer.UnpackUint8()

	if err != nil {
		return nil, err
	}

	header.ClientName, err = buffer.UnpackString()

	if err != nil {
		return nil, err
	}

	header.Username, err = buffer.UnpackString()

	if err != nil {
		return nil, err
	}

	header.Password, err = buffer.UnpackString()

	if err != nil {
		return nil, err
	}

	return &header, err
}

func (header Connect) Display() {
	fmt.Printf("Protocol Name: %s\n", header.ProtoName)
	fmt.Printf("Protocol Version: %d\n", header.ProtoVersion)
	fmt.Printf("Client Name: %s\n", header.ClientName)
	fmt.Printf("Username: %s\n", header.Username)
	fmt.Printf("Password: %s\n", header.Password)
}

// END OF SOURCE

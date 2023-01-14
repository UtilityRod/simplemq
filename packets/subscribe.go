package packets

import (
	"errors"
	"net"
	"smq/utilities"
)

type Subscribe struct {
	Topic string
}

func GetSubscribe(conn net.Conn, size uint32) (*Subscribe, error) {
	buffer := utilities.NewBuffer(size)
	nread, err := conn.Read(buffer.Bytes)

	if err != nil {
		return nil, err
	}

	if nread != int(size) {
		return nil, errors.New("invalid read size for subscribe packet")
	}

	var subscribe Subscribe
	subscribe.Topic, err = buffer.UnpackString()

	if err != nil {
		return nil, err
	}

	return &subscribe, nil
}

// END OF SOURCE

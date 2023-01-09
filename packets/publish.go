package packets

import (
	"errors"
	"net"
	"smq/utilities"
)

type Publish struct {
	Topic string
	Value string
}

func GetPublish(conn net.Conn, size uint32) (*Publish, error) {
	buffer := utilities.NewBuffer(size)
	nread, err := conn.Read(buffer.Bytes)

	if err != nil {
		return nil, err
	}

	if nread != int(size) {
		return nil, errors.New("invalid read size for publish packet")
	}

	topic, err := buffer.UnpackString()

	if err != nil {
		return nil, err
	}

	value, err := buffer.UnpackString()

	if err != nil {
		return nil, err
	}

	Publish := Publish{
		Topic: topic,
		Value: value,
	}

	return &Publish, nil
}

func (publish *Publish) Send(conn net.Conn) error {
	var buffer utilities.Buffer
	buffer.Pack(uint8(2))
	buffer.Pack(uint32(len(publish.Topic) + len(publish.Value) + 4))
	buffer.Pack(publish.Topic)
	buffer.Pack(publish.Value)
	nwrote, err := conn.Write(buffer.Bytes)

	if err != nil {
		return err
	}

	if nwrote != len(buffer.Bytes) {
		return errors.New("invalid write size for publish packet")
	}

	return nil
}

// END OF SOURCE

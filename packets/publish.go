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

func (publish *Publish) Set(buffer *utilities.Buffer) error {

	topic, err := buffer.UnpackString()

	if err != nil {
		return err
	}

	value, err := buffer.UnpackString()

	if err != nil {
		return err
	}

	publish.Topic = topic
	publish.Value = value

	return nil
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

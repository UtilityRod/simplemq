package packets

import (
	"errors"
	"net"
	"smq/utilities"
)

type Packet interface {
	Set(*utilities.Buffer) error
}

const (
	FixedHeaderLen = 5
)

type FixedHeader struct {
	PacketType   uint8
	RemainingLen uint32
}

func ReadFixedHeader(conn net.Conn) (*FixedHeader, error) {
	fixedBuffer := utilities.NewBuffer(FixedHeaderLen)
	nread, err := conn.Read(fixedBuffer.Bytes)

	if err != nil {
		return nil, err
	}

	if nread != FixedHeaderLen {
		return nil, errors.New("invalid read size for fixed header")
	}

	packetType, err := fixedBuffer.UnpackUint8()

	if err != nil {
		return nil, err
	}

	remaining, err := fixedBuffer.UnpackUint32()

	if err != nil {
		return nil, err
	}

	header := FixedHeader{
		PacketType:   packetType,
		RemainingLen: remaining,
	}

	return &header, nil
}

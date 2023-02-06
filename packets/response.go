package packets

import (
	"smq/utilities"
)

const (
	SUCCESS         = 0x01
	INVALID_AUTH    = 0x02
	GENERAL_FAILURE = 0xFF
)

const (
	RESPONSE_SZ = uint32(1)
)

type Response struct {
	ByteString []byte
}

func NewResponse(packet_type uint8, response_code uint8) *Response {
	buffer := utilities.NewBuffer(0)
	buffer.Pack(packet_type)
	buffer.Pack(RESPONSE_SZ)
	buffer.Pack(response_code)
	response := Response{
		ByteString: buffer.Bytes,
	}
	return &response
}

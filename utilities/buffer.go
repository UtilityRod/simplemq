package utilities

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

type Buffer struct {
	Bytes []byte
}

func NewBuffer(size uint32) *Buffer {
	var buffer Buffer
	buffer.Bytes = make([]byte, size)
	return &buffer
}

func (buffer *Buffer) Pack(value any) {
	var bytes []byte
	switch value.(type) {
	case uint8:
		bytes = make([]byte, 1)
		bytes[0] = byte(value.(uint8))
	case uint16:
		bytes = make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, value.(uint16))
	case uint32:
		bytes = make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, value.(uint32))
	case uint64:
		bytes = make([]byte, 8)
		binary.BigEndian.PutUint64(bytes, value.(uint64))
	case string:
		length := uint16(len(value.(string)))
		bytes = make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, length)
		strBytes := []byte(value.(string))
		bytes = append(bytes, strBytes...)
	default:
		fmt.Println("invalid variable type for pack", reflect.TypeOf(value))
	}

	buffer.Bytes = append(buffer.Bytes, bytes...)
}

func (buffer *Buffer) UnpackUint8() (uint8, error) {

	if len(buffer.Bytes) < 1 {
		return 0, errors.New("invalid size of buffer for uint8 unpack")
	}

	returnValue := uint8(buffer.Bytes[0])
	buffer.Bytes = buffer.Bytes[1:]
	return returnValue, nil
}

func (buffer *Buffer) UnpackUint16() (uint16, error) {

	if len(buffer.Bytes) < 2 {
		return 0, errors.New("invalid size of buffer for uint16 unpack")
	}

	returnValue := binary.BigEndian.Uint16(buffer.Bytes[:2])
	buffer.Bytes = buffer.Bytes[2:]
	return returnValue, nil
}

func (buffer *Buffer) UnpackUint32() (uint32, error) {

	if len(buffer.Bytes) < 4 {
		return 0, errors.New("invalid size of buffer for uint32 unpack")
	}

	returnValue := binary.BigEndian.Uint32(buffer.Bytes[:4])
	buffer.Bytes = buffer.Bytes[4:]
	return returnValue, nil
}

func (buffer *Buffer) UnpackUint64() (uint64, error) {

	if len(buffer.Bytes) < 8 {
		return 0, errors.New("invalid size of buffer for uint64 unpack")
	}

	returnValue := binary.BigEndian.Uint64(buffer.Bytes[:8])
	buffer.Bytes = buffer.Bytes[8:]
	return returnValue, nil
}

func (buffer *Buffer) UnpackString() (string, error) {
	if len(buffer.Bytes) < 2 {
		return "", errors.New("invalid size of buffer for string unpack")
	}

	stringLength := binary.BigEndian.Uint16(buffer.Bytes[:2])
	buffer.Bytes = buffer.Bytes[2:]
	if len(buffer.Bytes) < int(stringLength) {
		return "", errors.New("invalid size of buffer for string unpack")
	}

	returnValue := string(buffer.Bytes[:stringLength])
	buffer.Bytes = buffer.Bytes[stringLength:]
	return returnValue, nil
}

// END OF SOURCE

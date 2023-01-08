package packets

import (
	"fmt"
	"smq/utilities"
)

type Connect struct {
	ProtoName    string
	ProtoVersion uint8
	ClientName   string
	Username     string
	Password     string
}

func (header *Connect) Parse(buffer utilities.Buffer) (err error) {
	header.ProtoName, err = buffer.UnpackString()

	if err != nil {
		return err
	}

	header.ProtoVersion, err = buffer.UnpackUint8()

	if err != nil {
		return err
	}

	header.ClientName, err = buffer.UnpackString()

	if err != nil {
		return err
	}

	header.Username, err = buffer.UnpackString()

	if err != nil {
		return err
	}

	header.Password, err = buffer.UnpackString()

	if err != nil {
		return err
	}

	return err
}

func (header Connect) Display() {
	fmt.Printf("Protocol Name: %s\n", header.ProtoName)
	fmt.Printf("Protocol Version: %d\n", header.ProtoVersion)
	fmt.Printf("Client Name: %s\n", header.ClientName)
	fmt.Printf("Username: %s\n", header.Username)
	fmt.Printf("Password: %s\n", header.Password)
}

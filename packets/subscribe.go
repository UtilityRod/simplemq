package packets

import (
	"smq/utilities"
)

type Subscribe struct {
	Topic string
}

func (subscribe *Subscribe) Set(buffer *utilities.Buffer) error {
	topic, err := buffer.UnpackString()

	if err != nil {
		return err
	}

	subscribe.Topic = topic
	return nil
}

// END OF SOURCE

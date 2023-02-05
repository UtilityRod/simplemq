package packets

import (
	"smq/utilities"
)

type Unsubscribe struct {
	Topic string
}

func (unsub *Unsubscribe) Set(buffer *utilities.Buffer) error {
	topic, err := buffer.UnpackString()

	if err != nil {
		return err
	}

	unsub.Topic = topic
	return nil
}

// END OF SOURCE

package packets

import "smq/utilities"

type Disconnect struct {
	Flag bool
}

func (disconnect *Disconnect) Set(buffer *utilities.Buffer) error {
	disconnect.Flag = true
	return nil
}

// END OF SOURCE

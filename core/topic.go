package core

import (
	"errors"
	"fmt"
	"sync"
)

type Topic struct {
	TopicMutex sync.Mutex
	Clients    map[string]string
}

func (server *SMQServer) RegisterTopic(name string) error {
	if _, ok := server.Topics[name]; ok {
		errStr := fmt.Sprintf("could not register topic '%s': already exists", name)
		return errors.New(errStr)
	}

	topic := Topic{
		Clients: make(map[string]string, 0),
	}

	server.Topics[name] = &topic
	fmt.Printf("Topic registered: '%s'\n", name)
	return nil
}

// END OF SOURCE

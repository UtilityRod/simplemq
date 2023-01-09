package core

import (
	"fmt"
	"smq/packets"
	"sync"
)

const (
	Quit = 0xFF
)

type Event struct {
	EventType int
	Payload   interface{}
}

type EventHandler struct {
	Events        []Event
	EventChannel  chan Event
	Quit          bool
	EventSignaler sync.Cond
	Wg            sync.WaitGroup
	Topics        map[string]*Topic
}

func NewEventHandler() *EventHandler {
	handler := EventHandler{
		EventChannel: make(chan Event),
		Topics:       make(map[string]*Topic),
	}

	var mutex sync.Mutex
	mutex.Lock()
	handler.EventSignaler = *sync.NewCond(&mutex)

	return &handler
}

func (handler *EventHandler) Start() {
	fmt.Println("Starting event handler for SMQ server")
	handler.Wg.Add(2)
	go handler.EventReader()
	go handler.HandleEvents()
}

func (handler *EventHandler) Stop() {
	handler.Quit = true
	event := Event{
		EventType: Quit,
	}
	handler.EventChannel <- event
	handler.Wg.Wait()
}

func (handler *EventHandler) HandleEvents() {
	defer handler.Wg.Done()
	var event Event

	for handler.Quit != true {
		handler.EventSignaler.Wait()

		if len(handler.Events) <= 0 {
			continue
		}

		event, handler.Events = handler.Events[0], handler.Events[1:]

		switch event.EventType {
		case PublishType:
			message := event.Payload.(*packets.Publish)
			topic := handler.Topics[message.Topic]
			for _, client := range topic.Clients {
				err := message.Send(client.Conn)

				if err != nil {
					fmt.Println(err)
				}
			}
		case Quit:
			fmt.Println("Shutting down event handler")
		default:
			fmt.Printf("unknown event type: %d\n", event.EventType)
		}
	}
}

func (handler *EventHandler) EventReader() {
	defer handler.Wg.Done()

	for handler.Quit != true {
		select {
		case event := <-handler.EventChannel:
			handler.Events = append(handler.Events, event)
			handler.EventSignaler.Signal()
		}
	}
}

// END OF SOURCE

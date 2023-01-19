package core

import (
	"sync"
)

const (
	Quit = 0xFF
)

type EventFunc func(payload interface{})

type Event struct {
	Payload interface{}
	Func    EventFunc
}

type EventHandler struct {
	Events        []Event
	EventChannel  chan Event
	Quit          bool
	EventSignaler sync.Cond
	Wg            sync.WaitGroup
}

func NewEventHandler() *EventHandler {
	handler := EventHandler{
		EventChannel: make(chan Event),
	}

	var mutex sync.Mutex
	mutex.Lock()
	handler.EventSignaler = *sync.NewCond(&mutex)

	return &handler
}

func (handler *EventHandler) Start() {
	handler.Wg.Add(2)
	go handler.EventReader()
	go handler.HandleEvents()
}

func (handler *EventHandler) Stop() {
	handler.Quit = true
	event := Event{
		Func: nil,
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

		if event.Func != nil {
			event.Func(event.Payload)
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

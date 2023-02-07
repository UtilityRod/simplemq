package core

import (
	"fmt"
	"smq/packets"
)

func (server *SMQServer) publishHandler(inter interface{}) {
	fmt.Println("Publish")
	smqPayload := inter.(SMQPayload)
	publish := smqPayload.Payload.(*packets.Publish)
	topic := server.Topics[publish.Topic]
	for _, client := range topic.Clients {
		err := publish.Send(client.Conn)

		if err != nil {
			fmt.Println(err)
		}
	}
}

func (server *SMQServer) subscribeHandler(inter interface{}) {
	fmt.Println("Subscribe")
	smqPayload := inter.(SMQPayload)
	subscribe := smqPayload.Payload.(*packets.Subscribe)

	// Get the topic being subscribed too
	topicStr := subscribe.Topic
	topic := server.Topics[topicStr]

	if topic == nil {
		response := packets.NewResponse(SubackType, packets.INVALID_TOPIC)
		smqPayload.Client.Conn.Write(response.ByteString)
		return
	}

	// Check to see if client already subscribed to topic
	clientName := smqPayload.Client.ClientName
	if _, ok := topic.Clients[clientName]; ok {
		// Client already subscribed
		response := packets.NewResponse(SubackType, packets.REDUNDANT_SUB)
		smqPayload.Client.Conn.Write(response.ByteString)
		return
	}

	client := server.Clients[clientName]
	// Update topic to contain client information
	topic.Clients[clientName] = client
	// Update client's subscribed topics
	client.Topics[topicStr] = topic
	response := packets.NewResponse(SubackType, packets.SUCCESS)
	smqPayload.Client.Conn.Write(response.ByteString)
}

func (server *SMQServer) unsubscribeHandler(inter interface{}) {
	fmt.Println("Unsubscribe")
	payload := inter.(SMQPayload)
	client := payload.Client
	topicStr := payload.Payload.(*packets.Unsubscribe).Topic
	delete(client.Topics, topicStr)
	topic := server.Topics[topicStr]
	delete(topic.Clients, client.ClientName)
}

func (server *SMQServer) disconnectHandler(inter interface{}) {
	fmt.Println("Disconnect")
	smqPayload := inter.(SMQPayload)
	client := smqPayload.Client
	client.Quit = true

	// For every topic the client is subscribed too
	for _, topic := range client.Topics {
		delete(topic.Clients, client.ClientName)
	}

	delete(server.Clients, client.ClientName)
}

// END OF SOURCE

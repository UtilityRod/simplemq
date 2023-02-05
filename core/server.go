package core

import (
	"fmt"
	"net"
	"smq/packets"
	"smq/utilities"
)

const (
	ConnectType     = 1
	PublishType     = 2
	SubscribeType   = 3
	UnsubscribeType = 4
	DisconnectType  = 5
	DefaultUser     = "admin"
	DefaultPass     = "password"
)

type SMQClient struct {
	Conn       net.Conn
	ClientName string
	Topics     map[string]*Topic
	Quit       bool
}

type SMQServer struct {
	Addr     string
	Port     string
	Ln       net.Listener
	Handler  *EventHandler
	Topics   map[string]*Topic
	Clients  map[string]*SMQClient
	Auth     *utilities.Authenticator
	Payloads map[uint8]packets.Packet
	Handlers map[uint8]EventFunc
}

type SMQPayload struct {
	Payload packets.Packet
	Client  *SMQClient
}

func NewSMQServer(addr, port string) (*SMQServer, error) {
	fmt.Printf("Starting SMQ Server on '%s:%s'\n", addr, port)
	ln, err := net.Listen("tcp", addr+":"+port)

	if err != nil {
		return nil, err
	}

	server := SMQServer{
		Addr:    addr,
		Port:    port,
		Ln:      ln,
		Topics:  make(map[string]*Topic),
		Clients: make(map[string]*SMQClient),
		Auth:    utilities.NewAuthenticator(DefaultUser, DefaultPass),
	}

	server.Handlers = map[uint8]EventFunc{
		PublishType:     server.publishHandler,
		SubscribeType:   server.subscribeHandler,
		UnsubscribeType: server.unsubscribeHandler,
		DisconnectType:  server.disconnectHandler,
	}

	server.Payloads = map[uint8]packets.Packet{
		PublishType:     &packets.Publish{},
		SubscribeType:   &packets.Subscribe{},
		UnsubscribeType: &packets.Unsubscribe{},
		DisconnectType:  &packets.Disconnect{},
	}

	server.Handler = NewEventHandler()
	server.Handler.Start()

	return &server, nil
}

func NewSMQClient(name string, conn net.Conn) *SMQClient {
	client := SMQClient{
		Conn:       conn,
		ClientName: name,
		Topics:     make(map[string]*Topic),
	}

	return &client
}

func (server *SMQServer) ConnectionHandler(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Accepted connection for %s\n", conn.RemoteAddr())
	var err error
	fixedHeader, err := packets.ReadFixedHeader(conn)

	if err != nil {
		fmt.Println(err)
		return
	}

	connect, err := packets.ReadConnect(conn, fixedHeader.RemainingLen)

	if err != nil {
		fmt.Println(err)
		return
	}

	if server.Auth.Authenticate(connect.Username, connect.Password) != true {
		fmt.Printf("incorrect username or password from: %s\n", conn.RemoteAddr())
		return
	}

	// Create new client
	server.Clients[connect.ClientName] = NewSMQClient(connect.ClientName, conn)
	client := server.Clients[connect.ClientName]
	var event Event

	for client.Quit != true {
		fixedHeader, err := packets.ReadFixedHeader(client.Conn)

		if err != nil {
			fmt.Println(err)
			break
		}

		buffer := utilities.NewBuffer(fixedHeader.RemainingLen)
		nread, err := client.Conn.Read(buffer.Bytes)

		if err != nil {
			fmt.Println(err)
			continue
		}

		if nread != int(fixedHeader.RemainingLen) {
			fmt.Println("invalid read size for remaining bytes")
			continue
		}
		// Set function to handle payload
		event.Func = server.Handlers[fixedHeader.PacketType]
		// Set the payload for the packet
		payload := SMQPayload{
			Payload: server.Payloads[fixedHeader.PacketType],
			Client:  client,
		}

		err = payload.Payload.Set(buffer)
		event.Payload = payload

		if err != nil {
			fmt.Println(err)
		} else {
			server.Handler.EventChannel <- event
		}

		event.Func = nil
		event.Payload = nil
	}

	delete(server.Clients, client.ClientName)
	fmt.Printf("Closed connection for %s\n", conn.RemoteAddr())
}

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
	// Check to see if client already subscribed to topic
	clientName := smqPayload.Client.ClientName
	if _, ok := topic.Clients[clientName]; ok {
		// Client already subscribed
		return
	}

	client := server.Clients[clientName]
	// Update topic to contain client information
	topic.Clients[clientName] = client
	// Update client's subscribed topics
	client.Topics[topicStr] = topic
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
}

// END OF SOURCE
